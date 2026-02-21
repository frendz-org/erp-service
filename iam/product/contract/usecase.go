package contract

import (
	"context"

	"erp-service/entity"

	"github.com/google/uuid"
)

type Usecase interface {
	GetFrendzSaving(ctx context.Context, tenantID uuid.UUID) (*entity.Product, error)
}
