package boltdb

import (
	"database/sql"
	"encoding/json"

	"github.com/boltdb/bolt"
	"github.com/l2x/gopprof/common/structs"
	"github.com/l2x/gopprof/engine/database"
)

type TableConfig struct {
	db         *bolt.DB
	nodeID     string
	table      []byte
	defaultKey []byte
}

func NewTableConfig(db *bolt.DB, nodeID string) database.TableConfig {
	return &TableConfig{
		db:         db,
		nodeID:     nodeID,
		table:      []byte("config"),
		defaultKey: []byte("_default"),
	}
}

func (t *TableConfig) Table() []byte {
	return t.table
}

func (t *TableConfig) Save(data *structs.NodeConf) error {
	return t.db.Update(func(tx *bolt.Tx) error {
		v, err := json.Marshal(data)
		if err != nil {
			return err
		}
		b := tx.Bucket(t.Table())
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

func (t *TableConfig) SaveDefault(data *structs.NodeConf) error {
	return t.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(t.Table())
		v, err := json.Marshal(data)
		if err != nil {
			return err
		}
		return b.Put(t.defaultKey, v)
	})
}

func (t *TableConfig) GetDefault() (*structs.NodeConf, error) {
	nodeConf := &structs.NodeConf{
		EnableStats:   true,
		StatsCron:     "*/1 * * * *",
		EnableProfile: true,
		ProfileCron:   "*/1 * * * *",
		Profile:       []string{"heap", "cpu", "trace", "block"},
	}
	err := t.db.View(func(tx *bolt.Tx) error {
		v := tx.Bucket(t.Table()).Get([]byte(t.nodeID))
		if v == nil {
			return nil
		}
		if err := json.Unmarshal(v, &nodeConf); err != nil {
			return err
		}
		return nil
	})
	return nodeConf, err
}
