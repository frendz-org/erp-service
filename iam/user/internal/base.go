package internal

import (
	"erp-service/config"
	"erp-service/entity"
	"erp-service/iam/user/contract"
	"erp-service/iam/user/userdto"
)

type usecase struct {
	TxManager             contract.TransactionManager
	Config                *config.Config
	UserRepo              contract.UserRepository
	UserProfileRepo       contract.UserProfileRepository
	UserAuthMethodRepo    contract.UserAuthMethodRepository
	UserSecurityStateRepo contract.UserSecurityStateRepository
	TenantRepo            contract.TenantRepository
	RoleRepo              contract.RoleRepository
	UserRoleRepo          contract.UserRoleRepository
}

func NewUsecase(
	txManager contract.TransactionManager,
	cfg *config.Config,
	userRepo contract.UserRepository,
	userProfileRepo contract.UserProfileRepository,
	userAuthMethodRepo contract.UserAuthMethodRepository,
	userSecurityStateRepo contract.UserSecurityStateRepository,
	tenantRepo contract.TenantRepository,
	roleRepo contract.RoleRepository,
	userRoleRepo contract.UserRoleRepository,
) *usecase {
	return &usecase{
		TxManager:             txManager,
		Config:                cfg,
		UserRepo:              userRepo,
		UserProfileRepo:       userProfileRepo,
		UserAuthMethodRepo:    userAuthMethodRepo,
		UserSecurityStateRepo: userSecurityStateRepo,
		TenantRepo:            tenantRepo,
		RoleRepo:              roleRepo,
		UserRoleRepo:          userRoleRepo,
	}
}

func mapUserToDetailResponse(user *entity.User, profile *entity.UserProfile, authMethod *entity.UserAuthMethod, securityState *entity.UserSecurityState) *userdto.UserDetailResponse {
	resp := &userdto.UserDetailResponse{
		ID:        user.ID,
		Email:     user.Email,
		Status:    string(user.Status),
		IsActive:  user.IsActive(),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	if securityState != nil {
		resp.EmailVerified = securityState.EmailVerified
		resp.LastLoginAt = securityState.LastLoginAt
	}

	if profile != nil {
		resp.FirstName = profile.FirstName
		resp.LastName = profile.LastName
		resp.FullName = profile.FullName()
		resp.PhoneNumber = profile.PhoneNumber
		resp.Address = profile.Address
		resp.ProfilePictureURL = profile.ProfilePictureURL
		if profile.DateOfBirth != nil {
			formatted := profile.DateOfBirth.Format("2006-01-02")
			resp.DateOfBirth = &formatted
		}
	}

	if authMethod != nil {
		resp.PINSet = authMethod.MethodType == string(entity.AuthMethodPIN)
	}

	return resp
}

func mapUserToListItem(user *entity.User, profile *entity.UserProfile, securityState *entity.UserSecurityState) userdto.UserListItem {
	item := userdto.UserListItem{
		ID:        user.ID,
		Email:     user.Email,
		Status:    string(user.Status),
		IsActive:  user.IsActive(),
		CreatedAt: user.CreatedAt,
	}

	if securityState != nil {
		item.EmailVerified = securityState.EmailVerified
		item.LastLoginAt = securityState.LastLoginAt
	}

	if profile != nil {
		item.FirstName = profile.FirstName
		item.LastName = profile.LastName
		item.FullName = profile.FullName()
		item.PhoneNumber = profile.PhoneNumber
	}

	return item
}
