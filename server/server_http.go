package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/l2x/gopprof/common/structs"
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

	var response [][]*structs.ProfileData
	for _, nodeID := range req.Nodes {
		data, err := storeSaver.GetProfilesByTime(nodeID, req.Date.Start, req.Date.End)
		if err != nil {
			logger.Error(err)
			continue
		}
		response = append(response, data)
	}
	c.JSON(http.StatusOK, response)
}

func uploadHandler(c *gin.Context) {
	c.Request.ParseMultipartForm(10 * 1024 * 1024)
	file, handler, err := c.Request.FormFile("file")
	if err != nil {
		logger.Error(err)
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	defer file.Close()

	var data structs.ProfileData
	v := c.Request.FormValue("data")
	if err := json.Unmarshal([]byte(v), &data); err != nil {
		logger.Error(err)
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	logger.Debug(data)

	fname := filepath.Join(data.NodeID, data.Type, time.Now().Format("2006/01/02"), handler.Filename)
	if err = filesSaver.CopyTo(fname, file); err != nil {
		logger.Error(err)
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	data.File = fname
	evtReq := structs.NewEvent(structs.EventTypeProfile, data)
	if _, err = eventProfile(evtReq); err != nil {
		logger.Error(err)
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.Status(http.StatusOK)
}
