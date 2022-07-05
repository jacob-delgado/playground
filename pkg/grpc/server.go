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

	playgroundv1 "github.com/jacob-delgado/playground/gen/proto/go/playground/v1"
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

func (s *Server) GetFeature(context.Context, *playgroundv1.GetFeatureRequest) (*playgroundv1.GetFeatureResponse, error) {
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

	s.logger.Info("starting grpc server on localhost:8000")
	errCh <- grpcServer.Serve(lis)
}
