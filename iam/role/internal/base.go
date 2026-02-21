package internal

import (
	"erp-service/config"
	"erp-service/iam/role/contract"
)

type usecase struct {
	TxManager          contract.TransactionManager
	Config             *config.Config
	TenantRepo         contract.TenantRepository
	RoleRepo           contract.RoleRepository
	RolePermissionRepo contract.RolePermissionRepository
}

func NewUsecase(
	txManager contract.TransactionManager,
	cfg *config.Config,
	tenantRepo contract.TenantRepository,
	roleRepo contract.RoleRepository,
	rolePermissionRepo contract.RolePermissionRepository,
) *usecase {
	return &usecase{
		TxManager:          txManager,
		Config:             cfg,
		TenantRepo:         tenantRepo,
		RoleRepo:           roleRepo,
		RolePermissionRepo: rolePermissionRepo,
	}
}
