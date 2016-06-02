package boltstore

import (
	"database/sql"
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/boltdb/bolt"
	"github.com/l2x/gopprof/common/structs"
	"github.com/l2x/gopprof/engine/store"
)

func init() {
	store.Register("bolt", NewBoltstore)
}

// Boltstore use boltdb
type Boltstore struct {
	db        *bolt.DB
	tableConf []byte
}

// NewBoltstore return Boltstore
func NewBoltstore() store.Store {
	return &Boltstore{
		tableConf: []byte("store_conf"),
	}
}

// Open opens boltdb
func (b *Boltstore) Open(source string) error {
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

func (b *Boltstore) init() error {
	err := b.db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists(b.tableConf); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

// Close closes boltdb
func (b *Boltstore) Close() error {
	b.db.Close()
	return nil
}

// GetNode return NodeConf by nodeID
func (b *Boltstore) GetNode(nodeID string) (*structs.NodeConf, error) {
	var nodeConf *structs.NodeConf
	err := b.db.View(func(tx *bolt.Tx) error {
		v := tx.Bucket(b.tableConf).Get([]byte(nodeID))
		if v == nil {
			return sql.ErrNoRows
		}
		if err := json.Unmarshal(v, &nodeConf); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return nodeConf, nil
}

// GetNodeByTag return NodeConf slice by tag
func (b *Boltstore) GetNodeByTag(tag string) ([]*structs.NodeConf, error) {
	return nil, nil
}

// SetTags set tags
func (b *Boltstore) SetTags(nodeID string, tags []string) error {
	return nil
}

// GetDefault return default NodeConf
func (b *Boltstore) GetDefault() (*structs.NodeConf, error) {
	return &structs.NodeConf{
		EnableStat:   true,
		StatInterval: 30 * time.Second,
	}, nil
}
