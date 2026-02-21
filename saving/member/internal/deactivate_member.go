package internal

import (
	"context"

	"erp-service/entity"
	"erp-service/pkg/errors"
	"erp-service/saving/member/memberdto"
)

func (uc *usecase) DeactivateMember(ctx context.Context, req *memberdto.DeactivateRequest) (*memberdto.MemberDetailResponse, error) {
	reg, err := uc.utrRepo.GetByID(ctx, req.MemberID)
	if err != nil {
		return nil, err
	}

	if !uc.validateTenantBoundary(reg, req.TenantID, req.ProductID) {
		return nil, errors.ErrNotFound("member not found")
	}

	if reg.Status != entity.UTRStatusActive {
		return nil, errors.ErrBadRequest("only active members can be deactivated")
	}

	if reg.UserID == req.ActorID {
		return nil, errors.ErrForbidden("you cannot deactivate your own membership")
	}

	err = uc.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		reg.Status = entity.UTRStatusInactive

		if err := uc.utrRepo.UpdateStatus(txCtx, reg); err != nil {
			return err
		}

		return uc.userRole.SoftDeleteByUserAndProduct(txCtx, reg.UserID, req.ProductID)
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

	return mapToDetailResponse(reg, profile, email, nil, nil), nil
}
