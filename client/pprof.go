package client

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime/pprof"
	"runtime/trace"
	"time"
)

var pm = map[string]func(fn *os.File, option *ProfileOption) error{
	"cpu":   cpuProfile,
	"trace": traceProfile,
	"heap":  lookup,
	"block": lookup,
}

var pmChan = map[string]chan struct{}{}

func init() {
	for k := range pm {
		pmChan[k] = make(chan struct{}, 1)
	}
}

// StartProfile enables profiling for the current process.
// if success returns an profiling result file.
func StartProfile(option *ProfileOption) (string, error) {
	f, ok := pm[option.Name]
	if !ok {
		return "", fmt.Errorf("Unknown profile: %s", option.Name)
	}
	select {
	case pmChan[option.Name] <- struct{}{}:
		defer func() { <-pmChan[option.Name] }()
	default:
		return "", fmt.Errorf("profiling is already running: %s", option.Name)
	}

	fname := filepath.Join(option.tmp, fmt.Sprintf("%s_%v.pprof", option.Name, time.Now().Unix()))
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

func cpuProfile(fn *os.File, option *ProfileOption) error {
	if err := pprof.StartCPUProfile(fn); err != nil {
		return err
	}
	time.Sleep(time.Duration(option.Sleep) * time.Second)
	pprof.StopCPUProfile()
	return nil
}

func traceProfile(fn *os.File, option *ProfileOption) error {
	if err := trace.Start(fn); err != nil {
		return err
	}
	time.Sleep(time.Duration(option.Sleep) * time.Second)
	trace.Stop()
	return nil
}

func lookup(fn *os.File, option *ProfileOption) error {
	p := pprof.Lookup(option.Name)
	if p == nil {
		return fmt.Errorf("Unknown profile: %s", option.Name)
	}
	if err := p.WriteTo(fn, option.Debug); err != nil {
		return err
	}
	return nil
}
