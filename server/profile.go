package server

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/l2x/gopprof/common/structs"
)

// generate profiling graph to pdf file.
// go tool pprof -pdf /path/to/bin /path/to/pprof/file > /path/to/pdf/save
func pprofToPDF(data *structs.ProfileData) ([]byte, error) {
	var (
		err                                                 error
		tmpDir                                              = fmt.Sprintf("tmp/%d", time.Now().UnixNano())
		goRoot, goBin, tmpBinFile, tmpPprofFile, tmpPDFFile string
	)
	// TODO
	// auto check
	// get latest verison if not found
	goRoot, err = db.TableConfig(data.NodeID).GetGoroot(strings.TrimLeft(data.GoVersion, "go"))
	if err != nil {
		err = errors.New("go root not set")
		logger.Error(err)
		return nil, err
	}
	goBin = filepath.Join(goRoot, "bin/go")

	os.MkdirAll(tmpDir, 0755)
	defer os.RemoveAll(tmpDir)

	// get binary file
	// if the failure continues
	fname, err := db.TableBin(data.NodeID).Get(data.BinMD5)
	if err == nil {
		tmpBinFile = filepath.Join(tmpDir, filepath.Base(fname))
		if err = store.CopyToLocal(tmpBinFile, fname); err != nil {
			logger.Error(err)
			tmpBinFile = ""
		}
	}

	tmpPprofFile = filepath.Join(tmpDir, filepath.Base(data.File))
	if err = store.CopyToLocal(tmpPprofFile, data.File); err != nil {
		logger.Error(err)
		return nil, err
	}

	// set go root
	currentGoRoot := os.Getenv("GOROOT")
	os.Setenv("GOROOT", goRoot)
	defer func() {
		os.Setenv("GOROOT", currentGoRoot)
	}()

	tmpPDFFile = tmpPprofFile + ".pdf"
	cmd := fmt.Sprintf("%s tool pprof -pdf %s %s > %s", goBin, tmpBinFile, tmpPprofFile, tmpPDFFile)
	b, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		err = fmt.Errorf("%s,%s,%s", cmd, err.Error(), string(b))
		logger.Error(err)
		return nil, err
	}

	if b, err = ioutil.ReadFile(tmpPDFFile); err != nil {
		logger.Error(err)
		return nil, err
	}

	pdfFile := data.File + ".pdf"
	if err = store.Save(pdfFile, b); err != nil {
		logger.Error(err)
		return nil, err
	}
	return b, nil
}
