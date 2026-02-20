package internal

import (
	"context"
	"time"

	"iam-service/entity"
	"iam-service/iam/role/roledto"
	"iam-service/pkg/errors"
)

func (uc *usecase) Create(ctx context.Context, req *roledto.CreateRequest) (*roledto.CreateResponse, error) {
	tenant, err := uc.TenantRepo.GetByID(ctx, req.TenantID)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, errors.ErrTenantNotFound()
		}
		return nil, errors.ErrInternal("failed to verify tenant").WithError(err)
	}
	if !tenant.IsActive() {
		return nil, errors.ErrTenantInactive()
	}

	existingRole, err := uc.RoleRepo.GetByCode(ctx, req.TenantID, req.Code)
	if err != nil && !errors.IsNotFound(err) {
		return nil, errors.ErrInternal("failed to check role existence").WithError(err)
	}
	if existingRole != nil {
		return nil, errors.ErrConflict("Role with this code already exists in the tenant")
	}

	scopeLevel := entity.ScopeLevel(req.ScopeLevel)
	if scopeLevel != entity.ScopeLevelSystem && scopeLevel != entity.ScopeLevelTenant &&
		scopeLevel != entity.ScopeLevelBranch && scopeLevel != entity.ScopeLevelSelf {
		return nil, errors.ErrValidation("invalid scope level")
	}

	now := time.Now()
	role := &entity.Role{
		ProductID: req.TenantID,
		Code:      req.Code,
		Name:          req.Name,
		Description:   req.Description,
		IsSystem:      false,
		Status:        "ACTIVE",
	}

	err = uc.TxManager.WithTransaction(ctx, func(txCtx context.Context) error {
		if err := uc.RoleRepo.Create(txCtx, role); err != nil {
			return err
		}

		if len(req.Permissions) > 0 {
			for _, permissionID := range req.Permissions {
				rolePermission := &entity.RolePermission{
					RoleID:       role.ID,
					PermissionID: permissionID,
					CreatedAt:    now,
				}

				if err := uc.RolePermissionRepo.Create(txCtx, rolePermission); err != nil {
					return err
				}
			}
		}

		return nil
	})

	if err != nil {
		return nil, errors.ErrInternal("failed to create role").WithError(err)
	}

	response := &roledto.CreateResponse{
		RoleID:      role.ID,
		TenantID:    req.TenantID,
		Code:        role.Code,
		Name:        role.Name,
		Description: role.Description,
		ScopeLevel:  req.ScopeLevel,
		IsSystem:    role.IsSystem,
		IsActive:    role.IsActive(),
		CreatedAt:   role.CreatedAt,
	}

	return response, nil
}
