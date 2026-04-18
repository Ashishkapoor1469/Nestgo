package core

import (
	"log/slog"
)

// Option is a functional option for configuring a NestGoApp.
type Option func(*AppConfig)

// AppConfig holds all configuration for a NestGoApp.
type AppConfig struct {
	Address         string
	Logger          *slog.Logger
	GlobalPrefix    string
	CorsEnabled     bool
	CorsOrigins     []string
	ShutdownTimeout int // seconds
	EnableMetrics   bool
	EnableTracing   bool
}

// DefaultAppConfig returns sensible defaults.
func DefaultAppConfig() *AppConfig {
	return &AppConfig{
		Address:         ":3000",
		Logger:          slog.Default(),
		GlobalPrefix:    "",
		CorsEnabled:     false,
		ShutdownTimeout: 30,
		EnableMetrics:   false,
		EnableTracing:   false,
	}
}

// WithAddress sets the listen address.
func WithAddress(addr string) Option {
	return func(c *AppConfig) {
		c.Address = addr
	}
}

// WithLogger sets the logger.
func WithLogger(logger *slog.Logger) Option {
	return func(c *AppConfig) {
		c.Logger = logger
	}
}

// WithGlobalPrefix sets a global API prefix (e.g., "/api/v1").
func WithGlobalPrefix(prefix string) Option {
	return func(c *AppConfig) {
		c.GlobalPrefix = prefix
	}
}

// WithCors enables CORS with the specified origins.
func WithCors(origins ...string) Option {
	return func(c *AppConfig) {
		c.CorsEnabled = true
		c.CorsOrigins = origins
	}
}

// WithShutdownTimeout sets the graceful shutdown timeout in seconds.
func WithShutdownTimeout(seconds int) Option {
	return func(c *AppConfig) {
		c.ShutdownTimeout = seconds
	}
}

// WithMetrics enables Prometheus metrics endpoint.
func WithMetrics() Option {
	return func(c *AppConfig) {
		c.EnableMetrics = true
	}
}

// WithTracing enables distributed tracing.
func WithTracing() Option {
	return func(c *AppConfig) {
		c.EnableTracing = true
	}
}
