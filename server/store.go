package server

import (
	"github.com/l2x/gopprof/engine/database"
	_ "github.com/l2x/gopprof/engine/database/boltdb"
	st "github.com/l2x/gopprof/engine/store"
	_ "github.com/l2x/gopprof/engine/store/localfile"
)

var (
	db    database.Database
	store st.Store
)

func initDB(driver, source string) error {
	s, err := database.Open(driver, source)
	if err != nil {
		logger.Error(err)
		return err
	}
	db = s
	return nil
}

func initStore(driver, source string) error {
	s, err := st.Open(driver, source)
	if err != nil {
		logger.Error(err)
		return err
	}
	store = s
	return nil
}
