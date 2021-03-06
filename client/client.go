package client

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"net/rpc"
	"os"
	"reflect"
	"time"

	"gopkg.in/robfig/cron.v2"

	"github.com/l2x/gopprof/common/event"
	"github.com/l2x/gopprof/common/structs"
	"github.com/l2x/gopprof/common/util"
)

var (
	errInterval = 5 * time.Second
)

// Client is a gopprof client
type Client struct {
	rpc        *rpc.Client
	node       *structs.Node
	exInfo     structs.ExInfo
	serverAddr string
	job        *job
}

func NewClient(serverAddr, nodeID string) *Client {
	node := structs.NewNode(nodeID)
	if _, b, _ := util.GetBinFile(); b != nil {
		node.BinMD5 = fmt.Sprintf("%x", md5.Sum(b))
	}
	node.InternalIP, node.ExternalIP, _ = util.GetNetInterfaceIP()
	node.Hostname, _ = os.Hostname()
	c := &Client{
		node:       node,
		serverAddr: serverAddr,
	}
	return c
}

func (c *Client) Run() error {
	if err := c.connect(); err != nil {
		return err
	}
	if err := c.register(); err != nil {
		return err
	}
	c.job = &job{
		c:  c,
		cb: cron.New(),
	}
	c.job.cb.Start()

	go c.run()
	return nil
}

func (c *Client) register() error {
	evtReq := event.NewEvent(c.node.NodeID, event.EventTypeRegister, c.node.NodeBase)
	evtResp, err := c.sync(evtReq)
	if err != nil {
		log.Println("[gopprof/register] ", err)
		return err
	}
	// change extra info
	if evtResp == nil || evtResp.Type != event.EventTypeExInfo {
		return fmt.Errorf("[gopprof/register] incorrect response event: %#v", evtResp)
	}
	if _, err = eventExInfo(c, evtResp); err != nil {
		return err
	}
	return nil
}

func (c *Client) run() {
	var evtReq *event.Event
	for {
		if evtReq == nil {
			select {
			case evtReq = <-c.node.Event():
			default:
				evtReq = event.NewEvent(c.node.NodeID, event.EventTypeNone, nil)
			}
		}

		log.Println("[gopprof/evtReq] ", evtReq)
		evtResp, err := c.sync(evtReq)
		evtReq = nil
		if err != nil {
			log.Println("[gopprof/sync] ", err)
			time.Sleep(errInterval)
			continue
		}
		log.Println("[gopprof/evtResp] ", evtResp)

		evtReq, err = eventProxy(c, evtResp)
		if err != nil {
			log.Println("[gopprof/eventProxy] ", err)
			time.Sleep(errInterval)
			continue
		}
	}
}

func (c *Client) sync(evtReq *event.Event) (*event.Event, error) {
	evtResp := new(event.Event)
	if err := c.rpc.Call("RPCServer.Sync", evtReq, evtResp); err != nil {
		c.reconnect(err)
		return nil, err
	}
	return evtResp, nil
}

func (c *Client) connect() error {
	r, err := rpc.DialHTTP("tcp", c.serverAddr)
	if err != nil {
		return err
	}
	c.rpc = r
	return nil
}

func (c *Client) reconnect(e error) error {
	fmt.Println(reflect.TypeOf(e))
	if e != io.EOF && e != io.ErrUnexpectedEOF && e != rpc.ErrShutdown {
		return nil
	}
	return c.connect()
}
