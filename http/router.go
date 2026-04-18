package nesthttp

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/nestgo/nestgo/common"
)

// Router wraps chi and provides NestGo routing capabilities.
type Router struct {
	mux                *chi.Mux
	prefix             string
	globalMiddleware   []common.Middleware
	globalGuards       []common.Guard
	globalInterceptors []common.Interceptor
	errorHandler       ErrorHandler
	routes             []RouteInfo
}

// RouteInfo stores metadata about registered routes for introspection.
type RouteInfo struct {
	Method     string
	Path       string
	Controller string
	Handler    string
	Summary    string
}

// ErrorHandler handles errors returned by handlers.
type ErrorHandler func(ctx *common.Context, err error)

// NewRouter creates a new NestGo router.
func NewRouter() *Router {
	mux := chi.NewRouter()
	return &Router{
		mux:          mux,
		errorHandler: defaultErrorHandler,
	}
}

// Use adds global middleware.
func (r *Router) Use(middlewares ...func(http.Handler) http.Handler) {
	for _, mw := range middlewares {
		r.globalMiddleware = append(r.globalMiddleware, common.Middleware(mw))
		r.mux.Use(mw)
	}
}

// UseGuard adds a global guard.
func (r *Router) UseGuard(guards ...common.Guard) {
	r.globalGuards = append(r.globalGuards, guards...)
}

// UseInterceptor adds a global interceptor.
func (r *Router) UseInterceptor(interceptors ...common.Interceptor) {
	r.globalInterceptors = append(r.globalInterceptors, interceptors...)
}

// SetErrorHandler sets a custom error handler.
func (r *Router) SetErrorHandler(handler ErrorHandler) {
	r.errorHandler = handler
}

// RegisterController registers a controller's routes on the router.
func (r *Router) RegisterController(ctrl common.Controller) {
	prefix := normalizePrefix(r.prefix, ctrl.Prefix())
	routes := ctrl.Routes()

	r.mux.Route(prefix, func(sub chi.Router) {
		for _, route := range routes {
			handler := r.buildHandler(route, ctrl)

			path := route.Path
			if path == "" {
				path = "/"
			}

			// Register route info.
			fullPath := prefix + path
			if path == "/" {
				fullPath = prefix
			}
			r.routes = append(r.routes, RouteInfo{
				Method:     route.Method,
				Path:       fullPath,
				Controller: fmt.Sprintf("%T", ctrl),
				Summary:    route.Summary,
			})

			// Apply route-level middleware.
			if len(route.Middlewares) > 0 {
				chiMW := make([]func(http.Handler) http.Handler, len(route.Middlewares))
				for i, mw := range route.Middlewares {
					chiMW[i] = mw
				}
				sub.With(chiMW...).Method(route.Method, path, handler)
			} else {
				sub.Method(route.Method, path, handler)
			}
		}
	})
}

// SetGlobalPrefix sets a global prefix for all routes.
func (r *Router) SetGlobalPrefix(prefix string) {
	r.prefix = prefix
}

// Handler returns the underlying http.Handler.
func (r *Router) Handler() http.Handler {
	return r.mux
}

// Routes returns all registered route info for introspection.
func (r *Router) Routes() []RouteInfo {
	return r.routes
}

// ServeHTTP implements http.Handler.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}

// buildHandler wraps a route handler with guards, interceptors, and error handling.
func (r *Router) buildHandler(route common.Route, ctrl common.Controller) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := common.NewContext(w, req)

		// Run global guards.
		for _, guard := range r.globalGuards {
			allowed, err := guard.CanActivate(ctx)
			if err != nil {
				r.errorHandler(ctx, err)
				return
			}
			if !allowed {
				_ = ctx.Error(http.StatusForbidden, "access denied")
				return
			}
		}

		// Run route-level guards.
		for _, guard := range route.Guards {
			allowed, err := guard.CanActivate(ctx)
			if err != nil {
				r.errorHandler(ctx, err)
				return
			}
			if !allowed {
				_ = ctx.Error(http.StatusForbidden, "access denied")
				return
			}
		}

		// Build interceptor chain.
		finalHandler := route.Handler

		// Wrap with global interceptors (reverse order for proper nesting).
		for i := len(r.globalInterceptors) - 1; i >= 0; i-- {
			interceptor := r.globalInterceptors[i]
			next := finalHandler
			finalHandler = func(ctx *common.Context) error {
				return interceptor.Intercept(ctx, next)
			}
		}

		// Execute handler.
		if err := finalHandler(ctx); err != nil {
			r.errorHandler(ctx, err)
		}
	})
}

// defaultErrorHandler is the built-in error handler.
func defaultErrorHandler(ctx *common.Context, err error) {
	if ctx.IsWritten() {
		return
	}

	// Check if error is an HttpException.
	if httpErr, ok := err.(*common.HttpException); ok {
		_ = ctx.JSON(httpErr.Status, map[string]any{
			"error":      http.StatusText(httpErr.Status),
			"message":    httpErr.Message,
			"statusCode": httpErr.Status,
		})
		return
	}

	_ = ctx.Error(http.StatusInternalServerError, "internal server error")
}

func normalizePrefix(global, local string) string {
	path := global + "/" + strings.TrimPrefix(local, "/")
	path = strings.TrimSuffix(path, "/")
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	// Remove double slashes.
	for strings.Contains(path, "//") {
		path = strings.ReplaceAll(path, "//", "/")
	}
	return path
}
