package stats

import (
	"fmt"

	"github.com/l2x/gopprof/common/structs"
)

var (
	drivers = map[string]func() Stats{}
)

// Stats is the interface that stats information
type Stats interface {
	Open(souce string) error
	Close() error

	Save(data *structs.StatsData) error
	GetTimeRange(nodeID string, timeStart, timeEnd int64) ([]*structs.StatsData, error)
	GetLatest(nodeID string, num int) ([]*structs.StatsData, error)
}

// Register makes a database driver available by the provided name.
func Register(driver string, f func() Stats) {
	drivers[driver] = f
}

// Open opens a database specified by its database driver name and a
// driver-specific data source name,
func Open(driver, source string) (Stats, error) {
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
