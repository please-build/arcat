// Package ar provides an ar file archiver.
package ar

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/please-build/ar"
	"gopkg.in/op/go-logging.v1"
)

var log = logging.MustGetLogger("ar")

// mtime is the time we attach for the modification time of all files.
var mtime = time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)

// Create creates a new ar archive from the given sources.
// If combine is true they are treated as existing ar files and combined.
// If rename is true the srcs are renamed as gcc would (i.e. the extension is replaced by .o).
func Create(srcs []string, out string, combine, rename bool) error {
	// Rename the sources as gcc would.
	if rename {
		for i, src := range srcs {
			src = path.Base(src)
			if ext := path.Ext(src); ext != "" {
				src = src[:len(src)-len(ext)] + ".o"
			}
			srcs[i] = src
			log.Debug("renamed ar source to %s", src)
		}
	}

	log.Debug("Writing ar to %s", out)
	f, err := os.Create(out)
	if err != nil {
		return fmt.Errorf("create %s: %w", out, err)
	}
	defer f.Close()
	bw := bufio.NewWriter(f)
	defer bw.Flush()
	var w *ar.Writer
	// Write BSD-style names on OSX, GNU-style ones on Linux
	if runtime.GOOS == "darwin" {
		w = ar.NewWriter(bw, ar.BSD)
	} else {
		w = ar.NewWriter(bw, ar.GNU)
		allSrcs, err := allSourceNames(srcs, combine)
		if err != nil {
			return fmt.Errorf("find source file names: %w", err)
		}
		if err := w.WriteStringTable(allSrcs); err != nil {
			return fmt.Errorf("write ar string table: %w", err)
		}
	}
	for _, src := range srcs {
		log.Debug("ar source file: %s", src)
		f, err := os.Open(src)
		if err != nil {
			return fmt.Errorf("open %s: %w", src, err)
		}
		if combine {
			// Read archive & write its contents in
			r, err := ar.NewReader(f)
			if err != nil {
				return fmt.Errorf("read %s: %w", src, err)
			}
			for {
				hdr, err := r.Next()
				if err != nil {
					if err == io.EOF {
						break
					}
					return fmt.Errorf("read next file header: %w", err)
				}
				// Zero things out
				hdr.ModTime = mtime
				hdr.Uid = 0
				hdr.Gid = 0
				// Fix weird bug about octal numbers (looks like we're prepending 100 multiple times)
				hdr.Mode &= ^0100000
				log.Debug("copying '%s' in from %s, mode %x", hdr.Name, src, hdr.Mode)
				if err := w.WriteHeader(hdr); err != nil {
					return fmt.Errorf("write ar header: %w", err)
				} else if _, err = io.Copy(w, r); err != nil {
					return fmt.Errorf("write ar data: %w", err)
				}
			}
		} else {
			// Write in individual file
			info, err := os.Lstat(src)
			if err != nil {
				return fmt.Errorf("lstat %s: %w", src, err)
			}
			hdr := &ar.Header{
				Name:    src,
				ModTime: mtime,
				Mode:    int64(info.Mode()),
				Size:    info.Size(),
			}
			log.Debug("creating file %s", hdr.Name)
			if err := w.WriteHeader(hdr); err != nil {
				return fmt.Errorf("write ar header: %w", err)
			} else if _, err := io.Copy(w, f); err != nil {
				return fmt.Errorf("write ar data: %w", err)
			}
		}
		f.Close()
	}
	return nil
}

// Find finds all the .a files under the current directory and returns their names.
func Find() ([]string, error) {
	ret := []string{}
	return ret, filepath.WalkDir(".", func(name string, d fs.DirEntry, err error) error {
		if strings.HasSuffix(name, ".a") && !d.IsDir() {
			ret = append(ret, name)
		}
		return nil
	})
}

// allSourceNames returns the name of all source files that we will add to the archive.
func allSourceNames(srcs []string, combine bool) ([]string, error) {
	if !combine {
		return srcs, nil
	}
	ret := []string{}
	for _, src := range srcs {
		f, err := os.Open(src)
		if err == nil {
			r, err := ar.NewReader(f)
			if err != nil {
				return nil, err
			}
			for {
				hdr, err := r.Next()
				if err != nil {
					break
				}
				ret = append(ret, hdr.Name)
			}
		}
	}
	return ret, nil
}
