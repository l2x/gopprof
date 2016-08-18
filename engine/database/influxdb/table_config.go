package influxdb

import (
	"github.com/l2x/gopprof/common/structs"
	"github.com/l2x/gopprof/engine/database"
)

type TableConfig struct {
	b      *InfluxDB
	nodeID string
	table  string
}

func NewTableConfig(b *InfluxDB, nodeID string) database.TableConfig {
	return &TableConfig{
		b:      b,
		nodeID: nodeID,
		table:  "config",
	}
}

func (t *TableConfig) Table() []byte {
	return []byte(t.table)
}

func (t *TableConfig) Save(data *structs.NodeConf) error {
}

func (t *TableConfig) Get() (*structs.NodeConf, error) {
}

func (t *TableConfig) Goroots() ([]*structs.Goroot, error) {
}

func (t *TableConfig) GetGoroot(version string) (*structs.Goroot, error) {
}

func (t *TableConfig) SaveGoroot(goroot *structs.Goroot) error {
}

func (t *TableConfig) DelGoroot(goroot *structs.Goroot) error {
}
