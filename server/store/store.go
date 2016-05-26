package store

import "github.com/l2x/gopprof/common/structs"

// Store is the interface that storage information
type Store interface {
	GetNode(nodeID string) (structs.NodeConf, error)
	GetNodeByTag(tag string) (structs.NodeConf, error)
}
