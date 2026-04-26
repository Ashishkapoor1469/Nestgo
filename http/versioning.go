package nesthttp

import (
	"fmt"
	"net/http"

	"github.com/Ashishkapoor1469/Nestgo/common"
	"github.com/go-chi/chi/v5"
)

// VersionedRouter provides API version-based routing.
// It supports multiple API versions under a common prefix (e.g. /api/v1, /api/v2).
//
// Usage:
//
//	vr := nesthttp.NewVersionedRouter(router)
//	vr.SetPrefix("/api")             // default is /api
//	vr.SetDefaultVersion("v1")
//	vr.Version("v1", func(g *VersionGroup) {
//	    g.RegisterController(usersV1)
//	})
//	vr.Version("v2", func(g *VersionGroup) {
//	    g.RegisterController(usersV2)
//	})
type VersionedRouter struct {
	router         *Router
	prefix         string
	defaultVersion string
	versions       map[string]*VersionGroup
	deprecated     map[string]bool
}

// VersionGroup represents a versioned set of route registrations.
type VersionGroup struct {
	version    string
	prefix     string
	router     *Router
	mux        chi.Router
	deprecated bool
}

// NewVersionedRouter creates a new versioned router.
func NewVersionedRouter(router *Router) *VersionedRouter {
	return &VersionedRouter{
		router:     router,
		prefix:     "/api",
		versions:   make(map[string]*VersionGroup),
		deprecated: make(map[string]bool),
	}
}

// SetPrefix sets the base prefix for versioned routes.
func (vr *VersionedRouter) SetPrefix(prefix string) *VersionedRouter {
	vr.prefix = prefix
	return vr
}

// SetDefaultVersion sets the default API version.
func (vr *VersionedRouter) SetDefaultVersion(version string) *VersionedRouter {
	vr.defaultVersion = version
	return vr
}

// Version registers routes for a specific API version.
func (vr *VersionedRouter) Version(version string, fn func(g *VersionGroup)) {
	versionPrefix := vr.prefix + "/" + version

	group := &VersionGroup{
		version:    version,
		prefix:     versionPrefix,
		router:     vr.router,
		deprecated: vr.deprecated[version],
	}

	vr.versions[version] = group

	// Set up routes within the version group.
	vr.router.mux.Route(versionPrefix, func(r chi.Router) {
		// Add deprecation warning header for deprecated versions.
		if group.deprecated {
			r.Use(func(next http.Handler) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					w.Header().Set("X-API-Deprecation-Warning",
						fmt.Sprintf("API version %s is deprecated. Please migrate to a newer version.", version))
					w.Header().Set("Sunset", "")
					next.ServeHTTP(w, req)
				})
			})
		}

		group.mux = r
		fn(group)
	})
}

// Deprecate marks a version as deprecated.
// Requests to deprecated versions will include a deprecation warning header.
func (vr *VersionedRouter) Deprecate(version string) *VersionedRouter {
	vr.deprecated[version] = true
	return vr
}

// Versions returns all registered version names.
func (vr *VersionedRouter) Versions() []string {
	var versions []string
	for v := range vr.versions {
		versions = append(versions, v)
	}
	return versions
}

// RegisterController registers a controller within the version group.
func (g *VersionGroup) RegisterController(ctrl common.Controller) {
	prefix := normalizePrefix("", ctrl.Prefix())
	routes := ctrl.Routes()

	g.mux.Route(prefix, func(sub chi.Router) {
		for _, route := range routes {
			handler := g.buildHandler(route, ctrl)

			path := route.Path
			if path == "" {
				path = "/"
			}

			fullPath := g.prefix + prefix + path
			if path == "/" {
				fullPath = g.prefix + prefix
			}

			// Register in the main router's route table.
			g.router.routes = append(g.router.routes, RouteInfo{
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

// Use adds middleware to all routes in this version group.
func (g *VersionGroup) Use(middlewares ...func(http.Handler) http.Handler) {
	for _, mw := range middlewares {
		g.mux.Use(mw)
	}
}

// buildHandler wraps a route handler with guards and error handling.
func (g *VersionGroup) buildHandler(route common.Route, ctrl common.Controller) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := common.NewContext(w, req)

		// Run route-level guards.
		for _, guard := range route.Guards {
			allowed, err := guard.CanActivate(ctx)
			if err != nil {
				if httpErr, ok := err.(*common.HttpException); ok {
					_ = ctx.JSON(httpErr.Status, map[string]any{
						"error":      http.StatusText(httpErr.Status),
						"message":    httpErr.Message,
						"statusCode": httpErr.Status,
					})
				} else {
					_ = ctx.Error(http.StatusInternalServerError, "internal server error")
				}
				return
			}
			if !allowed {
				_ = ctx.Error(http.StatusForbidden, "access denied")
				return
			}
		}

		// Add version info to context.
		ctx.Set("api_version", g.version)

		// Execute handler.
		if err := route.Handler(ctx); err != nil {
			if httpErr, ok := err.(*common.HttpException); ok {
				_ = ctx.JSON(httpErr.Status, map[string]any{
					"error":   http.StatusText(httpErr.Status),
					"message": httpErr.Message,
				})
			} else {
				_ = ctx.Error(http.StatusInternalServerError, "internal server error")
			}
		}
	})
}


