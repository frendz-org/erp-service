package internal

import (
	"context"
	"time"

	"iam-service/entity"
	"iam-service/pkg/errors"
	"iam-service/saving/member/memberdto"
)

func (uc *usecase) ApproveMember(ctx context.Context, req *memberdto.ApproveRequest) (*memberdto.MemberDetailResponse, error) {
	reg, err := uc.utrRepo.GetByID(ctx, req.MemberID)
	if err != nil {
		return nil, err
	}

	if !uc.validateTenantBoundary(reg, req.TenantID, req.ProductID) {
		return nil, errors.ErrNotFound("member not found")
	}

	if reg.Status != entity.UTRStatusPendingApproval {
		return nil, errors.ErrBadRequest("only pending members can be approved")
	}

	if reg.UserID == req.ApproverID {
		return nil, errors.ErrForbidden("you cannot approve your own registration")
	}

	role, err := uc.roleRepo.GetByCodeAndProduct(ctx, req.ProductID, req.RoleCode)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, errors.ErrBadRequest("role not found for this product")
		}
		return nil, err
	}

	err = uc.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		now := time.Now()
		reg.Status = entity.UTRStatusActive
		reg.ApprovedBy = &req.ApproverID
		reg.ApprovedAt = &now

		if err := uc.utrRepo.UpdateStatus(txCtx, reg); err != nil {
			return err
		}

		userRole := &entity.UserRole{
			UserID:     reg.UserID,
			RoleID:     role.ID,
			ProductID:  &req.ProductID,
			AssignedAt: now,
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

	return mapToDetailResponse(reg, profile, email, &role.Code, &role.Name), nil
}
