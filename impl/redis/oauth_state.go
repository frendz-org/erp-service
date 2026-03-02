package redis

import (
	"context"
	"time"
)

const oauthStatePrefix = "oauth:state:"

func (r *Redis) StoreOAuthState(ctx context.Context, state string, ttl time.Duration) error {
	return r.client.Set(ctx, oauthStatePrefix+state, "1", ttl).Err()
}

func (r *Redis) GetAndDeleteOAuthState(ctx context.Context, state string) (bool, error) {
	result, err := r.client.GetDel(ctx, oauthStatePrefix+state).Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			return false, nil
		}
		return false, err
	}
	return result != "", nil
}
