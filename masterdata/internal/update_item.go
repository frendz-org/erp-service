package internal

import (
	"context"
	"fmt"
	"time"

	"iam-service/entity"
	"iam-service/masterdata/masterdatadto"
	"iam-service/pkg/errors"

	"github.com/google/uuid"
)

func (uc *usecase) UpdateItem(ctx context.Context, id uuid.UUID, req *masterdatadto.UpdateItemRequest) (*masterdatadto.ItemResponse, error) {
	item, err := uc.itemRepo.GetByID(ctx, id)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, errors.ErrNotFound("item not found")
		}
		return nil, err
	}

	if item.Version != req.Version {
		return nil, errors.ErrConflict(fmt.Sprintf("item has been modified by another user (current version: %d)", item.Version))
	}

	if item.IsSystem && req.Status != nil && *req.Status == string(entity.MasterdataItemStatusInactive) {
		return nil, errors.ErrValidation("system items cannot be deactivated")
	}

	if req.ParentItemID != nil {
		if *req.ParentItemID == id {
			return nil, errors.ErrValidation("item cannot be its own parent")
		}
		_, err := uc.itemRepo.GetByID(ctx, *req.ParentItemID)
		if err != nil {
			if errors.IsNotFound(err) {
				return nil, errors.ErrValidation("parent item not found")
			}
			return nil, err
		}
		item.ParentItemID = req.ParentItemID
	}
	if req.Name != nil {
		item.Name = *req.Name
	}
	if req.AltName != nil {
		item.AltName = req.AltName
	}
	if req.Description != nil {
		item.Description = req.Description
	}
	if req.SortOrder != nil {
		item.SortOrder = *req.SortOrder
	}
	if req.IsDefault != nil {
		item.IsDefault = *req.IsDefault
	}
	if req.Status != nil {
		item.Status = entity.MasterdataItemStatus(*req.Status)
	}
	if req.EffectiveFrom != nil {
		t, err := time.Parse(time.RFC3339, *req.EffectiveFrom)
		if err != nil {
			return nil, errors.ErrValidation("invalid effective_from date format (use ISO 8601)")
		}
		item.EffectiveFrom = &t
	}
	if req.EffectiveUntil != nil {
		t, err := time.Parse(time.RFC3339, *req.EffectiveUntil)
		if err != nil {
			return nil, errors.ErrValidation("invalid effective_until date format (use ISO 8601)")
		}
		item.EffectiveUntil = &t
	}
	if req.Metadata != nil {
		item.Metadata = req.Metadata
	}

	item.IncrementVersion()

	if err := item.Validate(); err != nil {
		return nil, errors.ErrValidation(err.Error())
	}

	if err := uc.itemRepo.Update(ctx, item); err != nil {
		return nil, err
	}

	_ = uc.cache.InvalidateItems(ctx)

	return masterdatadto.MapItemToResponse(item), nil
}
