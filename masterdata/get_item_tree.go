package masterdata

import (
	"context"

	"github.com/google/uuid"
)

func (uc *usecase) GetItemTree(ctx context.Context, categoryCode string, tenantID *uuid.UUID) ([]*ItemResponse, error) {
	if cached, _ := uc.cache.GetItemTree(ctx, categoryCode, tenantID); cached != nil {
		return cached, nil
	}

	items, err := uc.itemRepo.GetTree(ctx, categoryCode, tenantID)
	if err != nil {
		return nil, err
	}

	response := MapItemsToResponse(items)

	_ = uc.cache.SetItemTree(ctx, categoryCode, tenantID, response, uc.config.Masterdata.CacheTTLTree)

	return response, nil
}
