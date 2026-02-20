package internal

import (
	"context"

	"iam-service/masterdata/masterdatadto"
	"iam-service/pkg/errors"
)

func (uc *usecase) GetCategoryByCode(ctx context.Context, code string) (*masterdatadto.CategoryResponse, error) {
	if cached, _ := uc.cache.GetCategoryByCode(ctx, code); cached != nil {
		return cached, nil
	}

	category, err := uc.categoryRepo.GetByCode(ctx, code)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, errors.ErrNotFound("category not found")
		}
		return nil, err
	}

	response := masterdatadto.MapCategoryToResponse(category)

	_ = uc.cache.SetCategoryByCode(ctx, code, response, uc.config.Masterdata.CacheTTLCategories)

	return response, nil
}
