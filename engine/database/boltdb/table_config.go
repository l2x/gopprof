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

func (t *TableConfig) Goroots() ([]*structs.Goroot, error) {
	goroots := []*structs.Goroot{}
	err := t.db.View(func(tx *bolt.Tx) error {
		prefix := []byte("goroot_")
		c := tx.Bucket(t.Table()).Cursor()
		for k, v := c.Seek(prefix); bytes.HasPrefix(k, prefix); k, v = c.Next() {
			var goroot *structs.Goroot
			if err := json.Unmarshal(v, &goroot); err != nil {
				continue
			}
			goroots = append(goroots, goroot)
		}
		return nil
	})
	return goroots, err
}

func (t *TableConfig) GetGoroot(version string) (*structs.Goroot, error) {
	var goroot *structs.Goroot
	err := t.db.View(func(tx *bolt.Tx) error {
		k := fmt.Sprintf("goroot_%s", version)
		v := tx.Bucket(t.Table()).Get([]byte(k))
		if v == nil {
			return sql.ErrNoRows
		}
		if err := json.Unmarshal(v, &goroot); err != nil {
			return err
		}
		return nil
	})
	return goroot, err
}

func (t *TableConfig) SaveGoroot(goroot *structs.Goroot) error {
	return t.db.Update(func(tx *bolt.Tx) error {
		v, err := json.Marshal(goroot)
		if err != nil {
			return err
		}
		k := fmt.Sprintf("goroot_%s", goroot.Version)
		return tx.Bucket(t.Table()).Put([]byte(k), v)
	})
}

func (t *TableConfig) DelGoroot(goroot *structs.Goroot) error {
	return t.db.Update(func(tx *bolt.Tx) error {
		k := fmt.Sprintf("goroot_%s", goroot.Version)
		return tx.Bucket(t.Table()).Delete([]byte(k))
	})
}
