package http

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	playgroundv1 "github.com/jacob-delgado/playground/gen/proto/go/playground/v1"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

var (
	headersProcessed = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "headers_called",
			Help: "The total number of times headers routed was called",
		},
		[]string{"foo", "foobar"},
	)

	helloProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "hello_called",
		Help: "The total number of times hello route was called",
	})
)

type Server struct {
	logger *otelzap.Logger
	tracer trace.Tracer
}

func NewServer(logger *otelzap.Logger) *Server {
	return &Server{
		logger: logger,
		tracer: otel.Tracer("github.com/jacob-delgado/playground"),
	}
}

func (s *Server) GetFeature(ctx context.Context, in *playgroundv1.GetFeatureRequest, opts ...grpc.CallOption) (*playgroundv1.GetFeatureResponse, error) {
	return &playgroundv1.GetFeatureResponse{}, nil
}

func (s *Server) hello(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	_, span := s.tracer.Start(ctx, "headers")
	defer span.End()

	fmt.Fprintf(w, "hello\n")

	s.logger.Ctx(ctx).Info("hello")

	helloProcessed.Inc()
}

func (s *Server) headers(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	_, span := s.tracer.Start(ctx, "headers")
	defer span.End()

	s.logger.Ctx(ctx).Info("headers")

	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)

			// cardinality issues will be present, but this isn't production anyway
			if strings.ToUpper(name) == "FOO" {
				if strings.ToUpper(h) == "BAR" {
					headersProcessed.With(prometheus.Labels{"foobar": h}).Inc()
				} else {
					headersProcessed.With(prometheus.Labels{"foo": h}).Inc()
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

	errCh <- http.ListenAndServe(":8090", nil)
}
