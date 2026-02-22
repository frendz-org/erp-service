package masterdata

import (
	"context"

	"github.com/google/uuid"
)

func (uc *usecase) GetCategoryChildren(ctx context.Context, parentID uuid.UUID) ([]*CategoryResponse, error) {
	if cached, _ := uc.cache.GetCategoryChildren(ctx, parentID); cached != nil {
		return cached, nil
	}

	children, err := uc.categoryRepo.GetChildren(ctx, parentID)
	if err != nil {
		return nil, err
	}

	response := MapCategoriesToResponse(children)

	_ = uc.cache.SetCategoryChildren(ctx, parentID, response, uc.config.Masterdata.CacheTTLCategories)

	return response, nil
}
