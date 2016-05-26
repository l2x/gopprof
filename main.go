package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/l2x/gopprof/server"
)

func main() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signalChan
		server.Exit()
	}()

	server.Main()
}
