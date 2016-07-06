package server

import (
	"net"
	"net/http"
	"net/rpc"

	"github.com/l2x/gopprof/common/event"
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
	evt, err := EventProxy(evtReq)
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
