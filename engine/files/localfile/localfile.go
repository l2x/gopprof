package localfile

import (
	"os"

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
func (f *Localfile) Save(nodeID, typ string, data []byte) (string, error) {
	return "", nil
}

// Get file
func (f *Localfile) Get(fname string) ([]byte, error) {
	return nil, nil
}
