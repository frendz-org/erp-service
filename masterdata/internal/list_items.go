package internal

import (
	"context"

	"iam-service/masterdata/contract"
	"iam-service/masterdata/masterdatadto"
)

func (uc *usecase) ListItems(ctx context.Context, req *masterdatadto.ListItemsRequest) (*masterdatadto.ListItemsResponse, error) {
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
