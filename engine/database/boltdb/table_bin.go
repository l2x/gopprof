package boltdb

import (
	"database/sql"

	"github.com/boltdb/bolt"
	"github.com/l2x/gopprof/engine/database"
)

type TableBin struct {
	db     *bolt.DB
	nodeID string
	table  []byte
}

func NewTableBin(db *bolt.DB, nodeID string) database.TableBin {
	return &TableBin{
		db:     db,
		nodeID: nodeID,
		table:  []byte("bin"),
	}
}

func (t *TableBin) Table() []byte {
	return t.table
}

func (t *TableBin) Save(binMD5, file string) error {
	return t.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(t.Table()).Put([]byte(binMD5), []byte(file))
	})
}

func (t *TableBin) Get(binMD5 string) (string, error) {
	var file string
	err := t.db.View(func(tx *bolt.Tx) error {
		v := tx.Bucket(t.Table()).Get([]byte(binMD5))
		if v == nil {
			return sql.ErrNoRows
		}
		file = string(v)
		return nil
	})
	return file, err
}
