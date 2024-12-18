// Package main implements jarcat, a program to efficiently concatenate .zip files.
// Originally this was pretty simple and that was all it could do, over time it's
// gained a bunch more features on a more or less as needed basis.
//
// It's now used for most general-purpose zip and tar manipulation in Please, since
// the standard tools either differ between implementations (e.g. GNU tar vs. BSD tar)
// or introduce indeterminacy, often in regard to timestamps.
package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"slices"

	"github.com/peterebden/go-cli-init/v5/flags"
	"github.com/peterebden/go-cli-init/v5/logging"

	"github.com/please-build/arcat/ar"
	"github.com/please-build/arcat/tar"
	"github.com/please-build/arcat/unzip"
	"github.com/please-build/arcat/zip"
)

var javaExcludePrefixes = []string{
	"META-INF/LICENSE", "META-INF/NOTICE", "META-INF/maven/*", "META-INF/MANIFEST.MF",
	// Un-sign all jars by default, after concatenation the signatures will no longer be valid.
	"META-INF/*.SF", "META-INF/*.RSA", "META-INF/*.LIST",
}

var log = logging.MustGetLogger()

func must(err error) {
	if err != nil {
		log.Fatalf("%s", err)
	}
}

var opts = struct {
	Usage     string
	Verbosity logging.Verbosity `short:"v" long:"verbosity" default:"warning" description:"Verbosity of output (higher number = more output)"`

	Zip struct {
		In                    flags.StdinStrings  `short:"i" long:"input" description:"Input directory" required:"true"`
		Out                   string            `short:"o" long:"output" env:"OUT" description:"Output filename" required:"true"`
		Suffix                []string          `short:"s" long:"suffix" default:".jar" description:"Suffix of files to include"`
		ExcludeSuffix         []string          `short:"e" long:"exclude_suffix" description:"Suffix of files to exclude"`
		ExcludeTools          []string          `long:"exclude_tools" env:"TOOLS" env-delim:" " description:"Tools to exclude from the generated zipfile"`
		ExcludeInternalPrefix []string          `short:"x" long:"exclude_internal_prefix" description:"Prefix of files to exclude"`
		IncludeInternalPrefix []string          `short:"t" long:"include_internal_prefix" description:"Prefix of files to include"`
		StripPrefix           string            `long:"strip_prefix" description:"Prefix to strip off file names"`
		Preamble              string            `short:"p" long:"preamble" description:"Prepend this string to generated zip file"`
		PreambleFrom          string            `long:"preamble_from" description:"Prepend non-zip leading data from this file to generated zip file"`
		PreambleFile          string            `long:"preamble_file" description:"Prepend contents of this file to generated zip file"`
		MainClass             string            `short:"m" long:"main_class" description:"Write a Java manifest file containing the given main class."`
		Manifest              string            `long:"manifest" description:"Use the given file as a Java manifest"`
		ExcludeJavaPrefixes   bool              `short:"j" long:"exclude_java_prefixes" description:"Use default Java exclusions"`
		Strict                bool              `long:"strict" description:"Disallow duplicate files"`
		IncludeOther          bool              `long:"include_other" description:"Add files that are not jar files as well"`
		AddInitPy             bool              `long:"add_init_py" description:"Adds __init__.py files to all directories"`
		StripPy               bool              `long:"strip_py" description:"Strips .py files when there is a corresponding .pyc"`
		DumbMode              bool              `short:"d" long:"dumb" description:"Dumb mode, an alias for --suffix='' --exclude_suffix='' --include_other"`
		NoDirEntries          bool              `short:"n" long:"nodir_entries" description:"Don't add directory entries to zip"`
		RenameDirs            map[string]string `short:"r" long:"rename_dir" description:"Rename directories within zip file"`
		StoreSuffix           []string          `short:"u" long:"store_suffix" description:"Suffix of filenames to store instead of deflate (i.e. without compression). Note that this only affects files found with --include_other."`
		Prefix                string            `long:"prefix" description:"Prefix all entries with this directory name."`
	} `command:"zip" alias:"z" description:"Writes an output zipfile"`

	Tar struct {
		Gzip        bool     `short:"z" long:"gzip" description:"Apply gzip compression to the tar file."`
		Xzip        bool     `short:"x" long:"xzip" description:"Apply gzip compression to the tar file."`
		Flatten     bool     `long:"flatten" description:"Whether to flatten internal tar structure."`
		Out         string   `short:"o" long:"output" env:"OUT" description:"Output filename" required:"true"`
		Srcs        []string `long:"srcs" env:"SRCS" env-delim:" " description:"Source files for the tarball."`
		Prefix      string   `long:"prefix" description:"Prefix all entries with this directory name."`
		StripPrefix string   `long:"strip-prefix" description:"Prefix to remove from files. Only affects non-flattened tarballs."`
	} `command:"tar" alias:"t" description:"Builds a tarball instead of a zipfile."`

	Extract struct {
		Args struct {
			In   string `positional-arg-name:"input" required:"true" description:"Input archive"`
			File string `positional-arg-name:"file" description:"File to extract"`
		} `positional-args:"true"`
		StripPrefix string `short:"s" long:"strip_prefix" description:"Strip this prefix from extracted files"`
		OutDir      string `short:"o" long:"out" description:"Output directory"`
		Out         string `long:"out_file" hidden:"true" env:"OUT"`
	} `command:"extract" alias:"unzip" alias:"u" alias:"x" description:"Extracts a zipfile or tarball"`

	Ar struct {
		Srcs    []string `long:"srcs" env:"SRCS_SRCS" env-delim:" " description:"Source .ar files to combine"`
		Out     string   `long:"out" env:"OUT" description:"Output filename"`
		Rename  bool     `short:"r" long:"rename" description:"Rename source files as gcc would (i.e. change extension to .o)"`
		Combine bool     `short:"c" long:"combine" description:"Treat source files as .a files and combines them"`
		Find    bool     `short:"f" long:"find" description:"Find all .a files under the current directory & combine those (implies --combine)"`
	} `command:"ar" alias:"a" description:"Creates a new ar archive."`
}{
	Usage: `
arcat is a tool for creating, extracting and concatenating archives..

Its original and most unique feature is performing efficient concatenation of .jar files
when compiling Java code. This is possible with zip files because each file is compressed
individually so it's possible to combine them without decompressing and recompressing each one.

It now has a number of other features to help in compilation and serves as a general-purpose
zip manipulator for Please. To help us maintain reproduceability of builds it is able to strip
timestamps from files, and also has a bunch of Python-specific functionality to help with .pex files.

Typically you don't invoke this directly, Please will run it when individual rules need it.
You're welcome to use it separately if you find it useful, although be aware that we do not
aim to maintain compatibility very strongly.

Any apparent relationship between the name of this tool and bonsai kittens is completely coincidental.
`,
}

func main() {
	command := flags.ParseFlagsOrDie("arcat", &opts, nil)
	if opts.Zip.DumbMode {
		opts.Zip.Suffix = nil
		opts.Zip.ExcludeSuffix = nil
		opts.Zip.IncludeOther = true
	}
	logging.InitLogging(opts.Verbosity)

	if command == "tar" {
		if opts.Tar.Xzip && opts.Tar.Gzip {
			log.Fatalf("Can't pass --xzip and --gzip simultaneously")
		}
		if err := tar.Write(
			opts.Tar.Out, opts.Tar.Srcs, opts.Tar.Prefix,
			opts.Tar.Gzip, opts.Tar.Xzip, opts.Tar.Flatten, opts.Tar.StripPrefix); err != nil {
			log.Fatalf("Error writing tarball: %s\n", err)
		}
		os.Exit(0)
	} else if command == "extract" {
		// This comes up if we're in the root directory. Ignore it.
		if opts.Extract.StripPrefix == "." {
			opts.Extract.StripPrefix = ""
		}
		if opts.Extract.Args.File != "" && opts.Extract.OutDir == "" {
			opts.Extract.OutDir = opts.Extract.Out
		}
		if err := unzip.Extract(opts.Extract.Args.In, opts.Extract.OutDir, opts.Extract.Args.File, opts.Extract.StripPrefix); err != nil {
			log.Fatalf("Error extracting zipfile: %s", err)
		}
		os.Exit(0)
	} else if command == "ar" {
		if opts.Ar.Find {
			srcs, err := ar.Find()
			if err != nil {
				log.Fatalf("%s", err)
			}
			opts.Ar.Srcs = srcs
			opts.Ar.Combine = true
		}
		if err := ar.Create(opts.Ar.Srcs, opts.Ar.Out, opts.Ar.Combine, opts.Ar.Rename); err != nil {
			log.Fatalf("Error combining archives: %s", err)
		}
		os.Exit(0)
	}

	if len(slices.DeleteFunc([]string{
		opts.Zip.Preamble,
		opts.Zip.PreambleFrom,
		opts.Zip.PreambleFile,
	}, func(s string) bool { return s == "" })) > 1 {
		log.Fatal("Only one of --preamble, --preamble_from or --preamble_file may be specified.")
	}

	if opts.Zip.ExcludeJavaPrefixes {
		opts.Zip.ExcludeInternalPrefix = javaExcludePrefixes
	}

	tempFile, err := ioutil.TempFile(".", "jarcat-")
	must(err)
	filename := tempFile.Name()

	f := zip.NewFile(filename, opts.Zip.Strict)
	f.RenameDirs = opts.Zip.RenameDirs
	f.Include = opts.Zip.IncludeInternalPrefix
	f.Exclude = opts.Zip.ExcludeInternalPrefix
	f.StripPrefix = opts.Zip.StripPrefix
	f.Suffix = opts.Zip.Suffix
	f.ExcludeSuffix = opts.Zip.ExcludeSuffix
	f.StoreSuffix = opts.Zip.StoreSuffix
	f.IncludeOther = opts.Zip.IncludeOther
	f.AddInitPy = opts.Zip.AddInitPy
	f.StripPy = opts.Zip.StripPy
	f.DirEntries = !opts.Zip.NoDirEntries
	f.Prefix = opts.Zip.Prefix

	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to determine working directory: %s", err)
	}
	for _, excl := range opts.Zip.ExcludeTools {
		f.Exclude = append(f.Exclude, excl)
		if filepath.IsAbs(excl) {
			if rel, err := filepath.Rel(wd, excl); err == nil {
				f.Exclude = append(f.Exclude, rel)
			}
		}
	}
	// Never descend into the _please dir
	f.Exclude = append(f.Exclude, "_please")

	if opts.Zip.Preamble != "" {
		log.Debug("Adding preamble from command line option")
		must(f.WritePreambleBytes([]byte(opts.Zip.Preamble + "\n")))
	} else if opts.Zip.PreambleFrom != "" {
		log.Debugf("Adding preamble from non-zip data in %s", opts.Zip.PreambleFrom)
		pr, err := zip.Preamble(opts.Zip.PreambleFrom)
		must(err)
		defer pr.Close()
		must(f.WritePreamble(pr))
	} else if opts.Zip.PreambleFile != "" {
		log.Debugf("Adding preamble from %s", opts.Zip.PreambleFile)
		pf, err := os.Open(opts.Zip.PreambleFile)
		must(err)
		defer pf.Close()
		must(f.WritePreamble(pf))
	}
	if opts.Zip.MainClass != "" {
		must(f.AddManifest(opts.Zip.MainClass))
	}
	if opts.Zip.Manifest != "" {
		b, err := ioutil.ReadFile(opts.Zip.Manifest)
		must(err)
		must(f.WriteFile("META-INF/MANIFEST.MF", b, 0644))
	}
	for _, filename := range opts.Zip.In.Get() {
		must(f.AddFiles(filename))
	}

	f.Close()

	if err := os.Rename(filename, opts.Zip.Out); err != nil {
		panic(fmt.Sprintf("Failed to rename output file: %s", err))
	}
}
