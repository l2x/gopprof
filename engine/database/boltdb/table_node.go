package boltdb

import (
	"database/sql"
	"encoding/json"

	"github.com/boltdb/bolt"
	"github.com/l2x/gopprof/common/structs"
	"github.com/l2x/gopprof/engine/database"
)

type TableNode struct {
	db     *bolt.DB
	nodeID string
	table  []byte
}

func NewTableNode(db *bolt.DB, nodeID string) database.TableNode {
	return &TableNode{
		db:     db,
		nodeID: nodeID,
		table:  []byte("node"),
	}
}

func (t *TableNode) Table() []byte {
	return t.table
}

func (t *TableNode) Save(data *structs.NodeBase) error {
	return t.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(t.Table())
		v, err := json.Marshal(data)
		if err != nil {
			return err
		}
		return b.Put([]byte(t.nodeID), v)
	})
}

func (t *TableNode) Get() (*structs.NodeBase, error) {
	var nodeBase *structs.NodeBase
	err := t.db.View(func(tx *bolt.Tx) error {
		v := tx.Bucket(t.Table()).Get([]byte(t.nodeID))
		if v == nil {
			return sql.ErrNoRows
		}
		if err := json.Unmarshal(v, &nodeBase); err != nil {
			return err
		}
		return nil
	})
	return nodeBase, err
}
