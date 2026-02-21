package presenter

import (
	"erp-service/delivery/http/dto/response"
	"erp-service/saving/participant/participantdto"
)

func MapSelfRegisterResponse(resp *participantdto.SelfRegisterResponse) *participantdto.SelfRegisterResponse {
	return resp
}

func MapEmploymentResponse(dto *participantdto.EmploymentResponse) response.EmploymentResponse {
	return response.EmploymentResponse{
		ID:                 dto.ID,
		PersonnelNumber:    dto.PersonnelNumber,
		DateOfHire:         dto.DateOfHire,
		CorporateGroupName: dto.CorporateGroupName,
		LegalEntityCode:    dto.LegalEntityCode,
		LegalEntityName:    dto.LegalEntityName,
		BusinessUnitCode:   dto.BusinessUnitCode,
		BusinessUnitName:   dto.BusinessUnitName,
		TenantName:         dto.TenantName,
		EmploymentStatus:   dto.EmploymentStatus,
		PositionName:       dto.PositionName,
		JobLevel:           dto.JobLevel,
		LocationCode:       dto.LocationCode,
		LocationName:       dto.LocationName,
		SubLocationName:    dto.SubLocationName,
		RetirementDate:     dto.RetirementDate,
		RetirementTypeCode: dto.RetirementTypeCode,
		Version:            dto.Version,
		CreatedAt:          dto.CreatedAt,
		UpdatedAt:          dto.UpdatedAt,
	}
}

func MapPensionResponse(dto *participantdto.PensionResponse) response.PensionResponse {
	return response.PensionResponse{
		ID:                      dto.ID,
		ParticipantNumber:       dto.ParticipantNumber,
		PensionCategory:         dto.PensionCategory,
		PensionStatus:           dto.PensionStatus,
		EffectiveDate:           dto.EffectiveDate,
		EndDate:                 dto.EndDate,
		ProjectedRetirementDate: dto.ProjectedRetirementDate,
		Version:                 dto.Version,
		CreatedAt:               dto.CreatedAt,
		UpdatedAt:               dto.UpdatedAt,
	}
}

func MapParticipantResponse(dto *participantdto.ParticipantResponse) response.ParticipantResponse {
	identities := make([]response.IdentityResponse, 0, len(dto.Identities))
	for _, id := range dto.Identities {
		identities = append(identities, response.IdentityResponse{
			ID:                id.ID,
			IdentityType:      id.IdentityType,
			IdentityNumber:    id.IdentityNumber,
			IdentityAuthority: id.IdentityAuthority,
			IssueDate:         id.IssueDate,
			ExpiryDate:        id.ExpiryDate,
			PhotoFilePath:     id.PhotoFilePath,
			Version:           id.Version,
			CreatedAt:         id.CreatedAt,
			UpdatedAt:         id.UpdatedAt,
		})
	}

	addresses := make([]response.AddressResponse, 0, len(dto.Addresses))
	for _, addr := range dto.Addresses {
		addresses = append(addresses, response.AddressResponse{
			ID:              addr.ID,
			AddressType:     addr.AddressType,
			CountryCode:     addr.CountryCode,
			ProvinceCode:    addr.ProvinceCode,
			CityCode:        addr.CityCode,
			DistrictCode:    addr.DistrictCode,
			SubdistrictCode: addr.SubdistrictCode,
			PostalCode:      addr.PostalCode,
			RT:              addr.RT,
			RW:              addr.RW,
			AddressLine:     addr.AddressLine,
			IsPrimary:       addr.IsPrimary,
			Version:         addr.Version,
			CreatedAt:       addr.CreatedAt,
			UpdatedAt:       addr.UpdatedAt,
		})
	}

	bankAccounts := make([]response.BankAccountResponse, 0, len(dto.BankAccounts))
	for _, ba := range dto.BankAccounts {
		bankAccounts = append(bankAccounts, response.BankAccountResponse{
			ID:                ba.ID,
			BankCode:          ba.BankCode,
			AccountNumber:     ba.AccountNumber,
			AccountHolderName: ba.AccountHolderName,
			AccountType:       ba.AccountType,
			CurrencyCode:      ba.CurrencyCode,
			IsPrimary:         ba.IsPrimary,
			IssueDate:         ba.IssueDate,
			ExpiryDate:        ba.ExpiryDate,
			Version:           ba.Version,
			CreatedAt:         ba.CreatedAt,
			UpdatedAt:         ba.UpdatedAt,
		})
	}

	familyMembers := make([]response.FamilyMemberResponse, 0, len(dto.FamilyMembers))
	for _, fm := range dto.FamilyMembers {
		familyMembers = append(familyMembers, response.FamilyMemberResponse{
			ID:                    fm.ID,
			FullName:              fm.FullName,
			RelationshipType:      fm.RelationshipType,
			IsDependent:           fm.IsDependent,
			SupportingDocFilePath: fm.SupportingDocFilePath,
			Version:               fm.Version,
			CreatedAt:             fm.CreatedAt,
			UpdatedAt:             fm.UpdatedAt,
		})
	}

	var employment *response.EmploymentResponse
	if dto.Employment != nil {
		employment = &response.EmploymentResponse{
			ID:                 dto.Employment.ID,
			PersonnelNumber:    dto.Employment.PersonnelNumber,
			DateOfHire:         dto.Employment.DateOfHire,
			CorporateGroupName: dto.Employment.CorporateGroupName,
			LegalEntityCode:    dto.Employment.LegalEntityCode,
			LegalEntityName:    dto.Employment.LegalEntityName,
			BusinessUnitCode:   dto.Employment.BusinessUnitCode,
			BusinessUnitName:   dto.Employment.BusinessUnitName,
			TenantName:         dto.Employment.TenantName,
			EmploymentStatus:   dto.Employment.EmploymentStatus,
			PositionName:       dto.Employment.PositionName,
			JobLevel:           dto.Employment.JobLevel,
			LocationCode:       dto.Employment.LocationCode,
			LocationName:       dto.Employment.LocationName,
			SubLocationName:    dto.Employment.SubLocationName,
			RetirementDate:     dto.Employment.RetirementDate,
			RetirementTypeCode: dto.Employment.RetirementTypeCode,
			Version:            dto.Employment.Version,
			CreatedAt:          dto.Employment.CreatedAt,
			UpdatedAt:          dto.Employment.UpdatedAt,
		}
	}

	var pension *response.PensionResponse
	if dto.Pension != nil {
		pension = &response.PensionResponse{
			ID:                      dto.Pension.ID,
			ParticipantNumber:       dto.Pension.ParticipantNumber,
			PensionCategory:         dto.Pension.PensionCategory,
			PensionStatus:           dto.Pension.PensionStatus,
			EffectiveDate:           dto.Pension.EffectiveDate,
			EndDate:                 dto.Pension.EndDate,
			ProjectedRetirementDate: dto.Pension.ProjectedRetirementDate,
			Version:                 dto.Pension.Version,
			CreatedAt:               dto.Pension.CreatedAt,
			UpdatedAt:               dto.Pension.UpdatedAt,
		}
	}

	beneficiaries := make([]response.BeneficiaryResponse, 0, len(dto.Beneficiaries))
	for _, ben := range dto.Beneficiaries {
		beneficiaries = append(beneficiaries, response.BeneficiaryResponse{
			ID:                      ben.ID,
			FamilyMemberID:          ben.FamilyMemberID,
			IdentityPhotoFilePath:   ben.IdentityPhotoFilePath,
			FamilyCardPhotoFilePath: ben.FamilyCardPhotoFilePath,
			BankBookPhotoFilePath:   ben.BankBookPhotoFilePath,
			AccountNumber:           ben.AccountNumber,
			Version:                 ben.Version,
			CreatedAt:               ben.CreatedAt,
			UpdatedAt:               ben.UpdatedAt,
		})
	}

	return response.ParticipantResponse{
		ID:              dto.ID,
		TenantID:        dto.TenantID,
		ProductID:       dto.ProductID,
		UserID:          dto.UserID,
		FullName:        dto.FullName,
		Gender:          dto.Gender,
		PlaceOfBirth:    dto.PlaceOfBirth,
		DateOfBirth:     dto.DateOfBirth,
		MaritalStatus:   dto.MaritalStatus,
		Citizenship:     dto.Citizenship,
		Religion:        dto.Religion,
		KTPNumber:       dto.KTPNumber,
		EmployeeNumber:  dto.EmployeeNumber,
		PhoneNumber:     dto.PhoneNumber,
		Status:          dto.Status,
		CreatedBy:       dto.CreatedBy,
		SubmittedBy:     dto.SubmittedBy,
		SubmittedAt:     dto.SubmittedAt,
		ApprovedBy:      dto.ApprovedBy,
		ApprovedAt:      dto.ApprovedAt,
		RejectedBy:      dto.RejectedBy,
		RejectedAt:      dto.RejectedAt,
		RejectionReason: dto.RejectionReason,
		Version:         dto.Version,
		CreatedAt:       dto.CreatedAt,
		UpdatedAt:       dto.UpdatedAt,
		Identities:      identities,
		Addresses:       addresses,
		BankAccounts:    bankAccounts,
		FamilyMembers:   familyMembers,
		Employment:      employment,
		Pension:         pension,
		Beneficiaries:   beneficiaries,
	}
}
