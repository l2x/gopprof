package boltstore

import (
	"database/sql"
	"encoding/json"

	"github.com/boltdb/bolt"
	"github.com/l2x/gopprof/common/structs"
)

// TableConfName return table name
func (b *Boltstore) TableConfName() string {
	return "conf"
}

// GetConf return node conf
func (b *Boltstore) GetConf(nodeID string) (*structs.NodeConf, error) {
	var nodeConf *structs.NodeConf
	err := b.db.View(func(tx *bolt.Tx) error {
		v := tx.Bucket([]byte(b.TableConfName())).Get([]byte(nodeID))
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

// GetDefaultConf return default conf
func (b *Boltstore) GetDefaultConf() (*structs.NodeConf, error) {
	nodeConf := &structs.NodeConf{}
	err := b.db.View(func(tx *bolt.Tx) error {
		v := tx.Bucket([]byte(b.TableConfName())).Get(b.defaultConfKey)
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

// SaveConf save conf
func (b *Boltstore) SaveConf(nodeID string, nodeConf *structs.NodeConf) error {
	return nil
}