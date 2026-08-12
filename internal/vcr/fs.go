package vcr

import (
	"bytes"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
)

// GzipFS is a cassette.FS implementation that reads and writes gzip-compressed
// files. Note: not currently active. When in use the cassettes are not
// human-readable, so troubleshooting / reviewing is inconvenient. It is kept
// here (wired via a commented recorder.WithFS in recorder.go) for cases where a
// cassette grows large enough that on-disk size outweighs readability.
type GzipFS struct{}

func (fs *GzipFS) ReadFile(name string) ([]byte, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r, err := gzip.NewReader(f)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	return io.ReadAll(r)
}

func (fs *GzipFS) WriteFile(name string, data []byte) error {
	dir := filepath.Dir(name)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err = os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}

	f, err := os.Create(name)
	if err != nil {
		return err
	}
	defer f.Close()

	var buf bytes.Buffer
	w := gzip.NewWriter(&buf)
	if _, err := w.Write(data); err != nil {
		return err
	}
	if err := w.Close(); err != nil {
		return err
	}

	_, err = f.Write(buf.Bytes())
	return err
}

func (fs *GzipFS) IsFileExists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}
