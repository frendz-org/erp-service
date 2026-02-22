package masterdata

import (
	"context"

	"github.com/google/uuid"
)

func (uc *usecase) ListItemsByParent(ctx context.Context, categoryCode string, parentCode string, tenantID *uuid.UUID) ([]*ItemResponse, error) {
	items, err := uc.itemRepo.ListByParent(ctx, categoryCode, parentCode, tenantID)
	if err != nil {
		return nil, err
	}

	return MapItemsToResponse(items), nil
}
