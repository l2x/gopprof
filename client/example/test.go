package main

import (
	"log"

	"time"

	"github.com/l2x/gopprof/client"
)

func main() {
	go test()

	opt := client.NewProfileOption("cpu")
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
	}
}
