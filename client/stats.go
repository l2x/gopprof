package client

import "runtime"

// Stats records statistics about the memory allocator.
type Stats struct {
	NumGoroutine int
	runtime.MemStats
}

// StartStats enables stats for the current process.
func StartStats() Stats {
	m := runtime.MemStats{}
	runtime.ReadMemStats(&m)
	return Stats{
		MemStats:     m,
		NumGoroutine: runtime.NumGoroutine(),
	}
}
