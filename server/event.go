package server

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/l2x/gopprof/common/structs"
)

var eventFunc = map[structs.EventType]func(evtReq *structs.Event) (*structs.Event, error){
	structs.EventTypeNone:     eventNone,
	structs.EventTypeRegister: eventRegister,
	structs.EventTypeStat:     eventStat,
	structs.EventTypeBinCheck: eventBinCheck,
	structs.EventTypeCallback: eventCallback,
}

func eventProxy(evtReq *structs.Event) (*structs.Event, error) {
	f, ok := eventFunc[evtReq.Type]
	if !ok {
		err := fmt.Errorf("Unknown event: %v", evtReq.Type)
		logger.Error(err)
		return nil, err
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
		err := fmt.Errorf("Event data invalid: %#v", evtReq)
		logger.Error(err)
		return nil, err
	}
	nodeConf, err := storeSaver.GetConf(nodeBase.NodeID)
	if err == sql.ErrNoRows {
		nodeConf, err = storeSaver.GetDefaultConf()
	}
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	if err = storeSaver.SaveNode(&nodeBase); err != nil {
		return nil, err
	}

	node := nodesMap.Add(nodeBase.NodeID)
	node.NodeConf = *nodeConf
	node.LastSync = time.Now()
	node.Created = time.Now()

	info := structs.ExInfo{
		HTTPListen: conf.HTTPListen,
	}
	return structs.NewEvent(structs.EventTypeExInfo, info), nil
}

func eventNone(evtReq *structs.Event) (*structs.Event, error) {
	nodeID, ok := evtReq.Data.(string)
	if !ok {
		err := fmt.Errorf("Event data invalid: %#v", evtReq)
		logger.Error(err)
		return nil, err
	}
	node, ok := nodesMap.Get(nodeID)
	if !ok {
		logger.Warn("[eventNode] Node not registered, ", nodeID)
		return structs.NewEvent(structs.EventTypeRegister, nil), nil
	}
	node.LastSync = time.Now()

	select {
	case evt := <-node.Event():
		return evt, nil
	default:
	}

	evt, err := taskStats(node)
	if err != nil {
		return nil, err
	}
	if evt != nil {
		return evt, nil
	}

	evt, err = taskProfile(node)
	if err != nil {
		return nil, err
	}
	if evt != nil {
		return evt, nil
	}

	time.Sleep(conf.EventInterval)
	return nil, nil
}

func eventStat(evtReq *structs.Event) (*structs.Event, error) {
	data, ok := evtReq.Data.(structs.StatsData)
	if !ok {
		err := fmt.Errorf("Event data invalid: %#v", evtReq)
		logger.Error(err)
		return nil, err
	}
	if err := storeSaver.SaveStat(&data); err != nil {
		logger.Error(err)
		return nil, err
	}
	return nil, nil
}

func eventUploadProfile(evtReq *structs.Event) (*structs.Event, error) {
	data := evtReq.Data.(structs.ProfileData)
	if err := storeSaver.SaveProfile(&data); err != nil {
		logger.Error(err)
		return nil, err
	}
	return nil, nil
}

func eventCallback(evtReq *structs.Event) (*structs.Event, error) {
	return nil, nil
}

func taskProfile(node *structs.Node) (*structs.Event, error) {
	if node.EnableProfile == false {
		return nil, nil
	}
	if node.Created.Add(node.ProfileInterval).After(time.Now()) {
		return nil, nil
	}
	if node.LastProfile.Add(node.ProfileInterval).After(time.Now()) {
		return nil, nil
	}
	node.LastProfile = time.Now()
	return structs.NewEvent(structs.EventTypeProfile, node.Profile), nil
}

func taskStats(node *structs.Node) (*structs.Event, error) {
	if node.EnableStat == false {
		return nil, nil
	}
	if node.Created.Add(node.StatInterval).After(time.Now()) {
		return nil, nil
	}
	if node.LastStat.Add(node.StatInterval).After(time.Now()) {
		return nil, nil
	}
	node.LastStat = time.Now()
	return structs.NewEvent(structs.EventTypeStat, nil), nil
}

func eventBinCheck(evtReq *structs.Event) (*structs.Event, error) {
	exInfo, ok := evtReq.Data.(structs.ExInfo)
	if !ok {
		err := fmt.Errorf("Event data invalid: %#v", evtReq)
		logger.Error(err)
		return nil, err
	}
	_, err := storeSaver.GetBin(exInfo.NodeID, exInfo.MD5)
	if err != nil {
		return structs.NewEvent(structs.EventTypeUploadBin, nil), nil
	}
	return nil, nil
}
