package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"erp-service/entity"

	"github.com/google/uuid"
	goredis "github.com/redis/go-redis/v9"
)

const frendzSavingProductKey = "product:frendz-saving:%s"

func (r *Redis) GetFrendzSaving(ctx context.Context, tenantID uuid.UUID) (*entity.Product, error) {
	key := fmt.Sprintf(frendzSavingProductKey, tenantID.String())
	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == goredis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var product entity.Product
	if err := json.Unmarshal(data, &product); err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *Redis) SetFrendzSaving(ctx context.Context, tenantID uuid.UUID, product *entity.Product, ttl time.Duration) error {
	key := fmt.Sprintf(frendzSavingProductKey, tenantID.String())
	data, err := json.Marshal(product)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, key, data, ttl).Err()
}
