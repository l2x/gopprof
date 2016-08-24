package server

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/l2x/gopprof/common/structs"
)

var (
	errGetPDFLink = "https://github.com/l2x/gopprof#q-why-can-not-download-pdf"
)

// generate profiling graph to pdf file.
// go tool pprof -pdf /path/to/bin /path/to/pprof/file > /path/to/pdf/save
func pprofToPDF(data *structs.ProfileData) ([]byte, error) {
	var (
		err                                                 error
		tmpDir                                              = fmt.Sprintf("tmp/%d", time.Now().UnixNano())
		goRoot, goBin, tmpBinFile, tmpPprofFile, tmpPDFFile string
	)
	goVersion := strings.TrimLeft(data.GoVersion, "go")
	goroot, err := db.TableConfig("").GetGoroot(goVersion)
	if err != nil {
		if goroot, err = tryDownloadGo(goVersion); err != nil {
			err = fmt.Errorf("%s\nmore information: %s", err.Error(), errGetPDFLink)
			logger.Error(err)
			return nil, err
		}
	}
	goRoot = goroot.Path
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

// try to download go file and save GOROOT
func tryDownloadGo(goVersion string) (*structs.Goroot, error) {
	tmpFile := fmt.Sprintf("tmp/%v", time.Now().UnixNano())
	out, err := os.Create(tmpFile)
	if err != nil {
		return nil, err
	}
	defer out.Close()
	defer os.Remove(tmpFile)

	uri := fmt.Sprintf("https://storage.googleapis.com/golang/go%s.linux-amd64.tar.gz", goVersion)
	logger.Debug("download go file:", uri)
	resp, err := http.Get(uri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return nil, err
	}
	target := fmt.Sprintf("tmp/%v", time.Now().UnixNano())
	os.MkdirAll(target, 0755)
	if err = uncompress(tmpFile, target); err != nil {
		return nil, err
	}
	defer os.RemoveAll(target)

	newpath := fmt.Sprintf("%s/go%s", conf.GoFilePath, goVersion)
	os.RemoveAll(newpath)
	if err = os.Rename(fmt.Sprintf("%s/go", target), newpath); err != nil {
		return nil, err
	}
	goroot := &structs.Goroot{Version: goVersion, Path: newpath}
	db.TableConfig("").SaveGoroot(goroot)
	return db.TableConfig("").GetGoroot(goVersion)
}

func uncompress(source, target string) error {
	cmd := fmt.Sprintf("tar -C %s -xzf %s", target, source)
	b, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		err = fmt.Errorf("%s,%s,%s", cmd, err.Error(), string(b))
		return err
	}
	return nil
}
