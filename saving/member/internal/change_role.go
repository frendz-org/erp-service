package internal

import (
	"context"
	"time"

	"erp-service/entity"
	"erp-service/pkg/errors"
	"erp-service/saving/member/memberdto"
)

func (uc *usecase) ChangeRole(ctx context.Context, req *memberdto.ChangeRoleRequest) (*memberdto.MemberDetailResponse, error) {
	reg, err := uc.utrRepo.GetByID(ctx, req.MemberID)
	if err != nil {
		return nil, err
	}

	if !uc.validateTenantBoundary(reg, req.TenantID, req.ProductID) {
		return nil, errors.ErrNotFound("member not found")
	}

	if reg.Status != entity.UTRStatusActive {
		return nil, errors.ErrBadRequest("only active members can have their role changed")
	}

	if reg.UserID == req.ActorID {
		return nil, errors.ErrForbidden("you cannot change your own role")
	}

	_, err = uc.userRole.GetActiveByUserAndProduct(ctx, reg.UserID, req.ProductID)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, errors.ErrBadRequest("member has no active role to change")
		}
		return nil, err
	}

	newRole, err := uc.roleRepo.GetByCodeAndProduct(ctx, req.ProductID, req.RoleCode)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, errors.ErrBadRequest("role not found for this product")
		}
		return nil, err
	}

	err = uc.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		if err := uc.userRole.SoftDeleteByUserAndProduct(txCtx, reg.UserID, req.ProductID); err != nil {
			return err
		}

		userRole := &entity.UserRole{
			UserID:     reg.UserID,
			RoleID:     newRole.ID,
			ProductID:  &req.ProductID,
			AssignedAt: time.Now(),
		}

		return uc.userRole.Create(txCtx, userRole)
	})
	if err != nil {
		return nil, err
	}

	profile, profileErr := uc.profileRepo.GetByUserID(ctx, reg.UserID)
	if profileErr != nil && !errors.IsNotFound(profileErr) {
		return nil, profileErr
	}
	email, emailErr := uc.getUserEmail(ctx, reg.UserID)
	if emailErr != nil {
		return nil, emailErr
	}

	return mapToDetailResponse(reg, profile, email, &newRole.Code, &newRole.Name), nil
}
