package client

import (
	"fmt"
	"os"
	"runtime/pprof"
	"runtime/trace"
	"time"

	"github.com/l2x/gopprof/common/structs"
)

var pm = map[string]func(fn *os.File, option *structs.ProfileData) error{
	"cpu":       cpuProfile,
	"trace":     traceProfile,
	"heap":      lookup,
	"goroutine": lookup,
	"block":     lookup,
}

var pmChan = map[string]chan struct{}{}

func init() {
	for k := range pm {
		pmChan[k] = make(chan struct{}, 1)
	}
}

// StartProfile start profiling for the current process.
// if success returns an profiling file.
func StartProfile(option *structs.ProfileData) (string, error) {
	f, ok := pm[option.Type]
	if !ok {
		return "", fmt.Errorf("Unknown profile: %s", option.Type)
	}
	select {
	case pmChan[option.Type] <- struct{}{}:
		defer func() { <-pmChan[option.Type] }()
	default:
		return "", fmt.Errorf("Profiling is already running: %s", option.Type)
	}

	fname := fmt.Sprintf("/tmp/%s_%v.pprof", option.Type, time.Now().UnixNano())
	fn, err := os.Create(fname)
	if err != nil {
		return "", err
	}
	defer fn.Close()
	if err = f(fn, option); err != nil {
		return "", err
	}
	return fname, nil
}

func cpuProfile(fn *os.File, option *structs.ProfileData) error {
	if err := pprof.StartCPUProfile(fn); err != nil {
		return err
	}
	time.Sleep(time.Duration(option.Sleep) * time.Second)
	pprof.StopCPUProfile()
	return nil
}

func traceProfile(fn *os.File, option *structs.ProfileData) error {
	if err := trace.Start(fn); err != nil {
		return err
	}
	time.Sleep(time.Duration(option.Sleep) * time.Second)
	trace.Stop()
	return nil
}

func lookup(fn *os.File, option *structs.ProfileData) error {
	p := pprof.Lookup(option.Type)
	if p == nil {
		return fmt.Errorf("Unknown profile: %s", option.Type)
	}
	if err := p.WriteTo(fn, option.Debug); err != nil {
		return err
	}
	return nil
}
