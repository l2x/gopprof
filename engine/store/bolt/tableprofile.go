package boltstore

import "github.com/l2x/gopprof/common/structs"

// TableProfileName return table name
func (b *Boltstore) TableProfileName(nodeID string) string {
	return "profile"
}

// SaveProfile save profile data
func (b *Boltstore) SaveProfile(data *structs.ProfileData) error {
	return nil
}

// GetProfilesByTime return profile data
func (b *Boltstore) GetProfilesByTime(nodeID string, timeStart, timeEnd int64) ([]*structs.ProfileData, error) {
	return nil, nil
}

// GetProfilesLatest return latest profile data
func (b *Boltstore) GetProfilesLatest(nodeID string, num int) ([]*structs.ProfileData, error) {
	return nil, nil
}
