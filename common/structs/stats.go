package structs

import "encoding/gob"

func init() {
	gob.Register(StatsData{})
}

type StatsData struct {
	NodeID  string
	Created int64
}
