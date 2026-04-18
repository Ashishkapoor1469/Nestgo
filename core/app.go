package core

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Ashishkapoor1469/Nestgo/common"
	"github.com/Ashishkapoor1469/Nestgo/di"
	nesthttp "github.com/Ashishkapoor1469/Nestgo/http"
	"github.com/Ashishkapoor1469/Nestgo/middleware"
)

// NestGoApp is the central application container.
type NestGoApp struct {
	config    *AppConfig
	container *di.Container
	router    *nesthttp.Router
	registry  *ModuleRegistry
	lifecycle *LifecycleHost
	server    *http.Server
	plugins   []Plugin
	logger    *slog.Logger
}

// Plugin is the plugin interface.
type Plugin interface {
	Name() string
	Register(app *NestGoApp) error
}

// New creates a new NestGoApp with functional options.
func New(opts ...Option) *NestGoApp {
	config := DefaultAppConfig()
	for _, opt := range opts {
		opt(config)
	}

	app := &NestGoApp{
		config:    config,
		container: di.NewContainer(),
		router:    nesthttp.NewRouter(),
		registry:  NewModuleRegistry(),
		lifecycle: NewLifecycleHost(),
		logger:    config.Logger,
	}

	// Apply global prefix.
	if config.GlobalPrefix != "" {
		app.router.SetGlobalPrefix(config.GlobalPrefix)
	}

	// Register core middleware.
	app.router.Use(
		middleware.Recovery(app.logger),
		middleware.RequestID(),
		middleware.Logger(app.logger),
	)

	// CORS.
	if config.CorsEnabled {
		app.router.Use(middleware.CORS(config.CorsOrigins...))
	}

	// Security headers.
	app.router.Use(middleware.SecureHeaders())

	// Register the app itself in container.
	_ = app.container.ProvideValue(app)
	_ = app.container.ProvideValue(app.router)
	_ = app.container.ProvideValue(app.container)
	_ = app.container.ProvideValue(app.logger)

	return app
}

// RegisterModule registers a module with the application.
func (app *NestGoApp) RegisterModule(m common.Module) *NestGoApp {
	if err := app.registry.Register(m); err != nil {
		app.logger.Error("failed to register module", "error", err)
	}
	return app
}

// RegisterPlugin registers a plugin with the application.
func (app *NestGoApp) RegisterPlugin(p Plugin) *NestGoApp {
	app.plugins = append(app.plugins, p)
	return app
}

// UseMiddleware adds global middleware.
func (app *NestGoApp) UseMiddleware(mw ...func(http.Handler) http.Handler) *NestGoApp {
	app.router.Use(mw...)
	return app
}

// UseGuard adds a global guard.
func (app *NestGoApp) UseGuard(guards ...common.Guard) *NestGoApp {
	app.router.UseGuard(guards...)
	return app
}

// UseInterceptor adds a global interceptor.
func (app *NestGoApp) UseInterceptor(interceptors ...common.Interceptor) *NestGoApp {
	app.router.UseInterceptor(interceptors...)
	return app
}

// Container returns the DI container.
func (app *NestGoApp) Container() *di.Container {
	return app.container
}

// Router returns the HTTP router.
func (app *NestGoApp) Router() *nesthttp.Router {
	return app.router
}

// Logger returns the application logger.
func (app *NestGoApp) Logger() *slog.Logger {
	return app.logger
}

// Start boots the application and starts the HTTP server.
func (app *NestGoApp) Start(addr ...string) error {
	listenAddr := app.config.Address
	if len(addr) > 0 && addr[0] != "" {
		listenAddr = addr[0]
	}

	// Phase 1: Initialize plugins.
	for _, p := range app.plugins {
		app.logger.Info("registering plugin", "plugin", p.Name())
		if err := p.Register(app); err != nil {
			return fmt.Errorf("plugin %s registration failed: %w", p.Name(), err)
		}
	}

	// Phase 2: Resolve module dependencies.
	order, err := app.registry.ResolveOrder()
	if err != nil {
		return fmt.Errorf("module resolution failed: %w", err)
	}

	// Phase 3: Register providers and controllers from modules.
	for _, mod := range order {
		cfg := mod.Module()
		app.logger.Info("initializing module", "module", cfg.Name)

		// Register providers.
		for _, provider := range cfg.Providers {
			if provider.Factory != nil {
				if err := app.container.Provide(provider.Factory); err != nil {
					return fmt.Errorf("module %s: provider registration failed: %w", cfg.Name, err)
				}
			} else if provider.Instance != nil {
				if err := app.container.ProvideValue(provider.Instance); err != nil {
					return fmt.Errorf("module %s: provider value registration failed: %w", cfg.Name, err)
				}
			}
		}

		// Register the module itself.
		app.lifecycle.Register(mod)
	}

	// Phase 4: Validate and resolve all dependencies.
	if err := app.container.Validate(); err != nil {
		return fmt.Errorf("dependency validation failed: %w", err)
	}

	instances, err := app.container.ResolveAll()
	if err != nil {
		return fmt.Errorf("dependency resolution failed: %w", err)
	}

	// Register lifecycle hooks from all resolved instances.
	for _, inst := range instances {
		app.lifecycle.Register(inst)
	}

	// Phase 5: Register controllers.
	for _, mod := range order {
		cfg := mod.Module()
		for _, ctrl := range cfg.Controllers {
			if httpCtrl, ok := ctrl.(common.Controller); ok {
				app.router.RegisterController(httpCtrl)
			}
		}
	}

	// Phase 6: Run init hooks.
	if err := app.lifecycle.RunInitHooks(); err != nil {
		return fmt.Errorf("init hooks failed: %w", err)
	}

	// Phase 7: Start HTTP server.
	app.server = &http.Server{
		Addr:         listenAddr,
		Handler:      app.router.Handler(),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Print route table.
	app.printRoutes()

	// Run start hooks.
	if err := app.lifecycle.RunStartHooks(); err != nil {
		return fmt.Errorf("start hooks failed: %w", err)
	}

	// Graceful shutdown listener.
	go app.listenForShutdown()

	app.logger.Info("🚀 NestGo application started",
		"address", listenAddr,
		"modules", len(order),
		"routes", len(app.router.Routes()),
	)

	if err := app.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("server error: %w", err)
	}

	return nil
}

// Shutdown gracefully shuts down the application.
func (app *NestGoApp) Shutdown(ctx context.Context) error {
	app.logger.Info("shutting down NestGo application...")

	// Run shutdown hooks.
	if err := app.lifecycle.RunShutdownHooks(); err != nil {
		app.logger.Error("shutdown hook error", "error", err)
	}

	// Shutdown HTTP server.
	if app.server != nil {
		if err := app.server.Shutdown(ctx); err != nil {
			return fmt.Errorf("server shutdown error: %w", err)
		}
	}

	// Run destroy hooks.
	if err := app.lifecycle.RunDestroyHooks(); err != nil {
		app.logger.Error("destroy hook error", "error", err)
	}

	app.logger.Info("NestGo application stopped")
	return nil
}

// listenForShutdown waits for OS signals and triggers graceful shutdown.
func (app *NestGoApp) listenForShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(app.config.ShutdownTimeout)*time.Second)
	defer cancel()

	if err := app.Shutdown(ctx); err != nil {
		app.logger.Error("forced shutdown", "error", err)
		os.Exit(1)
	}
}

// printRoutes logs all registered routes.
func (app *NestGoApp) printRoutes() {
	routes := app.router.Routes()
	if len(routes) == 0 {
		return
	}

	app.logger.Info("registered routes", "count", len(routes))
	for _, r := range routes {
		app.logger.Info("route",
			"method", r.Method,
			"path", r.Path,
			"controller", r.Controller,
		)
	}
}
