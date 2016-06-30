package structs

import (
	"sync"
	"time"
)

const (
	EventTypeNone          EventType = 0
	EventTypeRegister      EventType = 1
	EventTypeCallback      EventType = 2
	EventTypeExInfo        EventType = 3
	EventTypeProfile       EventType = 4
	EventTypeStat          EventType = 5
	EventTypeUploadProfile EventType = 6
	EventTypeUploadBin     EventType = 7
)

var (
	mu    sync.Mutex
	start int64
)

// EventType defined event
type EventType int

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

type Event struct {
	Type EventType
	Data interface{}
	Seq  int64
	Ack  int64
}

// NewEvent return event with autoincrementing integer sequence
func NewEvent(typ EventType, data interface{}) *Event {
	evt := &Event{
		Type: typ,
		Data: data,
		Seq:  NextSequence(),
	}
	return evt
}

// NewCallback return event with Ack
func (e *Event) NewCallback(data interface{}) *Event {
	evt := &Event{
		Type: EventTypeCallback,
		Data: data,
		Ack:  e.Seq,
	}
	return evt
}
