package main

import (
	"io"
	"log"
	"net/rpc"
	"time"

	"github.com/l2x/gopprof/common/structs"
)

var (
	Host = "127.0.0.1:8081"
)

func main() {
	test()
}

func reconnect(c *rpc.Client, e error) {
	if e != io.EOF && e != io.ErrUnexpectedEOF && e != rpc.ErrShutdown {
		return
	}
	cl, err := rpc.DialHTTP("tcp", Host)
	if err != nil {
		time.Sleep(1 * time.Second)
		return
	}
	*c = *cl
}

func test() {
	c, err := rpc.DialHTTP("tcp", Host)
	if err != nil {
		log.Fatal(err)
	}

	for {
		time.Sleep(1 * time.Second)
		evtReq := &structs.Event{
			Type: structs.EventTypeRegister,
			Data: structs.NodeBase{
				NodeID: "nodeid1111",
			},
		}
		evtResp := new(structs.Event)
		if err := c.Call("RPCServer.Sync", evtReq, evtResp); err != nil {
			reconnect(c, err)
			log.Println(err)
			continue
		}

		log.Println(evtReq)
		log.Println(evtResp)
	}
}
