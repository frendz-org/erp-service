package participant

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"erp-service/entity"
	"erp-service/masterdata"
	apperrors "erp-service/pkg/errors"

	"github.com/google/uuid"
)

var (
	organizationCodeRegex  = regexp.MustCompile(`^TENANT_\d{3}$`)
	participantNumberRegex = regexp.MustCompile(`^([A-Z]{3}\d{5,8}|\d{8})$`)
	phoneNumberRegex       = regexp.MustCompile(`^(\+62\d{8,13}|08\d{8,12})$`)
)

func validateSelfRegisterRequest(req *SelfRegisterRequest) []apperrors.FieldError {

	errs := make([]apperrors.FieldError, 0, 4)

	if !organizationCodeRegex.MatchString(req.Organization) {
		errs = append(errs, apperrors.FieldError{
			Field:   "organization",
			Message: "must match pattern TENANT_XXX where XXX is 3 digits",
		})
	}

	if err := validateNIK(req.IdentityNumber); err != nil {
		errs = append(errs, apperrors.FieldError{
			Field:   "identity_number",
			Message: err.Error(),
		})
	}

	if !participantNumberRegex.MatchString(req.ParticipantNumber) {
		errs = append(errs, apperrors.FieldError{
			Field:   "participant_number",
			Message: "must be 3 uppercase letters followed by 5-8 digits, or exactly 8 digits",
		})
	}

	if !phoneNumberRegex.MatchString(req.PhoneNumber) {
		errs = append(errs, apperrors.FieldError{
			Field:   "phone_number",
			Message: "must be a valid Indonesian phone number (e.g. +628xxx or 08xxx)",
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

	day, err := strconv.Atoi(nik[6:8])
	if err != nil {
		return fmt.Errorf("invalid day in NIK")
	}

	if day > 40 {
		day -= 40
	}
	if day < 1 || day > 31 {
		return fmt.Errorf("invalid day in NIK: must be 01-31 (or 41-71 for female)")
	}

	month, err := strconv.Atoi(nik[8:10])
	if err != nil {
		return fmt.Errorf("invalid month in NIK")
	}
	if month < 1 || month > 12 {
		return fmt.Errorf("invalid month in NIK: must be 01-12")
	}

	return nil
}

func (uc *usecase) SelfRegister(ctx context.Context, req *SelfRegisterRequest) (*SelfRegisterResponse, error) {

	if fieldErrs := validateSelfRegisterRequest(req); len(fieldErrs) > 0 {
		return nil, apperrors.ErrValidationWithFields(fieldErrs)
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
		return nil, apperrors.ErrUnprocessable("invalid organization code")
	}

	tenant, err := uc.tenantRepo.GetByCode(ctx, req.Organization)
	if err != nil {
		return nil, err
	}
	if !tenant.IsActive() {
		return nil, apperrors.ErrUnprocessable("organization is inactive")
	}

	product, err := uc.productRepo.GetByCodeAndTenant(ctx, tenant.ID, "frendz-saving")
	if err != nil {
		return nil, err
	}

	config, err := uc.configRepo.GetByProductAndType(ctx, product.ID, "PARTICIPANT")
	if err != nil {
		return nil, apperrors.ErrUnprocessable("self-registration is not enabled")
	}
	if !config.IsActive {
		return nil, apperrors.ErrUnprocessable("self-registration is not enabled")
	}

	profile, err := uc.userProfileRepo.GetByUserID(ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	if profile.FullName() == "" || profile.DateOfBirth == nil || profile.Gender == nil {
		return nil, apperrors.ErrUnprocessable("user profile is incomplete: full name, date of birth, and gender are required")
	}

	existingParticipant, existingPension, err := uc.participantRepo.GetByKTPAndPensionNumber(
		ctx, req.IdentityNumber, req.ParticipantNumber, tenant.ID, product.ID,
	)
	if err != nil && !apperrors.IsNotFound(err) {
		return nil, fmt.Errorf("lookup participant: %w", err)
	}

	if existingParticipant != nil {

		switch {
		case existingParticipant.UserID != nil && *existingParticipant.UserID == req.UserID:
			return nil, apperrors.ErrConflict("registration not eligible")
		case existingParticipant.UserID != nil:
			return nil, apperrors.ErrConflict("registration not eligible")
		default:

			return uc.linkExistingParticipantToUser(ctx, req, existingParticipant, existingPension, tenant.ID, product.ID)
		}
	}

	_, err = uc.utrRepo.GetByUserAndProduct(ctx, req.UserID, tenant.ID, product.ID, "PARTICIPANT")
	if err == nil {

		return nil, apperrors.ErrConflict("registration not eligible")
	}
	if !apperrors.IsNotFound(err) {
		return nil, fmt.Errorf("check existing registration: %w", err)
	}

	return uc.createNewSelfRegisteredParticipant(ctx, req, profile, tenant.ID, product.ID)
}

func (uc *usecase) linkExistingParticipantToUser(
	ctx context.Context,
	req *SelfRegisterRequest,
	participant *entity.Participant,
	pension *entity.ParticipantPension,
	tenantID, productID uuid.UUID,
) (*SelfRegisterResponse, error) {
	var result *SelfRegisterResponse

	currentStatus := string(participant.Status)

	pensionNumber := ""
	if pension != nil && pension.ParticipantNumber != nil {
		pensionNumber = *pension.ParticipantNumber
	}

	err := uc.txManager.WithTransaction(ctx, func(txCtx context.Context) error {

		participant.UserID = &req.UserID
		if err := uc.participantRepo.Update(txCtx, participant); err != nil {
			return fmt.Errorf("update participant user link: %w", err)
		}

		utr := buildUTR(req.UserID, tenantID, productID)
		if err := uc.utrRepo.Create(txCtx, utr); err != nil {
			return fmt.Errorf("create user tenant registration: %w", err)
		}

		reason := "user linked via self-registration"
		history := &entity.ParticipantStatusHistory{
			ParticipantID: participant.ID,
			FromStatus:    &currentStatus,
			ToStatus:      currentStatus,
			ChangedBy:     req.UserID,
			Reason:        &reason,
			ChangedAt:     time.Now(),
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}
		if err := uc.statusHistoryRepo.Create(txCtx, history); err != nil {
			return fmt.Errorf("create status history: %w", err)
		}

		result = &SelfRegisterResponse{
			IsLinked:           true,
			RegistrationStatus: string(entity.UTRStatusPendingApproval),
			Data: &SelfRegisterParticipantData{
				ParticipantNumber: pensionNumber,
				Status:            currentStatus,
				CreatedAt:         participant.CreatedAt,
			},
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (uc *usecase) createNewSelfRegisteredParticipant(
	ctx context.Context,
	req *SelfRegisterRequest,
	profile *entity.UserProfile,
	tenantID, productID uuid.UUID,
) (*SelfRegisterResponse, error) {
	var result *SelfRegisterResponse

	err := uc.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		now := time.Now()
		genderStr := string(*profile.Gender)

		newParticipant := &entity.Participant{
			TenantID:       tenantID,
			ProductID:      productID,
			UserID:         &req.UserID,
			FullName:       profile.FullName(),
			DateOfBirth:    profile.DateOfBirth,
			Gender:         &genderStr,
			KTPNumber:      &req.IdentityNumber,
			PhoneNumber:    &req.PhoneNumber,
			Status:         entity.ParticipantStatusDraft,
			CreatedBy:      req.UserID,
			Version:        1,
			CreatedAt:      now,
			UpdatedAt:      now,
			StepsCompleted: map[string]bool{},
		}
		if err := uc.participantRepo.Create(txCtx, newParticipant); err != nil {

			if apperrors.IsConflict(err) {
				return apperrors.ErrConflict("registration not eligible")
			}
			return fmt.Errorf("create participant: %w", err)
		}

		pension := &entity.ParticipantPension{
			ParticipantID:     newParticipant.ID,
			ParticipantNumber: &req.ParticipantNumber,
			Version:           1,
			CreatedAt:         now,
			UpdatedAt:         now,
		}
		if err := uc.pensionRepo.Create(txCtx, pension); err != nil {
			return fmt.Errorf("create pension: %w", err)
		}

		reason := "self-registration"
		history := &entity.ParticipantStatusHistory{
			ParticipantID: newParticipant.ID,
			FromStatus:    nil,
			ToStatus:      string(entity.ParticipantStatusDraft),
			ChangedBy:     req.UserID,
			Reason:        &reason,
			ChangedAt:     now,
			CreatedAt:     now,
			UpdatedAt:     now,
		}
		if err := uc.statusHistoryRepo.Create(txCtx, history); err != nil {
			return fmt.Errorf("create status history: %w", err)
		}

		utr := buildUTR(req.UserID, tenantID, productID)
		if err := uc.utrRepo.Create(txCtx, utr); err != nil {
			return fmt.Errorf("create user tenant registration: %w", err)
		}

		result = &SelfRegisterResponse{
			IsLinked:           false,
			RegistrationStatus: string(entity.UTRStatusPendingApproval),
			Data: &SelfRegisterParticipantData{
				ParticipantNumber: req.ParticipantNumber,
				Status:            string(entity.ParticipantStatusDraft),
				CreatedAt:         now,
			},
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return result, nil
}

func buildUTR(userID, tenantID, productID uuid.UUID) *entity.UserTenantRegistration {
	pid := productID
	return &entity.UserTenantRegistration{
		UserID:           userID,
		TenantID:         tenantID,
		ProductID:        &pid,
		RegistrationType: "PARTICIPANT",
		Status:           entity.UTRStatusPendingApproval,
	}
}
