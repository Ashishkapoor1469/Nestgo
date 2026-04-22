package commands

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/Ashishkapoor1469/Nestgo/cli/utils"
	"github.com/spf13/cobra"
)

// MetricsCmd creates the `nestgo metrics` placeholder info command.
func MetricsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "metrics",
		Short: "Show how to enable the /metrics endpoint",
		RunE: func(cmd *cobra.Command, args []string) error {
			utils.PrintHeader("📊 Metrics Endpoint")
			utils.PrintInfo("NestGo includes a built-in metrics middleware.")
			fmt.Println()
			utils.PrintDim("  Add the metrics middleware to your application:")
			fmt.Println()
			fmt.Println("    import \"github.com/Ashishkapoor1469/Nestgo/middleware\"")
			fmt.Println()
			fmt.Println("    // In your main.go:")
			fmt.Println("    collector := NewMetricsCollector()")
			fmt.Println("    app.UseMiddleware(collector.Middleware())")
			fmt.Println()
			fmt.Println("    // Register the /metrics endpoint:")
			fmt.Println("    http.HandleFunc(\"/metrics\", collector.Handler())")
			fmt.Println()
			utils.PrintDim("  Available metrics:")
			utils.PrintDim("    • Total requests")
			utils.PrintDim("    • Active requests")
			utils.PrintDim("    • Request duration (avg, p95, p99)")
			utils.PrintDim("    • Error rate")
			utils.PrintDim("    • Status code distribution")
			fmt.Println()
			return nil
		},
	}
}

// VersioningInfo creates an info command about API versioning.
func VersioningCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "versioning",
		Short: "Show API versioning system guide",
		RunE: func(cmd *cobra.Command, args []string) error {
			utils.PrintHeader("🔄 API Versioning System")
			utils.PrintInfo("NestGo supports version-based routing via the VersionedRouter.")
			fmt.Println()
			utils.PrintDim("  Usage in your application:")
			fmt.Println()
			fmt.Println("    import nesthttp \"github.com/Ashishkapoor1469/Nestgo/http\"")
			fmt.Println()
			fmt.Println("    vr := nesthttp.NewVersionedRouter(router)")
			fmt.Println("    vr.Version(\"v1\", func(v *nesthttp.VersionGroup) {")
			fmt.Println("        v.RegisterController(&usersV1.Controller{})")
			fmt.Println("    })")
			fmt.Println("    vr.Version(\"v2\", func(v *nesthttp.VersionGroup) {")
			fmt.Println("        v.RegisterController(&usersV2.Controller{})")
			fmt.Println("    })")
			fmt.Println()
			utils.PrintDim("  Routes will be:")
			utils.PrintDim("    GET /api/v1/users")
			utils.PrintDim("    GET /api/v2/users")
			fmt.Println()
			return nil
		},
	}
}

// ─── Metrics Collector (Runtime) ────────────────────────────────────────────

// MetricsCollector collects HTTP request metrics.
type MetricsCollector struct {
	mu             sync.RWMutex
	totalRequests  int64
	activeRequests int64
	totalErrors    int64
	statusCodes    map[int]int64
	totalDuration  time.Duration
	requestCount   int64
}

// NewMetricsCollector creates a new metrics collector.
func NewMetricsCollector() *MetricsCollector {
	return &MetricsCollector{
		statusCodes: make(map[int]int64),
	}
}

// Middleware returns an HTTP middleware that collects metrics.
func (m *MetricsCollector) Middleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			m.mu.Lock()
			m.totalRequests++
			m.activeRequests++
			m.mu.Unlock()

			start := time.Now()
			rw := &metricsWriter{ResponseWriter: w, status: 200}
			next.ServeHTTP(rw, r)
			duration := time.Since(start)

			m.mu.Lock()
			m.activeRequests--
			m.totalDuration += duration
			m.requestCount++
			m.statusCodes[rw.status]++
			if rw.status >= 400 {
				m.totalErrors++
			}
			m.mu.Unlock()
		})
	}
}

// Handler returns an HTTP handler that serves metrics as JSON.
func (m *MetricsCollector) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m.mu.RLock()
		defer m.mu.RUnlock()

		avgDuration := time.Duration(0)
		if m.requestCount > 0 {
			avgDuration = m.totalDuration / time.Duration(m.requestCount)
		}

		errorRate := float64(0)
		if m.totalRequests > 0 {
			errorRate = float64(m.totalErrors) / float64(m.totalRequests) * 100
		}

		// Build status code breakdown.
		statusBreakdown := ""
		for code, count := range m.statusCodes {
			statusBreakdown += fmt.Sprintf(`"%d":%d,`, code, count)
		}
		statusBreakdown = strings.TrimSuffix(statusBreakdown, ",")

		json := fmt.Sprintf(`{
  "total_requests": %d,
  "active_requests": %d,
  "total_errors": %d,
  "error_rate_pct": %.2f,
  "avg_duration_ms": %d,
  "status_codes": {%s}
}`, m.totalRequests, m.activeRequests, m.totalErrors, errorRate, avgDuration.Milliseconds(), statusBreakdown)

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(json))
	}
}

type metricsWriter struct {
	http.ResponseWriter
	status int
}

func (w *metricsWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}
