package masterdata

import (
	"context"
	"fmt"

	"erp-service/pkg/errors"

	"github.com/google/uuid"
)

func (uc *usecase) GetItemDefault(ctx context.Context, categoryCode string, tenantID *uuid.UUID) (*ItemResponse, error) {
	category, err := uc.categoryRepo.GetByCode(ctx, categoryCode)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, errors.ErrNotFound("category not found")
		}
		return nil, err
	}

	if cached, _ := uc.cache.GetItemDefault(ctx, category.ID, tenantID); cached != nil {
		return cached, nil
	}

	item, err := uc.itemRepo.GetDefaultItem(ctx, category.ID, tenantID)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, errors.ErrNotFound(fmt.Sprintf("no default item found for category %s", categoryCode))
		}
		return nil, err
	}

	response := MapItemToResponse(item)

	_ = uc.cache.SetItemDefault(ctx, category.ID, tenantID, response, uc.config.Masterdata.CacheTTLItems)

	return response, nil
}
