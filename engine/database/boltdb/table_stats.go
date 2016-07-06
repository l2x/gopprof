package boltdb

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/l2x/gopprof/common/structs"
	"github.com/l2x/gopprof/engine/database"
)

type TableStats struct {
	db     *bolt.DB
	nodeID string
	table  []byte
}

func NewTableStats(db *bolt.DB, nodeID string) database.TableStats {
	return &TableStats{
		db:     db,
		nodeID: nodeID,
		table:  []byte("stats_" + nodeID),
	}
}

func (t *TableStats) Table() []byte {
	return t.table
}

func (t *TableStats) Save(data *structs.StatsData) error {
	return t.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(t.Table())
		if err != nil {
			return err
		}
		k := fmt.Sprintf("%d", data.Created)
		v, err := json.Marshal(data)
		if err != nil {
			return err
		}
		return b.Put([]byte(k), v)
	})
}

func (t *TableStats) GetRangeTime(start, end int64) ([]*structs.StatsData, error) {
	data := []*structs.StatsData{}
	err := t.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(t.Table())
		if b == nil {
			return nil
		}
		min := []byte(fmt.Sprintf("%d", start))
		max := []byte(fmt.Sprintf("%d", end))
		c := b.Cursor()
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
