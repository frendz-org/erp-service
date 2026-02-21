package member

import "erp-service/config"

type usecase struct {
	cfg         *config.Config
	txManager   TransactionManager
	utrRepo     UserTenantRegistrationRepository
	userRole    UserRoleRepository
	productRepo ProductRepository
	roleRepo    RoleRepository
	configRepo  ProductRegistrationConfigRepository
	profileRepo UserProfileRepository
	userRepo    UserRepository
}

func newUsecase(
	cfg *config.Config,
	txManager TransactionManager,
	utrRepo UserTenantRegistrationRepository,
	userRole UserRoleRepository,
	productRepo ProductRepository,
	roleRepo RoleRepository,
	configRepo ProductRegistrationConfigRepository,
	profileRepo UserProfileRepository,
	userRepo UserRepository,
) Usecase {
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
