package server

import (
	"log"
	"os"
	"time"
)

// Main func
func Main() {
	go ListenHTTP(conf.HTTPListen)
	go ListenRPC(conf.RPCListen)
	select {}
}

// Exit func
func Exit() {
	log.Println("exit")
	time.Sleep(1 * time.Second)
	os.Exit(1)
}

// Init at first
func Init(args []string) error {
	if err := initConfig(args); err != nil {
		return err
	}
	if err := initLogger(conf.LogPath); err != nil {
		return err
	}
	if err := initStoreSaver(conf.StoreDriver, conf.StoreSource); err != nil {
		return err
	}
	return nil
}

// Close at last
func Close() {
	if storeSaver != nil {
		storeSaver.Close()
	}
	if statsSaver != nil {
		statsSaver.Close()
	}
}
