package masterdata

import (
	"context"

	"erp-service/pkg/errors"

	"github.com/google/uuid"
)

func (uc *usecase) GetItemByCode(ctx context.Context, categoryCode string, tenantID *uuid.UUID, itemCode string) (*ItemResponse, error) {
	category, err := uc.categoryRepo.GetByCode(ctx, categoryCode)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, errors.ErrNotFound("category not found")
		}
		return nil, err
	}

	if cached, _ := uc.cache.GetItemByCode(ctx, category.ID, tenantID, itemCode); cached != nil {
		return cached, nil
	}

	item, err := uc.itemRepo.GetByCode(ctx, category.ID, tenantID, itemCode)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, errors.ErrNotFound("item not found")
		}
		return nil, err
	}

	response := MapItemToResponse(item)

	_ = uc.cache.SetItemByCode(ctx, category.ID, tenantID, itemCode, response, uc.config.Masterdata.CacheTTLItems)

	return response, nil
}
