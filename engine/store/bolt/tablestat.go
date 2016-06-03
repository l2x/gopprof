package boltstore

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/l2x/gopprof/common/structs"
)

// TableStatName return table name
func (b *Boltstore) TableStatName(nodeID string) string {
	return "stat_" + nodeID
}

// SaveStat save stat data
func (b *Boltstore) SaveStat(data *structs.StatsData) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		buc, err := tx.CreateBucketIfNotExists([]byte(b.TableStatName(data.NodeID)))
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

// GetStatsByTime return stat data
func (b *Boltstore) GetStatsByTime(nodeID string, timeStart, timeEnd int64) ([]*structs.StatsData, error) {
	data := []*structs.StatsData{}
	err := b.db.View(func(tx *bolt.Tx) error {
		buc := tx.Bucket([]byte(b.TableStatName(nodeID)))
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

// GetStatsLatest return latest stat data
func (b *Boltstore) GetStatsLatest(nodeID string, num int) ([]*structs.StatsData, error) {
	data := []*structs.StatsData{}
	err := b.db.View(func(tx *bolt.Tx) error {
		buc := tx.Bucket([]byte(b.TableStatName(nodeID)))
		if buc == nil {
			return nil
		}
		c := buc.Cursor()
		for k, v := c.Last(); k != nil && num > len(data); k, v = c.Prev() {
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
