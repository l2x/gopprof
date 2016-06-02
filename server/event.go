package server

import (
	"database/sql"
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

	nodeConf, err := storeSaver.GetNode(nodeBase.NodeID)
	if err == sql.ErrNoRows {
		nodeConf, err = storeSaver.GetDefault()
	}
	if err != nil {
		log.Println(err)
		return nil, err
	}

	node := nodesMap.Add(nodeBase.NodeID)
	node.NodeConf = *nodeConf

	return nil, nil
}

func eventNone(evtReq *structs.Event) (*structs.Event, error) {
	nodeID, ok := evtReq.Data.(string)
	if !ok {
		return nil, fmt.Errorf("Event data invalid: %#v", evtReq)
	}
	node, ok := nodesMap.Get(nodeID)
	if !ok {
		log.Println("[eventNode] Node not registered, ", nodeID)
		return &structs.Event{Type: structs.EventTypeRegister}, nil
	}

	select {
	case evt := <-node.Event():
		return evt, nil
	default:
	}

	evt, err := taskStats(node)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if evt != nil {
		return evt, nil
	}

	evt, err = taskProfile(node)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if evt != nil {
		return evt, nil
	}

	time.Sleep(conf.EventInterval)
	return nil, nil
}

func eventStat(evtReq *structs.Event) (*structs.Event, error) {
	return nil, nil
}

func eventProfile(evtReq *structs.Event) (*structs.Event, error) {
	return nil, nil
}

func taskProfile(node *structs.Node) (*structs.Event, error) {
	if node.EnableProfile == false {
		return nil, nil
	}
	if node.LastProfile.Add(node.ProfileInterval).After(time.Now()) {
		return nil, nil
	}
	node.LastProfile = time.Now()
	return &structs.Event{Type: structs.EventTypeProfile, Data: node.ProfileName}, nil
}

func taskStats(node *structs.Node) (*structs.Event, error) {
	if node.EnableStat == false {
		return nil, nil
	}
	if node.LastStat.Add(node.StatInterval).After(time.Now()) {
		return nil, nil
	}
	node.LastStat = time.Now()
	return &structs.Event{Type: structs.EventTypeStat}, nil
}
