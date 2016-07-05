package server

import (
	"os"
	"time"
)

// Init first
func Init(args []string) error {
	if err := initConfig(args); err != nil {
		return err
	}
	if err := initLogger(conf.LogPath, conf.Debug); err != nil {
		return err
	}
	if err := initDB(conf.DBDriver, conf.DBSource); err != nil {
		return err
	}
	if err := initStore(conf.StoreDriver, conf.StoreSource); err != nil {
		return err
	}
	return nil
}

// Main func
func Main() {
	go ListenHTTP(conf.HTTPListen)
	go ListenRPC(conf.RPCListen)
	select {}
}

// Exit func
func Exit(signal ...string) {
	logger.Info("exit:", signal)

	if db != nil {
		db.Close()
	}
	if store != nil {
		store.Close()
	}
	if logger != nil {
		logger.Flush()
	}
	time.Sleep(1 * time.Second)
	os.Exit(0)
}
