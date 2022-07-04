package http

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"net/http"
)

var tracer = otel.Tracer("github.com/jacob-delgado/playground")

func hello(w http.ResponseWriter, req *http.Request) {
	req.Context()
	fmt.Fprintf(w, "hello\n")
}

func headers(w http.ResponseWriter, req *http.Request) {
	_, span := tracer.Start(req.Context(), "sleep")
	defer span.End()

	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}

func Serve(errCh chan error) {
	helloHandler := http.HandlerFunc(hello)
	wrappedHelloHandler := otelhttp.NewHandler(helloHandler, "hello")
	http.Handle("/hello", wrappedHelloHandler)

	headerHandler := http.HandlerFunc(headers)
	wrappedHeaderHandler := otelhttp.NewHandler(headerHandler, "headers")
	http.Handle("/headers", wrappedHeaderHandler)

	http.Handle("/metrics", promhttp.Handler())

	errCh <- http.ListenAndServe(":8090", nil)
}
