package internal

import (
	"context"
	"fmt"
	"time"

	"erp-service/entity"
	"erp-service/masterdata/masterdatadto"
	"erp-service/pkg/errors"
)

func (uc *usecase) CreateItem(ctx context.Context, req *masterdatadto.CreateItemRequest) (*masterdatadto.ItemResponse, error) {
	category, err := uc.categoryRepo.GetByID(ctx, req.CategoryID)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, errors.ErrValidation("category not found")
		}
		return nil, err
	}

	if req.TenantID != nil && !category.IsTenantExtensible {
		return nil, errors.ErrValidation("category does not allow tenant-specific items")
	}

	exists, err := uc.itemRepo.ExistsByCode(ctx, req.CategoryID, req.TenantID, req.Code)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.ErrConflict(fmt.Sprintf("item with code '%s' already exists in this category", req.Code))
	}

	if req.ParentItemID != nil {
		_, err := uc.itemRepo.GetByID(ctx, *req.ParentItemID)
		if err != nil {
			if errors.IsNotFound(err) {
				return nil, errors.ErrValidation("parent item not found")
			}
			return nil, err
		}
	}

	var effectiveFrom, effectiveUntil *time.Time
	if req.EffectiveFrom != nil {
		t, err := time.Parse(time.RFC3339, *req.EffectiveFrom)
		if err != nil {
			return nil, errors.ErrValidation("invalid effective_from date format (use ISO 8601)")
		}
		effectiveFrom = &t
	}
	if req.EffectiveUntil != nil {
		t, err := time.Parse(time.RFC3339, *req.EffectiveUntil)
		if err != nil {
			return nil, errors.ErrValidation("invalid effective_until date format (use ISO 8601)")
		}
		effectiveUntil = &t
	}

	item := &entity.MasterdataItem{
		CategoryID:     req.CategoryID,
		TenantID:       req.TenantID,
		ParentItemID:   req.ParentItemID,
		Code:           req.Code,
		Name:           req.Name,
		AltName:        req.AltName,
		Description:    req.Description,
		SortOrder:      req.SortOrder,
		IsSystem:       req.IsSystem,
		IsDefault:      req.IsDefault,
		Status:         entity.MasterdataItemStatusActive,
		EffectiveFrom:  effectiveFrom,
		EffectiveUntil: effectiveUntil,
		Metadata:       req.Metadata,
		CreatedBy:      req.CreatedBy,
		Version:        1,
	}

	if err := item.Validate(); err != nil {
		return nil, errors.ErrValidation(err.Error())
	}

	if err := uc.itemRepo.Create(ctx, item); err != nil {
		return nil, err
	}

	_ = uc.cache.InvalidateItems(ctx)

	return masterdatadto.MapItemToResponse(item), nil
}
