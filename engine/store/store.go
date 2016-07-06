// Package store defines interfaces to be implemented by file storage.
//
// Store save profiling files.
package store

import (
	"fmt"
	"io"
)

var (
	drivers = map[string]func() Store{}
)

// Register makes a driver available by the provided name.
func Register(driver string, f func() Store) {
	drivers[driver] = f
}

// Open opens a file storage
func Open(driver, source string) (Store, error) {
	driveri, ok := drivers[driver]
	if !ok {
		return nil, fmt.Errorf("engine/store: unknown driver %q (forgotten import?)", driver)
	}
	s := driveri()
	if err := s.Open(source); err != nil {
		return nil, err
	}
	return s, nil
}

// Store is the interface that must be implemented by a file storage
type Store interface {
	Open(souce string) error
	Close() error

	Save(fname string, data []byte) error
	Copy(dst string, src io.Reader) error
	Get(fname string) ([]byte, error)
	CopyToLocal(dst, src string) error
}
