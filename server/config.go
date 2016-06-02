package server

import (
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/astaxie/beego/config"
)

var (
	conf *Config
)

// Config read from config file
type Config struct {
	HTTPListen    string
	RPCListen     string
	LogPath       string
	EventInterval time.Duration
	StoreDriver   string
	StoreSource   string
	StatsDriver   string
	StatsSource   string
}

func initConfig(args []string) error {
	var f string
	for _, arg := range args {
		if strings.HasPrefix(arg, "-f=") {
			f = strings.Trim(arg, "-f=")
			break
		}
	}
	var data []byte
	var err error
	if f != "" {
		if data, err = ioutil.ReadFile(f); err != nil {
			return err
		}
	}

	// read config file
	cnf, err := config.NewConfigData("ini", data)
	if err != nil {
		return err
	}
	for _, arg := range args {
		arg = strings.TrimLeft(arg, "-")
		a := strings.SplitN(arg, "=", 2)
		if len(a) < 2 {
			fmt.Println("ignore:", arg)
			continue
		}
		if _, err := cnf.DIY(a[0]); err == nil {
			fmt.Println("use ", a[0], a[1])
		}
		if err := cnf.Set(a[0], a[1]); err != nil {
			fmt.Println(err)
			continue
		}
	}

	conf = &Config{}
	conf.HTTPListen = ":" + strings.TrimLeft(cnf.DefaultString("http_listen", ":8670"), ":")
	conf.RPCListen = ":" + strings.TrimLeft(cnf.DefaultString("rpc_listen", ":8671"), ":")
	conf.LogPath = cnf.DefaultString("log_path", "./log")
	conf.EventInterval = time.Duration(cnf.DefaultInt("event_interval", 60)) * time.Second
	conf.StoreDriver = cnf.DefaultString("store_driver", "bolt")
	conf.StoreSource = cnf.DefaultString("store_driver", "./database/bolt_store.db")
	conf.StatsDriver = cnf.DefaultString("stats_driver", "bolt")
	conf.StatsSource = cnf.DefaultString("stats_driver", "./database/bolt_stats.db")

	fmt.Printf("%#v \n", conf)

	return nil
}

func initLogger(path string) error {
	return nil
}
