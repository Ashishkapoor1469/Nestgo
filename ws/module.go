package ws

import (
	"log/slog"

	"github.com/Ashishkapoor1469/Nestgo/common"
	"github.com/Ashishkapoor1469/Nestgo/di"
)

// WebSocketModule integrates the WebSocket gateway into the NestGo module system.
type WebSocketModule struct {
	gateway *Gateway
	path    string
}

// NewWebSocketModule creates a new WebSocket module.
func NewWebSocketModule(path string, logger *slog.Logger) *WebSocketModule {
	return &WebSocketModule{
		gateway: NewGateway(logger),
		path:    path,
	}
}

// Module returns the configuration for the WebSocketModule.
func (m *WebSocketModule) Module() common.ModuleConfig {
	return common.ModuleConfig{
		Name: "WebSocketModule",
		Controllers: []common.Controller{
			&wsController{gw: m.gateway, path: m.path},
		},
		Providers: []di.Provider{
			{Instance: m.gateway},
		},
		Exports: []any{
			m.gateway,
		},
	}
}

// Gateway returns the underlying Gateway instance.
func (m *WebSocketModule) Gateway() *Gateway {
	return m.gateway
}

type wsController struct {
	gw   *Gateway
	path string
}

func (c *wsController) Prefix() string {
	return c.path
}

func (c *wsController) Routes() []common.Route {
	return []common.Route{
		{
			Method: "GET",
			Path:   "",
			Handler: func(ctx *common.Context) error {
				c.gw.HandleHTTP(ctx.Writer, ctx.Request)
				return nil
			},
			Summary:     "WebSocket endpoint",
			Description: "Upgrades the connection to the WebSocket protocol.",
		},
	}
}
