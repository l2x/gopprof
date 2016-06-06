package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/l2x/gopprof/common/structs"
)

// ListenHTTP start http server
func ListenHTTP(port string) {
	log.Println("listen http:", port)
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
		log.Println(err)
	}
	endTime, err := strconv.ParseInt(end, 10, 64)
	if err != nil {
		log.Println(err)
	}

	resp := map[string][]*structs.StatsData{}
	for _, node := range strings.Split(nodes, ",") {
		data, err := storeSaver.GetStatsByTime(node, startTime, endTime)
		if err != nil {
			log.Println(err)
			continue
		}
		resp[node] = data
	}

	b, err := json.Marshal(resp)
	if err != nil {
		log.Println(err)
	}
	fmt.Fprint(w, string(b))
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("upload")

	r.ParseMultipartForm(10 * 1024 * 1024)
	file, handler, err := r.FormFile("file")
	if err != nil {
		log.Println("[upload]", err)
		fmt.Fprint(w, err.Error())
		return
	}
	defer file.Close()

	fmt.Println(handler.Filename)
}
