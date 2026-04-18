package common

import (
	"net/http"
)

// Route defines a single HTTP route.
type Route struct {
	Method      string
	Path        string
	Handler     HandlerFunc
	Middlewares []Middleware
	Guards      []Guard
	Summary     string // for OpenAPI
	Description string // for OpenAPI
}

// Controller is implemented by all NestGo controllers.
type Controller interface {
	Prefix() string
	Routes() []Route
}

// Guard is the authorization interface.
type Guard interface {
	CanActivate(ctx *Context) (bool, error)
}

// Interceptor can transform requests/responses.
type Interceptor interface {
	Intercept(ctx *Context, next HandlerFunc) error
}

// Middleware is a standard HTTP middleware.
type Middleware func(next http.Handler) http.Handler
