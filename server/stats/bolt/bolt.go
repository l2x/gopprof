package boltstats

import (
	"os"
	"path/filepath"

	"github.com/boltdb/bolt"
	"github.com/l2x/gopprof/common/structs"
	"github.com/l2x/gopprof/server/stats"
)

func init() {
	stats.Register("bolt", NewBoltstats)
}

// Boltstats use boltdb
type Boltstats struct {
	db *bolt.DB
}

// NewBoltstats return Boltstore
func NewBoltstats() stats.Stats {
	return &Boltstats{}
}

// Open opens boltdb
func (b *Boltstats) Open(source string) error {
	if err := os.MkdirAll(filepath.Dir(source), 0755); err != nil {
		return err
	}
	db, err := bolt.Open(source, 0600, nil)
	if err != nil {
		return err
	}
	b.db = db
	return nil
}

// Close closes boltdb
func (b *Boltstats) Close() error {
	b.db.Close()
	return nil
}

// Save data
func (b *Boltstats) Save(data *structs.StatsData) error {
	return nil
}

// GetTimeRange get data by time
func (b *Boltstats) GetTimeRange(nodeID string, timeStart, timeEnd int64) ([]*structs.StatsData, error) {
	return nil, nil
}

// GetLatest return latest rows
func (b *Boltstats) GetLatest(nodeID string, num int) ([]*structs.StatsData, error) {
	return nil, nil
}
