package server

import (
	"github.com/l2x/gopprof/engine/store"
	_ "github.com/l2x/gopprof/engine/store/bolt"
	"github.com/l2x/profile"
)

var (
	storeSaver   store.Store
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
