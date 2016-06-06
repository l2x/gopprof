package boltstore

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/l2x/gopprof/common/structs"
)

// TableProfileName return table name
func (b *Boltstore) TableProfileName(nodeID string) string {
	return "profile_" + nodeID
}

// SaveProfile save profile data
func (b *Boltstore) SaveProfile(data *structs.ProfileData) (int64, error) {
	return data.ID, b.db.Update(func(tx *bolt.Tx) error {
		buc, err := tx.CreateBucketIfNotExists([]byte(b.TableProfileName(data.NodeID)))
		if err != nil {
			return err
		}
		if data.ID == 0 {
			id, _ := buc.NextSequence()
			data.ID = int64(id)
		}

		k := fmt.Sprintf("%s_%d", data.NodeID, data.Created)
		v, err := json.Marshal(data)
		if err != nil {
			return err
		}
		return buc.Put([]byte(k), v)
	})
}

// GetProfilesByTime return profile data
func (b *Boltstore) GetProfilesByTime(nodeID string, timeStart, timeEnd int64) ([]*structs.ProfileData, error) {
	data := []*structs.ProfileData{}
	err := b.db.View(func(tx *bolt.Tx) error {
		buc := tx.Bucket([]byte(b.TableProfileName(nodeID)))
		if buc == nil {
			return nil
		}

		c := buc.Cursor()
		min := []byte(fmt.Sprintf("%s_%d", nodeID, timeStart))
		max := []byte(fmt.Sprintf("%s_%d", nodeID, timeEnd))
		for k, v := c.Seek(min); k != nil && bytes.Compare(k, max) <= 0; k, v = c.Next() {
			var d *structs.ProfileData
			if err := json.Unmarshal(v, &d); err != nil {
				return err
			}
			data = append(data, d)
		}
		return nil
	})
	return data, err
}

// GetProfilesLatest return latest profile data
func (b *Boltstore) GetProfilesLatest(nodeID string, num int) ([]*structs.ProfileData, error) {
	data := []*structs.ProfileData{}
	err := b.db.View(func(tx *bolt.Tx) error {
		buc := tx.Bucket([]byte(b.TableStatName(nodeID)))
		if buc == nil {
			return nil
		}
		c := buc.Cursor()
		for k, v := c.Last(); k != nil && num > len(data); k, v = c.Prev() {
			var d *structs.ProfileData
			if err := json.Unmarshal(v, &d); err != nil {
				return err
			}
			data = append(data, d)
		}
		return nil
	})
	return data, err
}
