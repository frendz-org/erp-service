package user

import (
	"context"

	"erp-service/entity"
	"erp-service/pkg/errors"

	"github.com/google/uuid"
)

func (uc *usecase) List(ctx context.Context, tenantID *uuid.UUID, req *ListRequest) (*ListResponse, error) {
	req.SetDefaults()

	filter := &UserListFilter{
		RoleID:    req.RoleID,
		Search:    req.Search,
		Page:      req.Page,
		PerPage:   req.PerPage,
		SortBy:    req.SortBy,
		SortOrder: req.SortOrder,
	}

	if req.Status != "" {
		status := entity.UserStatus(req.Status)
		filter.Status = &status
	}

	users, total, err := uc.UserRepo.List(ctx, filter)
	if err != nil {
		return nil, errors.ErrInternal("failed to list users").WithError(err)
	}

	items := make([]UserListItem, 0, len(users))
	for _, user := range users {
		profile, _ := uc.UserProfileRepo.GetByUserID(ctx, user.ID)
		securityState, _ := uc.UserSecurityStateRepo.GetByUserID(ctx, user.ID)
		items = append(items, mapUserToListItem(user, profile, securityState))
	}

	totalPages := int(total) / req.PerPage
	if int(total)%req.PerPage > 0 {
		totalPages++
	}

	return &ListResponse{
		Users: items,
		Pagination: Pagination{
			Total:      total,
			Page:       req.Page,
			PerPage:    req.PerPage,
			TotalPages: totalPages,
		},
	}, nil
}
