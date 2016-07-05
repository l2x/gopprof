package client

import (
	"fmt"

	"github.com/l2x/gopprof/common/event"
)

var eventFunc = map[event.EventType]func(c *Client, evtReq *event.Event) (*event.Event, error){
	event.EventTypeNone:     eventNone,
	event.EventTypeRegister: eventRegister,
}

func eventProxy(c *Client, evtReq *event.Event) (*event.Event, error) {
	f, ok := eventFunc[evtReq.Type]
	if !ok {
		return nil, fmt.Errorf("Unknown event: %v", evtReq.Type)
	}

	evt, err := f(c, evtReq)
	if err != nil {
		return nil, err
	}
	return evt, nil
}

func eventNone(c *Client, evtReq *event.Event) (*event.Event, error) {
	return nil, nil
}

func eventRegister(c *Client, evtReq *event.Event) (*event.Event, error) {
	return event.NewEvent(c.node.NodeID, event.EventTypeRegister, c.node.NodeBase), nil
}
