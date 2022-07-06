package http

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"github.com/jacob-delgado/playground/pkg/metrics"
)

type Server struct {
	logger *otelzap.Logger
	tracer trace.Tracer
}

func NewServer(logger *otelzap.Logger) *Server {
	return &Server{
		logger: logger,
		tracer: otel.Tracer("github.com/jacob-delgado/inventory/pkg/http"),
	}
}

func (s *Server) hello(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	fmt.Fprintf(w, "hello\n")

	s.logger.Ctx(ctx).Info("hello")

	metrics.HelloProcessed.Inc()
}

func (s *Server) headers(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	s.logger.Ctx(ctx).Info("headers")

	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)

			// cardinality issues will be present, but this isn't production anyway
			if strings.ToUpper(name) == "FOO" {
				if strings.ToUpper(h) == "BAR" {
					metrics.HeadersProcessed.With(prometheus.Labels{"foobar": h}).Inc()
				} else {
					metrics.HeadersProcessed.With(prometheus.Labels{"foo": h}).Inc()
				}
			}
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

	s.logger.Info("starting http server on localhost:8090")
	errCh <- http.ListenAndServe(":8090", nil)
}
