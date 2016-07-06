package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/l2x/gopprof/common/event"
	"github.com/l2x/gopprof/common/structs"
	"github.com/l2x/gopprof/common/util"
)

// ListenHTTP start http server
func ListenHTTP(port string) {
	logger.Infof("listen http %s", port)
	if conf.Debug == false {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type,Token")
		c.Next()
	})

	// router
	r.OPTIONS("/*cors", func(c *gin.Context) {})
	r.GET("/nodes", nodesHandler)
	r.POST("/stats", statsHandler)
	r.POST("/pprof", pprofHandler)
	r.POST("/upload", uploadHandler)
	r.GET("/download", downloadHandler)
	r.GET("/setting", settingHandler)
	r.POST("/setting/save", settingSaveHandler)

	if err := r.Run(port); err != nil {
		logger.Criticalf("Cannot start http server: %s", err)
	}
	Exit()
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

func nodesHandler(c *gin.Context) {
	nodes, err := db.TableNode("").GetAll()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, nodes)
}

func settingHandler(c *gin.Context) {
	var (
		nodeConf *structs.NodeConf
	)
	nodeID := c.Query("nodeid")
	if nodeID == "_default" {
		nodeConf, _ = db.TableConfig(nodeID).GetDefault()
	} else {
		nodeConf, _ = db.TableConfig(nodeID).Get()
	}
	c.JSON(http.StatusOK, nodeConf)
}

type settingReq struct {
	NodeConf structs.NodeConf `json:"conf"`
	Nodes    []string         `json:"nodes"`
}

func settingSaveHandler(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		logger.Error(err)
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	var req *settingReq
	if err = json.Unmarshal(body, &req); err != nil {
		logger.Error(err)
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	if req.Nodes == nil || len(req.Nodes) == 0 {
		err := errors.New("node empty")
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	if !req.NodeConf.EnableProfile {
		req.NodeConf.Profile = []string{}
		req.NodeConf.ProfileCron = ""
	}
	if !req.NodeConf.EnableStats {
		req.NodeConf.StatsCron = ""
	}
	var dfa bool
	if len(req.Nodes) == 1 && req.Nodes[0] == "_default" {
		dfa = true
	}

	if dfa {
		if err := db.TableConfig(req.Nodes[0]).SaveDefault(&req.NodeConf); err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		c.Status(http.StatusOK)
		return
	}

	for _, nodeID := range req.Nodes {
		if err = db.TableConfig(nodeID).Save(&req.NodeConf); err != nil {
			logger.Error(err)
			continue
		}
	}
	c.Status(http.StatusOK)
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

	statsDatas := []*statsData{}
	for _, nodeID := range req.Nodes {
		data, err := db.TableStats(nodeID).GetRangeTime(req.Date.Start, req.Date.End)
		if err != nil {
			logger.Error(err)
			continue
		}
		statsDatas = append(statsDatas, &statsData{nodeID: nodeID, data: data})
	}

	var response []*statsResp
	for _, opt := range req.Options {
		resp := &statsResp{Type: opt}
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

func statsParser(typ string, data *structs.StatsData) int64 {
	switch typ {
	case "goroutine":
		return int64(data.NumGoroutine)
	case "heap":
		return int64(data.HeapAlloc)
	case "gc":
		return int64(data.GCPauseNs)
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

	response := []*structs.ProfileData{}
	for _, nodeID := range req.Nodes {
		datas, err := db.TableProfile(nodeID).GetRangeTime(req.Date.Start, req.Date.End)
		if err != nil {
			logger.Error(err)
			continue
		}

		for _, data := range datas {
			if util.InStringSlice(data.Type, req.Options) {
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

	eventType := event.EventType(typ)
	file, handler, err := c.Request.FormFile("file")
	if err != nil {
		logger.Error(err)
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	defer file.Close()

	switch eventType {
	case event.EventTypeUploadProfile:
		err = uploadProfile(c, file, handler.Filename)
	case event.EventTypeUploadBin:
		err = uploadBin(c, file, handler.Filename)
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
	if err := store.Copy(fname, file); err != nil {
		logger.Error(err)
		return err
	}
	data.File = fname
	if err := eventUploadProfile(&data); err != nil {
		return err
	}
	return nil
}

func uploadBin(c *gin.Context, file multipart.File, filename string) error {
	nodeID := c.Request.FormValue("nodeid")
	binMD5 := c.Request.FormValue("bin_md5")
	fname := filepath.Join("bin", time.Now().Format("2006/01/02"), binMD5, filename)
	if err := store.Copy(fname, file); err != nil {
		logger.Error(err)
		return err
	}
	if err := eventUploadBin(nodeID, binMD5, fname); err != nil {
		return err
	}
	return nil
}

func downloadHandler(c *gin.Context) {
	var (
		typ        = c.Query("type")
		nodeID     = c.Query("nodeid")
		created, _ = strconv.ParseInt(c.Query("created"), 10, 64)
		b          []byte
		fname      string
		err        error
	)

	data, err := db.TableProfile(nodeID).GetCreated(created)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	switch typ {
	case "bin":
		b, fname, err = downloadBin(data)
	case "pdf":
		b, fname, err = downloadPDF(data)
	case "pprof":
		b, fname, err = downloadPprof(data)
	default:
		c.String(http.StatusBadRequest, "download type unsupport")
		return
	}

	c.Writer.Header().Set("Content-Disposition", "attachment; filename="+fname)
	c.Writer.Header().Set("Content-Type", c.Request.Header.Get("Content-Type"))
	io.Copy(c.Writer, bytes.NewBuffer(b))
}

func downloadBin(data *structs.ProfileData) ([]byte, string, error) {
	file, err := db.TableBin(data.NodeID).Get(data.BinMD5)
	if err != nil {
		return nil, "", err
	}
	b, err := store.Get(file)
	if err != nil {
		return nil, "", err
	}
	return b, filepath.Base(file), nil
}

func downloadPprof(data *structs.ProfileData) ([]byte, string, error) {
	b, err := store.Get(data.File)
	if err != nil {
		return nil, "", err
	}
	return b, filepath.Base(data.File), nil
}

func downloadPDF(data *structs.ProfileData) ([]byte, string, error) {
	pdfFile := data.File + ".pdf"
	if b, err := store.Get(pdfFile); err == nil {
		return b, filepath.Base(pdfFile), nil
	}
	b, err := pprofToPDF(data)
	if err != nil {
		return nil, "", err
	}
	return b, filepath.Base(pdfFile), nil
}
