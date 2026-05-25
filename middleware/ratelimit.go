package middleware

import (
	"context"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// RateLimitStore defines the interface for rate limiting backends.
type RateLimitStore interface {
	// Allow checks if the request is allowed under the rate limit.
	// Returns the remaining requests allowed, the time when the limit resets, and whether the request is allowed.
	Allow(ctx context.Context, key string, limit int, window time.Duration) (remaining int, resetAt time.Time, allowed bool)
}

// MemoryStore is an in-memory implementation of RateLimitStore.
type MemoryStore struct {
	mu       sync.Mutex
	visitors map[string]*visitor
}

type visitor struct {
	tokens    int
	lastReset time.Time
}

// NewMemoryStore creates a new MemoryStore.
func NewMemoryStore() *MemoryStore {
	store := &MemoryStore{
		visitors: make(map[string]*visitor),
	}
	go store.cleanupLoop()
	return store
}

func (s *MemoryStore) Allow(ctx context.Context, key string, limit int, window time.Duration) (int, time.Time, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	v, exists := s.visitors[key]
	if !exists {
		v = &visitor{
			tokens:    limit,
			lastReset: now,
		}
		s.visitors[key] = v
	}

	// Reset tokens if window has passed.
	if now.Sub(v.lastReset) > window {
		v.tokens = limit
		v.lastReset = now
	}

	resetAt := v.lastReset.Add(window)

	if v.tokens <= 0 {
		return 0, resetAt, false
	}

	v.tokens--
	return v.tokens, resetAt, true
}

func (s *MemoryStore) cleanupLoop() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		s.mu.Lock()
		now := time.Now()
		for key, v := range s.visitors {
			// If a visitor hasn't been active for 15 minutes, clean up
			if now.Sub(v.lastReset) > 15*time.Minute {
				delete(s.visitors, key)
			}
		}
		s.mu.Unlock()
	}
}

// RateLimiter handles rate limiting using a store.
type RateLimiter struct {
	store   RateLimitStore
	rate    int
	window  time.Duration
	keyFunc func(r *http.Request) string
}

// NewRateLimiter creates a new RateLimiter.
// By default, it uses MemoryStore and limits by IP.
func NewRateLimiter(rate int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		store:  NewMemoryStore(),
		rate:   rate,
		window: window,
		keyFunc: func(r *http.Request) string {
			// Extract IP from RemoteAddr, strip port
			ip := r.RemoteAddr
			if idx := lastIndex(ip, ':'); idx != -1 {
				ip = ip[:idx]
			}
			return ip
		},
	}
}

// WithStore sets a custom RateLimitStore (e.g. RedisStore).
func (rl *RateLimiter) WithStore(store RateLimitStore) *RateLimiter {
	rl.store = store
	return rl
}

// WithKeyFunc sets a custom key generator function (e.g. API key, User ID).
func (rl *RateLimiter) WithKeyFunc(keyFunc func(r *http.Request) string) *RateLimiter {
	rl.keyFunc = keyFunc
	return rl
}

// Middleware returns the rate limiting middleware.
func (rl *RateLimiter) Middleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := rl.keyFunc(r)

			remaining, resetAt, allowed := rl.store.Allow(r.Context(), key, rl.rate, rl.window)

			// Set Rate Limit headers
			w.Header().Set("X-RateLimit-Limit", strconv.Itoa(rl.rate))
			w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(remaining))
			w.Header().Set("X-RateLimit-Reset", strconv.FormatInt(resetAt.Unix(), 10))

			if !allowed {
				w.Header().Set("Content-Type", "application/json")
				// Calculate retry-after seconds
				retryAfter := int(time.Until(resetAt).Seconds())
				if retryAfter < 1 {
					retryAfter = 1
				}
				w.Header().Set("Retry-After", strconv.Itoa(retryAfter))
				w.WriteHeader(http.StatusTooManyRequests)
				_, _ = w.Write([]byte(`{"error":"Too Many Requests","statusCode":429,"message":"rate limit exceeded"}`))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func lastIndex(s string, c byte) int {
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == c {
			return i
		}
	}
	return -1
}
