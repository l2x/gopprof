// Package boltdb provides a boltDB driver for engine/database package
package boltdb

import (
	"os"
	"path/filepath"

	"github.com/boltdb/bolt"
	"github.com/l2x/gopprof/engine/database"
)

func init() {
	database.Register("bolt", NewBoltDB)
}

type BoltDB struct {
	db *bolt.DB
}

func NewBoltDB() database.Database {
	return &BoltDB{}
}

// Open opens boltdb
func (b *BoltDB) Open(source string) error {
	if err := os.MkdirAll(filepath.Dir(source), 0755); err != nil {
		return err
	}
	db, err := bolt.Open(source, 0600, nil)
	if err != nil {
		return err
	}
	b.db = db
	return b.init()
}

func (b *BoltDB) Close() error {
	return b.db.Close()
}

func (b *BoltDB) TableStats(nodeID string) database.TableStats {
	return NewTableStats(b.db, nodeID)
}

func (b *BoltDB) TableProfile(nodeID string) database.TableProfile {
	return NewTableProfile(b.db, nodeID)
}

func (b *BoltDB) TableConfig(nodeID string) database.TableConfig {
	return NewTableConfig(b.db, nodeID)
}

func (b *BoltDB) TableNode(nodeID string) database.TableNode {
	return NewTableNode(b.db, nodeID)
}

func (b *BoltDB) init() error {
	return b.db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists(NewTableNode(nil, "").Table()); err != nil {
			return err
		}
		if _, err := tx.CreateBucketIfNotExists(NewTableConfig(nil, "").Table()); err != nil {
			return err
		}
		return nil
	})
}
