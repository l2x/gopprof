package client

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/l2x/gopprof/common/structs"
)

var eventFunc = map[structs.EventType]func(client *Client, evtReq *structs.Event) (*structs.Event, error){
	structs.EventTypeNone:     eventNone,
	structs.EventTypeRegister: eventRegister,
	structs.EventTypeProfile:  eventProfile,
	structs.EventTypeStat:     eventStat,
	structs.EventTypeExInfo:   eventExInfo,
}

func eventProxy(client *Client, evtReq *structs.Event) (*structs.Event, error) {
	f, ok := eventFunc[evtReq.Type]
	if !ok {
		return nil, fmt.Errorf("Unknown event: %v", evtReq.Type)
	}

	evt, err := f(client, evtReq)
	if err != nil {
		return nil, err
	}
	return evt, nil
}

func eventRegister(client *Client, evtReq *structs.Event) (*structs.Event, error) {
	return structs.NewEvent(structs.EventTypeRegister, client.node.NodeBase), nil
}

func eventExInfo(client *Client, evtReq *structs.Event) (*structs.Event, error) {
	exInfo, ok := evtReq.Data.(structs.ExInfo)
	if !ok {
		return nil, fmt.Errorf("event data invalid: %#v", evtReq)
	}
	s := strings.Split(client.rpcServer, ":")[0]
	client.httpServer = fmt.Sprintf("http://%s:%s", s, strings.TrimLeft(exInfo.HTTPListen, ":"))
	return nil, nil
}

func eventNone(client *Client, evtReq *structs.Event) (*structs.Event, error) {
	return nil, nil
}

func eventProfile(client *Client, evtReq *structs.Event) (*structs.Event, error) {
	opts, ok := evtReq.Data.([]structs.ProfileData)
	if !ok {
		return nil, fmt.Errorf("event data invalid: %#v", evtReq)
	}
	for _, opt := range opts {
		file, err := StartProfile(&opt)
		if err != nil {
			opt.Status = 2
			opt.ErrMsg = err.Error()
			return structs.NewEvent(structs.EventTypeProfile, opt), err
		}
		defer func() {
			os.Remove(file)
		}()

		opt.Version = runtime.Version()
		opt.Created = time.Now().Unix()
		opt.Status = 1
		opt.NodeID = client.node.NodeID

		data, err := json.Marshal(opt)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		params := map[string]string{
			"data": string(data),
			"type": strconv.Itoa(int(structs.EventTypeUploadProfile)),
		}
		_, err = fileUpload(fmt.Sprintf("%s/upload", client.httpServer), file, params)
		if err != nil {
			log.Println(err)
			return nil, err
		}
	}
	return nil, nil
}

func eventStat(client *Client, evtReq *structs.Event) (*structs.Event, error) {
	data := StartStats()
	data.NodeID = client.node.NodeID
	return structs.NewEvent(structs.EventTypeStat, data), nil
}
