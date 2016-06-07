package server

import (
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
	logger.Info("exit")
	time.Sleep(1 * time.Second)
	os.Exit(1)
}

// Init at first
func Init(args []string) error {
	if err := initConfig(args); err != nil {
		return err
	}
	if err := initLogger(conf.LogPath, conf.Debug); err != nil {
		return err
	}
	if err := initStoreSaver(conf.StoreDriver, conf.StoreSource); err != nil {
		return err
	}
	if err := initFilesSaver(conf.FilesDriver, conf.FilesSource); err != nil {
		return err
	}
	return nil
}

// Close at last
func Close() {
	if storeSaver != nil {
		storeSaver.Close()
	}
	if filesSaver != nil {
		filesSaver.Close()
	}
	if logger != nil {
		logger.Flush()
	}
}
