package server

import "time"

var (
	conf *Config
)

// Config read from config file
type Config struct {
	StoreDriver   string
	StoreSource   string
	LogPath       string
	HTTPListen    string
	RPCListen     string
	EventInterval time.Duration
}

func initConfig(cfg string) error {
	conf = &Config{
		EventInterval: 60 * time.Second,
		HTTPListen:    ":8081",
		RPCListen:     ":8082",
	}
	return nil
}

func initLogger(path string) error {
	return nil
}
