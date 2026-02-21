package internal

import (
	"context"
	"time"

	"erp-service/entity"
	"erp-service/iam/product/contract"

	"github.com/google/uuid"
)

const (
	frendzSavingCode = "frendz-saving"
	cacheTTL         = 24 * time.Hour
)

type usecase struct {
	repo  contract.ProductRepository
	cache contract.Cache
}

func NewUsecase(repo contract.ProductRepository, cache contract.Cache) contract.Usecase {
	return &usecase{
		repo:  repo,
		cache: cache,
	}
}

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
