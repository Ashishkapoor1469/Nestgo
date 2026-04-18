package core

import (
	"context"
	"sync"
)

// RequestContext wraps a standard context.Context with request-scoped values.
type RequestContext struct {
	context.Context
	mu     sync.RWMutex
	values map[string]any
}

// NewRequestContext creates a new request context from a parent context.
func NewRequestContext(parent context.Context) *RequestContext {
	return &RequestContext{
		Context: parent,
		values:  make(map[string]any),
	}
}

// Set stores a value in the request context.
func (rc *RequestContext) Set(key string, value any) {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	rc.values[key] = value
}

// Get retrieves a value from the request context.
func (rc *RequestContext) Get(key string) (any, bool) {
	rc.mu.RLock()
	defer rc.mu.RUnlock()
	v, ok := rc.values[key]
	return v, ok
}

// MustGet retrieves a value from the request context, panics if not found.
func (rc *RequestContext) MustGet(key string) any {
	v, ok := rc.Get(key)
	if !ok {
		panic("nestgo: request context key not found: " + key)
	}
	return v
}

// GetString retrieves a string value from the request context.
func (rc *RequestContext) GetString(key string) string {
	v, ok := rc.Get(key)
	if !ok {
		return ""
	}
	s, _ := v.(string)
	return s
}

// UserID returns the authenticated user ID from the context.
func (rc *RequestContext) UserID() string {
	return rc.GetString("user_id")
}

// TraceID returns the trace ID from the context.
func (rc *RequestContext) TraceID() string {
	return rc.GetString("trace_id")
}

// RequestID returns the request ID from the context.
func (rc *RequestContext) RequestID() string {
	return rc.GetString("request_id")
}

// WithValue returns a new RequestContext with the given key-value pair.
func (rc *RequestContext) WithValue(key string, value any) *RequestContext {
	newCtx := &RequestContext{
		Context: rc.Context,
		values:  make(map[string]any, len(rc.values)+1),
	}
	rc.mu.RLock()
	for k, v := range rc.values {
		newCtx.values[k] = v
	}
	rc.mu.RUnlock()
	newCtx.values[key] = value
	return newCtx
}
