package boltstore

import (
	"database/sql"
	"encoding/json"

	"github.com/boltdb/bolt"
	"github.com/l2x/gopprof/common/structs"
)

// TableConfName return table name
func (b *Boltstore) TableNodeName() string {
	return "node"
}

func (b *Boltstore) SaveNode(node *structs.NodeBase) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		v, err := json.Marshal(node)
		if err != nil {
			return err
		}
		return tx.Bucket([]byte(b.TableNodeName())).Put([]byte(node.NodeID), v)
	})
}

func (b *Boltstore) GetNodes() ([]*structs.NodeBase, error) {
	nodes := []*structs.NodeBase{}
	err := b.db.View(func(tx *bolt.Tx) error {
		tx.Bucket([]byte(b.TableNodeName())).ForEach(func(k, v []byte) error {
			var nodeBase *structs.NodeBase
			if err := json.Unmarshal(v, &nodeBase); err != nil {
				return err
			}
			nodes = append(nodes, nodeBase)
			return nil
		})
		return nil
	})
	return nodes, err
}

func (b *Boltstore) GetNode(nodeID string) (*structs.NodeBase, error) {
	var nodeBase *structs.NodeBase
	err := b.db.View(func(tx *bolt.Tx) error {
		v := tx.Bucket([]byte(b.TableNodeName())).Get([]byte(nodeID))
		if v == nil {
			return sql.ErrNoRows
		}
		if err := json.Unmarshal(v, &nodeBase); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return nodeBase, nil
}
