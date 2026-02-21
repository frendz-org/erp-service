package internal

import (
	"context"
	"fmt"

	"erp-service/entity"
	"erp-service/masterdata/masterdatadto"
	"erp-service/pkg/errors"
)

func (uc *usecase) CreateCategory(ctx context.Context, req *masterdatadto.CreateCategoryRequest) (*masterdatadto.CategoryResponse, error) {
	exists, err := uc.categoryRepo.ExistsByCode(ctx, req.Code)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.ErrConflict(fmt.Sprintf("category with code '%s' already exists", req.Code))
	}

	if req.ParentCategoryID != nil {
		_, err := uc.categoryRepo.GetByID(ctx, *req.ParentCategoryID)
		if err != nil {
			if errors.IsNotFound(err) {
				return nil, errors.ErrValidation("parent category not found")
			}
			return nil, err
		}
	}

	category := &entity.MasterdataCategory{
		Code:               req.Code,
		Name:               req.Name,
		Description:        req.Description,
		ParentCategoryID:   req.ParentCategoryID,
		IsSystem:           req.IsSystem,
		IsTenantExtensible: req.IsTenantExtensible,
		SortOrder:          req.SortOrder,
		Status:             entity.MasterdataCategoryStatusActive,
		Metadata:           req.Metadata,
		Version:            1,
	}

	if err := category.Validate(); err != nil {
		return nil, errors.ErrValidation(err.Error())
	}

	if err := uc.categoryRepo.Create(ctx, category); err != nil {
		return nil, err
	}

	_ = uc.cache.InvalidateCategories(ctx)

	return masterdatadto.MapCategoryToResponse(category), nil
}
