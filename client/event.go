package client

import (
	"fmt"

	"github.com/l2x/gopprof/common/structs"
)

var eventFunc = map[structs.EventType]func(client *Client, evtReq *structs.Event) (*structs.Event, error){
	structs.EventTypeNone:     eventNone,
	structs.EventTypeRegister: eventRegister,
	structs.EventTypeProfile:  eventProfile,
	structs.EventTypeStat:     eventStat,
}

func eventProxy(client *Client, evtReq *structs.Event) (*structs.Event, error) {
	f, ok := eventFunc[evtReq.Type]
	if !ok {
		return nil, fmt.Errorf("Unknown event: %v", evtReq.Type)
	}

	evt, err := f(client, evtReq)
	if err != nil {
		return nil, err
	}
	return evt, nil
}

func eventRegister(client *Client, evtReq *structs.Event) (*structs.Event, error) {
	evt := &structs.Event{
		Type: structs.EventTypeRegister,
		Data: client.node.NodeBase,
	}
	return evt, nil
}

func eventNone(client *Client, evtReq *structs.Event) (*structs.Event, error) {
	return nil, nil
}

func eventProfile(client *Client, evtReq *structs.Event) (*structs.Event, error) {
	return nil, nil
}

func eventStat(client *Client, evtReq *structs.Event) (*structs.Event, error) {
	data := StartStats()
	data.NodeID = client.node.NodeID
	return &structs.Event{Type: structs.EventTypeStat, Data: data}, nil
}
