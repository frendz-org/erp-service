package masterdata

import (
	"context"

	"erp-service/pkg/errors"
)

func (uc *usecase) GetItemChildrenByCode(ctx context.Context, itemCode string) ([]*ItemResponse, error) {
	item, err := uc.itemRepo.GetByCodeOnly(ctx, itemCode)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, errors.ErrNotFound("item not found")
		}
		return nil, err
	}

	return uc.GetItemChildren(ctx, item.ID)
}
