package internal

import (
	"context"

	"erp-service/saving/member/contract"
	"erp-service/saving/member/memberdto"
)

func (uc *usecase) ListMembers(ctx context.Context, req *memberdto.ListRequest) (*memberdto.ListResponse, error) {
	filter := &contract.MemberListFilter{
		TenantID:  req.TenantID,
		ProductID: req.ProductID,
		Status:    req.Status,
		Search:    req.Search,
		Page:      req.Page,
		PerPage:   req.PerPage,
		SortBy:    req.SortBy,
		SortOrder: req.SortOrder,
	}

	rows, total, err := uc.utrRepo.ListByProductWithFilters(ctx, filter)
	if err != nil {
		return nil, err
	}

	members := make([]memberdto.MemberListItem, 0, len(rows))
	for _, row := range rows {
		members = append(members, mapRowToListItem(row))
	}

	return &memberdto.ListResponse{
		Members: members,
		Pagination: memberdto.PaginationMeta{
			Page:    req.Page,
			PerPage: req.PerPage,
			Total:   total,
		},
	}, nil
}
