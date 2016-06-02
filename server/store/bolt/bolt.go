package boltstore

import (
	"github.com/boltdb/bolt"
	"github.com/l2x/gopprof/common/structs"
	"github.com/l2x/gopprof/server/store"
)

func init() {
	store.Register("bolt", NewBoltstore)
}

// Boltstore use boltdb
type Boltstore struct {
	db *bolt.DB
}

// NewBoltstore return Boltstore
func NewBoltstore() store.Store {
	return &Boltstore{}
}

// Open opens boltdb
func (b *Boltstore) Open(source string) error {
	db, err := bolt.Open(source, 0600, nil)
	if err != nil {
		return err
	}
	b.db = db
	return nil
}

// Close closes boltdb
func (b *Boltstore) Close() error {
	b.db.Close()
	return nil
}

// GetNode return NodeConf by nodeID
func (b *Boltstore) GetNode(nodeID string) (*structs.NodeConf, error) {
	return nil, nil
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
	return nil, nil
}
