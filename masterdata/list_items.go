package masterdata

import (
	"context"
)

func (uc *usecase) ListItems(ctx context.Context, req *ListItemsRequest) (*ListItemsResponse, error) {
	req.Page, req.PerPage = normalizePageParams(req.Page, req.PerPage)

	filterHash := hashFilter(req)
	if cached, _ := uc.cache.GetItemsList(ctx, filterHash); cached != nil {
		return cached, nil
	}

	filter := &ItemFilter{
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

	response := &ListItemsResponse{
		Items: MapItemsToResponse(items),
		Pagination: Pagination{
			Total:      total,
			Page:       req.Page,
			PerPage:    req.PerPage,
			TotalPages: CalculateTotalPages(total, req.PerPage),
		},
	}

	_ = uc.cache.SetItemsList(ctx, filterHash, response, uc.config.Masterdata.CacheTTLItems)

	return response, nil
}
