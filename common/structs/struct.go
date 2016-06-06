package structs

import (
	"encoding/gob"
	"runtime"
)

func init() {
	gob.Register(Event{})
	gob.Register(StatsData{})
	gob.Register(ProfileData{})
	gob.Register([]ProfileData{})
	gob.Register(Node{})
	gob.Register(NodeBase{})
	gob.Register(NodeConf{})
}

// StatsData records statistics about the memory allocator.
type StatsData struct {
	ID           int64
	NodeID       string
	Created      int64
	NumGoroutine int
	runtime.MemStats
}
