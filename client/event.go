package client

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/l2x/gopprof/common/event"
	"github.com/l2x/gopprof/common/structs"
	"github.com/l2x/gopprof/common/util"
)

var eventFunc = map[event.EventType]func(c *Client, evtReq *event.Event) (*event.Event, error){
	event.EventTypeNone:          eventNone,
	event.EventTypeRegister:      eventRegister,
	event.EventTypeConf:          eventConf,
	event.EventTypeUploadProfile: eventUploadProfile,
	event.EventTypeUploadBin:     eventUploadBin,
	event.EventTypeStats:         eventStats,
	event.EventTypeExInfo:        eventExInfo,
}

func eventProxy(c *Client, evtReq *event.Event) (*event.Event, error) {
	f, ok := eventFunc[evtReq.Type]
	if !ok {
		return nil, fmt.Errorf("Unknown event: %v", evtReq.Type)
	}

	evt, err := f(c, evtReq)
	if err != nil {
		return nil, err
	}
	return evt, nil
}

func eventNone(c *Client, evtReq *event.Event) (*event.Event, error) {
	return nil, nil
}

func eventRegister(c *Client, evtReq *event.Event) (*event.Event, error) {
	return event.NewEvent(c.node.NodeID, event.EventTypeRegister, c.node.NodeBase), nil
}

func eventConf(c *Client, evtReq *event.Event) (*event.Event, error) {
	nodeConf, ok := evtReq.Data.(structs.NodeConf)
	if !ok {
		err := fmt.Errorf("Event data invalid: %#v", evtReq)
		log.Println(err)
		return nil, err
	}
	c.node.NodeConf = nodeConf
	c.job.reload()
	return nil, nil
}

func eventStats(c *Client, evtReq *event.Event) (*event.Event, error) {
	return event.NewEvent(c.node.NodeID, event.EventTypeStats, StartStats()), nil
}

func eventUploadProfile(c *Client, evtReq *event.Event) (*event.Event, error) {
	return nil, nil
}

func uploadProfile(c *Client, nodeConf structs.NodeConf) error {
	if !nodeConf.EnableProfile {
		return nil
	}

	for _, typ := range nodeConf.Profile {
		f, err := StartProfile(typ)
		if err != nil {
			log.Println(err)
			continue
		}
		defer func() {
			os.Remove(f)
		}()

		data := structs.ProfileData{
			Created:   time.Now().Unix(),
			NodeID:    c.node.NodeID,
			Type:      typ,
			File:      f,
			BinMD5:    c.node.BinMD5,
			GoVersion: runtime.Version(),
		}
		b, _ := json.Marshal(data)
		params := map[string]string{
			"data": string(b),
			"type": strconv.Itoa(int(event.EventTypeUploadProfile)),
		}
		if _, err = util.PostFile(fmt.Sprintf("%s/upload", c.exInfo.HTTPListen), f, params); err != nil {
			log.Println(err)
			continue
		}
		time.Sleep(1 * time.Second)
	}
	return nil
}

func eventUploadBin(c *Client, evtReq *event.Event) (*event.Event, error) {
	f, b, err := util.GetBinFile()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	params := map[string]string{
		"nodeid":  c.node.NodeID,
		"bin_md5": c.node.BinMD5,
		"type":    strconv.Itoa(int(event.EventTypeUploadBin)),
	}
	_, err = util.PostData(fmt.Sprintf("%s/upload", c.exInfo.HTTPListen), filepath.Base(f), b, params)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return nil, nil
}

func eventExInfo(c *Client, evtReq *event.Event) (*event.Event, error) {
	exInfo, ok := evtReq.Data.(structs.ExInfo)
	if !ok {
		return nil, fmt.Errorf("[gopprof/register] response event data invalid: %#v", evtReq)
	}
	exInfo.HTTPListen = fmt.Sprintf("http://%s:%s", strings.Split(c.serverAddr, ":")[0], strings.TrimLeft(exInfo.HTTPListen, ":"))
	c.exInfo = exInfo
	return nil, nil
}
