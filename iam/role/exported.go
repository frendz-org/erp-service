package role

import (
	"iam-service/config"
	"iam-service/iam/role/contract"
	"iam-service/iam/role/internal"
)

type Usecase = contract.Usecase

func NewUsecase(
	txManager contract.TransactionManager,
	cfg *config.Config,
	tenantRepo contract.TenantRepository,
	roleRepo contract.RoleRepository,
	rolePermissionRepo contract.RolePermissionRepository,
) Usecase {
	return internal.NewUsecase(
		txManager,
		cfg,
		tenantRepo,
		roleRepo,
		rolePermissionRepo,
	)
}
