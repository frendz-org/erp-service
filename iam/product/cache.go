package product

import (
	"context"
	"time"

	"erp-service/entity"

	"github.com/google/uuid"
)

type Cache interface {
	GetFrendzSaving(ctx context.Context, tenantID uuid.UUID) (*entity.Product, error)
	SetFrendzSaving(ctx context.Context, tenantID uuid.UUID, product *entity.Product, ttl time.Duration) error
}
