package influxdb

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/influxdata/influxdb/client/v2"
	"github.com/l2x/gopprof/engine/database"
)

func init() {
	database.Register("influxdb", NewInfluxDB)
}

type InfluxDB struct {
	db     client.Client
	dbname string
}

func NewInfluxDB() database.Database {
	return &InfluxDB{}
}

// Open opens influxdb
// source format:
//  http client - http:addr=localhost:port&database=dbname&username=user&password=pwd
//  udp client  - udp:addr=localhost:port&database=dbname
func (b *InfluxDB) Open(source string) error {
	s := strings.SplitN(source, ":", 2)
	if len(s) != 2 {
		return errors.New("invalid influxdb source")
	}

	var err error
	proto := strings.ToLower(s[0])
	q, err := url.Parse(s[1])
	if err != nil {
		return errors.New("invalid influxdb source")
	}
	addr := q.Query().Get("addr")
	if addr == "" {
		return errors.New("invalid influxdb source: addr empty")
	}
	b.dbname = q.Query().Get("database")
	if b.dbname == "" {
		return errors.New("invalid influxdb source: database empty")
	}

	switch proto {
	case "udp":
		cfg := client.UDPConfig{
			Addr: addr,
		}
		if b.db, err = client.NewUDPClient(cfg); err != nil {
			return err
		}
	case "http":
		cfg := client.HTTPConfig{
			Addr:     addr,
			Username: q.Query().Get("username"),
			Password: q.Query().Get("password"),
		}
		if b.db, err = client.NewHTTPClient(cfg); err != nil {
			return err
		}
	default:
		return errors.New("invalid influxdb source: unsupport protocol")
	}

	return b.init()
}

func (b *InfluxDB) Close() error {
	b.db.Close()
	return nil
}

func (b *InfluxDB) TableStats(nodeID string) database.TableStats {
	return NewTableStats(b, nodeID)
}

func (b *InfluxDB) TableProfile(nodeID string) database.TableProfile {
	return NewTableProfile(b, nodeID)
}

func (b *InfluxDB) TableConfig(nodeID string) database.TableConfig {
	return NewTableConfig(b, nodeID)
}

func (b *InfluxDB) TableNode(nodeID string) database.TableNode {
	return NewTableNode(b, nodeID)
}

func (b *InfluxDB) TableBin(nodeID string) database.TableBin {
	return NewTableBin(b, nodeID)
}

func (b *InfluxDB) init() error {
	_, err := b.queryDB(fmt.Sprintf("CREATE DATABASE %s", b.dbname))
	if err != nil {
		return err
	}
	return nil
}

func (b *InfluxDB) queryDB(cmd string) (res []client.Result, err error) {
	q := client.Query{
		Command:  cmd,
		Database: b.dbname,
	}
	if response, err := b.db.Query(q); err == nil {
		if response.Error() != nil {
			return res, response.Error()
		}
		res = response.Results
	} else {
		return res, err
	}
	return res, nil
}
