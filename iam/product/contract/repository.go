package contract

import (
	"context"

	"iam-service/entity"

	"github.com/google/uuid"
)

type ProductRepository interface {
	GetByCodeAndTenant(ctx context.Context, tenantID uuid.UUID, code string) (*entity.Product, error)
}
