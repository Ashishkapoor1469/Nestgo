package modules

import (
	"net/http"

	"github.com/Ashishkapoor1469/Nestgo/common"
)

// AppController handles incoming root-level HTTP requests.
type AppController struct{}

// NewAppController acts as a factory function.
func NewAppController() *AppController {
	return &AppController{}
}

// Prefix returns the route prefix for this controller.
func (c *AppController) Prefix() string {
	return "/"
}

// Routes connects the paths to their handlers.
func (c *AppController) Routes() []common.Route {
	return []common.Route{
		{
			Method:  http.MethodGet,
			Path:    "/",
			Handler: c.welcome,
			Summary: "Welcome to the API",
		},
		{
			Method:  http.MethodGet,
			Path:    "/ping",
			Handler: c.ping,
			Summary: "Health check",
		},
	}
}

func (c *AppController) welcome(ctx *common.Context) error {
	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"message":   "Welcome to the NestGo Sample App API!",
		"version":   "1.0.0",
		"endpoints": []string{"/ping", "/api/auth/login", "/api/auth/profile"},
	})
}

func (c *AppController) ping(ctx *common.Context) error {
	return ctx.JSON(http.StatusOK, map[string]string{
		"status": "pong",
	})
}
