package http

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type Server struct {
	tracer trace.Tracer
	logger *otelzap.Logger
}

func NewServer(logger *otelzap.Logger) *Server {
	return &Server{
		logger: logger,
		tracer: otel.Tracer("github.com/jacob-delgado/playground"),
	}
}

func (s *Server) hello(w http.ResponseWriter, req *http.Request) {
	req.Context()
	fmt.Fprintf(w, "hello\n")
}

func (s *Server) headers(w http.ResponseWriter, req *http.Request) {
	_, span := s.tracer.Start(req.Context(), "sleep")
	defer span.End()

	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}

func (s *Server) Serve(errCh chan error) {
	helloHandler := http.HandlerFunc(s.hello)
	wrappedHelloHandler := otelhttp.NewHandler(helloHandler, "hello")
	http.Handle("/hello", wrappedHelloHandler)

	headerHandler := http.HandlerFunc(s.headers)
	wrappedHeaderHandler := otelhttp.NewHandler(headerHandler, "headers")
	http.Handle("/headers", wrappedHeaderHandler)

	http.Handle("/metrics", promhttp.Handler())

	errCh <- http.ListenAndServe(":8090", nil)
}
