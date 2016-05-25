package structs

import (
	"encoding/gob"
	"runtime"
)

func init() {
	gob.Register(Event{})
	gob.Register(Stats{})
	gob.Register(ProfileOption{})
	gob.Register(Node{})
	gob.Register(NodeBase{})
	gob.Register(NodeConf{})
}

// Stats records statistics about the memory allocator.
type Stats struct {
	Timestamp    int64
	NumGoroutine int
	runtime.MemStats
}

// ProfileOption is options for profiling.
type ProfileOption struct {
	Name  string
	Sleep int
	Debug int
	GC    bool
	Tmp   string
}

// NewProfileOption return ProfileOption with default value.
func NewProfileOption(name string) ProfileOption {
	return ProfileOption{
		Name:  name,
		Sleep: 30,
		Debug: 1,
		Tmp:   "/tmp",
	}
}

// Node storage node information
type Node struct {
	NodeBase
	NodeConf
}

// NodeBase storage node base information
type NodeBase struct {
	NodeID string
}

// NodeConf storage node config
type NodeConf struct {
}
