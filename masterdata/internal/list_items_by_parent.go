package internal

import (
	"context"

	"erp-service/masterdata/masterdatadto"

	"github.com/google/uuid"
)

func (uc *usecase) ListItemsByParent(ctx context.Context, categoryCode string, parentCode string, tenantID *uuid.UUID) ([]*masterdatadto.ItemResponse, error) {
	items, err := uc.itemRepo.ListByParent(ctx, categoryCode, parentCode, tenantID)
	if err != nil {
		return nil, err
	}

	return masterdatadto.MapItemsToResponse(items), nil
}
