package server

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/l2x/gopprof/common/structs"
)

// generate profiling pdf file.
// go tool pprof -pdf /path/to/bin /path/to/pprof/file > /path/to/pdf/save
// TODO
func pprofToPDF(data *structs.ProfileData) ([]byte, error) {
	var (
		tmpDir       = fmt.Sprintf("tmp/%d", time.Now().UnixNano())
		goBin        = "go"
		tmpBinFile   = ""
		tmpPprofFile = ""
		tmpPDFFile   = ""
	)
	os.MkdirAll(tmpDir, 0755)
	defer os.RemoveAll(tmpDir)

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

	tmpPDFFile = tmpPprofFile + ".pdf"
	cmd := fmt.Sprintf("%s tool pprof -pdf %s %s > %s", goBin, tmpBinFile, tmpPprofFile, tmpPDFFile)
	b, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		err = fmt.Errorf("%s,%s,%s", cmd, err.Error(), string(b))
		logger.Error(err)
		return nil, err
	}

	b, err = ioutil.ReadFile(tmpPDFFile)
	if err != nil {
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
