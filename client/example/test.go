package main

import (
	"log"

	"time"

	"github.com/l2x/gopprof/client"
)

func main() {
	go test()
	go test3()

	time.Sleep(1 * time.Second)

	opt := client.NewProfileOption("block")
	fname, err := client.StartProfile(opt)
	if err != nil {
		log.Println(err)
	}
	log.Println(fname)
	time.Sleep(2 * time.Second)
}

func test() {
	m := map[string]string{}
	for {
		m[time.Now().String()] = time.Now().String()
		time.Sleep(1 * time.Millisecond)
		go test2(time.Now().Unix())
	}
}

func test2(a int64) int64 {
	time.Sleep(5 * time.Second)
	return a + 1
}

func test3() {
	ch := make(chan int64, 100)
	for {
		ch <- time.Now().UnixNano()
		go func() {
			i := <-ch
			i = i + 1
			time.Sleep(1 * time.Millisecond)
		}()
	}
}
