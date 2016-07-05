package boltdb

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/l2x/gopprof/common/structs"
	"github.com/l2x/gopprof/engine/database"
)

type TableProfile struct {
	db     *bolt.DB
	nodeID string
	table  []byte
}

func NewTableProfile(db *bolt.DB, nodeID string) database.TableProfile {
	return &TableProfile{
		db:     db,
		nodeID: nodeID,
		table:  []byte("profile_" + nodeID),
	}
}

func (t *TableProfile) Table() []byte {
	return t.table
}

func (t *TableProfile) Save(data *structs.ProfileData) error {
	return t.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(t.Table())
		if err != nil {
			return err
		}
		k := fmt.Sprintf("%s_%d", data.NodeID, data.Created)
		v, err := json.Marshal(data)
		if err != nil {
			return err
		}
		return b.Put([]byte(k), v)
	})
}

func (t *TableProfile) GetRangeTime(start, end int64) ([]*structs.ProfileData, error) {
	data := []*structs.ProfileData{}
	err := t.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(t.Table())
		if b == nil {
			return nil
		}
		min := []byte(fmt.Sprintf("%d", t.nodeID, start))
		max := []byte(fmt.Sprintf("%d", t.nodeID, end))
		c := b.Cursor()
		for k, v := c.Seek(min); k != nil && bytes.Compare(k, max) <= 0; k, v = c.Next() {
			var d *structs.ProfileData
			if err := json.Unmarshal(v, &d); err != nil {
				return err
			}
			data = append(data, d)
		}
		return nil
	})
	return data, err
}

func (t *TableProfile) GetCreated(created int64) (*structs.ProfileData, error) {
	var data *structs.ProfileData
	err := t.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(t.Table())
		if b == nil {
			return sql.ErrNoRows
		}
		v := b.Get([]byte(fmt.Sprintf("%d", created)))
		if v == nil {
			return sql.ErrNoRows
		}
		if err := json.Unmarshal(v, &data); err != nil {
			return err
		}
		return nil
	})
	return data, err
}
