package middleware

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisStore is a Redis-backed implementation of RateLimitStore.
type RedisStore struct {
	client *redis.Client
	prefix string
}

// NewRedisStore creates a new RedisStore.
func NewRedisStore(client *redis.Client, prefix string) *RedisStore {
	if prefix == "" {
		prefix = "ratelimit"
	}
	return &RedisStore{
		client: client,
		prefix: prefix,
	}
}

// Lua script for atomic sliding window rate limiting.
// KEYS[1]: rate limit key
// ARGV[1]: current timestamp in milliseconds
// ARGV[2]: window duration in milliseconds
// ARGV[3]: limit (max requests)
const slidingWindowLua = `
local key = KEYS[1]
local now = tonumber(ARGV[1])
local window = tonumber(ARGV[2])
local limit = tonumber(ARGV[3])

local clearBefore = now - window

-- Remove old requests
redis.call('ZREMRANGEBYSCORE', key, '-inf', clearBefore)

-- Count remaining requests
local count = redis.call('ZCARD', key)
local allowed = 1
local remaining = limit - count

if count >= limit then
    allowed = 0
    remaining = 0
else
    -- Add current request timestamp
    redis.call('ZADD', key, now, now)
    remaining = remaining - 1
end

-- Get oldest request timestamp to calculate reset time
local oldest = redis.call('ZRANGE', key, 0, 0, 'WITHSCORES')
local oldestTime = now
if oldest and oldest[2] then
    oldestTime = tonumber(oldest[2])
end
local resetAt = oldestTime + window

-- Set TTL to prevent leaks
local ttl = math.ceil(window / 1000) * 2
redis.call('EXPIRE', key, ttl)

return {allowed, remaining, resetAt}
`

func (s *RedisStore) Allow(ctx context.Context, key string, limit int, window time.Duration) (int, time.Time, bool) {
	fullKey := s.prefix + ":" + key
	nowMs := time.Now().UnixNano() / int64(time.Millisecond)
	windowMs := window.Milliseconds()

	res, err := s.client.Eval(ctx, slidingWindowLua, []string{fullKey}, nowMs, windowMs, limit).Result()
	if err != nil {
		// Fallback to allowing request in case of Redis errors to not block users
		return limit, time.Now().Add(window), true
	}

	results, ok := res.([]any)
	if !ok || len(results) < 3 {
		return limit, time.Now().Add(window), true
	}

	allowedVal := results[0].(int64)
	remainingVal := results[1].(int64)
	resetAtMs := results[2].(int64)

	allowed := allowedVal == 1
	remaining := int(remainingVal)
	resetAt := time.UnixMilli(resetAtMs)

	return remaining, resetAt, allowed
}
