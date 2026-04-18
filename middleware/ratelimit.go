package middleware

import (
	"net/http"
	"sync"
	"time"
)

// RateLimiter implements a token bucket rate limiter per IP.
type RateLimiter struct {
	mu       sync.Mutex
	visitors map[string]*visitor
	rate     int           // requests per window
	window   time.Duration // time window
	cleanup  time.Duration // cleanup interval
}

type visitor struct {
	tokens    int
	lastReset time.Time
}

// NewRateLimiter creates a new rate limiter.
// rate is the number of requests allowed per window duration.
func NewRateLimiter(rate int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*visitor),
		rate:     rate,
		window:   window,
		cleanup:  window * 3,
	}

	go rl.cleanupLoop()
	return rl
}

// Middleware returns the rate limiting middleware.
func (rl *RateLimiter) Middleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := r.RemoteAddr

			rl.mu.Lock()
			v, exists := rl.visitors[ip]
			if !exists {
				v = &visitor{
					tokens:    rl.rate,
					lastReset: time.Now(),
				}
				rl.visitors[ip] = v
			}

			// Reset tokens if window has passed.
			if time.Since(v.lastReset) > rl.window {
				v.tokens = rl.rate
				v.lastReset = time.Now()
			}

			if v.tokens <= 0 {
				rl.mu.Unlock()
				w.Header().Set("Content-Type", "application/json")
				w.Header().Set("Retry-After", rl.window.String())
				w.WriteHeader(http.StatusTooManyRequests)
				_, _ = w.Write([]byte(`{"error":"Too Many Requests","statusCode":429,"message":"rate limit exceeded"}`))
				return
			}

			v.tokens--
			rl.mu.Unlock()

			next.ServeHTTP(w, r)
		})
	}
}

func (rl *RateLimiter) cleanupLoop() {
	ticker := time.NewTicker(rl.cleanup)
	defer ticker.Stop()
	for range ticker.C {
		rl.mu.Lock()
		for ip, v := range rl.visitors {
			if time.Since(v.lastReset) > rl.cleanup {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}
