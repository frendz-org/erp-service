package role

import (
	"context"

	"erp-service/config"
)

type Usecase interface {
	Create(ctx context.Context, req *CreateRequest) (*CreateResponse, error)
}

func NewUsecase(
	txManager TransactionManager,
	cfg *config.Config,
	tenantRepo TenantRepository,
	roleRepo RoleRepository,
	rolePermissionRepo RolePermissionRepository,
) Usecase {
	return newUsecase(txManager, cfg, tenantRepo, roleRepo, rolePermissionRepo)
}
