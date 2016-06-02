package client

import (
	"runtime"
	"time"

	"github.com/l2x/gopprof/common/structs"
)

// StartStats enables stats for the current process.
func StartStats() structs.StatsData {
	m := runtime.MemStats{}
	runtime.ReadMemStats(&m)
	return structs.StatsData{
		Created:      time.Now().Unix(),
		MemStats:     m,
		NumGoroutine: runtime.NumGoroutine(),
	}
}
