package grpc

import (
	"log/slog"
	"net"

	"google.golang.org/grpc"
)

// GRPCServer wraps the standard grpc.Server with logger and registered services.
type GRPCServer struct {
	server     *grpc.Server
	registrars []ServiceRegistrar
	logger     *slog.Logger
}

// ServiceRegistrar is an interface that gRPC controllers/services implement
// to register themselves with the gRPC server.
type ServiceRegistrar interface {
	RegisterGRPC(server *grpc.Server)
}

// NewGRPCServer creates a new GRPCServer wrapper.
func NewGRPCServer(logger *slog.Logger, opts ...grpc.ServerOption) *GRPCServer {
	return &GRPCServer{
		server: grpc.NewServer(opts...),
		logger: logger,
	}
}

// Register registers a gRPC service registrar.
func (s *GRPCServer) Register(registrar ServiceRegistrar) {
	s.registrars = append(s.registrars, registrar)
}

// Start starts the gRPC server listening on the given address.
func (s *GRPCServer) Start(addr string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	// Register all services
	for _, registrar := range s.registrars {
		registrar.RegisterGRPC(s.server)
	}

	s.logger.Info("gRPC server started", "address", addr)
	return s.server.Serve(lis)
}

// Stop gracefully stops the gRPC server.
func (s *GRPCServer) Stop() {
	s.server.GracefulStop()
	s.logger.Info("gRPC server stopped")
}

// Server returns the underlying grpc.Server.
func (s *GRPCServer) Server() *grpc.Server {
	return s.server
}
