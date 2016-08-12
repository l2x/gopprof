package influxdb

import (
	"fmt"
	"time"

	"github.com/influxdata/influxdb/client/v2"
	"github.com/l2x/gopprof/engine/database"
)

type TableBin struct {
	b      *InfluxDB
	nodeID string
	table  string
}

func NewTableBin(b *InfluxDB, nodeID string) database.TableBin {
	return &TableBin{
		b:      b,
		nodeID: nodeID,
		table:  "bin",
	}
}

func (t *TableBin) Table() []byte {
	return []byte(t.table)
}

func (t *TableBin) Save(binMD5, file string) error {
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{})
	if err != nil {
		return err
	}
	tags := map[string]string{
		"nodeID": t.nodeID,
	}
	fields := map[string]interface{}{
		"md5":  binMD5,
		"file": file,
	}
	pt, err := client.NewPoint(t.table, tags, fields, time.Now())
	if err != nil {
		return err
	}
	bp.AddPoint(pt)
	return t.b.db.Write(bp)
}

func (t *TableBin) Get(binMD5 string) (string, error) {
	q := fmt.Sprintf("SELECT * FROM %s WHERE nodeID='%s' and md5='%s' LIMIT 1", t.table, t.nodeID, binMD5)
	res, err := t.b.queryDB(q)
	if err != nil {
		return "", err
	}

}
