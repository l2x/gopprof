package server

import (
	"fmt"
	"time"

	"github.com/l2x/gopprof/common/event"
	"github.com/l2x/gopprof/common/structs"
)

// register event handler function
var eventFunc = map[event.EventType]func(node *structs.Node, evtReq *event.Event) (*event.Event, error){
	event.EventTypeNone:     eventNone,
	event.EventTypeRegister: eventRegister,
	event.EventTypeStats:    eventStats,
}

func EventProxy(evt *event.Event) (*event.Event, error) {
	f, ok := eventFunc[evt.Type]
	if !ok {
		err := fmt.Errorf("Unknown event: %v", evt.Type)
		logger.Error(err)
		return nil, err
	}

	var node *structs.Node
	if evt.Type != event.EventTypeRegister {
		if node, ok = nodesMap.Get(evt.NodeID); !ok {
			logger.Warnf("Node not registered: %s", evt.NodeID)
			return event.NewEvent(evt.NodeID, event.EventTypeRegister, nil), nil
		}
	}

	evt, err := f(node, evt)
	if err != nil {
		return nil, err
	}
	return evt, nil
}

func eventNone(node *structs.Node, evtReq *event.Event) (*event.Event, error) {
	node.LastSync = time.Now()
	timer := time.NewTimer(conf.EventInterval)
	select {
	case <-timer.C:
		return nil, nil
	case evt := <-node.Event():
		return evt, nil
	}
	return nil, nil
}

func eventRegister(node *structs.Node, evtReq *event.Event) (*event.Event, error) {
	nodeBase, ok := evtReq.Data.(structs.NodeBase)
	if !ok {
		err := fmt.Errorf("Event data invalid: %#v", evtReq)
		logger.Error(err)
		return nil, err
	}

	// get node profile and stats configs
	nodeConf, err := db.TableConfig(nodeBase.NodeID).Get()
	if err != nil {
		nodeConf, _ = db.TableConfig(nodeBase.NodeID).GetDefault()
	}

	// save node
	if err = db.TableNode(nodeBase.NodeID).Save(&nodeBase); err != nil {
		return nil, err
	}

	node = nodesMap.Add(nodeBase.NodeID)
	node.NodeBase = nodeBase
	node.NodeConf = *nodeConf
	node.LastSync = time.Now()
	node.Created = time.Now()

	if !checkBinFileExist(node.NodeID, node.BinMD5) {
		node.AddEvent(event.NewEvent(node.NodeID, event.EventTypeUploadBin, nil))
	}
	if node.EnableProfile || node.EnableStats {
		node.AddEvent(event.NewEvent(node.NodeID, event.EventTypeConf, node.NodeConf))
	}
	exInfo := structs.ExInfo{
		HTTPListen: conf.HTTPListen,
	}
	return event.NewEvent(node.NodeID, event.EventTypeExInfo, exInfo), nil
}

func eventStats(node *structs.Node, evtReq *event.Event) (*event.Event, error) {
	data, ok := evtReq.Data.(structs.StatsData)
	if !ok {
		err := fmt.Errorf("Event data invalid: %#v", evtReq)
		logger.Error(err)
		return nil, err
	}

	if err := db.TableStats(evtReq.NodeID).Save(&data); err != nil {
		logger.Error(err)
		return nil, err
	}
	return nil, nil
}

func eventUploadProfile(data *structs.ProfileData) error {
	if err := db.TableProfile(data.NodeID).Save(data); err != nil {
		logger.Error(err)
		return err
	}
	return nil
}

func eventUploadBin(nodeID, binMD5, file string) error {
	if err := db.TableBin(nodeID).Save(binMD5, file); err != nil {
		logger.Error(err)
		return err
	}
	return nil
}

// Check binary file is already exists.
// It will be used in profiling data analysis.
func checkBinFileExist(nodeID string, binMD5 string) bool {
	_, err := db.TableBin(nodeID).Get(binMD5)
	if err != nil {
		return false
	}
	return true
}
