package presenter

import (
	"erp-service/delivery/http/dto/response"
	"erp-service/iam/user/userdto"
)

func ToUserResponse(resp *userdto.UserDetailResponse) *response.UserResponse {
	if resp == nil {
		return nil
	}

	roles := make([]response.UserRoleResponse, len(resp.Roles))
	for i, r := range resp.Roles {
		roles[i] = response.UserRoleResponse{
			Code: r.Code,
			Name: r.Name,
		}
	}

	return &response.UserResponse{
		ID:          resp.ID,
		Email:       resp.Email,
		FirstName:   resp.FirstName,
		LastName:    resp.LastName,
		FullName:    resp.FullName,
		PhoneNumber: resp.PhoneNumber,
		DateOfBirth: resp.DateOfBirth,
		Address:     resp.Address,
		Status:      resp.Status,
		IsActive:    resp.IsActive,
		Roles:       roles,
	}
}

func ToUserListItemResponse(item *userdto.UserListItem) *response.UserListItemResponse {
	if item == nil {
		return nil
	}
	return &response.UserListItemResponse{
		ID:          item.ID,
		Email:       item.Email,
		FullName:    item.FullName,
		PhoneNumber: item.PhoneNumber,
		Status:      item.Status,
		IsActive:    item.IsActive,
		CreatedAt:   item.CreatedAt,
	}
}

func ToUserListResponse(items []userdto.UserListItem) []*response.UserListItemResponse {
	if items == nil {
		return nil
	}
	result := make([]*response.UserListItemResponse, len(items))
	for i := range items {
		result[i] = ToUserListItemResponse(&items[i])
	}
	return result
}

func ToCreateUserResponse(resp *userdto.CreateResponse) *response.CreateUserResponse {
	if resp == nil {
		return nil
	}
	return &response.CreateUserResponse{
		UserID:   resp.UserID,
		Email:    resp.Email,
		FullName: resp.FullName,
		RoleCode: resp.RoleCode,
	}
}

func ToUpdateUserResponse(resp *userdto.UpdateResponse) *response.UpdateUserResponse {
	if resp == nil {
		return nil
	}

	roles := make([]response.UserRoleResponse, len(resp.Roles))
	for i, r := range resp.Roles {
		roles[i] = response.UserRoleResponse{
			Code: r.Code,
			Name: r.Name,
		}
	}

	return &response.UpdateUserResponse{
		UserResponse: response.UserResponse{
			ID:          resp.ID,
			Email:       resp.Email,
			FirstName:   resp.FirstName,
			LastName:    resp.LastName,
			FullName:    resp.FullName,
			PhoneNumber: resp.PhoneNumber,
			DateOfBirth: resp.DateOfBirth,
			Address:     resp.Address,
			Status:      resp.Status,
			IsActive:    resp.IsActive,
			Roles:       roles,
		},
		Message: resp.Message,
	}
}

func ToApproveUserResponse(resp *userdto.ApproveResponse) *response.ApproveUserResponse {
	if resp == nil {
		return nil
	}
	return &response.ApproveUserResponse{
		UserID:  resp.UserID,
		Message: resp.Message,
	}
}

func ToRejectUserResponse(resp *userdto.RejectResponse) *response.RejectUserResponse {
	if resp == nil {
		return nil
	}
	return &response.RejectUserResponse{
		UserID:  resp.UserID,
		Message: resp.Message,
	}
}

func ToUnlockUserResponse(resp *userdto.UnlockResponse) *response.UnlockUserResponse {
	if resp == nil {
		return nil
	}
	return &response.UnlockUserResponse{
		UserID:  resp.UserID,
		Message: resp.Message,
	}
}

func ToResetUserPINResponse(resp *userdto.ResetPINResponse) *response.ResetUserPINResponse {
	if resp == nil {
		return nil
	}
	return &response.ResetUserPINResponse{
		UserID:  resp.UserID,
		Message: resp.Message,
	}
}
