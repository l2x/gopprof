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

// ProfileData is data of profiling
type ProfileData struct {
	ID      int64
	NodeID  string
	Type    string
	Created int64
	File    string

	// status
	Status int // 0 - pending, 1 - success, 2 - failed
	ErrMsg string

	// option
	Sleep int
	Debug int
	GC    bool
}

// NewProfileData .
func NewProfileData() *ProfileData {
	return &ProfileData{}
}
