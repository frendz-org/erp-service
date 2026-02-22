package member

import (
	"context"

	"erp-service/entity"
	"erp-service/pkg/errors"

	"github.com/google/uuid"
)

func (uc *usecase) getUserEmail(ctx context.Context, userID uuid.UUID) (string, error) {
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.IsNotFound(err) {
			return "", nil
		}
		return "", err
	}
	return user.Email, nil
}

func (uc *usecase) validateTenantBoundary(reg *entity.UserTenantRegistration, tenantID, productID uuid.UUID) bool {
	if reg.TenantID != tenantID {
		return false
	}
	if reg.ProductID == nil || *reg.ProductID != productID {
		return false
	}
	return true
}

func mapRowToListItem(row MemberListRow) MemberListItem {
	return MemberListItem{
		ID:               row.Registration.ID,
		UserID:           row.Registration.UserID,
		FirstName:        row.FirstName,
		LastName:         row.LastName,
		Email:            row.Email,
		Status:           string(row.Registration.Status),
		RegistrationType: row.Registration.RegistrationType,
		RoleCode:         row.RoleCode,
		RoleName:         row.RoleName,
		CreatedAt:        row.Registration.CreatedAt,
	}
}

func mapToDetailResponse(reg *entity.UserTenantRegistration, profile *entity.UserProfile, email string, roleCode, roleName *string) *MemberDetailResponse {
	resp := &MemberDetailResponse{
		ID:               reg.ID,
		UserID:           reg.UserID,
		Email:            email,
		Status:           string(reg.Status),
		RegistrationType: reg.RegistrationType,
		RoleCode:         roleCode,
		RoleName:         roleName,
		ApprovedBy:       reg.ApprovedBy,
		ApprovedAt:       reg.ApprovedAt,
		Version:          reg.Version,
		CreatedAt:        reg.CreatedAt,
		UpdatedAt:        reg.UpdatedAt,
	}

	if profile != nil {
		resp.FirstName = profile.FirstName
		resp.LastName = profile.LastName
	}

	return resp
}
