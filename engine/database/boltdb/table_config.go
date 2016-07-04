package boltdb

import (
	"database/sql"
	"encoding/json"

	"github.com/boltdb/bolt"
	"github.com/l2x/gopprof/common/structs"
	"github.com/l2x/gopprof/engine/database"
)

type TableConfig struct {
	db     *bolt.DB
	nodeID string
	table  []byte
}

func NewTableConfig(db *bolt.DB, nodeID string) database.TableConfig {
	return &TableConfig{
		db:     db,
		nodeID: nodeID,
		table:  []byte("config"),
	}
}

func (t *TableConfig) Table() []byte {
	return t.table
}

func (t *TableConfig) Save(data *structs.NodeConf) error {
	return t.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(t.Table())
		v, err := json.Marshal(data)
		if err != nil {
			return err
		}
		return b.Put([]byte(t.nodeID), v)
	})
}

func (t *TableConfig) Get() (*structs.NodeConf, error) {
	var nodeConf *structs.NodeConf
	err := t.db.View(func(tx *bolt.Tx) error {
		v := tx.Bucket(t.Table()).Get([]byte(t.nodeID))
		if v == nil {
			return sql.ErrNoRows
		}
		if err := json.Unmarshal(v, &nodeConf); err != nil {
			return err
		}
		return nil
	})
	return nodeConf, err
}
