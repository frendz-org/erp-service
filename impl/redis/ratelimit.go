package redis

import (
	"context"
	"erp-service/pkg/errors"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

func (r *Redis) RateLimitAllow(ctx context.Context, key string, limit int64, window time.Duration) (*RateLimitResult, error) {
	fullKey := rateLimitKey(key)

	script := goredis.NewScript(`
		local current = redis.call("INCR", KEYS[1])
		if current == 1 then
			redis.call("PEXPIRE", KEYS[1], ARGV[1])
		end
		local ttl = redis.call("PTTL", KEYS[1])
		return {current, ttl}
	`)

	result, err := script.Run(ctx, r.client, []string{fullKey}, int64(window/time.Millisecond)).Slice()
	if err != nil {
		return nil, errors.ErrInternal("failed to check rate limit").WithError(err)
	}

	current := result[0].(int64)
	ttl := result[1].(int64)

	remaining := limit - current
	if remaining < 0 {
		remaining = 0
	}

	return &RateLimitResult{
		Allowed:   current <= limit,
		Remaining: remaining,
		Total:     limit,
		ResetIn:   time.Duration(ttl) * time.Millisecond,
	}, nil
}

func (r *Redis) RateLimitAllowSlidingWindow(ctx context.Context, key string, limit int64, window time.Duration) (*RateLimitResult, error) {
	fullKey := rateLimitKey(key)
	now := time.Now()
	windowStart := now.Add(-window).UnixMilli()

	script := goredis.NewScript(`
		-- Remove old entries
		redis.call("ZREMRANGEBYSCORE", KEYS[1], "-inf", ARGV[1])

		-- Add current request
		redis.call("ZADD", KEYS[1], ARGV[2], ARGV[2])

		-- Set expiry
		redis.call("PEXPIRE", KEYS[1], ARGV[3])

		-- Count requests in window
		local count = redis.call("ZCARD", KEYS[1])

		return count
	`)

	count, err := script.Run(ctx, r.client, []string{fullKey}, windowStart, now.UnixMilli(), int64(window/time.Millisecond)).Int64()
	if err != nil {
		return nil, errors.ErrInternal("failed to check rate limit").WithError(err)
	}

	remaining := limit - count
	if remaining < 0 {
		remaining = 0
	}

	return &RateLimitResult{
		Allowed:   count <= limit,
		Remaining: remaining,
		Total:     limit,
		ResetIn:   window,
	}, nil
}

func (r *Redis) RateLimitAllowTokenBucket(ctx context.Context, key string, capacity int64, refillRate float64, refillInterval time.Duration) (*RateLimitResult, error) {
	fullKey := rateLimitKey(key)
	now := time.Now().UnixMilli()

	script := goredis.NewScript(`
		local tokens_key = KEYS[1] .. ":tokens"
		local timestamp_key = KEYS[1] .. ":ts"

		local capacity = tonumber(ARGV[1])
		local refill_rate = tonumber(ARGV[2])
		local refill_interval = tonumber(ARGV[3])
		local now = tonumber(ARGV[4])

		-- Get current tokens and last update time
		local tokens = tonumber(redis.call("GET", tokens_key) or capacity)
		local last_update = tonumber(redis.call("GET", timestamp_key) or now)

		-- Calculate tokens to add based on time elapsed
		local elapsed = now - last_update
		local tokens_to_add = math.floor(elapsed / refill_interval) * refill_rate
		tokens = math.min(capacity, tokens + tokens_to_add)

		-- Try to consume a token
		local allowed = 0
		if tokens >= 1 then
			tokens = tokens - 1
			allowed = 1
		end

		-- Update state
		redis.call("SET", tokens_key, tokens)
		redis.call("SET", timestamp_key, now)
		redis.call("PEXPIRE", tokens_key, refill_interval * capacity / refill_rate * 2)
		redis.call("PEXPIRE", timestamp_key, refill_interval * capacity / refill_rate * 2)

		return {allowed, tokens}
	`)

	result, err := script.Run(ctx, r.client, []string{fullKey}, capacity, refillRate, int64(refillInterval/time.Millisecond), now).Slice()
	if err != nil {
		return nil, errors.ErrInternal("failed to check rate limit").WithError(err)
	}

	allowed := result[0].(int64) == 1
	remaining := int64(result[1].(int64))

	return &RateLimitResult{
		Allowed:   allowed,
		Remaining: remaining,
		Total:     capacity,
		ResetIn:   time.Duration(float64(capacity-remaining)/refillRate) * refillInterval,
	}, nil
}

func (r *Redis) RateLimitReset(ctx context.Context, key string) error {
	fullKey := rateLimitKey(key)
	pipe := r.client.Pipeline()
	pipe.Del(ctx, fullKey)
	pipe.Del(ctx, fullKey+":tokens")
	pipe.Del(ctx, fullKey+":ts")
	_, err := pipe.Exec(ctx)
	return err
}

func (r *Redis) RateLimitGetCount(ctx context.Context, key string) (int64, error) {
	fullKey := rateLimitKey(key)
	val, err := r.client.Get(ctx, fullKey).Int64()
	if err != nil {
		if err == goredis.Nil {
			return 0, nil
		}
		return 0, err
	}
	return val, nil
}
