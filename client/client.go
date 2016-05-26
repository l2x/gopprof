package client

import (
	"io"
	"log"
	"net/rpc"
	"time"

	"github.com/l2x/gopprof/common/structs"
)

// Client is a gopprof client
type Client struct {
	rpc    *rpc.Client
	server string
	node   *structs.Node
}

// NewClient return client
func NewClient(server, nodeID string) *Client {
	node := structs.NewNode(nodeID)
	c := &Client{
		server: server,
		node:   node,
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
	_, err := c.sync(evtReq)
	if err != nil {
		log.Println("[register]", err)
		return err
	}
	return nil
}

func (c *Client) run() {
	for {
		var evtReq *structs.Event
		select {
		case evtReq = <-c.node.Event():
		default:
			evtReq = structs.NewEvent()
			evtReq.Data = c.node.NodeID
		}
		evtResp, err := c.sync(evtReq)
		if err != nil {
			log.Println("[sync]", err)
			time.Sleep(5 * time.Second)
			continue
		}
		evt, err := eventProxy(c, evtResp)
		if err != nil {
			log.Println("[eventProxy]", err)
			continue
		}
		if evt != nil {
			c.node.AddEvent(evt)
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
	r, err := rpc.DialHTTP("tcp", c.server)
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
	cl, err := rpc.DialHTTP("tcp", c.server)
	if err != nil {
		return
	}
	*c.rpc = *cl
}
