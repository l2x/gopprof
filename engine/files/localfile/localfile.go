package localfile

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/l2x/gopprof/engine/files"
)

func init() {
	files.Register("localfile", NewLocalfile)
}

// Localfile store file in local file system
type Localfile struct {
	base string
}

// NewLocalfile return localfile storage
func NewLocalfile() files.Files {
	return &Localfile{}
}

// Open init file storage
func (f *Localfile) Open(souce string) error {
	f.base = souce
	return os.MkdirAll(f.base, 0755)
}

// Close file
func (f *Localfile) Close() error {
	return nil
}

// Save file
func (f *Localfile) Save(fname string, data []byte) error {
	fname = filepath.Join(f.base, fname)
	if err := os.MkdirAll(filepath.Dir(fname), 0755); err != nil {
		return err
	}
	return ioutil.WriteFile(fname, data, 0644)
}

// CopyTo copy src to dst
func (f *Localfile) CopyTo(dst string, src io.Reader) error {
	fname := filepath.Join(f.base, dst)
	if err := os.MkdirAll(filepath.Dir(fname), 0755); err != nil {
		return err
	}
	fn, err := os.OpenFile(fname, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer fn.Close()
	_, err = io.Copy(fn, src)
	if err != nil {
		return err
	}
	return nil
}

// Get file
func (f *Localfile) Get(fname string) ([]byte, error) {
	return ioutil.ReadFile(filepath.Join(f.base, fname))
}
