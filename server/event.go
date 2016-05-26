package server

import (
	"fmt"
	"log"
	"time"

	"github.com/l2x/gopprof/common/structs"
)

var eventFunc = map[structs.EventType]func(evtReq *structs.Event) (*structs.Event, error){
	structs.EventTypeNone:     eventNone,
	structs.EventTypeRegister: eventRegister,
}

func eventProxy(evtReq *structs.Event) (*structs.Event, error) {
	f, ok := eventFunc[evtReq.Type]
	if !ok {
		return nil, fmt.Errorf("Unknown event: %v", evtReq.Type)
	}

	evt, err := f(evtReq)
	if err != nil {
		return nil, err
	}
	return evt, nil
}

func eventRegister(evtReq *structs.Event) (*structs.Event, error) {
	nodeBase, ok := evtReq.Data.(structs.NodeBase)
	if !ok {
		return nil, fmt.Errorf("Event data invalid: %#v", evtReq)
	}
	NodesMap.Add(nodeBase.NodeID)

	// TODO get node conf

	return nil, nil
}

func eventNone(evtReq *structs.Event) (*structs.Event, error) {
	nodeBase, ok := evtReq.Data.(structs.NodeBase)
	if !ok {
		return nil, fmt.Errorf("Event data invalid: %#v", evtReq)
	}
	node, ok := NodesMap.Get(nodeBase.NodeID)
	if !ok {
		log.Println("please register")
		return &structs.Event{Type: structs.EventTypeRegister}, nil
	}
	// TODO check task

	_ = node

	time.Sleep(60 * time.Second)
	return nil, nil
}

func eventStat(evtReq *structs.Event) (*structs.Event, error) {
	return nil, nil
}
