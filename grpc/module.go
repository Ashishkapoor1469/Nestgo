package grpc

import (
	"log/slog"

	"github.com/Ashishkapoor1469/Nestgo/common"
	"github.com/Ashishkapoor1469/Nestgo/di"
)

// GRPCModule integrates gRPC server into NestGo's module system.
type GRPCModule struct {
	server *GRPCServer
}

// NewGRPCModule creates a new GRPCModule.
func NewGRPCModule(logger *slog.Logger) *GRPCModule {
	return &GRPCModule{
		server: NewGRPCServer(logger),
	}
}

// Module returns the common.ModuleConfig.
func (m *GRPCModule) Module() common.ModuleConfig {
	return common.ModuleConfig{
		Name: "GRPCModule",
		Providers: []di.Provider{
			{
				Instance: m.server,
			},
		},
	}
}

// Server returns the GRPCServer instance.
func (m *GRPCModule) Server() *GRPCServer {
	return m.server
}
