package main

import (
	"flag"
	"log"
	"math/rand"
	"sync"
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
		log.Fatal(err)
	}
	go test()

	select {}
}

func test() {
	var wg sync.WaitGroup
	for i := 1; i < 100000; i++ {
		wg.Add(1)
		go func(i int) {
			runMem(i)
			wg.Done()
		}(i)
		wg.Add(1)
		go func(i int) {
			runGoroutine(i)
			wg.Done()
		}(i)
		wg.Wait()
	}
}

func runMem(i int) {
	b := make([]byte, i)
	b[i-1] = 'a'
	rand.Seed(time.Now().UnixNano())
	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
}

func runGoroutine(n int) {
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < rand.Intn(n); i++ {
		go func() {
			time.Sleep(1 * time.Second)
		}()
	}
	time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
}
