package client

import (
	"fmt"
	"io"
	"log"
	"net/rpc"
	"time"

	"github.com/l2x/gopprof/common/structs"
	"github.com/l2x/gopprof/common/utils"
)

// Client is a gopprof client
type Client struct {
	rpc        *rpc.Client
	rpcServer  string
	httpServer string
	node       *structs.Node
}

// NewClient return client
func NewClient(rpcServer, nodeID string) *Client {
	node := structs.NewNode(nodeID)
	node.InternalIP, _ = utils.GetInternalIP()
	c := &Client{
		rpcServer: rpcServer,
		node:      node,
	}
	return c
}

// Run an client
func (c *Client) Run() error {
	if err := c.connect(); err != nil {
		return err
	}
	if err := c.register(); err != nil {
		return err
	}
	go c.run()
	return nil
}

func (c *Client) register() error {
	evtReq := &structs.Event{
		Type: structs.EventTypeRegister,
		Data: c.node.NodeBase,
	}
	evtResp, err := c.sync(evtReq)
	if err != nil {
		log.Println("[register]", err)
		return err
	}
	if evtResp.Type != structs.EventTypeExInfo {
		return fmt.Errorf("incorrect response event: %d", evtResp.Type)
	}
	if _, err = eventExInfo(c, evtResp); err != nil {
		return err
	}
	return nil
}

func (c *Client) run() {
	var evtReq *structs.Event
	for {
		if evtReq == nil {
			select {
			case evtReq = <-c.node.Event():
			default:
				evtReq = structs.NewEvent(structs.EventTypeNone, c.node.NodeID)
			}
		}

		log.Println("[evtReq]", evtReq)
		evtResp, err := c.sync(evtReq)
		evtReq = nil
		if err != nil {
			log.Println("[sync]", err)
			time.Sleep(5 * time.Second)
			continue
		}
		log.Println("[evtResp]", evtResp)

		evtReq, err = eventProxy(c, evtResp)
		if err != nil {
			time.Sleep(5 * time.Second)
			log.Println("[eventProxy]", err)
			continue
		}
	}
}

func (c *Client) sync(evtReq *structs.Event) (*structs.Event, error) {
	evtResp := new(structs.Event)
	if err := c.rpc.Call("RPCServer.Sync", evtReq, evtResp); err != nil {
		c.reconnect(err)
		return nil, err
	}
	return evtResp, nil
}

func (c *Client) connect() error {
	r, err := rpc.DialHTTP("tcp", c.rpcServer)
	if err != nil {
		return err
	}
	c.rpc = r
	return nil
}

func (c *Client) reconnect(e error) {
	if e != io.EOF && e != io.ErrUnexpectedEOF && e != rpc.ErrShutdown {
		return
	}
	cl, err := rpc.DialHTTP("tcp", c.rpcServer)
	if err != nil {
		return
	}
	*c.rpc = *cl
}
