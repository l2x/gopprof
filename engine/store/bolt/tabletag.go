package boltstore

import "github.com/l2x/gopprof/common/structs"

// TableTagName return table name
func (b *Boltstore) TableTagName() string {
	return "tag"
}

// GetTags return tags
func (b *Boltstore) GetTags() ([]string, error) {
	return nil, nil
}

// GetNodeByTag return nodes by tag
func (b *Boltstore) GetNodeByTag(tag string) ([]*structs.NodeConf, error) {
	return nil, nil
}

// SaveTags save tags
func (b *Boltstore) SaveTags(nodeID string, tags []string) error {
	return nil
}