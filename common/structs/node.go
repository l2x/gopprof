package structs

import (
	"errors"
	"time"
)

// Node storage node information
type Node struct {
	event chan *Event

	NodeBase
	NodeConf
}

// NewNode return an Node
func NewNode(NodeID string) *Node {
	return &Node{
		event: make(chan *Event, 20),
		NodeBase: NodeBase{
			NodeID: NodeID,
		},
	}
}

// Event return event chan
func (n *Node) Event() chan *Event {
	return n.event
}

// AddEvent add event
func (n *Node) AddEvent(evt *Event) error {
	select {
	case n.event <- evt:
	default:
		return errors.New("event chan is full")
	}
	return nil
}

// NodeBase storage node base information
type NodeBase struct {
	NodeID   string
	LastSync time.Time
}

// NodeConf storage node config
type NodeConf struct {
}
