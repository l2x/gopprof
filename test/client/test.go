package main

import (
	"log"

	"github.com/l2x/gopprof/client"
)

var (
	Host   = "127.0.0.1:8981"
	NodeID = "node1"
)

func main() {
	if err := client.NewClient(Host, NodeID).Run(); err != nil {
		log.Fatal(err)
	}

	select {}
}
