package internal

import (
	"iam-service/config"
	"iam-service/saving/member/contract"
)

type usecase struct {
	cfg         *config.Config
	txManager   contract.TransactionManager
	utrRepo     contract.UserTenantRegistrationRepository
	userRole    contract.UserRoleRepository
	productRepo contract.ProductRepository
	roleRepo    contract.RoleRepository
	configRepo  contract.ProductRegistrationConfigRepository
	profileRepo contract.UserProfileRepository
	userRepo    contract.UserRepository
}

func NewUsecase(
	cfg *config.Config,
	txManager contract.TransactionManager,
	utrRepo contract.UserTenantRegistrationRepository,
	userRole contract.UserRoleRepository,
	productRepo contract.ProductRepository,
	roleRepo contract.RoleRepository,
	configRepo contract.ProductRegistrationConfigRepository,
	profileRepo contract.UserProfileRepository,
	userRepo contract.UserRepository,
) contract.Usecase {
	return &usecase{
		cfg:         cfg,
		txManager:   txManager,
		utrRepo:     utrRepo,
		userRole:    userRole,
		productRepo: productRepo,
		roleRepo:    roleRepo,
		configRepo:  configRepo,
		profileRepo: profileRepo,
		userRepo:    userRepo,
	}
}
