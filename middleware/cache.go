package middleware

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"sync"
	"time"
)

// CacheEntry holds a cached response.
type CacheEntry struct {
	Body       []byte
	Status     int
	Headers    http.Header
	Expiration time.Time
}

// Cache is an in-memory HTTP response cache.
type Cache struct {
	mu      sync.RWMutex
	entries map[string]*CacheEntry
	ttl     time.Duration
}

// NewCache creates a new cache with the given TTL.
func NewCache(ttl time.Duration) *Cache {
	c := &Cache{
		entries: make(map[string]*CacheEntry),
		ttl:     ttl,
	}
	go c.cleanupLoop()
	return c
}

// Middleware returns a caching middleware for GET requests.
func (c *Cache) Middleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Only cache GET requests.
			if r.Method != http.MethodGet {
				next.ServeHTTP(w, r)
				return
			}

			key := c.cacheKey(r)

			// Check cache.
			c.mu.RLock()
			entry, exists := c.entries[key]
			c.mu.RUnlock()

			if exists && time.Now().Before(entry.Expiration) {
				// Serve from cache.
				for k, v := range entry.Headers {
					w.Header()[k] = v
				}
				w.Header().Set("X-Cache", "HIT")
				w.WriteHeader(entry.Status)
				_, _ = w.Write(entry.Body)
				return
			}

			// Cache miss, execute handler and capture response.
			rec := &cacheRecorder{
				ResponseWriter: w,
				status:         http.StatusOK,
				headers:        make(http.Header),
			}

			next.ServeHTTP(rec, r)

			// Only cache successful responses.
			if rec.status >= 200 && rec.status < 300 {
				c.mu.Lock()
				c.entries[key] = &CacheEntry{
					Body:       rec.body,
					Status:     rec.status,
					Headers:    rec.headers,
					Expiration: time.Now().Add(c.ttl),
				}
				c.mu.Unlock()
			}

			w.Header().Set("X-Cache", "MISS")
		})
	}
}

// Invalidate removes a cached entry.
func (c *Cache) Invalidate(path string) {
	hash := sha256.Sum256([]byte(path))
	key := hex.EncodeToString(hash[:])
	c.mu.Lock()
	delete(c.entries, key)
	c.mu.Unlock()
}

// Clear removes all cached entries.
func (c *Cache) Clear() {
	c.mu.Lock()
	c.entries = make(map[string]*CacheEntry)
	c.mu.Unlock()
}

func (c *Cache) cacheKey(r *http.Request) string {
	hash := sha256.Sum256([]byte(r.URL.String()))
	return hex.EncodeToString(hash[:])
}

func (c *Cache) cleanupLoop() {
	ticker := time.NewTicker(c.ttl)
	defer ticker.Stop()
	for range ticker.C {
		now := time.Now()
		c.mu.Lock()
		for key, entry := range c.entries {
			if now.After(entry.Expiration) {
				delete(c.entries, key)
			}
		}
		c.mu.Unlock()
	}
}

// cacheRecorder captures the response for caching.
type cacheRecorder struct {
	http.ResponseWriter
	status  int
	body    []byte
	headers http.Header
}

func (r *cacheRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func (r *cacheRecorder) Write(b []byte) (int, error) {
	r.body = append(r.body, b...)
	return r.ResponseWriter.Write(b)
}
