package boltstore

import (
	"os"
	"path/filepath"

	"github.com/boltdb/bolt"
	"github.com/l2x/gopprof/engine/store"
)

func init() {
	store.Register("bolt", NewBoltstore)
}

// Boltstore use boltdb
type Boltstore struct {
	db             *bolt.DB
	defaultConfKey []byte
}

// NewBoltstore return Store
func NewBoltstore() store.Store {
	return &Boltstore{
		defaultConfKey: []byte("default_conf"),
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
		if _, err := tx.CreateBucketIfNotExists([]byte(b.TableConfName())); err != nil {
			return err
		}
		if _, err := tx.CreateBucketIfNotExists([]byte(b.TableTagName())); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

// Close db
func (b *Boltstore) Close() error {
	b.db.Close()
	return nil
}
