package client

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/l2x/gopprof/common/structs"
)

var eventFunc = map[structs.EventType]func(client *Client, evtReq *structs.Event) (*structs.Event, error){
	structs.EventTypeNone:     eventNone,
	structs.EventTypeRegister: eventRegister,
	structs.EventTypeProfile:  eventProfile,
	structs.EventTypeStat:     eventStat,
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
	evt := &structs.Event{
		Type: structs.EventTypeRegister,
		Data: client.node.NodeBase,
	}
	return evt, nil
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
			return structs.NewEvent(structs.EventTypeProfile, opt), err
		}
		defer func() {
			os.Remove(file)
		}()

		opt.File = file
		opt.Created = time.Now().Unix()
		opt.Status = 1
		data, err := json.Marshal(opt)
		if err != nil {
			return nil, err
		}
		params := map[string]string{
			"data": string(data),
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
