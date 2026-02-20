package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"iam-service/entity"
	"iam-service/pkg/errors"

	"github.com/google/uuid"
	goredis "github.com/redis/go-redis/v9"
)

const (
	registrationSessionPrefix  = "reg:%s"
	registrationEmailPrefix    = "reg_email:%s"
	registrationRatePrefix     = "reg_rate:%s"
	registrationPasswordPrefix = "reg_pw:%s"
)

func (r *Redis) registrationSessionKey(sessionID uuid.UUID) string {
	return fmt.Sprintf(registrationSessionPrefix, sessionID.String())
}

func (r *Redis) registrationEmailLockKey(email string) string {
	return fmt.Sprintf(registrationEmailPrefix, strings.ToLower(email))
}

func (r *Redis) registrationRateLimitKey(email string) string {
	return fmt.Sprintf(registrationRatePrefix, strings.ToLower(email))
}

func (r *Redis) registrationPasswordKey(sessionID uuid.UUID) string {
	return fmt.Sprintf(registrationPasswordPrefix, sessionID.String())
}

func (r *Redis) CreateRegistrationSession(ctx context.Context, session *entity.RegistrationSession, ttl time.Duration) error {
	key := r.registrationSessionKey(session.ID)

	data, err := json.Marshal(session)
	if err != nil {
		return errors.ErrInternal("failed to marshal registration session").WithError(err)
	}

	if err := r.client.Set(ctx, key, data, ttl).Err(); err != nil {
		return errors.ErrInternal("failed to store registration session").WithError(err)
	}

	return nil
}

func (r *Redis) GetRegistrationSession(ctx context.Context, sessionID uuid.UUID) (*entity.RegistrationSession, error) {
	key := r.registrationSessionKey(sessionID)

	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == goredis.Nil {
			return nil, errors.ErrNotFound("registration session not found or expired")
		}
		return nil, errors.ErrInternal("failed to get registration session").WithError(err)
	}

	var session entity.RegistrationSession
	if err := json.Unmarshal(data, &session); err != nil {
		return nil, errors.ErrInternal("failed to unmarshal registration session").WithError(err)
	}

	return &session, nil
}

var updateSessionScript = goredis.NewScript(`
	local current = redis.call("GET", KEYS[1])
	if not current then
		return {err = "NOT_FOUND"}
	end
	local pttl = redis.call("PTTL", KEYS[1])
	if pttl <= 0 then
		pttl = -1
	end
	redis.call("SET", KEYS[1], ARGV[1], "PX", pttl)
	return current
`)

func (r *Redis) updateSessionAtomically(
	ctx context.Context,
	sessionID uuid.UUID,
	mutate func(session *entity.RegistrationSession) error,
) (*entity.RegistrationSession, error) {
	key := r.registrationSessionKey(sessionID)

	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == goredis.Nil {
			return nil, errors.ErrNotFound("registration session not found or expired")
		}
		return nil, errors.ErrInternal("failed to get registration session").WithError(err)
	}

	var session entity.RegistrationSession
	if err := json.Unmarshal(data, &session); err != nil {
		return nil, errors.ErrInternal("failed to unmarshal registration session").WithError(err)
	}

	if err := mutate(&session); err != nil {
		return nil, err
	}

	newData, err := json.Marshal(&session)
	if err != nil {
		return nil, errors.ErrInternal("failed to marshal registration session").WithError(err)
	}

	result, err := updateSessionScript.Run(ctx, r.client, []string{key}, newData).Text()
	if err != nil {
		if strings.Contains(err.Error(), "NOT_FOUND") {
			return nil, errors.ErrNotFound("registration session not found or expired")
		}
		return nil, errors.ErrInternal("failed to update registration session").WithError(err)
	}
	_ = result

	return &session, nil
}

func (r *Redis) UpdateRegistrationSession(ctx context.Context, session *entity.RegistrationSession, ttl time.Duration) error {
	key := r.registrationSessionKey(session.ID)

	data, err := json.Marshal(session)
	if err != nil {
		return errors.ErrInternal("failed to marshal registration session").WithError(err)
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
			return errors.ErrNotFound("registration session not found or expired")
		}
		return errors.ErrInternal("failed to update registration session").WithError(err)
	}

	return nil
}

func (r *Redis) DeleteRegistrationSession(ctx context.Context, sessionID uuid.UUID) error {
	key := r.registrationSessionKey(sessionID)
	return r.client.Del(ctx, key).Err()
}

func (r *Redis) IncrementRegistrationAttempts(ctx context.Context, sessionID uuid.UUID) (int, error) {
	session, err := r.updateSessionAtomically(ctx, sessionID, func(s *entity.RegistrationSession) error {
		s.Attempts++
		if s.Attempts >= s.MaxAttempts {
			s.Status = entity.RegistrationSessionStatusFailed
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return session.Attempts, nil
}

func (r *Redis) UpdateRegistrationOTP(ctx context.Context, sessionID uuid.UUID, otpHash string, expiresAt time.Time) error {
	_, err := r.updateSessionAtomically(ctx, sessionID, func(s *entity.RegistrationSession) error {
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

func (r *Redis) MarkRegistrationVerified(ctx context.Context, sessionID uuid.UUID, tokenHash string) error {
	_, err := r.updateSessionAtomically(ctx, sessionID, func(s *entity.RegistrationSession) error {
		now := time.Now()
		s.Status = entity.RegistrationSessionStatusVerified
		s.VerifiedAt = &now
		s.RegistrationTokenHash = &tokenHash
		return nil
	})
	return err
}

func (r *Redis) MarkRegistrationPasswordSet(ctx context.Context, sessionID uuid.UUID, passwordHash string, tokenHash string) error {

	pwKey := r.registrationPasswordKey(sessionID)
	sessionKey := r.registrationSessionKey(sessionID)

	pttl, err := r.client.PTTL(ctx, sessionKey).Result()
	if err != nil {
		return errors.ErrInternal("failed to get session TTL").WithError(err)
	}
	if pttl <= 0 {
		return errors.ErrNotFound("registration session not found or expired")
	}

	if err := r.client.Set(ctx, pwKey, passwordHash, pttl).Err(); err != nil {
		return errors.ErrInternal("failed to store password hash").WithError(err)
	}

	_, err = r.updateSessionAtomically(ctx, sessionID, func(s *entity.RegistrationSession) error {
		now := time.Now()
		s.Status = entity.RegistrationSessionStatusPasswordSet
		s.PasswordSetAt = &now
		s.RegistrationTokenHash = &tokenHash
		return nil
	})
	return err
}

func (r *Redis) GetRegistrationPasswordHash(ctx context.Context, sessionID uuid.UUID) (string, error) {
	key := r.registrationPasswordKey(sessionID)

	hash, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == goredis.Nil {
			return "", errors.ErrNotFound("password hash not found or expired")
		}
		return "", errors.ErrInternal("failed to get password hash").WithError(err)
	}

	return hash, nil
}

func (r *Redis) LockRegistrationEmail(ctx context.Context, email string, ttl time.Duration) (bool, error) {
	key := r.registrationEmailLockKey(email)
	return r.client.SetNX(ctx, key, "1", ttl).Result()
}

func (r *Redis) UnlockRegistrationEmail(ctx context.Context, email string) error {
	key := r.registrationEmailLockKey(email)
	return r.client.Del(ctx, key).Err()
}

func (r *Redis) IsRegistrationEmailLocked(ctx context.Context, email string) (bool, error) {
	key := r.registrationEmailLockKey(email)
	exists, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, errors.ErrInternal("failed to check email lock").WithError(err)
	}
	return exists > 0, nil
}

var rateLimitScript = goredis.NewScript(`
	local current = redis.call("INCR", KEYS[1])
	if current == 1 then
		redis.call("PEXPIRE", KEYS[1], ARGV[1])
	end
	return current
`)

func (r *Redis) IncrementRegistrationRateLimit(ctx context.Context, email string, ttl time.Duration) (int64, error) {
	key := r.registrationRateLimitKey(email)

	count, err := rateLimitScript.Run(ctx, r.client, []string{key}, int64(ttl/time.Millisecond)).Int64()
	if err != nil {
		return 0, errors.ErrInternal("failed to increment rate limit").WithError(err)
	}

	return count, nil
}

func (r *Redis) GetRegistrationRateLimitCount(ctx context.Context, email string) (int64, error) {
	key := r.registrationRateLimitKey(email)

	count, err := r.client.Get(ctx, key).Int64()
	if err != nil {
		if err == goredis.Nil {
			return 0, nil
		}
		return 0, errors.ErrInternal("failed to get rate limit count").WithError(err)
	}

	return count, nil
}
