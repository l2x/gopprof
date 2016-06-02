package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/l2x/gopprof/server"
)

var (
	cfg string
)

func main() {
	flag.StringVar(&cfg, "f", "gopprof.conf", "config file")
	flag.Parse()

	if err := server.Init(cfg); err != nil {
		log.Fatal(err)
	}
	defer server.Close()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signalChan
		server.Exit()
	}()

	server.Main()
}
