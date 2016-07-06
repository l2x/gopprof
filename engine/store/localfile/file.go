package localfile

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/l2x/gopprof/engine/store"
)

var (
	permDir  os.FileMode = 0755
	permFile os.FileMode = 0644
)

func init() {
	store.Register("localfile", NewLocalfile)
}

// Localfile store file in local file system
type Localfile struct {
	base string
}

// NewLocalfile return localfile storage
func NewLocalfile() store.Store {
	return &Localfile{}
}

// Open init file storage
func (f *Localfile) Open(souce string) error {
	f.base = souce
	return os.MkdirAll(f.base, permDir)
}

// Close file
func (f *Localfile) Close() error {
	return nil
}

// Save file
func (f *Localfile) Save(fname string, data []byte) error {
	fname = filepath.Join(f.base, fname)
	if err := os.MkdirAll(filepath.Dir(fname), permDir); err != nil {
		return err
	}
	return ioutil.WriteFile(fname, data, permFile)
}

// Copy src to dst
func (f *Localfile) Copy(dst string, src io.Reader) error {
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

func (f *Localfile) CopyToLocal(dst, src string) error {
	b, err := f.Get(src)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(dst, b, 0755)
	if err != nil {
		return err
	}
	return nil
}
