package boltstore

import (
	"database/sql"

	"github.com/boltdb/bolt"
)

// TableBinName return table name
func (b *Boltstore) TableBinName(nodeID string) string {
	return "bin_" + nodeID
}

func (b *Boltstore) SaveBin(nodeID, binMd5, file string) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		buc, err := tx.CreateBucketIfNotExists([]byte(b.TableBinName(nodeID)))
		if err != nil {
			return err
		}
		return buc.Put([]byte(binMd5), []byte(file))
	})
}

func (b *Boltstore) GetBin(nodeID, binMd5 string) (string, error) {
	var file string
	return file, b.db.View(func(tx *bolt.Tx) error {
		buc := tx.Bucket([]byte(b.TableBinName(nodeID)))
		if buc == nil {
			return sql.ErrNoRows
		}
		v := buc.Get([]byte(binMd5))
		if v == nil {
			return sql.ErrNoRows
		}
		file = string(v)
		return nil
	})
}
