package member

import (
	"context"

	"erp-service/pkg/errors"
)

func (uc *usecase) GetMyMember(ctx context.Context, req *GetMyMemberRequest) (*MyMemberResponse, error) {
	reg, err := uc.utrRepo.GetByUserAndProduct(ctx, req.UserID, req.TenantID, req.ProductID, "MEMBER")
	if err != nil {
		return nil, err
	}

	profile, profileErr := uc.profileRepo.GetByUserID(ctx, reg.UserID)
	if profileErr != nil && !errors.IsNotFound(profileErr) {
		return nil, profileErr
	}

	email, err := uc.getUserEmail(ctx, reg.UserID)
	if err != nil {
		return nil, err
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

	resp := &MyMemberResponse{
		ID:               reg.ID,
		UserID:           reg.UserID,
		Email:            email,
		Status:           string(reg.Status),
		RegistrationType: reg.RegistrationType,
		RoleCode:         roleCode,
		RoleName:         roleName,
		CreatedAt:        reg.CreatedAt,
		UpdatedAt:        reg.UpdatedAt,
	}

	if profile != nil {
		resp.FirstName = profile.FirstName
		resp.LastName = profile.LastName
	}

	m, mErr := uc.memberRepo.GetByRegistrationID(ctx, reg.ID)
	if mErr != nil && !errors.IsNotFound(mErr) {
		return nil, mErr
	}
	if m != nil {
		resp.ParticipantNumber = &m.ParticipantNumber
		resp.IdentityNumber = &m.IdentityNumber
		resp.OrganizationCode = &m.OrganizationCode
	}

	return resp, nil
}
