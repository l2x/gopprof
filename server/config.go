package server

import (
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/astaxie/beego/config"
)

var (
	conf *Config
)

// Config read from config file
type Config struct {
	Debug         bool
	LogPath       string
	HTTPListen    string
	RPCListen     string
	EventInterval time.Duration
	StoreDriver   string
	StoreSource   string
	DBDriver      string
	DBSource      string
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
			log.Println(err)
			return err
		}
	}

	// read config file
	cnf, err := config.NewConfigData("ini", data)
	if err != nil {
		log.Println(err)
		return err
	}
	for _, arg := range args {
		arg = strings.TrimLeft(arg, "-")
		a := strings.SplitN(arg, "=", 2)
		if len(a) < 2 {
			log.Println("ignore:", arg)
			continue
		}
		if _, err := cnf.DIY(a[0]); err == nil {
			log.Println("use ", a[0], a[1])
		}
		if err := cnf.Set(a[0], a[1]); err != nil {
			log.Println(err)
			continue
		}
	}

	conf = &Config{}
	conf.Debug = cnf.DefaultBool("debug", false)
	conf.LogPath = cnf.DefaultString("log_path", "./logs")
	conf.HTTPListen = ":" + strings.TrimLeft(cnf.DefaultString("http_listen", ":8980"), ":")
	conf.RPCListen = ":" + strings.TrimLeft(cnf.DefaultString("rpc_listen", ":8981"), ":")
	conf.EventInterval = time.Duration(cnf.DefaultInt("event_interval", 10)) * time.Second

	conf.DBDriver = cnf.DefaultString("db_driver", "bolt")
	conf.DBSource = cnf.DefaultString("db_source", "./data/db/bolt.db")
	conf.StoreDriver = cnf.DefaultString("store_driver", "localfile")
	conf.StoreSource = cnf.DefaultString("store_source", "./data/files")

	// TODO
	conf.Debug = true

	log.Printf("%#v \n", conf)
	return nil
}
