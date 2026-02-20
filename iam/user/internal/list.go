package internal

import (
	"context"

	"iam-service/entity"
	"iam-service/iam/user/contract"
	"iam-service/iam/user/userdto"
	"iam-service/pkg/errors"

	"github.com/google/uuid"
)

func (uc *usecase) List(ctx context.Context, tenantID *uuid.UUID, req *userdto.ListRequest) (*userdto.ListResponse, error) {
	req.SetDefaults()

	filter := &contract.UserListFilter{
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

	items := make([]userdto.UserListItem, 0, len(users))
	for _, user := range users {
		profile, _ := uc.UserProfileRepo.GetByUserID(ctx, user.ID)
		securityState, _ := uc.UserSecurityStateRepo.GetByUserID(ctx, user.ID)
		items = append(items, mapUserToListItem(user, profile, securityState))
	}

	totalPages := int(total) / req.PerPage
	if int(total)%req.PerPage > 0 {
		totalPages++
	}

	return &userdto.ListResponse{
		Users: items,
		Pagination: userdto.Pagination{
			Total:      total,
			Page:       req.Page,
			PerPage:    req.PerPage,
			TotalPages: totalPages,
		},
	}, nil
}
