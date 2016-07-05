package server

import (
	"fmt"

	"github.com/l2x/gopprof/common/event"
	"github.com/valyala/gorpc"
)

// ListenRPC start RPC server
func ListenRPC(port string) {
	s := &gorpc.Server{
		Addr:    port,
		Handler: sync,
	}
	if err := s.Serve(); err != nil {
		logger.Criticalf("Cannot start rpc server: %s", err)
	}
	Exit()
}

func sync(clientAddr string, request interface{}) interface{} {
	logger.Debugf("evtReq[%s][%#v]", clientAddr, request)

	evt, ok := request.(*event.Event)
	if !ok {
		logger.Errorf("event invalid[%#v]", request)
		return nil
	}
	response, err := eventProxy(evt)
	if err != nil {
		return nil
	}
	if response == nil {
		response = event.NewEvent("", event.EventTypeNone, nil)
	}

	logger.Debugf("evtResp[%#v]", response)
	return response
}

var eventFunc = map[event.EventType]func(evtReq *event.Event) (*event.Event, error){
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
	evt, err := f(evt)
	if err != nil {
		return nil, err
	}
	return evt, nil
}

func eventNone(evtReq *event.Event) (*event.Event, error) {
	return nil, nil
}

func eventRegister(evtReq *event.Event) (*event.Event, error) {
	return nil, nil
}
