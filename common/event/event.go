package event

import (
	"encoding/gob"
	"sync/atomic"
	"time"
)

type EventType int

const (
	EventTypeNone          EventType = 0
	EventTypeRegister      EventType = 1
	EventTypeStats         EventType = 2
	EventTypeUploadProfile EventType = 3
	EventTypeUploadBin     EventType = 4
	EventTypeConf          EventType = 5
	EventTypeExInfo        EventType = 6
)

var (
	seq = time.Now().UnixNano()
)

func init() {
	gob.Register(&Event{})
}

type Event struct {
	NodeID string
	Type   EventType
	Data   interface{}

	Ack int64
	Seq int64
}

func NewEvent(nodeID string, typ EventType, data interface{}) *Event {
	return &Event{
		NodeID: nodeID,
		Type:   typ,
		Data:   data,
		Seq:    atomic.AddInt64(&seq, 1),
	}
}
