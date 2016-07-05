package server

import (
	"bytes"
	"encoding/json"
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

func uploadBin(c *gin.Context, file multipart.File) error {
	nodeID := c.Request.FormValue("nodeid")
	binMD5 := c.Request.FormValue("bin_md5")
	fname := filepath.Join(nodeID, "bin", time.Now().Format("2006/01/02"), binMD5)
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

	switch data.Type {
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
	fname, err := db.TableBin(data.NodeID).Get(data.BinMD5)
	if err != nil {
		return nil, "", err
	}
	b, err := store.Get(fname)
	if err != nil {
		return nil, "", err
	}
	return b, filepath.Base(fname), nil
}

func downloadPprof(data *structs.ProfileData) ([]byte, string, error) {
	b, err := store.Get(data.File)
	if err != nil {
		return nil, "", err
	}
	return b, filepath.Base(data.File), nil
}

func downloadPDF(data *structs.ProfileData) ([]byte, string, error) {
	fname := data.File + ".pdf"
	if b, err := store.Get(fname); err == nil {
		return b, "", nil
	}
	b, err := pprofToPDF(data)
	if err != nil {
		return nil, "", err
	}
	return b, filepath.Base(fname), nil
}
