package server

import (
	"github.com/l2x/gopprof/server/stats"
	_ "github.com/l2x/gopprof/server/stats/bolt"
	"github.com/l2x/gopprof/server/store"
	_ "github.com/l2x/gopprof/server/store/bolt"
	"github.com/l2x/profile"
)

var (
	storeSaver   store.Store
	statsSaver   stats.Stats
	profileSaver profile.Profile
)

func initStoreSaver(driver, source string) error {
	s, err := store.Open(driver, source)
	if err != nil {
		return err
	}
	storeSaver = s
	return nil
}

func initStatsSaver(driver, source string) error {
	s, err := stats.Open(driver, source)
	if err != nil {
		return err
	}
	statsSaver = s
	return nil
}
