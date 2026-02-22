package member

import "context"

func (uc *usecase) ListMembers(ctx context.Context, req *ListRequest) (*ListResponse, error) {
	filter := &MemberListFilter{
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

	members := make([]MemberListItem, 0, len(rows))
	for _, row := range rows {
		members = append(members, mapRowToListItem(row))
	}

	return &ListResponse{
		Members: members,
		Pagination: PaginationMeta{
			Page:    req.Page,
			PerPage: req.PerPage,
			Total:   total,
		},
	}, nil
}
