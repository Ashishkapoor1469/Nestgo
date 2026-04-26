// Package logger provides structured logging for the NestGo framework.
package logger

import (
	"context"
	"io"
	"log/slog"
	"os"
	"strings"
)

// Level represents the log level.
type Level = slog.Level

const (
	LevelDebug = slog.LevelDebug
	LevelInfo  = slog.LevelInfo
	LevelWarn  = slog.LevelWarn
	LevelError = slog.LevelError
)

// Logger is the NestGo structured logger built on slog.
type Logger struct {
	*slog.Logger
	level *slog.LevelVar
}

// Config holds logger configuration.
type Config struct {
	Level      string `env:"LOG_LEVEL" default:"info"`
	Format     string `env:"LOG_FORMAT" default:"json"` // "json" or "text"
	Output     io.Writer
	AddSource  bool
	TimeFormat string
}

// New creates a new structured logger.
func New(cfg Config) *Logger {
	if cfg.Output == nil {
		cfg.Output = os.Stdout
	}

	level := &slog.LevelVar{}
	level.Set(parseLevel(cfg.Level))

	opts := &slog.HandlerOptions{
		Level:     level,
		AddSource: cfg.AddSource,
	}

	var handler slog.Handler
	switch strings.ToLower(cfg.Format) {
	case "text":
		handler = slog.NewTextHandler(cfg.Output, opts)
	default:
		handler = slog.NewJSONHandler(cfg.Output, opts)
	}

	return &Logger{
		Logger: slog.New(handler),
		level:  level,
	}
}

// Default creates a default JSON logger.
func Default() *Logger {
	return New(Config{
		Level:  "info",
		Format: "json",
	})
}

// SetLevel dynamically changes the log level.
func (l *Logger) SetLevel(level string) {
	l.level.Set(parseLevel(level))
}

// WithContext returns a logger with context values (trace_id, request_id).
func (l *Logger) WithContext(ctx context.Context) *Logger {
	logger := l.Logger
	if traceID, ok := ctx.Value("trace_id").(string); ok && traceID != "" {
		logger = logger.With("trace_id", traceID)
	}
	if requestID, ok := ctx.Value("request_id").(string); ok && requestID != "" {
		logger = logger.With("request_id", requestID)
	}
	return &Logger{Logger: logger, level: l.level}
}

// WithFields returns a new logger with additional fields.
func (l *Logger) WithFields(fields map[string]any) *Logger {
	args := make([]any, 0, len(fields)*2)
	for k, v := range fields {
		args = append(args, k, v)
	}
	return &Logger{
		Logger: l.With(args...),
		level:  l.level,
	}
}

// WithComponent returns a logger scoped to a specific component.
func (l *Logger) WithComponent(name string) *Logger {
	return &Logger{
		Logger: l.With("component", name),
		level:  l.level,
	}
}

// WithModule returns a logger scoped to a specific module.
func (l *Logger) WithModule(name string) *Logger {
	return &Logger{
		Logger: l.With("module", name),
		level:  l.level,
	}
}

// SlogLogger returns the underlying slog.Logger for stdlib compatibility.
func (l *Logger) SlogLogger() *slog.Logger {
	return l.Logger
}

func parseLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
