// Package errors provides centralized error handling and exception filters
// for the NestGo framework.
package errors

import (
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/nestgo/nestgo/common"
)

// ExceptionFilter handles specific types of errors.
type ExceptionFilter interface {
	// Catch handles an error and writes the response.
	Catch(err error, ctx *common.Context) bool
}

// ErrorHandler is the centralized error handler with filter chain.
type ErrorHandler struct {
	filters []ExceptionFilter
	logger  *slog.Logger
}

// NewErrorHandler creates a new error handler.
func NewErrorHandler(logger *slog.Logger) *ErrorHandler {
	return &ErrorHandler{
		logger: logger,
	}
}

// AddFilter adds an exception filter to the chain.
func (h *ErrorHandler) AddFilter(f ExceptionFilter) {
	h.filters = append(h.filters, f)
}

// Handle processes an error through the filter chain.
func (h *ErrorHandler) Handle(ctx *common.Context, err error) {
	if ctx.IsWritten() {
		return
	}

	// Try each filter in order.
	for _, f := range h.filters {
		if f.Catch(err, ctx) {
			return
		}
	}

	// Default error handling.
	if httpErr, ok := err.(*common.HttpException); ok {
		_ = ctx.JSON(httpErr.Status, map[string]any{
			"error":      http.StatusText(httpErr.Status),
			"message":    httpErr.Message,
			"statusCode": httpErr.Status,
		})
		return
	}

	// Log unexpected errors.
	h.logger.Error("unhandled error",
		"error", err.Error(),
	)

	_ = ctx.Error(http.StatusInternalServerError, "internal server error")
}

// HandlePanic recovers from a panic and returns an error response.
func (h *ErrorHandler) HandlePanic(ctx *common.Context, recovered any) {
	stack := string(debug.Stack())
	h.logger.Error("panic recovered",
		"panic", fmt.Sprintf("%v", recovered),
		"stack", stack,
	)

	if !ctx.IsWritten() {
		_ = ctx.Error(http.StatusInternalServerError, "internal server error")
	}
}

// --- Built-in Exception Filters ---

// ValidationExceptionFilter handles validation errors.
type ValidationExceptionFilter struct{}

func (f *ValidationExceptionFilter) Catch(err error, ctx *common.Context) bool {
	if valErr, ok := err.(*common.ValidationErrors); ok {
		_ = ctx.JSON(http.StatusUnprocessableEntity, map[string]any{
			"error":      "Validation Failed",
			"message":    valErr.Error(),
			"statusCode": http.StatusUnprocessableEntity,
			"details":    valErr.Errors,
		})
		return true
	}
	return false
}

// NotFoundExceptionFilter handles not found errors.
type NotFoundExceptionFilter struct{}

func (f *NotFoundExceptionFilter) Catch(err error, ctx *common.Context) bool {
	if httpErr, ok := err.(*common.HttpException); ok && httpErr.Status == http.StatusNotFound {
		_ = ctx.JSON(http.StatusNotFound, map[string]any{
			"error":      "Not Found",
			"message":    httpErr.Message,
			"statusCode": http.StatusNotFound,
		})
		return true
	}
	return false
}

// --- Common Error Types ---

// AppError is a general application error with context.
type AppError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"statusCode"`
	Details any    `json:"details,omitempty"`
}

func (e *AppError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// NewAppError creates a new application error.
func NewAppError(code string, status int, message string) *AppError {
	return &AppError{
		Code:    code,
		Status:  status,
		Message: message,
	}
}

// WithDetails adds details to the error.
func (e *AppError) WithDetails(details any) *AppError {
	e.Details = details
	return e
}

// AppErrorFilter handles AppError types.
type AppErrorFilter struct{}

func (f *AppErrorFilter) Catch(err error, ctx *common.Context) bool {
	if appErr, ok := err.(*AppError); ok {
		_ = ctx.JSON(appErr.Status, map[string]any{
			"error":      appErr.Code,
			"message":    appErr.Message,
			"statusCode": appErr.Status,
			"details":    appErr.Details,
		})
		return true
	}
	return false
}
