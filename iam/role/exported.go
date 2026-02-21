package role

import (
	"erp-service/config"
	"erp-service/iam/role/contract"
	"erp-service/iam/role/internal"
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
