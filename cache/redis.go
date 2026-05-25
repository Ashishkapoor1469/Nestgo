package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Ashishkapoor1469/Nestgo/middleware"
	"github.com/redis/go-redis/v9"
)

// RedisCacheStore is a Redis-backed implementation of CacheStore.
type RedisCacheStore struct {
	client *redis.Client
	prefix string
}

// NewRedisCacheStore creates a new RedisCacheStore.
func NewRedisCacheStore(client *redis.Client, prefix string) *RedisCacheStore {
	if prefix == "" {
		prefix = "cache"
	}
	return &RedisCacheStore{
		client: client,
		prefix: prefix,
	}
}

// Get retrieves a cache entry from Redis.
func (s *RedisCacheStore) Get(ctx context.Context, key string) (*middleware.CacheEntry, bool) {
	fullKey := s.prefix + ":" + key
	data, err := s.client.Get(ctx, fullKey).Bytes()
	if err != nil {
		return nil, false
	}

	var entry middleware.CacheEntry
	if err := json.Unmarshal(data, &entry); err != nil {
		return nil, false
	}

	return &entry, true
}

// Set stores a cache entry in Redis and links it with any tags.
func (s *RedisCacheStore) Set(ctx context.Context, key string, entry *middleware.CacheEntry, ttl time.Duration, tags ...string) {
	fullKey := s.prefix + ":" + key
	data, err := json.Marshal(entry)
	if err != nil {
		return
	}

	pipe := s.client.Pipeline()
	pipe.Set(ctx, fullKey, data, ttl)

	// Keep track of tags
	for _, tag := range tags {
		tagKey := s.prefix + ":tag:" + tag
		pipe.SAdd(ctx, tagKey, key)
		pipe.Expire(ctx, tagKey, ttl*2) // Keep tag mapping alive slightly longer than TTL
	}

	_, _ = pipe.Exec(ctx)
}

// Delete removes a cache entry from Redis.
func (s *RedisCacheStore) Delete(ctx context.Context, key string) {
	fullKey := s.prefix + ":" + key
	_ = s.client.Del(ctx, fullKey).Err()
}

// InvalidateTags deletes all cache entries associated with the specified tags.
func (s *RedisCacheStore) InvalidateTags(ctx context.Context, tags ...string) {
	for _, tag := range tags {
		tagKey := s.prefix + ":tag:" + tag
		keys, err := s.client.SMembers(ctx, tagKey).Result()
		if err != nil || len(keys) == 0 {
			continue
		}

		pipe := s.client.Pipeline()
		for _, key := range keys {
			fullKey := s.prefix + ":" + key
			pipe.Del(ctx, fullKey)
		}
		pipe.Del(ctx, tagKey)
		_, _ = pipe.Exec(ctx)
	}
}

// Clear flushes all cache entries matching the prefix.
func (s *RedisCacheStore) Clear(ctx context.Context) {
	var cursor uint64
	for {
		keys, nextCursor, err := s.client.Scan(ctx, cursor, s.prefix+":*", 100).Result()
		if err != nil {
			break
		}

		if len(keys) > 0 {
			_ = s.client.Del(ctx, keys...).Err()
		}

		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}
}
