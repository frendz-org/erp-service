package internal

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"iam-service/entity"
	"iam-service/saving/participant/participantdto"
	"iam-service/pkg/errors"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

func validateParticipantOwnership(participant *entity.Participant, tenantID, productID uuid.UUID) error {
	if participant.TenantID != tenantID {
		return errors.ErrForbidden("participant does not belong to this tenant")
	}
	if participant.ProductID != productID {
		return errors.ErrForbidden("participant does not belong to this product")
	}
	return nil
}

func validateEditableState(participant *entity.Participant) error {
	if !participant.CanBeEdited() {
		return errors.ErrBadRequest(fmt.Sprintf("participant in %s status cannot be edited", participant.Status))
	}
	return nil
}

var allowedFieldNames = map[string]bool{
	"ktp_photo":        true,
	"passport_photo":   true,
	"family_card":      true,
	"identity_photo":   true,
	"bank_book_photo":  true,
	"supporting_doc":   true,
	"profile_photo":    true,
}

func sanitizeFieldName(fieldName string) string {
	safe := filepath.Base(fieldName)
	if safe == "." || safe == ".." || strings.ContainsAny(safe, "/\\") {
		return "unknown"
	}
	if !allowedFieldNames[safe] {
		return "unknown"
	}
	return safe
}

func sanitizeFilename(filename string) string {
	safe := filepath.Base(filename)
	if safe == "." || safe == ".." || strings.ContainsAny(safe, "/\\") {
		return "upload"
	}
	return safe
}

func generateObjectKey(tenantID, productID, participantID uuid.UUID, fieldName, filename string) string {
	safeField := sanitizeFieldName(fieldName)
	safeFile := sanitizeFilename(filename)
	return fmt.Sprintf("participants/%s/%s/%s/%s/%s", tenantID.String(), productID.String(), participantID.String(), safeField, safeFile)
}

func mapIdentityToResponse(identity *entity.ParticipantIdentity) participantdto.IdentityResponse {
	return participantdto.IdentityResponse{
		ID:                identity.ID,
		IdentityType:      identity.IdentityType,
		IdentityNumber:    identity.IdentityNumber,
		IdentityAuthority: identity.IdentityAuthority,
		IssueDate:         identity.IssueDate,
		ExpiryDate:        identity.ExpiryDate,
		PhotoFilePath:     identity.PhotoFilePath,
		Version:           identity.Version,
		CreatedAt:         identity.CreatedAt,
		UpdatedAt:         identity.UpdatedAt,
	}
}

func mapAddressToResponse(address *entity.ParticipantAddress) participantdto.AddressResponse {
	return participantdto.AddressResponse{
		ID:              address.ID,
		AddressType:     address.AddressType,
		CountryCode:     address.CountryCode,
		ProvinceCode:    address.ProvinceCode,
		CityCode:        address.CityCode,
		DistrictCode:    address.DistrictCode,
		SubdistrictCode: address.SubdistrictCode,
		PostalCode:      address.PostalCode,
		RT:              address.RT,
		RW:              address.RW,
		AddressLine:     address.AddressLine,
		IsPrimary:       address.IsPrimary,
		Version:         address.Version,
		CreatedAt:       address.CreatedAt,
		UpdatedAt:       address.UpdatedAt,
	}
}

func mapBankAccountToResponse(account *entity.ParticipantBankAccount) participantdto.BankAccountResponse {
	return participantdto.BankAccountResponse{
		ID:                account.ID,
		BankCode:          account.BankCode,
		AccountNumber:     account.AccountNumber,
		AccountHolderName: account.AccountHolderName,
		AccountType:       account.AccountType,
		CurrencyCode:      account.CurrencyCode,
		IsPrimary:         account.IsPrimary,
		IssueDate:         account.IssueDate,
		ExpiryDate:        account.ExpiryDate,
		Version:           account.Version,
		CreatedAt:         account.CreatedAt,
		UpdatedAt:         account.UpdatedAt,
	}
}

func mapFamilyMemberToResponse(member *entity.ParticipantFamilyMember) participantdto.FamilyMemberResponse {
	return participantdto.FamilyMemberResponse{
		ID:                    member.ID,
		FullName:              member.FullName,
		RelationshipType:      member.RelationshipType,
		IsDependent:           member.IsDependent,
		SupportingDocFilePath: member.SupportingDocFilePath,
		Version:               member.Version,
		CreatedAt:             member.CreatedAt,
		UpdatedAt:             member.UpdatedAt,
	}
}

func mapEmploymentToResponse(employment *entity.ParticipantEmployment) participantdto.EmploymentResponse {
	return participantdto.EmploymentResponse{
		ID:                 employment.ID,
		PersonnelNumber:    employment.PersonnelNumber,
		DateOfHire:         employment.DateOfHire,
		CorporateGroupName: employment.CorporateGroupName,
		LegalEntityCode:    employment.LegalEntityCode,
		LegalEntityName:    employment.LegalEntityName,
		BusinessUnitCode:   employment.BusinessUnitCode,
		BusinessUnitName:   employment.BusinessUnitName,
		TenantName:         employment.TenantName,
		EmploymentStatus:   employment.EmploymentStatus,
		PositionName:       employment.PositionName,
		JobLevel:           employment.JobLevel,
		LocationCode:       employment.LocationCode,
		LocationName:       employment.LocationName,
		SubLocationName:    employment.SubLocationName,
		RetirementDate:     employment.RetirementDate,
		RetirementTypeCode: employment.RetirementTypeCode,
		Version:            employment.Version,
		CreatedAt:          employment.CreatedAt,
		UpdatedAt:          employment.UpdatedAt,
	}
}

func mapPensionToResponse(pension *entity.ParticipantPension) participantdto.PensionResponse {
	return participantdto.PensionResponse{
		ID:                      pension.ID,
		ParticipantNumber:       pension.ParticipantNumber,
		PensionCategory:         pension.PensionCategory,
		PensionStatus:           pension.PensionStatus,
		EffectiveDate:           pension.EffectiveDate,
		EndDate:                 pension.EndDate,
		ProjectedRetirementDate: pension.ProjectedRetirementDate,
		Version:                 pension.Version,
		CreatedAt:               pension.CreatedAt,
		UpdatedAt:               pension.UpdatedAt,
	}
}

func mapBeneficiaryToResponse(beneficiary *entity.ParticipantBeneficiary) participantdto.BeneficiaryResponse {
	return participantdto.BeneficiaryResponse{
		ID:                      beneficiary.ID,
		FamilyMemberID:          beneficiary.FamilyMemberID,
		IdentityPhotoFilePath:   beneficiary.IdentityPhotoFilePath,
		FamilyCardPhotoFilePath: beneficiary.FamilyCardPhotoFilePath,
		BankBookPhotoFilePath:   beneficiary.BankBookPhotoFilePath,
		AccountNumber:           beneficiary.AccountNumber,
		Version:                 beneficiary.Version,
		CreatedAt:               beneficiary.CreatedAt,
		UpdatedAt:               beneficiary.UpdatedAt,
	}
}

func mapStatusHistoryToResponse(history *entity.ParticipantStatusHistory) participantdto.StatusHistoryResponse {
	return participantdto.StatusHistoryResponse{
		ID:         history.ID,
		FromStatus: history.FromStatus,
		ToStatus:   history.ToStatus,
		ChangedBy:  history.ChangedBy,
		Reason:     history.Reason,
		ChangedAt:  history.ChangedAt,
		CreatedAt:  history.CreatedAt,
	}
}

func (uc *usecase) buildFullParticipantResponse(ctx context.Context, participant *entity.Participant, concurrent bool) (*participantdto.ParticipantResponse, error) {
	resp := &participantdto.ParticipantResponse{
		ID:              participant.ID,
		TenantID:        participant.TenantID,
		ProductID:       participant.ProductID,
		UserID:          participant.UserID,
		FullName:        participant.FullName,
		Gender:          participant.Gender,
		PlaceOfBirth:    participant.PlaceOfBirth,
		DateOfBirth:     participant.DateOfBirth,
		MaritalStatus:   participant.MaritalStatus,
		Citizenship:     participant.Citizenship,
		Religion:        participant.Religion,
		KTPNumber:       participant.KTPNumber,
		EmployeeNumber:  participant.EmployeeNumber,
		PhoneNumber:     participant.PhoneNumber,
		Status:          string(participant.Status),
		CreatedBy:       participant.CreatedBy,
		SubmittedBy:     participant.SubmittedBy,
		SubmittedAt:     participant.SubmittedAt,
		ApprovedBy:      participant.ApprovedBy,
		ApprovedAt:      participant.ApprovedAt,
		RejectedBy:      participant.RejectedBy,
		RejectedAt:      participant.RejectedAt,
		RejectionReason: participant.RejectionReason,
		Version:         participant.Version,
		CreatedAt:       participant.CreatedAt,
		UpdatedAt:       participant.UpdatedAt,
	}

	var (
		identities    []*entity.ParticipantIdentity
		addresses     []*entity.ParticipantAddress
		bankAccounts  []*entity.ParticipantBankAccount
		familyMembers []*entity.ParticipantFamilyMember
		employment    *entity.ParticipantEmployment
		pension       *entity.ParticipantPension
		beneficiaries []*entity.ParticipantBeneficiary
	)

	if concurrent {
		if err := uc.loadChildEntitiesConcurrent(ctx, participant.ID, &identities, &addresses, &bankAccounts, &familyMembers, &employment, &pension, &beneficiaries); err != nil {
			return nil, err
		}
	} else {
		if err := uc.loadChildEntitiesSequential(ctx, participant.ID, &identities, &addresses, &bankAccounts, &familyMembers, &employment, &pension, &beneficiaries); err != nil {
			return nil, err
		}
	}

	resp.Identities = make([]participantdto.IdentityResponse, 0, len(identities))
	for _, identity := range identities {
		resp.Identities = append(resp.Identities, mapIdentityToResponse(identity))
	}

	resp.Addresses = make([]participantdto.AddressResponse, 0, len(addresses))
	for _, address := range addresses {
		resp.Addresses = append(resp.Addresses, mapAddressToResponse(address))
	}

	resp.BankAccounts = make([]participantdto.BankAccountResponse, 0, len(bankAccounts))
	for _, account := range bankAccounts {
		resp.BankAccounts = append(resp.BankAccounts, mapBankAccountToResponse(account))
	}

	resp.FamilyMembers = make([]participantdto.FamilyMemberResponse, 0, len(familyMembers))
	for _, member := range familyMembers {
		resp.FamilyMembers = append(resp.FamilyMembers, mapFamilyMemberToResponse(member))
	}

	if employment != nil {
		empResp := mapEmploymentToResponse(employment)
		resp.Employment = &empResp
	}

	if pension != nil {
		penResp := mapPensionToResponse(pension)
		resp.Pension = &penResp
	}

	resp.Beneficiaries = make([]participantdto.BeneficiaryResponse, 0, len(beneficiaries))
	for _, beneficiary := range beneficiaries {
		resp.Beneficiaries = append(resp.Beneficiaries, mapBeneficiaryToResponse(beneficiary))
	}

	return resp, nil
}

func (uc *usecase) loadChildEntitiesSequential(ctx context.Context, participantID uuid.UUID,
	identities *[]*entity.ParticipantIdentity,
	addresses *[]*entity.ParticipantAddress,
	bankAccounts *[]*entity.ParticipantBankAccount,
	familyMembers *[]*entity.ParticipantFamilyMember,
	employment **entity.ParticipantEmployment,
	pension **entity.ParticipantPension,
	beneficiaries *[]*entity.ParticipantBeneficiary,
) error {
	var err error

	*identities, err = uc.identityRepo.ListByParticipantID(ctx, participantID)
	if err != nil {
		return fmt.Errorf("load identities: %w", err)
	}

	*addresses, err = uc.addressRepo.ListByParticipantID(ctx, participantID)
	if err != nil {
		return fmt.Errorf("load addresses: %w", err)
	}

	*bankAccounts, err = uc.bankAccountRepo.ListByParticipantID(ctx, participantID)
	if err != nil {
		return fmt.Errorf("load bank accounts: %w", err)
	}

	*familyMembers, err = uc.familyMemberRepo.ListByParticipantID(ctx, participantID)
	if err != nil {
		return fmt.Errorf("load family members: %w", err)
	}

	*employment, err = uc.employmentRepo.GetByParticipantID(ctx, participantID)
	if err != nil && !errors.IsNotFound(err) {
		return fmt.Errorf("load employment: %w", err)
	}

	*pension, err = uc.pensionRepo.GetByParticipantID(ctx, participantID)
	if err != nil && !errors.IsNotFound(err) {
		return fmt.Errorf("load pension: %w", err)
	}

	*beneficiaries, err = uc.beneficiaryRepo.ListByParticipantID(ctx, participantID)
	if err != nil {
		return fmt.Errorf("load beneficiaries: %w", err)
	}

	return nil
}

func (uc *usecase) loadChildEntitiesConcurrent(ctx context.Context, participantID uuid.UUID,
	identities *[]*entity.ParticipantIdentity,
	addresses *[]*entity.ParticipantAddress,
	bankAccounts *[]*entity.ParticipantBankAccount,
	familyMembers *[]*entity.ParticipantFamilyMember,
	employment **entity.ParticipantEmployment,
	pension **entity.ParticipantPension,
	beneficiaries *[]*entity.ParticipantBeneficiary,
) error {
	g, gctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		var err error
		*identities, err = uc.identityRepo.ListByParticipantID(gctx, participantID)
		if err != nil {
			return fmt.Errorf("load identities: %w", err)
		}
		return nil
	})

	g.Go(func() error {
		var err error
		*addresses, err = uc.addressRepo.ListByParticipantID(gctx, participantID)
		if err != nil {
			return fmt.Errorf("load addresses: %w", err)
		}
		return nil
	})

	g.Go(func() error {
		var err error
		*bankAccounts, err = uc.bankAccountRepo.ListByParticipantID(gctx, participantID)
		if err != nil {
			return fmt.Errorf("load bank accounts: %w", err)
		}
		return nil
	})

	g.Go(func() error {
		var err error
		*familyMembers, err = uc.familyMemberRepo.ListByParticipantID(gctx, participantID)
		if err != nil {
			return fmt.Errorf("load family members: %w", err)
		}
		return nil
	})

	g.Go(func() error {
		var err error
		*employment, err = uc.employmentRepo.GetByParticipantID(gctx, participantID)
		if err != nil && !errors.IsNotFound(err) {
			return fmt.Errorf("load employment: %w", err)
		}
		return nil
	})

	g.Go(func() error {
		var err error
		*pension, err = uc.pensionRepo.GetByParticipantID(gctx, participantID)
		if err != nil && !errors.IsNotFound(err) {
			return fmt.Errorf("load pension: %w", err)
		}
		return nil
	})

	g.Go(func() error {
		var err error
		*beneficiaries, err = uc.beneficiaryRepo.ListByParticipantID(gctx, participantID)
		if err != nil {
			return fmt.Errorf("load beneficiaries: %w", err)
		}
		return nil
	})

	return g.Wait()
}
