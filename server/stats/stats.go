package stats

import "github.com/l2x/gopprof/common/structs"

// Stats is the interface that stats information
type Stats interface {
	Save(nodeID string, data *structs.StatsData) error
	GetTimeRange(nodeID string, timeStart, timeEnd int64) ([]*structs.StatsData, error)
	GetLatest(nodeID string, num int) ([]*structs.StatsData, error)
}
