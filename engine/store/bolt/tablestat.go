package boltstore

import "github.com/l2x/gopprof/common/structs"

// TableStatName return table name
func (b *Boltstore) TableStatName(nodeID string) string {
	return "stat"
}

// SaveStat save stat data
func (b *Boltstore) SaveStat(data *structs.StatsData) error {
	return nil
}

// GetStatsByTime return stat data
func (b *Boltstore) GetStatsByTime(nodeID string, timeStart, timeEnd int64) ([]*structs.StatsData, error) {
	return nil, nil
}

// GetStatsLatest return latest stat data
func (b *Boltstore) GetStatsLatest(nodeID string, num int) ([]*structs.StatsData, error) {
	return nil, nil
}
