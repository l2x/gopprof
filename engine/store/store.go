package store

import (
	"fmt"

	"github.com/l2x/gopprof/common/structs"
)

var (
	drivers = map[string]func() Store{}
)

// Store is the interface that wraps the tables
type Store interface {
	Open(source string) error
	Close() error

	TableNode
	TableConf
	TableTag
	TableProfile
	TableStat
	TableBin
}

// TableConf is the interface defined table conf
type TableConf interface {
	TableConfName() string
	GetConf(nodeID string) (*structs.NodeConf, error)
	GetDefaultConf() (*structs.NodeConf, error)
	SaveConf(nodeID string, nodeConf *structs.NodeConf) error
	SaveDefaultConf(nodeConf *structs.NodeConf) error
}

// TableNode is the interface defined table node
type TableNode interface {
	SaveNode(node *structs.NodeBase) error
	GetNodes() ([]*structs.NodeBase, error)
	GetNode(nodeID string) (*structs.NodeBase, error)
}

// TableTag is the interface defined table tags
type TableTag interface {
	TableTagName() string
	GetTags() ([]string, error)
	SaveTags(nodeID string, tags []string) error
	DelTag(nodeID, tag string) error
}

// TableProfile is the interface defined table profile
type TableProfile interface {
	TableProfileName(nodeID string) string
	SaveProfile(data *structs.ProfileData) error
	GetProfilesByTime(nodeID string, timeStart, timeEnd int64) ([]*structs.ProfileData, error)
	GetProfilesLatest(nodeID string, num int) ([]*structs.ProfileData, error)
}

// TableStat is the interface defined table stat
type TableStat interface {
	TableStatName(nodeID string) string
	SaveStat(data *structs.StatsData) error
	GetStatsByTime(nodeID string, timeStart, timeEnd int64) ([]*structs.StatsData, error)
	GetStatsLatest(nodeID string, num int) ([]*structs.StatsData, error)
}

// TableBin is the interface defined table bin
type TableBin interface {
	TableBinName(nodeID string) string
	SaveBin(nodeID, binMd5, file string) error
	GetBin(nodeID, binMd5 string) (string, error)
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
