package event

import (
	"encoding/gob"
	"sync/atomic"
	"time"
)

const (
	EventTypeNone     = 0
	EventTypeRegister = 1
)

func init() {
	gob.Register(&Event{})
}

type EventType int

type Event struct {
	NodeID string
	Type   EventType
	Data   interface{}

	Ack int64
	Seq int64
}

var (
	seq = time.Now().UnixNano()
)

func NewEvent(nodeID string, typ EventType, data interface{}) *Event {
	return &Event{
		NodeID: nodeID,
		Type:   typ,
		Data:   data,
		Seq:    atomic.AddInt64(&seq, 1),
	}
}
