package masterdata

import (
	"context"

	"erp-service/pkg/errors"

	"github.com/google/uuid"
)

func (uc *usecase) GetCategoryByID(ctx context.Context, id uuid.UUID) (*CategoryResponse, error) {
	if cached, _ := uc.cache.GetCategoryByID(ctx, id); cached != nil {
		return cached, nil
	}

	category, err := uc.categoryRepo.GetByID(ctx, id)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, errors.ErrNotFound("category not found")
		}
		return nil, err
	}

	response := MapCategoryToResponse(category)

	_ = uc.cache.SetCategoryByID(ctx, id, response, uc.config.Masterdata.CacheTTLCategories)

	return response, nil
}
