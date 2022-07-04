package http

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

func hello(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "hello\n")
}

func headers(w http.ResponseWriter, req *http.Request) {
	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}

func Serve() {
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/headers", headers)
	http.Handle("/metrics", promhttp.Handler())

	http.ListenAndServe(":8090", nil)
}
