package common

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

// HandlerFunc is the NestGo HTTP handler signature.
type HandlerFunc func(ctx *Context) error

// Context wraps an HTTP request and response with convenience methods.
type Context struct {
	Writer  http.ResponseWriter
	Request *http.Request
	params  map[string]string
	query   map[string][]string
	values  map[string]any
	status  int
	written bool
}

// NewContext creates a new HTTP context.
func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Writer:  w,
		Request: r,
		params:  make(map[string]string),
		query:   r.URL.Query(),
		values:  make(map[string]any),
		status:  http.StatusOK,
	}
}

// --- Request Binding ---

// Bind decodes the JSON request body into the given struct.
func (c *Context) Bind(v any) error {
	if c.Request.Body == nil {
		return fmt.Errorf("request body is empty")
	}
	defer c.Request.Body.Close()
	decoder := json.NewDecoder(c.Request.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(v); err != nil {
		return fmt.Errorf("invalid request body: %w", err)
	}

	// Run validation if the struct implements Validatable.
	if val, ok := v.(Validatable); ok {
		if err := val.Validate(); err != nil {
			return err
		}
	}

	return nil
}

// Param returns a URL path parameter.
func (c *Context) Param(key string) string {
	if v := chi.URLParam(c.Request, key); v != "" {
		return v
	}
	return c.params[key]
}

// ParamInt returns a URL path parameter as int.
func (c *Context) ParamInt(key string) (int, error) {
	return strconv.Atoi(c.Param(key))
}

// Query returns a query string parameter.
func (c *Context) Query(key string) string {
	vals := c.query[key]
	if len(vals) > 0 {
		return vals[0]
	}
	return ""
}

// QueryDefault returns a query string parameter with a default value.
func (c *Context) QueryDefault(key, defaultVal string) string {
	v := c.Query(key)
	if v == "" {
		return defaultVal
	}
	return v
}

// QueryInt returns a query string parameter as int.
func (c *Context) QueryInt(key string, defaultVal int) int {
	v := c.Query(key)
	if v == "" {
		return defaultVal
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return defaultVal
	}
	return n
}

// Header returns a request header value.
func (c *Context) Header(key string) string {
	return c.Request.Header.Get(key)
}

// SetHeader sets a response header.
func (c *Context) SetHeader(key, value string) {
	c.Writer.Header().Set(key, value)
}

// BearerToken extracts the bearer token from the Authorization header.
func (c *Context) BearerToken() string {
	auth := c.Header("Authorization")
	if strings.HasPrefix(auth, "Bearer ") {
		return strings.TrimPrefix(auth, "Bearer ")
	}
	return ""
}

// --- Context Values ---

// Set stores a value in the context.
func (c *Context) Set(key string, value any) {
	c.values[key] = value
}

// Get retrieves a value from the context.
func (c *Context) Get(key string) (any, bool) {
	v, ok := c.values[key]
	return v, ok
}

// MustGet retrieves a value from the context, panics if not found.
func (c *Context) MustGet(key string) any {
	v, ok := c.values[key]
	if !ok {
		panic("nestgo: context key not found: " + key)
	}
	return v
}

// --- Response Methods ---

// Status sets the HTTP status code.
func (c *Context) Status(code int) *Context {
	c.status = code
	return c
}

// JSON sends a JSON response.
func (c *Context) JSON(status int, v any) error {
	c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	c.Writer.WriteHeader(status)
	c.written = true
	return json.NewEncoder(c.Writer).Encode(v)
}

// OK sends a 200 JSON response.
func (c *Context) OK(v any) error {
	return c.JSON(http.StatusOK, v)
}

// Created sends a 201 JSON response.
func (c *Context) Created(v any) error {
	return c.JSON(http.StatusCreated, v)
}

// NoContent sends a 204 response.
func (c *Context) NoContent() error {
	c.Writer.WriteHeader(http.StatusNoContent)
	c.written = true
	return nil
}

// Error sends an error response.
func (c *Context) Error(status int, message string) error {
	return c.JSON(status, map[string]any{
		"error":      http.StatusText(status),
		"message":    message,
		"statusCode": status,
	})
}

// Success sends a standard success envelope.
func (c *Context) Success(data any) error {
	return c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    data,
	})
}

// Paginated sends a paginated response.
func (c *Context) Paginated(data any, total, page, perPage int) error {
	return c.JSON(http.StatusOK, PaginatedResponse{
		Success: true,
		Data:    data,
		Meta: PaginationMeta{
			Total:   total,
			Page:    page,
			PerPage: perPage,
			Pages:   (total + perPage - 1) / perPage,
		},
	})
}

// IsWritten returns true if the response has already been written.
func (c *Context) IsWritten() bool {
	return c.written
}

// --- Response Types ---

// Response is the standard API response envelope.
type Response struct {
	Success bool   `json:"success"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
}

// PaginatedResponse is a response with pagination metadata.
type PaginatedResponse struct {
	Success bool           `json:"success"`
	Data    any            `json:"data"`
	Meta    PaginationMeta `json:"meta"`
}

// PaginationMeta holds pagination information.
type PaginationMeta struct {
	Total   int `json:"total"`
	Page    int `json:"page"`
	PerPage int `json:"perPage"`
	Pages   int `json:"pages"`
}

// Validatable can be implemented by DTOs for automatic validation.
type Validatable interface {
	Validate() error
}
