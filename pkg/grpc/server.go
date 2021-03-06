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

	healthv1 "github.com/jacob-delgado/playground/gen/proto/go/grpc/health/v1"
	inventoryv1 "github.com/jacob-delgado/playground/gen/proto/go/inventory/v1"
	"github.com/jacob-delgado/playground/pkg/metrics"
)

type Server struct {
	logger *otelzap.Logger
	tracer trace.Tracer
}

func NewServer(logger *otelzap.Logger) *Server {
	return &Server{
		logger: logger,
		tracer: otel.Tracer("github.com/jacob-delgado/inventory/pkg/grpc"),
	}
}

func (s *Server) Check(ctx context.Context, req *healthv1.HealthCheckRequest) (*healthv1.HealthCheckResponse, error) {
	return &healthv1.HealthCheckResponse{}, nil
}

func (s *Server) Watch(req *healthv1.HealthCheckRequest, srv healthv1.Health_WatchServer) error {
	return srv.Send(&healthv1.HealthCheckResponse{
		Status: healthv1.HealthCheckResponse_SERVING,
	})
}

func (s *Server) GetInventory(ctx context.Context, req *inventoryv1.GetInventoryRequest) (*inventoryv1.GetInventoryResponse, error) {
	s.logger.Ctx(ctx).Info("GetFeature rpc call")
	metrics.GetFeatureProcessed.Inc()

	return &inventoryv1.GetInventoryResponse{}, nil
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
	healthv1.RegisterHealthServer(grpcServer, s)
	inventoryv1.RegisterInventoryServiceServer(grpcServer, s)
	reflection.Register(grpcServer)

	s.logger.Info("starting grpc server on localhost:8000")
	errCh <- grpcServer.Serve(lis)
}
