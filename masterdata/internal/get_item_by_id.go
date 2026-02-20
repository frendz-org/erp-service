package internal

import (
	"context"

	"iam-service/masterdata/masterdatadto"
	"iam-service/pkg/errors"

	"github.com/google/uuid"
)

func (uc *usecase) GetItemByID(ctx context.Context, id uuid.UUID) (*masterdatadto.ItemResponse, error) {
	if cached, _ := uc.cache.GetItemByID(ctx, id); cached != nil {
		return cached, nil
	}

	item, err := uc.itemRepo.GetByID(ctx, id)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, errors.ErrNotFound("item not found")
		}
		return nil, err
	}

	response := masterdatadto.MapItemToResponse(item)

	_ = uc.cache.SetItemByID(ctx, id, response, uc.config.Masterdata.CacheTTLItems)

	return response, nil
}
