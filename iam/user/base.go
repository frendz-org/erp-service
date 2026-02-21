package user

import (
	"erp-service/config"
	"erp-service/entity"
)

type usecase struct {
	TxManager             TransactionManager
	Config                *config.Config
	UserRepo              UserRepository
	UserProfileRepo       UserProfileRepository
	UserAuthMethodRepo    UserAuthMethodRepository
	UserSecurityStateRepo UserSecurityStateRepository
	TenantRepo            TenantRepository
	RoleRepo              RoleRepository
	UserRoleRepo          UserRoleRepository
}

func newUsecase(
	txManager TransactionManager,
	cfg *config.Config,
	userRepo UserRepository,
	userProfileRepo UserProfileRepository,
	userAuthMethodRepo UserAuthMethodRepository,
	userSecurityStateRepo UserSecurityStateRepository,
	tenantRepo TenantRepository,
	roleRepo RoleRepository,
	userRoleRepo UserRoleRepository,
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

func mapUserToDetailResponse(user *entity.User, profile *entity.UserProfile, authMethod *entity.UserAuthMethod, securityState *entity.UserSecurityState) *UserDetailResponse {
	resp := &UserDetailResponse{
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

func mapUserToListItem(user *entity.User, profile *entity.UserProfile, securityState *entity.UserSecurityState) UserListItem {
	item := UserListItem{
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
