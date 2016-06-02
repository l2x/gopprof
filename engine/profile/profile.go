package profile

import "github.com/l2x/gopprof/common/structs"

// Profile is the interface that profiling information
type Profile interface {
	Open(souce string) error
	Close() error

	Save(data *structs.ProfileData) error
	GetTimeRange(nodeID string, timeStart, timeEnd int64) ([]*structs.ProfileData, error)
	GetLatest(nodeID string, num int) ([]*structs.ProfileData, error)
}
