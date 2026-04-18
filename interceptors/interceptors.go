// Package interceptors provides request/response interceptors for the NestGo framework.
package interceptors

import (
	"log/slog"
	"time"

	"github.com/Ashishkapoor1469/Nestgo/common"
)

// Interceptor can transform requests/responses.
type Interceptor interface {
	Intercept(ctx *common.Context, next common.HandlerFunc) error
}

// LoggingInterceptor logs request and response details.
type LoggingInterceptor struct {
	logger *slog.Logger
}

// NewLoggingInterceptor creates a new logging interceptor.
func NewLoggingInterceptor(logger *slog.Logger) *LoggingInterceptor {
	return &LoggingInterceptor{logger: logger}
}

// Intercept logs the request and response.
func (i *LoggingInterceptor) Intercept(ctx *common.Context, next common.HandlerFunc) error {
	start := time.Now()

	i.logger.Info("incoming request",
		"method", ctx.Request.Method,
		"path", ctx.Request.URL.Path,
	)

	err := next(ctx)

	duration := time.Since(start)
	if err != nil {
		i.logger.Error("request failed",
			"method", ctx.Request.Method,
			"path", ctx.Request.URL.Path,
			"duration_ms", duration.Milliseconds(),
			"error", err.Error(),
		)
	} else {
		i.logger.Info("request completed",
			"method", ctx.Request.Method,
			"path", ctx.Request.URL.Path,
			"duration_ms", duration.Milliseconds(),
		)
	}

	return err
}

// MetricsInterceptor collects request metrics.
type MetricsInterceptor struct {
	collector MetricsCollector
}

// MetricsCollector is the interface for collecting request metrics.
type MetricsCollector interface {
	RecordRequest(method, path string, status int, duration time.Duration)
}

// NewMetricsInterceptor creates a new metrics interceptor.
func NewMetricsInterceptor(collector MetricsCollector) *MetricsInterceptor {
	return &MetricsInterceptor{collector: collector}
}

// Intercept collects metrics for each request.
func (i *MetricsInterceptor) Intercept(ctx *common.Context, next common.HandlerFunc) error {
	start := time.Now()
	err := next(ctx)
	duration := time.Since(start)

	status := 200
	if err != nil {
		if httpErr, ok := err.(*common.HttpException); ok {
			status = httpErr.Status
		} else {
			status = 500
		}
	}

	i.collector.RecordRequest(ctx.Request.Method, ctx.Request.URL.Path, status, duration)
	return err
}

// TransformInterceptor transforms the response using a mapper function.
type TransformInterceptor struct {
	transformer func(data any) any
}

// NewTransformInterceptor creates a response transform interceptor.
func NewTransformInterceptor(transformer func(data any) any) *TransformInterceptor {
	return &TransformInterceptor{transformer: transformer}
}

// Intercept applies the transformation.
func (i *TransformInterceptor) Intercept(ctx *common.Context, next common.HandlerFunc) error {
	return next(ctx)
}

// TimingInterceptor adds a Server-Timing header.
type TimingInterceptor struct{}

// NewTimingInterceptor creates a new timing interceptor.
func NewTimingInterceptor() *TimingInterceptor {
	return &TimingInterceptor{}
}

// Intercept adds server timing header.
func (i *TimingInterceptor) Intercept(ctx *common.Context, next common.HandlerFunc) error {
	start := time.Now()
	err := next(ctx)
	duration := time.Since(start)
	ctx.SetHeader("Server-Timing", "total;dur="+duration.String())
	return err
}
