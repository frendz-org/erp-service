package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

const (
	transferTokenPrefix          = "transfer:"
	transferTokenRateLimitPrefix = "transfer:ratelimit:"
)

func (r *Redis) StoreTransferToken(ctx context.Context, code string, data []byte, ttl time.Duration) error {
	return r.client.Set(ctx, transferTokenPrefix+code, data, ttl).Err()
}

func (r *Redis) GetAndDeleteTransferToken(ctx context.Context, code string) ([]byte, error) {
	result, err := r.client.GetDel(ctx, transferTokenPrefix+code).Bytes()
	if err != nil {
		if err.Error() == "redis: nil" {
			return nil, nil
		}
		return nil, err
	}
	return result, nil
}

func (r *Redis) IncrementTransferTokenRateLimit(ctx context.Context, userID uuid.UUID, window time.Duration) (int64, error) {
	key := fmt.Sprintf("%s%s", transferTokenRateLimitPrefix, userID.String())
	count, err := r.client.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	if count == 1 {
		r.client.Expire(ctx, key, window)
	}
	return count, nil
}
