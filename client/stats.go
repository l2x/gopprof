package client

import (
	"runtime"
	"time"

	"github.com/l2x/gopprof/common/structs"
)

// StartStats start stats for the current process.
func StartStats() structs.StatsData {
	m := runtime.MemStats{}
	runtime.ReadMemStats(&m)
	return structs.StatsData{
		Created:      time.Now().Unix(),
		NumGoroutine: uint(runtime.NumGoroutine()),
		HeapAlloc:    m.HeapAlloc,
		GCPauseNs:    m.PauseNs[(m.NumGC+255)%256],
	}
}
