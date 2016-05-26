package server

import (
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

// RPCServer is rpc server
type RPCServer struct{}

// Sync event
func (r *RPCServer) Sync(evtReq *structs.Event, evtResp *structs.Event) error {
	log.Println(evtReq)

	evt, err := eventProxy(evtReq)
	if err != nil {
		return err
	}
	if evt != nil {
		*evtResp = *evt
	}
	return nil
}
