package zip

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"reflect"
)

// PreambleReader is an io.ReadSeekCloser that facilitates the reading of the data preceding the zip
// data in a zip file.
type PreambleReader struct {
	rc io.ReadCloser
	sr *io.SectionReader
}

func (r *PreambleReader) Read(p []byte) (int, error) {
	return r.sr.Read(p)
}

func (r *PreambleReader) Seek(offset int64, whence int) (int64, error) {
	return r.sr.Seek(offset, whence)
}

func (r *PreambleReader) ReadAt(p []byte, off int64) (int, error) {
	return r.sr.ReadAt(p, off)
}

func (r *PreambleReader) Size() int64 {
	return r.sr.Size()
}

func (r *PreambleReader) Close() error {
	return r.rc.Close()
}

// Preamble returns a PreambleReader enabling the preamble to be read from the given zip file.
func Preamble(path string) (*PreambleReader, error) {
	zr, err := zip.OpenReader(path)
	if err != nil {
		return nil, fmt.Errorf("open zip: %w", err)
	}
	// Vile, but archive/zip identifies the byte offset of the start of the zip data within the file
	// when the Reader is created, so this is the fastest (and easiest) way to find out where the zip
	// data begins in the underlying file.
	zipOffset := reflect.ValueOf(zr).Elem().FieldByName("baseOffset").Int()
	log.Debugf("%s: zip data begins at byte offset %d", path, zipOffset)
	zr.Close()
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open: %w", err)
	}
	return &PreambleReader{
		rc: f,
		sr: io.NewSectionReader(f, 0, zipOffset),
	}, nil
}
