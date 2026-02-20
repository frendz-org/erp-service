package internal

import (
	"context"
	"fmt"

	"iam-service/masterdata/masterdatadto"
)

func (uc *usecase) ValidateItemCode(ctx context.Context, req *masterdatadto.ValidateCodeRequest) (*masterdatadto.ValidateCodeResponse, error) {
	valid, err := uc.itemRepo.ValidateCode(ctx, req.CategoryCode, req.ItemCode, req.TenantID)
	if err != nil {
		return nil, err
	}

	response := &masterdatadto.ValidateCodeResponse{
		Valid:        valid,
		CategoryCode: req.CategoryCode,
		ItemCode:     req.ItemCode,
	}

	if !valid {
		response.Message = fmt.Sprintf("Item code '%s' not found in category '%s'", req.ItemCode, req.CategoryCode)
		return response, nil
	}

	if !req.RequireActive && req.ParentItemCode == "" {
		return response, nil
	}

	category, err := uc.categoryRepo.GetByCode(ctx, req.CategoryCode)
	if err != nil {
		return nil, fmt.Errorf("get category: %w", err)
	}

	item, err := uc.itemRepo.GetByCode(ctx, category.ID, req.TenantID, req.ItemCode)
	if err != nil {
		return nil, fmt.Errorf("get item: %w", err)
	}

	if req.RequireActive && !item.IsActive() {
		response.Valid = false
		response.Message = fmt.Sprintf("Item code '%s' is not active", req.ItemCode)
		return response, nil
	}

	if req.ParentItemCode != "" {
		parentItem, err := uc.itemRepo.GetByCode(ctx, category.ID, req.TenantID, req.ParentItemCode)
		if err != nil {
			response.Valid = false
			response.Message = fmt.Sprintf("Parent item code '%s' not found", req.ParentItemCode)
			return response, nil
		}
		if item.ParentItemID == nil || *item.ParentItemID != parentItem.ID {
			response.Valid = false
			response.Message = fmt.Sprintf("Item code '%s' does not belong to parent '%s'", req.ItemCode, req.ParentItemCode)
			return response, nil
		}
	}

	return response, nil
}
