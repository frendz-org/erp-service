package member

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"

	"erp-service/entity"
	"erp-service/masterdata"
	"erp-service/pkg/errors"
)

var (
	organizationCodeRegex  = regexp.MustCompile(`^TENANT_\d{3}$`)
	participantNumberRegex = regexp.MustCompile(`^([A-Z]{3}\d{5,8}|\d{8})$`)
)

func validateRegisterRequest(req *RegisterRequest) []errors.FieldError {
	errs := make([]errors.FieldError, 0, 3)

	if !organizationCodeRegex.MatchString(req.Organization) {
		errs = append(errs, errors.FieldError{
			Field:   "organization",
			Message: "must match pattern TENANT_XXX where XXX is 3 digits",
		})
	}

	if err := validateNIK(req.IdentityNumber); err != nil {
		errs = append(errs, errors.FieldError{
			Field:   "identity_number",
			Message: err.Error(),
		})
	}

	if !participantNumberRegex.MatchString(req.ParticipantNumber) {
		errs = append(errs, errors.FieldError{
			Field:   "participant_number",
			Message: "must be 3 uppercase letters followed by 5-8 digits, or exactly 8 digits",
		})
	}

	return errs
}

func validateNIK(nik string) error {
	if len(nik) != 16 {
		return fmt.Errorf("must be exactly 16 digits")
	}
	for _, c := range nik {
		if c < '0' || c > '9' {
			return fmt.Errorf("must contain only digits")
		}
	}

	day := int(nik[6]-'0')*10 + int(nik[7]-'0')
	if day > 40 {
		day -= 40
	}
	if day < 1 || day > 31 {
		return fmt.Errorf("invalid day in NIK: must be 01-31 (or 41-71 for female)")
	}

	month := int(nik[8]-'0')*10 + int(nik[9]-'0')
	if month < 1 || month > 12 {
		return fmt.Errorf("invalid month in NIK: must be 01-12")
	}

	return nil
}

func (uc *usecase) RegisterMember(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error) {
	if fieldErrs := validateRegisterRequest(req); len(fieldErrs) > 0 {
		return nil, errors.ErrValidationWithFields(fieldErrs)
	}

	validateResp, err := uc.masterdataUsecase.ValidateItemCode(ctx, &masterdata.ValidateCodeRequest{
		CategoryCode:  "TENANT",
		ItemCode:      req.Organization,
		RequireActive: true,
	})
	if err != nil {
		return nil, fmt.Errorf("validate organization: %w", err)
	}
	if !validateResp.Valid {
		return nil, errors.ErrUnprocessable("invalid organization code")
	}

	tenant, err := uc.tenantRepo.GetByCode(ctx, req.Organization)
	if err != nil {
		return nil, err
	}
	if !tenant.IsActive() {
		return nil, errors.ErrUnprocessable("organization is inactive")
	}

	product, err := uc.productRepo.GetByCodeAndTenant(ctx, tenant.ID, "frendz-saving")
	if err != nil {
		return nil, err
	}

	regConfig, err := uc.configRepo.GetByProductAndType(ctx, product.ID, "MEMBER")
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, errors.ErrBadRequest("member registration is not configured for this product")
		}
		return nil, err
	}
	if !regConfig.IsActive {
		return nil, errors.ErrBadRequest("member registration is currently not accepting new registrations")
	}

	existing, err := uc.utrRepo.GetByUserAndProduct(ctx, req.UserID, tenant.ID, product.ID, "MEMBER")
	if err != nil && !errors.IsNotFound(err) {
		return nil, err
	}
	if existing != nil {
		return nil, errors.ErrConflict("you already have a member registration for this product")
	}

	_, err = uc.csiEmployeeRepo.GetByEmployeeNo(ctx, req.ParticipantNumber)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, errors.ErrNotFound("employee number not found")
		}
		return nil, fmt.Errorf("validate employee number: %w", err)
	}

	existingMember, err := uc.memberRepo.GetByParticipantNumber(ctx, tenant.ID, product.ID, req.ParticipantNumber)
	if err != nil && !errors.IsNotFound(err) {
		return nil, fmt.Errorf("check participant number: %w", err)
	}
	if existingMember != nil {
		return nil, errors.ErrConflict("participant number is already used")
	}

	profile, err := uc.profileRepo.GetByUserID(ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	if profile.FullName() == "" || profile.DateOfBirth == nil || profile.Gender == nil {
		return nil, errors.ErrUnprocessable("user profile is incomplete: full name, date of birth, and gender are required")
	}

	var reg *entity.UserTenantRegistration
	var memberRecord *entity.Member

	err = uc.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		pid := product.ID
		reg = &entity.UserTenantRegistration{
			UserID:           req.UserID,
			TenantID:         tenant.ID,
			ProductID:        &pid,
			RegistrationType: "MEMBER",
			Status:           entity.UTRStatusPendingApproval,
			Metadata:         json.RawMessage(`{}`),
		}

		if err := uc.utrRepo.Create(txCtx, reg); err != nil {
			return err
		}

		genderStr := string(*profile.Gender)
		memberRecord = &entity.Member{
			UserTenantRegistrationID: reg.ID,
			TenantID:                 tenant.ID,
			ProductID:                product.ID,
			UserID:                   req.UserID,
			ParticipantNumber:        req.ParticipantNumber,
			IdentityNumber:           req.IdentityNumber,
			OrganizationCode:         req.Organization,
			FullName:                 profile.FullName(),
			Gender:                   &genderStr,
			DateOfBirth:              profile.DateOfBirth,
		}

		return uc.memberRepo.Create(txCtx, memberRecord)
	})
	if err != nil {
		return nil, err
	}

	return &RegisterResponse{
		ID:                reg.ID,
		TenantID:          tenant.ID,
		Status:            string(reg.Status),
		RegistrationType:  reg.RegistrationType,
		ParticipantNumber: memberRecord.ParticipantNumber,
		IdentityNumber:    memberRecord.IdentityNumber,
		OrganizationCode:  memberRecord.OrganizationCode,
		FullName:          memberRecord.FullName,
		CreatedAt:         reg.CreatedAt,
	}, nil
}
