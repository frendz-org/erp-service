package redis

import (
	"context"
	"fmt"
	"strconv"
	"time"

	goredis "github.com/redis/go-redis/v9"

	"github.com/google/uuid"
)

const (
	tokenBlacklistKeyPrefix = "blacklist:token:"
	userBlacklistKeyPrefix  = "blacklist:user:"
)

func (r *Redis) BlacklistToken(ctx context.Context, jti string, ttl time.Duration) error {
	key := tokenBlacklistKeyPrefix + jti
	return r.client.Set(ctx, key, "1", ttl).Err()
}

func (r *Redis) IsTokenBlacklisted(ctx context.Context, jti string) (bool, error) {
	key := tokenBlacklistKeyPrefix + jti
	result, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("check token blacklist: %w", err)
	}
	return result > 0, nil
}

func (r *Redis) BlacklistUser(ctx context.Context, userID uuid.UUID, timestamp time.Time, ttl time.Duration) error {
	key := userBlacklistKeyPrefix + userID.String()
	return r.client.Set(ctx, key, strconv.FormatInt(timestamp.Unix(), 10), ttl).Err()
}

func (r *Redis) GetUserBlacklistTimestamp(ctx context.Context, userID uuid.UUID) (*time.Time, error) {
	key := userBlacklistKeyPrefix + userID.String()
	result, err := r.client.Get(ctx, key).Int64()
	if err != nil {
		if err == goredis.Nil {
			return nil, nil
		}
		return nil, fmt.Errorf("get user blacklist timestamp: %w", err)
	}
	t := time.Unix(result, 0)
	return &t, nil
}
