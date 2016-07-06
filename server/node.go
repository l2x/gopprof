package server

import (
	"sync"

	"github.com/l2x/gopprof/common/structs"
)

var (
	nodesMap = NewNodes()
)

// Nodes is node map
type Nodes struct {
	nodes map[string]*structs.Node
	mu    sync.RWMutex
}

func NewNodes() *Nodes {
	return &Nodes{
		nodes: map[string]*structs.Node{},
	}
}

// Add node
func (n *Nodes) Add(nodeID string) *structs.Node {
	node := structs.NewNode(nodeID)
	n.mu.Lock()
	n.nodes[nodeID] = node
	n.mu.Unlock()
	return node
}

// Get node by nodeID
func (n *Nodes) Get(nodeID string) (*structs.Node, bool) {
	n.mu.RLock()
	node, ok := n.nodes[nodeID]
	n.mu.RUnlock()
	return node, ok
}
