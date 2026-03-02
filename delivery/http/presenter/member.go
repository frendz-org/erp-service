package presenter

import (
	"erp-service/delivery/http/dto/response"
	"erp-service/saving/member"
)

func MapMemberRegisterResponse(dto *member.RegisterResponse) response.MemberRegisterResponse {
	return response.MemberRegisterResponse{
		ID:                dto.ID,
		TenantID:          dto.TenantID,
		Status:            dto.Status,
		RegistrationType:  dto.RegistrationType,
		ParticipantNumber: dto.ParticipantNumber,
		IdentityNumber:    dto.IdentityNumber,
		OrganizationCode:  dto.OrganizationCode,
		FullName:          dto.FullName,
		CreatedAt:         dto.CreatedAt,
	}
}

func MapMyMemberResponse(dto *member.MyMemberResponse) response.MyMemberResponse {
	return response.MyMemberResponse{
		ID:                dto.ID,
		UserID:            dto.UserID,
		FirstName:         dto.FirstName,
		LastName:          dto.LastName,
		Email:             dto.Email,
		Status:            dto.Status,
		RegistrationType:  dto.RegistrationType,
		RoleCode:          dto.RoleCode,
		RoleName:          dto.RoleName,
		ParticipantNumber: dto.ParticipantNumber,
		IdentityNumber:    dto.IdentityNumber,
		OrganizationCode:  dto.OrganizationCode,
		CreatedAt:         dto.CreatedAt,
		UpdatedAt:         dto.UpdatedAt,
	}
}

func MapMemberDetailResponse(dto *member.MemberDetailResponse) response.MemberDetailResponse {
	return response.MemberDetailResponse{
		ID:                dto.ID,
		UserID:            dto.UserID,
		FirstName:         dto.FirstName,
		LastName:          dto.LastName,
		Email:             dto.Email,
		Status:            dto.Status,
		RegistrationType:  dto.RegistrationType,
		RoleCode:          dto.RoleCode,
		RoleName:          dto.RoleName,
		ParticipantNumber: dto.ParticipantNumber,
		IdentityNumber:    dto.IdentityNumber,
		OrganizationCode:  dto.OrganizationCode,
		ApprovedBy:        dto.ApprovedBy,
		ApprovedAt:        dto.ApprovedAt,
		Version:           dto.Version,
		CreatedAt:         dto.CreatedAt,
		UpdatedAt:         dto.UpdatedAt,
	}
}

func MapMemberListResponse(dto *member.ListResponse) response.MemberListResponse {
	members := make([]response.MemberListItemResponse, 0, len(dto.Members))
	for _, m := range dto.Members {
		members = append(members, response.MemberListItemResponse{
			ID:                m.ID,
			UserID:            m.UserID,
			FirstName:         m.FirstName,
			LastName:          m.LastName,
			Email:             m.Email,
			Status:            m.Status,
			RegistrationType:  m.RegistrationType,
			RoleCode:          m.RoleCode,
			RoleName:          m.RoleName,
			ParticipantNumber: m.ParticipantNumber,
			IdentityNumber:    m.IdentityNumber,
			OrganizationCode:  m.OrganizationCode,
			CreatedAt:         m.CreatedAt,
		})
	}

	totalPages := 0
	if dto.Pagination.PerPage > 0 {
		totalPages = int((dto.Pagination.Total + int64(dto.Pagination.PerPage) - 1) / int64(dto.Pagination.PerPage))
	}

	return response.MemberListResponse{
		Members: members,
		Pagination: response.Pagination{
			Page:       dto.Pagination.Page,
			Limit:      dto.Pagination.PerPage,
			Total:      dto.Pagination.Total,
			TotalPages: totalPages,
		},
	}
}
