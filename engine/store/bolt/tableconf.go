package boltstore

import "github.com/l2x/gopprof/common/structs"

// TableConfName return table name
func (b *Boltstore) TableConfName() string {
	return "conf"
}

// GetConf return node conf
func (b *Boltstore) GetConf(nodeID string) (*structs.NodeConf, error) {
	return nil, nil
}

// GetDefaultConf return default conf
func (b *Boltstore) GetDefaultConf() (*structs.NodeConf, error) {
	return nil, nil
}

// SaveConf save conf
func (b *Boltstore) SaveConf(nodeID string, nodeConf *structs.NodeConf) error {
	return nil
}
