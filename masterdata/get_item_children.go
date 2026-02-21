package masterdata

import (
	"context"

	"github.com/google/uuid"
)

func (uc *usecase) GetItemChildren(ctx context.Context, parentID uuid.UUID) ([]*ItemResponse, error) {
	if cached, _ := uc.cache.GetItemChildren(ctx, parentID); cached != nil {
		return cached, nil
	}

	children, err := uc.itemRepo.GetChildren(ctx, parentID)
	if err != nil {
		return nil, err
	}

	response := MapItemsToResponse(children)

	_ = uc.cache.SetItemChildren(ctx, parentID, response, uc.config.Masterdata.CacheTTLItems)

	return response, nil
}
