package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/l2x/gopprof/server"
)

func main() {
	if err := server.Init(args()); err != nil {
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

func args() []string {
	args := os.Args[1:]
	if len(args) == 1 && (args[0] == "-v" || args[0] == "--version") {
		fmt.Println(version)
		os.Exit(0)
	}
	if len(args) == 1 && (args[0] == "-h" || args[0] == "-help") {
		fmt.Println("help")
		os.Exit(0)
	}
	return args
}
