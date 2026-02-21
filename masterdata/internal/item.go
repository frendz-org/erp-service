package internal

import (
	"context"
	"fmt"
	"time"

	"erp-service/entity"
	"erp-service/masterdata/contract"
	"erp-service/masterdata/masterdatadto"
	"erp-service/pkg/errors"

	"github.com/google/uuid"
)

func (uc *usecase) ItemList(ctx context.Context, req *masterdatadto.ListItemsRequest) (*masterdatadto.ListItemsResponse, error) {
	req.Page, req.PerPage = normalizePageParams(req.Page, req.PerPage)

	filterHash := hashFilter(req)
	if cached, _ := uc.cache.GetItemsList(ctx, filterHash); cached != nil {
		return cached, nil
	}

	filter := &contract.ItemFilter{
		CategoryID:   req.CategoryID,
		CategoryCode: req.CategoryCode,
		TenantID:     req.TenantID,
		ParentID:     req.ParentID,
		ParentCode:   req.ParentCode,
		Status:       req.Status,
		Search:       req.Search,
		IsDefault:    req.IsDefault,
		IsSystem:     req.IsSystem,
		Page:         req.Page,
		PerPage:      req.PerPage,
		SortBy:       req.SortBy,
		SortOrder:    req.SortOrder,
	}

	items, total, err := uc.itemRepo.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	response := &masterdatadto.ListItemsResponse{
		Items: masterdatadto.MapItemsToResponse(items),
		Pagination: masterdatadto.Pagination{
			Total:      total,
			Page:       req.Page,
			PerPage:    req.PerPage,
			TotalPages: masterdatadto.CalculateTotalPages(total, req.PerPage),
		},
	}

	_ = uc.cache.SetItemsList(ctx, filterHash, response, uc.config.Masterdata.CacheTTLItems)

	return response, nil
}

func (uc *usecase) ItemGetByID(ctx context.Context, id uuid.UUID) (*masterdatadto.ItemResponse, error) {

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

func (uc *usecase) ItemGetByCode(ctx context.Context, categoryCode string, tenantID *uuid.UUID, itemCode string) (*masterdatadto.ItemResponse, error) {

	category, err := uc.categoryRepo.GetByCode(ctx, categoryCode)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, errors.ErrNotFound("category not found")
		}
		return nil, err
	}

	if cached, _ := uc.cache.GetItemByCode(ctx, category.ID, tenantID, itemCode); cached != nil {
		return cached, nil
	}

	item, err := uc.itemRepo.GetByCode(ctx, category.ID, tenantID, itemCode)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, errors.ErrNotFound("item not found")
		}
		return nil, err
	}

	response := masterdatadto.MapItemToResponse(item)

	_ = uc.cache.SetItemByCode(ctx, category.ID, tenantID, itemCode, response, uc.config.Masterdata.CacheTTLItems)

	return response, nil
}

func (uc *usecase) ItemCreate(ctx context.Context, req *masterdatadto.CreateItemRequest) (*masterdatadto.ItemResponse, error) {

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

func (uc *usecase) ItemUpdate(ctx context.Context, id uuid.UUID, req *masterdatadto.UpdateItemRequest) (*masterdatadto.ItemResponse, error) {

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

func (uc *usecase) ItemDelete(ctx context.Context, id uuid.UUID) error {

	item, err := uc.itemRepo.GetByID(ctx, id)
	if err != nil {
		if errors.IsNotFound(err) {
			return errors.ErrNotFound("item not found")
		}
		return err
	}

	if item.IsSystem {
		return errors.ErrValidation("system items cannot be deleted")
	}

	if err := uc.itemRepo.Delete(ctx, id); err != nil {
		return err
	}

	_ = uc.cache.InvalidateItems(ctx)

	return nil
}

func (uc *usecase) ItemGetChildren(ctx context.Context, parentID uuid.UUID) ([]*masterdatadto.ItemResponse, error) {

	if cached, _ := uc.cache.GetItemChildren(ctx, parentID); cached != nil {
		return cached, nil
	}

	children, err := uc.itemRepo.GetChildren(ctx, parentID)
	if err != nil {
		return nil, err
	}

	response := masterdatadto.MapItemsToResponse(children)

	_ = uc.cache.SetItemChildren(ctx, parentID, response, uc.config.Masterdata.CacheTTLItems)

	return response, nil
}

func (uc *usecase) ItemGetTree(ctx context.Context, categoryCode string, tenantID *uuid.UUID) ([]*masterdatadto.ItemResponse, error) {

	if cached, _ := uc.cache.GetItemTree(ctx, categoryCode, tenantID); cached != nil {
		return cached, nil
	}

	items, err := uc.itemRepo.GetTree(ctx, categoryCode, tenantID)
	if err != nil {
		return nil, err
	}

	response := masterdatadto.MapItemsToResponse(items)

	_ = uc.cache.SetItemTree(ctx, categoryCode, tenantID, response, uc.config.Masterdata.CacheTTLTree)

	return response, nil
}

func (uc *usecase) ItemListByParent(ctx context.Context, categoryCode string, parentCode string, tenantID *uuid.UUID) ([]*masterdatadto.ItemResponse, error) {

	items, err := uc.itemRepo.ListByParent(ctx, categoryCode, parentCode, tenantID)
	if err != nil {
		return nil, err
	}

	return masterdatadto.MapItemsToResponse(items), nil
}

func (uc *usecase) ItemGetDefault(ctx context.Context, categoryCode string, tenantID *uuid.UUID) (*masterdatadto.ItemResponse, error) {

	category, err := uc.categoryRepo.GetByCode(ctx, categoryCode)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, errors.ErrNotFound("category not found")
		}
		return nil, err
	}

	if cached, _ := uc.cache.GetItemDefault(ctx, category.ID, tenantID); cached != nil {
		return cached, nil
	}

	item, err := uc.itemRepo.GetDefaultItem(ctx, category.ID, tenantID)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, errors.ErrNotFound(fmt.Sprintf("no default item found for category %s", categoryCode))
		}
		return nil, err
	}

	response := masterdatadto.MapItemToResponse(item)

	_ = uc.cache.SetItemDefault(ctx, category.ID, tenantID, response, uc.config.Masterdata.CacheTTLItems)

	return response, nil
}

func (uc *usecase) ItemValidateCode(ctx context.Context, req *masterdatadto.ValidateCodeRequest) (*masterdatadto.ValidateCodeResponse, error) {
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
		response.Message = "Item code '" + req.ItemCode + "' not found in category '" + req.CategoryCode + "'"
	}

	return response, nil
}

func (uc *usecase) ItemValidateCodes(ctx context.Context, req *masterdatadto.ValidateCodesRequest) (*masterdatadto.ValidateCodesResponse, error) {
	results := make([]masterdatadto.ValidationResult, len(req.Validations))
	allValid := true

	for i, v := range req.Validations {
		valid, err := uc.itemRepo.ValidateCode(ctx, v.CategoryCode, v.ItemCode, v.TenantID)
		if err != nil {
			return nil, err
		}

		result := masterdatadto.ValidationResult{
			CategoryCode: v.CategoryCode,
			ItemCode:     v.ItemCode,
			Valid:        valid,
		}

		if !valid {
			allValid = false
			result.Message = "Item code not found"
		}

		results[i] = result
	}

	return &masterdatadto.ValidateCodesResponse{
		AllValid: allValid,
		Results:  results,
	}, nil
}
