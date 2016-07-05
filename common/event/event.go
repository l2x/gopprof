package event

import (
	"encoding/gob"
	"sync/atomic"
	"time"
)

const (
	EventTypeNone          = 0
	EventTypeRegister      = 1
	EventTypeStats         = 2
	EventTypeUploadProfile = 3
	EventTypeUploadBin     = 4
	EventTypeConf          = 5
	EventTypeExInfo        = 6
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
