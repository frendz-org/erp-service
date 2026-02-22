package masterdata

import (
	"context"
)

func (uc *usecase) ListCategories(ctx context.Context, req *ListCategoriesRequest) (*ListCategoriesResponse, error) {
	req.Page, req.PerPage = normalizePageParams(req.Page, req.PerPage)

	filterHash := hashFilter(req)
	if cached, _ := uc.cache.GetCategoriesList(ctx, filterHash); cached != nil {
		return cached, nil
	}

	filter := &CategoryFilter{
		Status:    req.Status,
		IsSystem:  req.IsSystem,
		ParentID:  req.ParentID,
		Page:      req.Page,
		PerPage:   req.PerPage,
		SortBy:    req.SortBy,
		SortOrder: req.SortOrder,
	}

	categories, total, err := uc.categoryRepo.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	response := &ListCategoriesResponse{
		Categories: MapCategoriesToResponse(categories),
		Pagination: Pagination{
			Total:      total,
			Page:       req.Page,
			PerPage:    req.PerPage,
			TotalPages: CalculateTotalPages(total, req.PerPage),
		},
	}

	_ = uc.cache.SetCategoriesList(ctx, filterHash, response, uc.config.Masterdata.CacheTTLCategories)

	return response, nil
}
