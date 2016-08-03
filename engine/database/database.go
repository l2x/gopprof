// Package database defines interfaces to be implemented by data storage.
package database

import (
	"fmt"

	"github.com/l2x/gopprof/common/structs"
)

var (
	drivers = map[string]func() Database{}
)

// Register makes a database driver available by the provided name.
func Register(driver string, f func() Database) {
	drivers[driver] = f
}

// Open opens a database specified by its database driver name and a
// driver-specific data source name,
func Open(driver, source string) (Database, error) {
	driveri, ok := drivers[driver]
	if !ok {
		return nil, fmt.Errorf("engine/database: unknown driver %q (forgotten import?)", driver)
	}
	s := driveri()
	if err := s.Open(source); err != nil {
		return nil, err
	}
	return s, nil
}

// Database is the interface that must be implemented by a data storage
type Database interface {
	Open(source string) error
	Close() error

	TableStats(nodeID string) TableStats
	TableProfile(nodeID string) TableProfile
	TableConfig(nodeID string) TableConfig
	TableNode(nodeID string) TableNode
	TableBin(nodeID string) TableBin
}

// TableStats save stats data
type TableStats interface {
	Table() []byte
	Save(data *structs.StatsData) error
	GetRangeTime(start, end int64) ([]*structs.StatsData, error)
}

// TableProfile save profile data
type TableProfile interface {
	Table() []byte
	Save(data *structs.ProfileData) error
	GetRangeTime(start, end int64) ([]*structs.ProfileData, error)
	GetCreated(created int64) (*structs.ProfileData, error)
}

// TableConfig save all configure
type TableConfig interface {
	Table() []byte
	Save(data *structs.NodeConf) error
	Get() (*structs.NodeConf, error)
	Goroots() ([]*structs.Goroot, error)
	GetGoroot(version string) (*structs.Goroot, error)
	SaveGoroot(goroot *structs.Goroot) error
}

// TableNode save all node info
type TableNode interface {
	Table() []byte
	Save(data *structs.NodeBase) error
	Get() (*structs.NodeBase, error)
	GetAll() ([]*structs.NodeBase, error)
}

// TableBin save binary file info
type TableBin interface {
	Table() []byte
	Save(binMD5 string, file string) error
	Get(binMD5 string) (string, error)
}
