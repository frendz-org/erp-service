package redis

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"erp-service/pkg/errors"
	"fmt"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

type Lock struct {
	redis    *Redis
	name     string
	token    string
	expiry   time.Duration
	acquired bool
}

func generateToken() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func (r *Redis) AcquireLock(ctx context.Context, name string, expiry time.Duration) (*Lock, error) {
	token, err := generateToken()
	if err != nil {
		return nil, errors.ErrInternal("failed to generate lock token").WithError(err)
	}

	lock := &Lock{
		redis:  r,
		name:   name,
		token:  token,
		expiry: expiry,
	}

	lockKey := fmt.Sprintf(LockPrefix, name)
	acquired, err := r.client.SetNX(ctx, lockKey, token, expiry).Result()
	if err != nil {
		return nil, errors.ErrInternal("failed to acquire lock").WithError(err)
	}

	if !acquired {
		return nil, errors.SentinelLockNotAcquired
	}

	lock.acquired = true
	return lock, nil
}

func (r *Redis) AcquireLockWithRetry(ctx context.Context, name string, expiry time.Duration, maxRetries int, retryDelay time.Duration) (*Lock, error) {
	ticker := time.NewTicker(retryDelay)
	defer ticker.Stop()

	for i := 0; i < maxRetries; i++ {
		lock, err := r.AcquireLock(ctx, name, expiry)
		if err == nil {
			return lock, nil
		}
		if err != errors.SentinelLockNotAcquired {
			return nil, err
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-ticker.C:
			continue
		}
	}
	return nil, errors.SentinelLockNotAcquired
}

func (r *Redis) AcquireLockWithWait(ctx context.Context, name string, expiry time.Duration, pollInterval time.Duration) (*Lock, error) {
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	for {
		lock, err := r.AcquireLock(ctx, name, expiry)
		if err == nil {
			return lock, nil
		}
		if err != errors.SentinelLockNotAcquired {
			return nil, err
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-ticker.C:
			continue
		}
	}
}

func (lock *Lock) Release(ctx context.Context) error {
	if !lock.acquired {
		return errors.SentinelLockNotHeld
	}

	script := goredis.NewScript(`
		if redis.call("get", KEYS[1]) == ARGV[1] then
			return redis.call("del", KEYS[1])
		else
			return 0
		end
	`)

	lockKey := fmt.Sprintf(LockPrefix, lock.name)
	result, err := script.Run(ctx, lock.redis.client, []string{lockKey}, lock.token).Int64()
	if err != nil {
		return errors.ErrInternal("failed to release lock").WithError(err)
	}

	if result == 0 {
		return errors.SentinelLockNotHeld
	}

	lock.acquired = false
	return nil
}

func (lock *Lock) Extend(ctx context.Context, expiry time.Duration) error {
	if !lock.acquired {
		return errors.SentinelLockNotHeld
	}

	script := goredis.NewScript(`
		if redis.call("get", KEYS[1]) == ARGV[1] then
			return redis.call("pexpire", KEYS[1], ARGV[2])
		else
			return 0
		end
	`)

	lockKey := fmt.Sprintf(LockPrefix, lock.name)
	result, err := script.Run(ctx, lock.redis.client, []string{lockKey}, lock.token, int64(expiry/time.Millisecond)).Int64()
	if err != nil {
		return errors.ErrInternal("failed to extend lock").WithError(err)
	}

	if result == 0 {
		lock.acquired = false
		return errors.SentinelLockNotHeld
	}

	lock.expiry = expiry
	return nil
}

func (lock *Lock) IsHeld(ctx context.Context) (bool, error) {
	lockKey := fmt.Sprintf(LockPrefix, lock.name)
	val, err := lock.redis.client.Get(ctx, lockKey).Result()
	if err != nil {
		if err == goredis.Nil {
			lock.acquired = false
			return false, nil
		}
		return false, errors.ErrInternal("failed to check lock").WithError(err)
	}
	held := val == lock.token
	if !held {
		lock.acquired = false
	}
	return held, nil
}

func (lock *Lock) TTL(ctx context.Context) (time.Duration, error) {
	lockKey := fmt.Sprintf(LockPrefix, lock.name)
	ttl, err := lock.redis.client.TTL(ctx, lockKey).Result()
	if err != nil {
		return 0, errors.ErrInternal("failed to get lock TTL").WithError(err)
	}
	return ttl, nil
}

func (r *Redis) WithLock(ctx context.Context, name string, expiry time.Duration, fn func(ctx context.Context) error) error {
	lock, err := r.AcquireLock(ctx, name, expiry)
	if err != nil {
		return err
	}
	defer lock.Release(ctx)

	return fn(ctx)
}

func (r *Redis) WithLockRetry(ctx context.Context, name string, expiry time.Duration, maxRetries int, retryDelay time.Duration, fn func(ctx context.Context) error) error {
	lock, err := r.AcquireLockWithRetry(ctx, name, expiry, maxRetries, retryDelay)
	if err != nil {
		return err
	}
	defer lock.Release(ctx)

	return fn(ctx)
}

func (r *Redis) IsLocked(ctx context.Context, name string) (bool, error) {
	lockKey := fmt.Sprintf(LockPrefix, name)
	exists, err := r.client.Exists(ctx, lockKey).Result()
	if err != nil {
		return false, errors.ErrInternal("failed to check lock").WithError(err)
	}
	return exists > 0, nil
}

func (r *Redis) ForceReleaseLock(ctx context.Context, name string) error {
	lockKey := fmt.Sprintf(LockPrefix, name)
	return r.client.Del(ctx, lockKey).Err()
}

func (r *Redis) NewSemaphore(name string, maxCount int64) *Semaphore {
	return &Semaphore{
		redis:    r,
		name:     name,
		maxCount: maxCount,
	}
}

func (s *Semaphore) Acquire(ctx context.Context, expiry time.Duration) (string, error) {
	token, err := generateToken()
	if err != nil {
		return "", errors.ErrInternal("failed to generate token").WithError(err)
	}

	script := goredis.NewScript(`
		local current = redis.call("ZCARD", KEYS[1])
		if current < tonumber(ARGV[1]) then
			redis.call("ZADD", KEYS[1], ARGV[2], ARGV[3])
			return 1
		else
			return 0
		end
	`)

	semKey := fmt.Sprintf(SemaphorePrefix, s.name)
	expireAt := float64(time.Now().Add(expiry).Unix())
	result, err := script.Run(ctx, s.redis.client, []string{semKey}, s.maxCount, expireAt, token).Int64()
	if err != nil {
		return "", errors.ErrInternal("failed to acquire semaphore").WithError(err)
	}

	if result == 0 {
		return "", errors.SentinelSemaphoreFull
	}

	return token, nil
}

func (s *Semaphore) Release(ctx context.Context, token string) error {
	semKey := fmt.Sprintf(SemaphorePrefix, s.name)
	removed, err := s.redis.client.ZRem(ctx, semKey, token).Result()
	if err != nil {
		return errors.ErrInternal("failed to release semaphore").WithError(err)
	}
	if removed == 0 {
		return errors.SentinelSemaphoreTokenNotFound
	}
	return nil
}

func (s *Semaphore) Cleanup(ctx context.Context) error {
	semKey := fmt.Sprintf(SemaphorePrefix, s.name)
	now := float64(time.Now().Unix())
	return s.redis.client.ZRemRangeByScore(ctx, semKey, "-inf", fmt.Sprintf("%f", now)).Err()
}

func (s *Semaphore) Count(ctx context.Context) (int64, error) {
	semKey := fmt.Sprintf(SemaphorePrefix, s.name)
	return s.redis.client.ZCard(ctx, semKey).Result()
}

func (s *Semaphore) Available(ctx context.Context) (int64, error) {
	count, err := s.Count(ctx)
	if err != nil {
		return 0, err
	}
	return s.maxCount - count, nil
}
