package structs

import "encoding/gob"

func init() {
	gob.Register(StatsData{})
}

// StatsData records statistics about the memory allocator.
type StatsData struct {
	NodeID       string
	Created      int64
	NumGoroutine uint
	HeapAlloc    uint64
	GCPauseNs    uint64
}
