// Package plugins provides the plugin system and built-in plugins
// for the NestGo framework.
package plugins

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Plugin is the NestGo plugin interface.
type Plugin interface {
	Name() string
	Register(router chi.Router, logger *slog.Logger) error
}

// --- Health Check Plugin ---

// HealthPlugin provides /health and /ready endpoints.
type HealthPlugin struct {
	checks map[string]HealthChecker
	mu     sync.RWMutex
}

// HealthChecker performs a health check.
type HealthChecker func() error

// HealthStatus represents the health status response.
type HealthStatus struct {
	Status    string                 `json:"status"` // "healthy" or "unhealthy"
	Timestamp string                 `json:"timestamp"`
	Checks    map[string]CheckResult `json:"checks"`
	Uptime    string                 `json:"uptime"`
}

// CheckResult is an individual check result.
type CheckResult struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

var startTime = time.Now()

// NewHealthPlugin creates a new health plugin.
func NewHealthPlugin() *HealthPlugin {
	return &HealthPlugin{
		checks: make(map[string]HealthChecker),
	}
}

// AddCheck adds a named health check.
func (p *HealthPlugin) AddCheck(name string, checker HealthChecker) *HealthPlugin {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.checks[name] = checker
	return p
}

// Name returns the plugin name.
func (p *HealthPlugin) Name() string { return "health" }

// Register registers health endpoints.
func (p *HealthPlugin) Register(router chi.Router, logger *slog.Logger) error {
	router.Get("/health", p.healthHandler)
	router.Get("/ready", p.readyHandler)
	logger.Info("health plugin registered", "endpoints", []string{"/health", "/ready"})
	return nil
}

func (p *HealthPlugin) healthHandler(w http.ResponseWriter, r *http.Request) {
	status := HealthStatus{
		Status:    "healthy",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Checks:    make(map[string]CheckResult),
		Uptime:    time.Since(startTime).String(),
	}

	p.mu.RLock()
	defer p.mu.RUnlock()

	for name, checker := range p.checks {
		if err := checker(); err != nil {
			status.Status = "unhealthy"
			status.Checks[name] = CheckResult{Status: "unhealthy", Message: err.Error()}
		} else {
			status.Checks[name] = CheckResult{Status: "healthy"}
		}
	}

	code := http.StatusOK
	if status.Status != "healthy" {
		code = http.StatusServiceUnavailable
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(status)
}

func (p *HealthPlugin) readyHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"status": "ready",
		"uptime": time.Since(startTime).String(),
	})
}

// --- Metrics Plugin (Prometheus) ---

// MetricsPlugin provides Prometheus metrics at /metrics.
type MetricsPlugin struct {
	requestsTotal   *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
	requestsActive  prometheus.Gauge
	responseSize    *prometheus.HistogramVec
}

// NewMetricsPlugin creates a new metrics plugin.
func NewMetricsPlugin() *MetricsPlugin {
	p := &MetricsPlugin{
		requestsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "nestgo_http_requests_total",
				Help: "Total number of HTTP requests",
			},
			[]string{"method", "path", "status"},
		),
		requestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "nestgo_http_request_duration_seconds",
				Help:    "HTTP request duration in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "path"},
		),
		requestsActive: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "nestgo_http_requests_active",
				Help: "Number of active HTTP requests",
			},
		),
		responseSize: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "nestgo_http_response_size_bytes",
				Help:    "HTTP response size in bytes",
				Buckets: []float64{100, 1000, 10000, 100000, 1000000},
			},
			[]string{"method", "path"},
		),
	}

	prometheus.MustRegister(p.requestsTotal, p.requestDuration, p.requestsActive, p.responseSize)

	// System metrics.
	prometheus.MustRegister(prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name: "nestgo_goroutines",
			Help: "Number of goroutines",
		},
		func() float64 { return float64(runtime.NumGoroutine()) },
	))

	return p
}

// Name returns the plugin name.
func (p *MetricsPlugin) Name() string { return "metrics" }

// Register registers the /metrics endpoint and middleware.
func (p *MetricsPlugin) Register(router chi.Router, logger *slog.Logger) error {
	router.Handle("/metrics", promhttp.Handler())
	router.Use(p.Middleware)
	logger.Info("metrics plugin registered", "endpoint", "/metrics")
	return nil
}

// Middleware is the metrics collection middleware.
func (p *MetricsPlugin) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		p.requestsActive.Inc()

		ww := &metricsWriter{ResponseWriter: w, status: 200}
		next.ServeHTTP(ww, r)

		duration := time.Since(start).Seconds()
		p.requestsActive.Dec()

		p.requestsTotal.WithLabelValues(r.Method, r.URL.Path, fmt.Sprintf("%d", ww.status)).Inc()
		p.requestDuration.WithLabelValues(r.Method, r.URL.Path).Observe(duration)
		p.responseSize.WithLabelValues(r.Method, r.URL.Path).Observe(float64(ww.bytes))
	})
}

// RecordRequest records request metrics (implements interceptors.MetricsCollector).
func (p *MetricsPlugin) RecordRequest(method, path string, status int, duration time.Duration) {
	p.requestsTotal.WithLabelValues(method, path, fmt.Sprintf("%d", status)).Inc()
	p.requestDuration.WithLabelValues(method, path).Observe(duration.Seconds())
}

type metricsWriter struct {
	http.ResponseWriter
	status int
	bytes  int
}

func (w *metricsWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *metricsWriter) Write(b []byte) (int, error) {
	n, err := w.ResponseWriter.Write(b)
	w.bytes += n
	return n, err
}

// --- Feature Flags Plugin ---

// FeatureFlagsPlugin provides runtime feature flags.
type FeatureFlagsPlugin struct {
	mu    sync.RWMutex
	flags map[string]*atomic.Bool
}

// NewFeatureFlagsPlugin creates a new feature flags plugin.
func NewFeatureFlagsPlugin() *FeatureFlagsPlugin {
	return &FeatureFlagsPlugin{
		flags: make(map[string]*atomic.Bool),
	}
}

// Name returns the plugin name.
func (p *FeatureFlagsPlugin) Name() string { return "feature-flags" }

// Register registers the feature flags endpoints.
func (p *FeatureFlagsPlugin) Register(router chi.Router, logger *slog.Logger) error {
	router.Get("/features", p.listHandler)
	logger.Info("feature flags plugin registered")
	return nil
}

// Set sets a feature flag.
func (p *FeatureFlagsPlugin) Set(name string, enabled bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	flag, ok := p.flags[name]
	if !ok {
		flag = &atomic.Bool{}
		p.flags[name] = flag
	}
	flag.Store(enabled)
}

// IsEnabled checks if a feature flag is enabled.
func (p *FeatureFlagsPlugin) IsEnabled(name string) bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	flag, ok := p.flags[name]
	if !ok {
		return false
	}
	return flag.Load()
}

func (p *FeatureFlagsPlugin) listHandler(w http.ResponseWriter, r *http.Request) {
	p.mu.RLock()
	result := make(map[string]bool, len(p.flags))
	for name, flag := range p.flags {
		result[name] = flag.Load()
	}
	p.mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(result)
}
