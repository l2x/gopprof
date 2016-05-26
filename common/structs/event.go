package structs

// EventType defined event type
type EventType int

const (
	// EventTypeNone is default event
	EventTypeNone EventType = iota
	// EventTypeRegister is register event
	EventTypeRegister
)

const (
	// EventTypeProfile is profiling event
	EventTypeProfile EventType = iota + 100
)

const (
	// EventTypeStat is stats event
	EventTypeStat EventType = iota + 200
)

// Event struct
type Event struct {
	Type EventType
	Data interface{}
}

// NewEvent return default event
func NewEvent() *Event {
	return &Event{Type: EventTypeNone}
}
