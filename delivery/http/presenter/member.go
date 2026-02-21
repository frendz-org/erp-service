package presenter

import (
	"erp-service/delivery/http/dto/response"
	"erp-service/saving/member/memberdto"
)

func MapMemberRegisterResponse(dto *memberdto.RegisterResponse) response.MemberRegisterResponse {
	return response.MemberRegisterResponse{
		ID:               dto.ID,
		Status:           dto.Status,
		RegistrationType: dto.RegistrationType,
		CreatedAt:        dto.CreatedAt,
	}
}

func MapMemberDetailResponse(dto *memberdto.MemberDetailResponse) response.MemberDetailResponse {
	return response.MemberDetailResponse{
		ID:               dto.ID,
		UserID:           dto.UserID,
		FirstName:        dto.FirstName,
		LastName:         dto.LastName,
		Email:            dto.Email,
		Status:           dto.Status,
		RegistrationType: dto.RegistrationType,
		RoleCode:         dto.RoleCode,
		RoleName:         dto.RoleName,
		ApprovedBy:       dto.ApprovedBy,
		ApprovedAt:       dto.ApprovedAt,
		Version:          dto.Version,
		CreatedAt:        dto.CreatedAt,
		UpdatedAt:        dto.UpdatedAt,
	}
}

func MapMemberListResponse(dto *memberdto.ListResponse) response.MemberListResponse {
	members := make([]response.MemberListItemResponse, 0, len(dto.Members))
	for _, m := range dto.Members {
		members = append(members, response.MemberListItemResponse{
			ID:               m.ID,
			UserID:           m.UserID,
			FirstName:        m.FirstName,
			LastName:         m.LastName,
			Email:            m.Email,
			Status:           m.Status,
			RegistrationType: m.RegistrationType,
			RoleCode:         m.RoleCode,
			RoleName:         m.RoleName,
			CreatedAt:        m.CreatedAt,
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
