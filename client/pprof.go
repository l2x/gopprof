package client

import (
	"fmt"
	"os"
	"runtime/pprof"
	"runtime/trace"
	"time"
)

var pm = map[string]func(fn *os.File, typ string) error{
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
func StartProfile(typ string) (string, error) {
	f, ok := pm[typ]
	if !ok {
		return "", fmt.Errorf("Unknown profile: %s", typ)
	}
	select {
	case pmChan[typ] <- struct{}{}:
		defer func() { <-pmChan[typ] }()
	default:
		return "", fmt.Errorf("Profiling is already running: %s", typ)
	}

	fname := fmt.Sprintf("%s_%v.pprof", typ, time.Now().UnixNano())
	fn, err := os.Create(fname)
	if err != nil {
		return "", err
	}
	defer fn.Close()
	if err = f(fn, typ); err != nil {
		return "", err
	}
	return fname, nil
}

func cpuProfile(fn *os.File, typ string) error {
	if err := pprof.StartCPUProfile(fn); err != nil {
		return err
	}
	time.Sleep(time.Duration(30) * time.Second)
	pprof.StopCPUProfile()
	return nil
}

func traceProfile(fn *os.File, typ string) error {
	if err := trace.Start(fn); err != nil {
		return err
	}
	time.Sleep(time.Duration(30) * time.Second)
	trace.Stop()
	return nil
}

func lookup(fn *os.File, typ string) error {
	p := pprof.Lookup(typ)
	if p == nil {
		return fmt.Errorf("Unknown profile: %s", typ)
	}
	// TODO debug value
	if err := p.WriteTo(fn, 0); err != nil {
		return err
	}
	return nil
}
