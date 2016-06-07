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
	if len(args) == 1 {
		switch args[0] {
		case "-v", "--v", "-version", "--version":
			fmt.Println(version)
			os.Exit(0)
		case "-h", "--h", "-help", "--help":
			fmt.Println("help")
			os.Exit(0)
		}
	}
	return args
}
