package influxdb

import (
	"database/sql"
	"encoding/json"
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
	_, err := b.query(fmt.Sprintf("CREATE DATABASE %s", b.dbname))
	if err != nil {
		return err
	}
	return nil
}

func (b *InfluxDB) query(cmd string) ([][]interface{}, error) {
	q := client.Query{
		Command:  cmd,
		Database: b.dbname,
	}
	response, err := b.db.Query(q)
	if err != nil {
		return nil, err
	}
	if response.Error() != nil {
		return nil, response.Error()
	}
	if len(response.Results) == 0 || len(response.Results[0].Series) == 0 || len(response.Results[0].Series[0].Values) == 0 {
		return nil, sql.ErrNoRows
	}
	res := response.Results[0].Series[0].Values
	return res, nil
}

func (b *InfluxDB) queryRow(cmd string) ([]interface{}, error) {
	rows, err := b.query(cmd)
	if err != nil {
		return nil, err
	}
	return rows[0], nil
}

func scan(row []interface{}, dest ...interface{}) error {
	if len(row)-1 < len(dest) {
		return fmt.Errorf("expected %d destination arguments in Scan, not %d", len(row), len(dest))
	}

	var r json.Number
	for i := 0; i < len(dest); i++ {
		switch d := dest[i].(type) {
		case *string:
			*d = row[i+1].(string)
		case *int64:
			r, _ = row[i+1].(json.Number)
			*d, _ = r.Int64()
		case *float64:
			r, _ = row[i+1].(json.Number)
			*d, _ = r.Float64()
		default:
			return fmt.Errorf("unsupported type: %v", d)
		}
	}
	return nil
}
