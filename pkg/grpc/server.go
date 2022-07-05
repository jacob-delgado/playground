package grpc

import (
	"context"
	"log"
	"net"

	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	playgroundv1 "github.com/jacob-delgado/playground/gen/proto/go/playground/v1"
	"github.com/jacob-delgado/playground/pkg/metrics"
)

type Server struct {
	logger *otelzap.Logger
	tracer trace.Tracer
}

func NewServer(logger *otelzap.Logger) *Server {
	return &Server{
		logger: logger,
		tracer: otel.Tracer("github.com/jacob-delgado/playground/pkg/grpc"),
	}
}

func (s *Server) GetFeature(ctx context.Context, req *playgroundv1.GetFeatureRequest) (*playgroundv1.GetFeatureResponse, error) {
	s.logger.Ctx(ctx).Info("GetFeature rpc call")
	metrics.GetFeatureProcessed.Inc()

	return &playgroundv1.GetFeatureResponse{}, nil
}

func (s *Server) Serve(errCh chan error) {
	lis, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
		grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()),
	)
	playgroundv1.RegisterPlaygroundServiceServer(grpcServer, s)
	reflection.Register(grpcServer)

	s.logger.Info("starting grpc server on localhost:8000")
	errCh <- grpcServer.Serve(lis)
}
