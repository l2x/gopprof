package server

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"

	"github.com/l2x/gopprof/common/structs"
)

// ListenRPC start rpc server
func ListenRPC(port string) {
	log.Println("listen rpc:", port)
	rpcServer := new(RPCServer)
	rpc.Register(rpcServer)
	rpc.HandleHTTP()
	l, err := net.Listen("tcp", port)
	if err != nil {
		panic(err)
	}
	if err = http.Serve(l, nil); err != nil {
		panic(err)
	}
}

var eventFunc = map[structs.EventType]func(evtReq *structs.Event) (*structs.Event, error){
	structs.EventTypeNone:     eventNone,
	structs.EventTypeRegister: eventRegister,
}

// RPCServer is rpc server
type RPCServer struct{}

// Sync event
func (r *RPCServer) Sync(evtReq *structs.Event, evtResp *structs.Event) error {
	log.Println(evtReq)

	f, ok := eventFunc[evtReq.Type]
	if !ok {
		return fmt.Errorf("Unknown event: %v", evtReq.Type)
	}

	evt, err := f(evtReq)
	if err != nil {
		log.Println(err)
		return err
	}
	if evt != nil {
		*evtResp = *evt
	}
	return nil
}

func eventRegister(evtReq *structs.Event) (*structs.Event, error) {
	nodeBase, ok := evtReq.Data.(structs.NodeBase)
	if !ok {
		return nil, fmt.Errorf("Event data invalid: %#v", evtReq)
	}
	NodesMap.Add(nodeBase)
	return nil, nil
}

func eventNone(evtReq *structs.Event) (*structs.Event, error) {
	return nil, nil
}
