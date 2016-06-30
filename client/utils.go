package client

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

// GetBinFile return current running process file
func GetBinFile() ([]byte, error) {
	bf, err := filepath.Abs(os.Args[0])
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadFile(bf)
	if err != nil {
		return nil, err
	}
	return b, nil
}
