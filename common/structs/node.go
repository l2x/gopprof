package structs

import (
	"encoding/gob"
	"errors"
	"time"

	"github.com/l2x/gopprof/common/event"
)

var (
	nodeEventChanSize = 10
)

func init() {
	gob.Register(NodeBase{})
	gob.Register(NodeConf{})
	gob.Register(ExInfo{})
}

type Node struct {
	evt chan *event.Event

	NodeConf
	NodeBase
}

func NewNode(nodeID string) *Node {
	return &Node{
		evt: make(chan *event.Event, nodeEventChanSize),
		NodeBase: NodeBase{
			NodeID: nodeID,
		},
	}
}

func (n *Node) Event() chan *event.Event {
	return n.evt
}

func (n *Node) AddEvent(evt *event.Event) error {
	select {
	case n.evt <- evt:
		return nil
	default:
		return errors.New("event chan full")
	}
}

type NodeBase struct {
	NodeID     string
	Hostname   string
	InternalIP []string
	ExternalIP []string
	BinMD5     string
	LastSync   time.Time
	Created    time.Time
}

type NodeConf struct {
	EnableProfile bool
	ProfileCron   string
	Profile       []string

	EnableStats bool
	StatsCron   string
}

type ExInfo struct {
	HTTPListen string
}
