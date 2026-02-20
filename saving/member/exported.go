package member

import (
	"iam-service/config"
	"iam-service/saving/member/contract"
	"iam-service/saving/member/internal"
)

type Usecase = contract.Usecase

func NewUsecase(
	cfg *config.Config,
	txManager contract.TransactionManager,
	utrRepo contract.UserTenantRegistrationRepository,
	userRoleRepo contract.UserRoleRepository,
	productRepo contract.ProductRepository,
	roleRepo contract.RoleRepository,
	configRepo contract.ProductRegistrationConfigRepository,
	profileRepo contract.UserProfileRepository,
	userRepo contract.UserRepository,
) Usecase {
	return internal.NewUsecase(
		cfg,
		txManager,
		utrRepo,
		userRoleRepo,
		productRepo,
		roleRepo,
		configRepo,
		profileRepo,
		userRepo,
	)
}
