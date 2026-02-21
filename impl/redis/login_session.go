package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"erp-service/entity"
	"erp-service/pkg/errors"

	"github.com/google/uuid"
	goredis "github.com/redis/go-redis/v9"
)

const (
	loginSessionPrefix = "login:%s"
	loginRatePrefix    = "login_rate:%s"
)

func (r *Redis) loginSessionKey(sessionID uuid.UUID) string {
	return fmt.Sprintf(loginSessionPrefix, sessionID.String())
}

func (r *Redis) loginRateLimitKey(email string) string {
	return fmt.Sprintf(loginRatePrefix, strings.ToLower(email))
}

func (r *Redis) CreateLoginSession(ctx context.Context, session *entity.LoginSession, ttl time.Duration) error {
	key := r.loginSessionKey(session.ID)

	data, err := json.Marshal(session)
	if err != nil {
		return errors.ErrInternal("failed to marshal login session").WithError(err)
	}

	if err := r.client.Set(ctx, key, data, ttl).Err(); err != nil {
		return errors.ErrInternal("failed to store login session").WithError(err)
	}

	return nil
}

func (r *Redis) GetLoginSession(ctx context.Context, sessionID uuid.UUID) (*entity.LoginSession, error) {
	key := r.loginSessionKey(sessionID)

	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == goredis.Nil {
			return nil, errors.ErrNotFound("login session not found or expired")
		}
		return nil, errors.ErrInternal("failed to get login session").WithError(err)
	}

	var session entity.LoginSession
	if err := json.Unmarshal(data, &session); err != nil {
		return nil, errors.ErrInternal("failed to unmarshal login session").WithError(err)
	}

	return &session, nil
}

func (r *Redis) updateLoginSessionAtomically(
	ctx context.Context,
	sessionID uuid.UUID,
	mutate func(session *entity.LoginSession) error,
) (*entity.LoginSession, error) {
	key := r.loginSessionKey(sessionID)

	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == goredis.Nil {
			return nil, errors.ErrNotFound("login session not found or expired")
		}
		return nil, errors.ErrInternal("failed to get login session").WithError(err)
	}

	var session entity.LoginSession
	if err := json.Unmarshal(data, &session); err != nil {
		return nil, errors.ErrInternal("failed to unmarshal login session").WithError(err)
	}

	if err := mutate(&session); err != nil {
		return nil, err
	}

	newData, err := json.Marshal(&session)
	if err != nil {
		return nil, errors.ErrInternal("failed to marshal login session").WithError(err)
	}

	result, err := updateSessionScript.Run(ctx, r.client, []string{key}, newData).Text()
	if err != nil {
		if strings.Contains(err.Error(), "NOT_FOUND") {
			return nil, errors.ErrNotFound("login session not found or expired")
		}
		return nil, errors.ErrInternal("failed to update login session").WithError(err)
	}
	_ = result

	return &session, nil
}

func (r *Redis) UpdateLoginSession(ctx context.Context, session *entity.LoginSession, ttl time.Duration) error {
	key := r.loginSessionKey(session.ID)

	data, err := json.Marshal(session)
	if err != nil {
		return errors.ErrInternal("failed to marshal login session").WithError(err)
	}

	script := goredis.NewScript(`
		local exists = redis.call("EXISTS", KEYS[1])
		if exists == 0 then
			return redis.error_reply("NOT_FOUND")
		end
		local ttl_ms = tonumber(ARGV[2])
		if ttl_ms == 0 then
			local pttl = redis.call("PTTL", KEYS[1])
			if pttl > 0 then
				ttl_ms = pttl
			end
		end
		if ttl_ms > 0 then
			redis.call("SET", KEYS[1], ARGV[1], "PX", ttl_ms)
		else
			redis.call("SET", KEYS[1], ARGV[1])
		end
		return "OK"
	`)

	ttlMs := int64(ttl / time.Millisecond)
	_, err = script.Run(ctx, r.client, []string{key}, data, ttlMs).Text()
	if err != nil {
		if strings.Contains(err.Error(), "NOT_FOUND") {
			return errors.ErrNotFound("login session not found or expired")
		}
		return errors.ErrInternal("failed to update login session").WithError(err)
	}

	return nil
}

func (r *Redis) DeleteLoginSession(ctx context.Context, sessionID uuid.UUID) error {
	key := r.loginSessionKey(sessionID)
	return r.client.Del(ctx, key).Err()
}

func (r *Redis) IncrementLoginAttempts(ctx context.Context, sessionID uuid.UUID) (int, error) {
	session, err := r.updateLoginSessionAtomically(ctx, sessionID, func(s *entity.LoginSession) error {
		s.Attempts++
		if s.Attempts >= s.MaxAttempts {
			s.Status = entity.LoginSessionStatusFailed
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return session.Attempts, nil
}

func (r *Redis) UpdateLoginOTP(ctx context.Context, sessionID uuid.UUID, otpHash string, expiresAt time.Time) error {
	_, err := r.updateLoginSessionAtomically(ctx, sessionID, func(s *entity.LoginSession) error {
		now := time.Now()
		s.OTPHash = otpHash
		s.OTPCreatedAt = now
		s.OTPExpiresAt = expiresAt
		s.ResendCount++
		s.LastResentAt = &now
		return nil
	})
	return err
}

func (r *Redis) MarkLoginVerified(ctx context.Context, sessionID uuid.UUID) error {
	_, err := r.updateLoginSessionAtomically(ctx, sessionID, func(s *entity.LoginSession) error {
		now := time.Now()
		s.Status = entity.LoginSessionStatusVerified
		s.VerifiedAt = &now
		return nil
	})
	return err
}

func (r *Redis) IncrementLoginRateLimit(ctx context.Context, email string, ttl time.Duration) (int64, error) {
	key := r.loginRateLimitKey(email)

	count, err := rateLimitScript.Run(ctx, r.client, []string{key}, int64(ttl/time.Millisecond)).Int64()
	if err != nil {
		return 0, errors.ErrInternal("failed to increment login rate limit").WithError(err)
	}

	return count, nil
}

func (r *Redis) GetLoginRateLimitCount(ctx context.Context, email string) (int64, error) {
	key := r.loginRateLimitKey(email)

	count, err := r.client.Get(ctx, key).Int64()
	if err != nil {
		if err == goredis.Nil {
			return 0, nil
		}
		return 0, errors.ErrInternal("failed to get login rate limit count").WithError(err)
	}

	return count, nil
}
