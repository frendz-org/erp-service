package internal

import (
	"context"

	"iam-service/masterdata/contract"
	"iam-service/masterdata/masterdatadto"
)

func (uc *usecase) ListCategories(ctx context.Context, req *masterdatadto.ListCategoriesRequest) (*masterdatadto.ListCategoriesResponse, error) {
	req.Page, req.PerPage = normalizePageParams(req.Page, req.PerPage)

	filterHash := hashFilter(req)
	if cached, _ := uc.cache.GetCategoriesList(ctx, filterHash); cached != nil {
		return cached, nil
	}

	filter := &contract.CategoryFilter{
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

	response := &masterdatadto.ListCategoriesResponse{
		Categories: masterdatadto.MapCategoriesToResponse(categories),
		Pagination: masterdatadto.Pagination{
			Total:      total,
			Page:       req.Page,
			PerPage:    req.PerPage,
			TotalPages: masterdatadto.CalculateTotalPages(total, req.PerPage),
		},
	}

	_ = uc.cache.SetCategoriesList(ctx, filterHash, response, uc.config.Masterdata.CacheTTLCategories)

	return response, nil
}
