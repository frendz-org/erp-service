package internal

import (
	"context"
	"fmt"

	"erp-service/entity"
	"erp-service/masterdata/masterdatadto"
	"erp-service/pkg/errors"

	"github.com/google/uuid"
)

func (uc *usecase) UpdateCategory(ctx context.Context, id uuid.UUID, req *masterdatadto.UpdateCategoryRequest) (*masterdatadto.CategoryResponse, error) {
	category, err := uc.categoryRepo.GetByID(ctx, id)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, errors.ErrNotFound("category not found")
		}
		return nil, err
	}

	if category.Version != req.Version {
		return nil, errors.ErrConflict(fmt.Sprintf("category has been modified by another user (current version: %d)", category.Version))
	}

	if category.IsSystem && req.Status != nil && *req.Status == string(entity.MasterdataCategoryStatusInactive) {
		return nil, errors.ErrValidation("system categories cannot be deactivated")
	}

	if req.Name != nil {
		category.Name = *req.Name
	}
	if req.Description != nil {
		category.Description = req.Description
	}
	if req.ParentCategoryID != nil {

		if *req.ParentCategoryID == id {
			return nil, errors.ErrValidation("category cannot be its own parent")
		}
		_, err := uc.categoryRepo.GetByID(ctx, *req.ParentCategoryID)
		if err != nil {
			if errors.IsNotFound(err) {
				return nil, errors.ErrValidation("parent category not found")
			}
			return nil, err
		}
		category.ParentCategoryID = req.ParentCategoryID
	}
	if req.IsTenantExtensible != nil {
		category.IsTenantExtensible = *req.IsTenantExtensible
	}
	if req.SortOrder != nil {
		category.SortOrder = *req.SortOrder
	}
	if req.Status != nil {
		category.Status = entity.MasterdataCategoryStatus(*req.Status)
	}
	if req.Metadata != nil {
		category.Metadata = req.Metadata
	}

	category.IncrementVersion()

	if err := category.Validate(); err != nil {
		return nil, errors.ErrValidation(err.Error())
	}

	if err := uc.categoryRepo.Update(ctx, category); err != nil {
		return nil, err
	}

	_ = uc.cache.InvalidateCategories(ctx)

	return masterdatadto.MapCategoryToResponse(category), nil
}
