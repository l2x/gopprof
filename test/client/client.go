package main

import (
	"flag"
	"log"
	"math/rand"
	"time"

	"github.com/l2x/gopprof/client"
)

var (
	Host   = ":8981"
	NodeID = "node1"
)

func main() {
	flag.StringVar(&Host, "h", ":8981", "host")
	flag.StringVar(&NodeID, "n", "node1", "nodeid")
	flag.Parse()

	if err := client.NewClient(Host, NodeID).Run(); err != nil {
		log.Println(err)
	}
	go test()

	select {}
}

func test() {
	for {
		runMem()
		runGoroutine()
	}
}

func runMem() {
	rand.Seed(time.Now().UnixNano())
	i := rand.Intn(1024 * 1024 * 10)

	b := make([]byte, i)
	b[i-1] = 'a'

	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
}

func runGoroutine() {
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(50)
	for i := 0; i < n; i++ {
		go func() {
			time.Sleep(1 * time.Second)
		}()
	}

	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
}
