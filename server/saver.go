package server

import (
	"github.com/l2x/gopprof/engine/store"
	_ "github.com/l2x/gopprof/engine/store/bolt"

	"github.com/l2x/gopprof/engine/files"
	_ "github.com/l2x/gopprof/engine/files/localfile"
)

var (
	storeSaver store.Store
	filesSaver files.Files
)

func initStoreSaver(driver, source string) error {
	s, err := store.Open(driver, source)
	if err != nil {
		logger.Error(err)
		return err
	}
	storeSaver = s
	return nil
}

func initFilesSaver(driver, source string) error {
	s, err := files.Open(driver, source)
	if err != nil {
		logger.Error(err)
		return err
	}
	filesSaver = s
	return nil
}
