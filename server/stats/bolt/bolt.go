package boltstats

import (
	"bytes"
	"encoding/json"
	"fmt"
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

func (b *Boltstats) tableStats(nodeID string) []byte {
	return []byte("stats_" + nodeID)
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
	return b.db.Update(func(tx *bolt.Tx) error {
		buc, err := tx.CreateBucketIfNotExists(b.tableStats(data.NodeID))
		if err != nil {
			return err
		}
		k := fmt.Sprintf("%s_%d", data.NodeID, data.Created)
		v, err := json.Marshal(data)
		if err != nil {
			return err
		}
		return buc.Put([]byte(k), v)
	})
}

// GetTimeRange get data by time
func (b *Boltstats) GetTimeRange(nodeID string, timeStart, timeEnd int64) ([]*structs.StatsData, error) {
	data := []*structs.StatsData{}
	err := b.db.View(func(tx *bolt.Tx) error {
		buc := tx.Bucket(b.tableStats(nodeID))
		if buc == nil {
			return nil
		}
		c := buc.Cursor()
		min := []byte(fmt.Sprintf("%s_%d", nodeID, timeStart))
		max := []byte(fmt.Sprintf("%s_%d", nodeID, timeEnd))

		for k, v := c.Seek(min); k != nil && bytes.Compare(k, max) <= 0; k, v = c.Next() {
			var d *structs.StatsData
			if err := json.Unmarshal(v, &d); err != nil {
				return err
			}
			data = append(data, d)
		}
		return nil
	})
	return data, err
}

// GetLatest return latest rows
func (b *Boltstats) GetLatest(nodeID string, num int) ([]*structs.StatsData, error) {
	data := []*structs.StatsData{}
	err := b.db.View(func(tx *bolt.Tx) error {
		buc := tx.Bucket(b.tableStats(nodeID))
		if buc == nil {
			return nil
		}
		c := buc.Cursor()
		for k, v := c.Last(); k != nil; k, v = c.Prev() {
			if len(data) == num {
				return nil
			}
			var d *structs.StatsData
			if err := json.Unmarshal(v, &d); err != nil {
				return err
			}
			data = append(data, d)
		}
		return nil
	})
	return data, err
}
