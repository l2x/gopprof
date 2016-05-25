package client

import (
	"runtime"

	"github.com/l2x/gopprof/common/structs"
)

// StartStats enables stats for the current process.
func StartStats() structs.Stats {
	m := runtime.MemStats{}
	runtime.ReadMemStats(&m)
	return structs.Stats{
		MemStats:     m,
		NumGoroutine: runtime.NumGoroutine(),
	}
}
