package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/l2x/gopprof/common/structs"
	"github.com/l2x/gopprof/common/utils"
)

// ListenHTTP start http server
func ListenHTTP(port string) {
	logger.Infof("listen http %s", port)

	//gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type,Token")
		c.Next()
	})

	r.OPTIONS("/*cors", func(c *gin.Context) {})
	r.GET("/nodes", nodesHandler)
	r.POST("/stats", statsHandler)
	r.POST("/pprof", pprofHandler)
	r.POST("/upload", uploadHandler)
	r.GET("/download", downloadHandler)
	r.Run(port)
}

func nodesHandler(c *gin.Context) {
	nodes, err := storeSaver.GetNodes()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, nodes)
}

type statsReq struct {
	Nodes   []string `json:"nodes"`
	Options []string `json:"options"`
	Date    struct {
		Start int64 `json:"start"`
		End   int64 `json:"end"`
	} `json:"date"`
}

type statsResp struct {
	Type string           `json:"type"`
	Data []*statsRespData `json:"data"`
}

type statsRespData struct {
	NodeID string     `json:"name"`
	Data   [][2]int64 `json:"data"`
}

type statsData struct {
	nodeID string
	data   []*structs.StatsData
}

func statsHandler(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		logger.Error(err)
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	var req *statsReq
	if err = json.Unmarshal(body, &req); err != nil {
		logger.Error(err)
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	if len(req.Nodes) == 0 || len(req.Options) == 0 {
		c.String(http.StatusBadRequest, "nodes or options invalid")
		return
	}

	logger.Debug(string(body))

	statsDatas := []*statsData{}
	for _, nodeID := range req.Nodes {
		data, err := storeSaver.GetStatsByTime(nodeID, req.Date.Start, req.Date.End)
		if err != nil {
			logger.Error(err)
			continue
		}
		statsDatas = append(statsDatas, &statsData{nodeID: nodeID, data: data})
	}

	var response []*statsResp
	for _, opt := range req.Options {
		resp := &statsResp{Type: statsTitle(opt)}
		for _, datas := range statsDatas {
			srd := &statsRespData{NodeID: datas.nodeID, Data: [][2]int64{}}
			for _, data := range datas.data {
				d := [2]int64{data.Created * 1000, statsParser(opt, data)}
				srd.Data = append(srd.Data, d)
			}
			resp.Data = append(resp.Data, srd)
		}
		response = append(response, resp)
	}

	c.JSON(http.StatusOK, response)
}

func statsTitle(typ string) string {
	switch typ {
	case "heap":
		return "heap (byte)"
	case "gc":
		return "gc pause (ns)"
	}
	return typ
}

func statsParser(typ string, data *structs.StatsData) int64 {
	switch typ {
	case "goroutine":
		return int64(data.NumGoroutine)
	case "heap":
		return int64(data.HeapAlloc)
	case "gc":
		return int64(data.PauseNs[(data.NumGC+255)%256])
	}
	return 0
}

func pprofHandler(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		logger.Error(err)
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	var req *statsReq
	if err = json.Unmarshal(body, &req); err != nil {
		logger.Error(err)
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	if len(req.Nodes) == 0 || len(req.Options) == 0 {
		c.String(http.StatusBadRequest, "nodes or options invalid")
		return
	}

	logger.Debug(string(body))

	response := []*structs.ProfileData{}
	for _, nodeID := range req.Nodes {
		datas, err := storeSaver.GetProfilesByTime(nodeID, req.Date.Start, req.Date.End)
		if err != nil {
			logger.Error(err)
			continue
		}

		for _, data := range datas {
			if utils.InStringSlice(data.Type, req.Options) {
				response = append(response, data)
			}
		}
	}
	c.JSON(http.StatusOK, response)
}

func uploadHandler(c *gin.Context) {
	c.Request.ParseMultipartForm(10 * 1024 * 1024)

	v := c.Request.FormValue("type")
	typ, err := strconv.Atoi(v)
	if err != nil {
		logger.Error(err)
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	eventType := structs.EventType(typ)
	file, handler, err := c.Request.FormFile("file")
	if err != nil {
		logger.Error(err)
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	defer file.Close()

	switch eventType {
	case structs.EventTypeUploadProfile:
		err = uploadProfile(c, file, handler.Filename)
	case structs.EventTypeUploadBin:
		err = uploadBin(c, file)
	default:
		err = fmt.Errorf("Unknown event: %v", eventType)
	}
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

func uploadProfile(c *gin.Context, file multipart.File, filename string) error {
	var data structs.ProfileData
	v := c.Request.FormValue("data")
	if err := json.Unmarshal([]byte(v), &data); err != nil {
		logger.Error(err)
		return err
	}

	fname := filepath.Join(data.NodeID, data.Type, time.Now().Format("2006/01/02"), filename)
	if err := filesSaver.CopyTo(fname, file); err != nil {
		logger.Error(err)
		return err
	}
	data.File = fname
	if err := eventUploadProfile(data); err != nil {
		return err
	}
	return nil
}

func uploadBin(c *gin.Context, file multipart.File) error {
	var data structs.ExInfo
	v := c.Request.FormValue("data")
	if err := json.Unmarshal([]byte(v), &data); err != nil {
		logger.Error(err)
		return err
	}
	fname := filepath.Join(data.NodeID, "bin", time.Now().Format("2006/01/02"), data.MD5)
	if err := filesSaver.CopyTo(fname, file); err != nil {
		logger.Error(err)
		return err
	}
	if err := eventUploadBin(data, fname); err != nil {
		return err
	}
	return nil
}

func downloadHandler(c *gin.Context) {
	typ := c.Query("type")
	var data structs.ProfileData
	if err := json.Unmarshal([]byte(c.Query("data")), &data); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	var (
		b     []byte
		fname string
		err   error
	)
	switch typ {
	case "bin":
	case "pdf":
		pdfFile := data.File + ".pdf"
		fname = filepath.Base(pdfFile)
		if b, err = filesSaver.Get(pdfFile); err != nil {
			if b, err = pprofToPDF(&data); err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				return
			}
		}
	case "pprof":
		fname = filepath.Base(data.File)
		if b, err = filesSaver.Get(data.File); err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
	default:
		c.String(http.StatusBadRequest, "download type unsupport")
		return
	}

	c.Writer.Header().Set("Content-Disposition", "attachment; filename="+fname)
	c.Writer.Header().Set("Content-Type", c.Request.Header.Get("Content-Type"))
	io.Copy(c.Writer, bytes.NewBuffer(b))
}

// TODO
func pprofToPDF(data *structs.ProfileData) ([]byte, error) {
	var (
		tmpDir       = fmt.Sprintf("tmp/%d", time.Now().UnixNano())
		goBin        = "go"
		tmpBinFile   = ""
		tmpPprofFile = ""
		tmpPdfFile   = ""
	)
	os.MkdirAll(tmpDir, 0755)
	defer os.RemoveAll(tmpDir)

	fname, err := storeSaver.GetBin(data.NodeID, data.BinMd5)
	if err != nil {
		logger.Error(err)
	}
	tmpBinFile = filepath.Join(tmpDir, filepath.Base(fname))
	if err = filesSaver.CopyFile(fname, tmpBinFile); err != nil {
		logger.Error(err)
		tmpBinFile = ""
	}

	tmpPprofFile = filepath.Join(tmpDir, filepath.Base(data.File))
	if err = filesSaver.CopyFile(data.File, tmpPprofFile); err != nil {
		logger.Error(err)
		return nil, err
	}
	tmpPdfFile = tmpPdfFile + ".pdf"

	cmd := fmt.Sprintf("%s tool pprof -pdf %s %s > %s", goBin, tmpBinFile, tmpPprofFile, tmpPdfFile)
	b, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		err = fmt.Errorf("%s,%s,%s", cmd, err.Error(), string(b))
		logger.Error(err)
		return nil, err
	}
	b, err = ioutil.ReadFile(tmpPdfFile)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	pdfFile := data.File + ".pdf"
	if err = filesSaver.Save(pdfFile, b); err != nil {
		logger.Error(err)
		return nil, err
	}
	return b, nil
}
