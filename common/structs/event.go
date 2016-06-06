package structs

import (
	"sync"
	"time"
)

// EventType defined event type
type EventType int

const (
	EventTypeNone EventType = iota
	EventTypeRegister
	EventTypeCallback
	EventTypeInfo
)

const (
	EventTypeProfile EventType = iota + 100
)

const (
	EventTypeStat EventType = iota + 200
)

var (
	mu    sync.Mutex
	start int64
)

func init() {
	start = time.Now().UnixNano()
}

// NextSequence return autoincrementing integer
func NextSequence() int64 {
	mu.Lock()
	start++
	s := start
	mu.Unlock()
	return s
}

// Event struct
type Event struct {
	Type EventType
	Data interface{}
	Seq  int64
	Ack  int64
}

// NewEvent return default event
func NewEvent(typ EventType, data interface{}) *Event {
	evt := &Event{
		Type: typ,
		Data: data,
		Seq:  NextSequence(),
	}
	return evt
}

// NewCallback return callback event
func (e *Event) NewCallback(data interface{}) *Event {
	evt := &Event{
		Type: EventTypeCallback,
		Data: data,
		Ack:  e.Seq,
	}
	return evt
}
