package localfile

import (
	"os"

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
