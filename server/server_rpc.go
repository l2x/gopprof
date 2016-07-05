package server

import (
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"net/rpc"
	"time"

	"github.com/l2x/gopprof/common/event"
	"github.com/l2x/gopprof/common/structs"
)

// ListenRPC start RPC server
func ListenRPC(port string) {
	logger.Infof("listen rpc %s", port)
	rpcServer := new(RPCServer)
	rpc.Register(rpcServer)
	rpc.HandleHTTP()
	l, err := net.Listen("tcp", port)
	if err != nil {
		logger.Criticalf("Cannot start rpc server: %s", err)
		Exit()
	}
	if err = http.Serve(l, nil); err != nil {
		logger.Criticalf("Cannot start rpc server: %s", err)
	}
	Exit()
}

type RPCServer struct{}

func (r *RPCServer) Sync(evtReq *event.Event, evtResp *event.Event) error {
	logger.Debugf("evtReq[%#v]", evtReq)
	evt, err := eventProxy(evtReq)
	if err != nil {
		logger.Error(err)
		return err
	}
	logger.Debugf("evtResp[%#v]", evt)

	if evt != nil {
		*evtResp = *evt
	}
	return nil
}

// register event handler function
var eventFunc = map[event.EventType]func(node *structs.Node, evtReq *event.Event) (*event.Event, error){
	event.EventTypeNone:     eventNone,
	event.EventTypeRegister: eventRegister,
}

func eventProxy(evt *event.Event) (*event.Event, error) {
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
			return event.NewEvent("", event.EventTypeRegister, nil), nil
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
	if err == sql.ErrNoRows {
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

// Check binary file is already exists.
// It will be used in profiling data analysis.
func checkBinFileExist(nodeID string, binMD5 string) bool {
	return false
}
