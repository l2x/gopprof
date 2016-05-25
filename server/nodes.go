package server

import (
	"sync"

	"github.com/l2x/gopprof/common/structs"
)

var (
	NodesMap = NewNodes()
)

type Nodes struct {
	nodes map[string]*structs.Node
	mu    sync.RWMutex
}

func NewNodes() *Nodes {
	return &Nodes{
		nodes: map[string]*structs.Node{},
	}
}

func (n *Nodes) Add(nodeBase structs.NodeBase) {
	n.mu.Lock()
	n.nodes[nodeBase.NodeID] = &structs.Node{
		NodeBase: nodeBase,
	}
	n.mu.Unlock()
}
