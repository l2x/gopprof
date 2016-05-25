package server

import (
	"fmt"
	"html"
	"log"
	"net/http"
)

// ListenHTTP start http server
func ListenHTTP(port string) {
	log.Println("listen http:", port)
	http.HandleFunc("/", indexHandler)
	if err := http.ListenAndServe(port, nil); err != nil {
		panic(err)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}
