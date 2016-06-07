package server

import (
	"net"
	"net/http"
	"net/rpc"

	"github.com/l2x/gopprof/common/structs"
)

// ListenRPC start rpc server
func ListenRPC(port string) {
	logger.Infof("listen rpc %s", port)
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

type RPCServer struct{}

func (r *RPCServer) Sync(evtReq *structs.Event, evtResp *structs.Event) error {
	logger.Debugf("evtReq[%#v]", evtReq)
	evt, err := eventProxy(evtReq)
	if err != nil {
		logger.Error(err)
		return err
	}
	if evt != nil {
		*evtResp = *evt
	}
	logger.Debugf("evtResp[%#v]", evtResp)
	return nil
}
