package role

import "erp-service/config"

type usecase struct {
	TxManager          TransactionManager
	Config             *config.Config
	TenantRepo         TenantRepository
	RoleRepo           RoleRepository
	RolePermissionRepo RolePermissionRepository
}

func newUsecase(
	txManager TransactionManager,
	cfg *config.Config,
	tenantRepo TenantRepository,
	roleRepo RoleRepository,
	rolePermissionRepo RolePermissionRepository,
) *usecase {
	return &usecase{
		TxManager:          txManager,
		Config:             cfg,
		TenantRepo:         tenantRepo,
		RoleRepo:           roleRepo,
		RolePermissionRepo: rolePermissionRepo,
	}
}
