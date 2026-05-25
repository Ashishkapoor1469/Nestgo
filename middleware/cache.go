package middleware

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"sync"
	"time"
)

// CacheEntry holds a cached response.
type CacheEntry struct {
	Body       []byte      `json:"body"`
	Status     int         `json:"status"`
	Headers    http.Header `json:"headers"`
	Expiration time.Time   `json:"expiration"`
}

// CacheStore defines the interface for HTTP response cache backends.
type CacheStore interface {
	Get(ctx context.Context, key string) (*CacheEntry, bool)
	Set(ctx context.Context, key string, entry *CacheEntry, ttl time.Duration, tags ...string)
	Delete(ctx context.Context, key string)
	InvalidateTags(ctx context.Context, tags ...string)
	Clear(ctx context.Context)
}

// MemoryCacheStore is an in-memory implementation of CacheStore.
type MemoryCacheStore struct {
	mu      sync.RWMutex
	entries map[string]*CacheEntry
	tagMap  map[string]map[string]bool // tag -> set of keys
}

// NewMemoryCacheStore creates a new MemoryCacheStore.
func NewMemoryCacheStore() *MemoryCacheStore {
	store := &MemoryCacheStore{
		entries: make(map[string]*CacheEntry),
		tagMap:  make(map[string]map[string]bool),
	}
	go store.cleanupLoop()
	return store
}

func (s *MemoryCacheStore) Get(ctx context.Context, key string) (*CacheEntry, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	entry, exists := s.entries[key]
	if !exists {
		return nil, false
	}
	if time.Now().After(entry.Expiration) {
		return nil, false
	}
	return entry, true
}

func (s *MemoryCacheStore) Set(ctx context.Context, key string, entry *CacheEntry, ttl time.Duration, tags ...string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries[key] = entry

	for _, tag := range tags {
		if _, exists := s.tagMap[tag]; !exists {
			s.tagMap[tag] = make(map[string]bool)
		}
		s.tagMap[tag][key] = true
	}
}

func (s *MemoryCacheStore) Delete(ctx context.Context, key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.entries, key)
}

func (s *MemoryCacheStore) InvalidateTags(ctx context.Context, tags ...string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, tag := range tags {
		if keys, exists := s.tagMap[tag]; exists {
			for key := range keys {
				delete(s.entries, key)
			}
			delete(s.tagMap, tag)
		}
	}
}

func (s *MemoryCacheStore) Clear(ctx context.Context) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries = make(map[string]*CacheEntry)
	s.tagMap = make(map[string]map[string]bool)
}

func (s *MemoryCacheStore) cleanupLoop() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		now := time.Now()
		s.mu.Lock()
		for key, entry := range s.entries {
			if now.After(entry.Expiration) {
				delete(s.entries, key)
			}
		}
		s.mu.Unlock()
	}
}

// Cache is a middleware handler that wraps an HTTP request with caching.
type Cache struct {
	store CacheStore
	ttl   time.Duration
	tags  []string
}

// NewCache creates a new cache middleware with default MemoryCacheStore.
func NewCache(ttl time.Duration) *Cache {
	return &Cache{
		store: NewMemoryCacheStore(),
		ttl:   ttl,
	}
}

// WithStore configures a custom CacheStore.
func (c *Cache) WithStore(store CacheStore) *Cache {
	c.store = store
	return c
}

// WithTags adds static invalidation tags to the cached responses.
func (c *Cache) WithTags(tags ...string) *Cache {
	c.tags = tags
	return c
}

// Middleware returns the caching middleware.
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
			entry, exists := c.store.Get(r.Context(), key)

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
			}

			next.ServeHTTP(rec, r)

			// Only cache successful responses.
			if rec.status >= 200 && rec.status < 300 {
				headers := make(http.Header)
				for k, v := range w.Header() {
					headers[k] = v
				}

				entry = &CacheEntry{
					Body:       rec.body,
					Status:     rec.status,
					Headers:    headers,
					Expiration: time.Now().Add(c.ttl),
				}
				c.store.Set(r.Context(), key, entry, c.ttl, c.tags...)
			}

			w.Header().Set("X-Cache", "MISS")
		})
	}
}

// Invalidate invalidates a cache key based on the URL path.
func (c *Cache) Invalidate(ctx context.Context, path string) {
	hash := sha256.Sum256([]byte(path))
	key := hex.EncodeToString(hash[:])
	c.store.Delete(ctx, key)
}

// InvalidateTags invalidates all cache entries matching any of the given tags.
func (c *Cache) InvalidateTags(ctx context.Context, tags ...string) {
	c.store.InvalidateTags(ctx, tags...)
}

// Clear clears all cache entries.
func (c *Cache) Clear(ctx context.Context) {
	c.store.Clear(ctx)
}

func (c *Cache) cacheKey(r *http.Request) string {
	hash := sha256.Sum256([]byte(r.URL.String()))
	return hex.EncodeToString(hash[:])
}

// cacheRecorder captures the response for caching.
type cacheRecorder struct {
	http.ResponseWriter
	status int
	body   []byte
}

func (r *cacheRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func (r *cacheRecorder) Write(b []byte) (int, error) {
	r.body = append(r.body, b...)
	return r.ResponseWriter.Write(b)
}
