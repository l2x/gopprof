package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/l2x/gopprof/common/structs"
)

// ListenHTTP start http server
func ListenHTTP(port string) {
	logger.Infof("listen http %s", port)
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/stats", statsHandler)
	http.HandleFunc("/upload", uploadHandler)
	if err := http.ListenAndServe(port, nil); err != nil {
		panic(err)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, time.Now().Unix())
}

func statsHandler(w http.ResponseWriter, r *http.Request) {
	nodes := r.URL.Query().Get("nodes")
	start := r.URL.Query().Get("start")
	end := r.URL.Query().Get("end")

	startTime, err := strconv.ParseInt(start, 10, 64)
	if err != nil {
		logger.Error(err)
		fmt.Fprint(w, err)
		return
	}
	endTime, err := strconv.ParseInt(end, 10, 64)
	if err != nil {
		logger.Error(err)
		fmt.Fprint(w, err)
		return
	}

	resp := map[string][]*structs.StatsData{}
	for _, node := range strings.Split(nodes, ",") {
		data, err := storeSaver.GetStatsByTime(node, startTime, endTime)
		if err != nil {
			logger.Error(err)
			continue
		}
		resp[node] = data
	}

	b, err := json.Marshal(resp)
	if err != nil {
		logger.Error(err)
	}
	fmt.Fprint(w, string(b))
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 * 1024 * 1024)
	file, handler, err := r.FormFile("file")
	if err != nil {
		logger.Error(err)
		fmt.Fprint(w, err.Error())
		return
	}
	defer file.Close()

	var data structs.ProfileData
	v := r.FormValue("data")
	if err := json.Unmarshal([]byte(v), &data); err != nil {
		logger.Error(err)
		fmt.Fprint(w, err.Error())
		return
	}
	logger.Debug(data)

	fname := filepath.Join(data.NodeID, data.Type, time.Now().Format("2006/01/02"), handler.Filename)
	if err = filesSaver.CopyTo(fname, file); err != nil {
		logger.Error(err)
		fmt.Fprint(w, err.Error())
		return
	}
	data.File = fname
	evtReq := structs.NewEvent(structs.EventTypeProfile, data)
	if _, err = eventProfile(evtReq); err != nil {
		fmt.Fprint(w, err.Error())
		return
	}
	fmt.Fprint(w, "")
}
