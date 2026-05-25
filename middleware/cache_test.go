package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCache_Memory(t *testing.T) {
	c := NewCache(100 * time.Millisecond)

	count := 0
	handler := c.Middleware()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count++
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("Hello"))
	}))

	// Request 1: Miss
	req := httptest.NewRequest("GET", "/test", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", rec.Code)
	}
	if rec.Body.String() != "Hello" {
		t.Errorf("Expected Hello, got %s", rec.Body.String())
	}
	if rec.Header().Get("X-Cache") != "MISS" {
		t.Errorf("Expected X-Cache: MISS, got %s", rec.Header().Get("X-Cache"))
	}
	if count != 1 {
		t.Errorf("Expected handler execution count to be 1, got %d", count)
	}

	// Request 2: Hit
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", rec.Code)
	}
	if rec.Body.String() != "Hello" {
		t.Errorf("Expected Hello, got %s", rec.Body.String())
	}
	if rec.Header().Get("X-Cache") != "HIT" {
		t.Errorf("Expected X-Cache: HIT, got %s", rec.Header().Get("X-Cache"))
	}
	if count != 1 {
		t.Errorf("Expected handler execution count to be 1, got %d", count)
	}

	// Invalidate key
	c.Invalidate(context.Background(), "/test")

	// Request 3: Miss after invalidation
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Header().Get("X-Cache") != "MISS" {
		t.Errorf("Expected X-Cache: MISS, got %s", rec.Header().Get("X-Cache"))
	}
	if count != 2 {
		t.Errorf("Expected handler execution count to be 2, got %d", count)
	}
}
