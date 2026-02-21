package member

import (
	"context"

	"erp-service/pkg/errors"
)

func (uc *usecase) GetMember(ctx context.Context, req *GetMemberRequest) (*MemberDetailResponse, error) {
	reg, err := uc.utrRepo.GetByID(ctx, req.MemberID)
	if err != nil {
		return nil, err
	}

	if !uc.validateTenantBoundary(reg, req.TenantID, req.ProductID) {
		return nil, errors.ErrNotFound("member not found")
	}

	profile, profileErr := uc.profileRepo.GetByUserID(ctx, reg.UserID)
	if profileErr != nil && !errors.IsNotFound(profileErr) {
		return nil, profileErr
	}
	email, emailErr := uc.getUserEmail(ctx, reg.UserID)
	if emailErr != nil {
		return nil, emailErr
	}

	var roleCode, roleName *string
	activeRole, err := uc.userRole.GetActiveByUserAndProduct(ctx, reg.UserID, req.ProductID)
	if err == nil && activeRole != nil {
		role, rErr := uc.roleRepo.GetByID(ctx, activeRole.RoleID)
		if rErr == nil && role != nil {
			roleCode = &role.Code
			roleName = &role.Name
		}
	}

	return mapToDetailResponse(reg, profile, email, roleCode, roleName), nil
}
