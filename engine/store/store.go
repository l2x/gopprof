package store

import (
	"fmt"

	"github.com/l2x/gopprof/common/structs"
)

var (
	drivers = map[string]func() Store{}
)

// Store is the interface that storage information
type Store interface {
	Open(souce string) error
	Close() error

	GetNode(nodeID string) (*structs.NodeConf, error)
	GetNodeByTag(tag string) ([]*structs.NodeConf, error)
	SetTags(nodeID string, tags []string) error
	GetDefault() (*structs.NodeConf, error)
}

// Register makes a database driver available by the provided name.
func Register(driver string, f func() Store) {
	drivers[driver] = f
}

// Open opens a database specified by its database driver name and a
// driver-specific data source name,
func Open(driver, source string) (Store, error) {
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
