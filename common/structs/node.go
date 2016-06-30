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
	NodeID     string `json:"nodeid"`
	Hostname   string `json:"hostname"`
	InternalIP string `json:"internal_ip"`
	LastSync   time.Time
	Created    time.Time
	Status     uint8 `json:"status"`
}

// NodeConf storage node config
type NodeConf struct {
	Tags []string

	EnableProfile   bool
	Profile         []ProfileData
	LastProfile     time.Time
	ProfileInterval time.Duration

	EnableStat   bool
	LastStat     time.Time
	StatInterval time.Duration
}
