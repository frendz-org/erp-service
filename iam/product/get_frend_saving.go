package product

import (
	"context"
	"time"

	"erp-service/entity"

	"github.com/google/uuid"
)

const (
	frendzSavingCode = "frendz-saving"
	cacheTTL         = 24 * time.Hour
)

func (u *usecase) GetFrendzSaving(ctx context.Context, tenantID uuid.UUID) (*entity.Product, error) {
	if cached, err := u.cache.GetFrendzSaving(ctx, tenantID); err == nil && cached != nil {
		return cached, nil
	}

	product, err := u.repo.GetByCodeAndTenant(ctx, tenantID, frendzSavingCode)
	if err != nil {
		return nil, err
	}

	_ = u.cache.SetFrendzSaving(ctx, tenantID, product, cacheTTL)

	return product, nil
}
