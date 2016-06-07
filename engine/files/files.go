package files

import (
	"fmt"
	"io"
)

var (
	drivers = map[string]func() Files{}
)

// Files is the interface store profiling files
type Files interface {
	Open(souce string) error
	Close() error

	Save(fname string, data []byte) error
	CopyTo(dst string, src io.Reader) error
	Get(fname string) ([]byte, error)
}

// Register makes a database driver available by the provided name.
func Register(driver string, f func() Files) {
	drivers[driver] = f
}

// Open opens a file storage
func Open(driver, source string) (Files, error) {
	driveri, ok := drivers[driver]
	if !ok {
		return nil, fmt.Errorf("sql: unknown driver %q (forgotten import?)", driver)
	}
	s := driveri()
	if err := s.Open(source); err != nil {
		return nil, err
	}
	return s, nil
}
