package postgres

import (
	"context"

	"iam-service/entity"

	"gorm.io/gorm"
)

type rolePermissionRepository struct {
	baseRepository
}

func NewRolePermissionRepository(db *gorm.DB) *rolePermissionRepository {
	return &rolePermissionRepository{
		baseRepository: baseRepository{db: db},
	}
}

func (r *rolePermissionRepository) Create(ctx context.Context, rolePermission *entity.RolePermission) error {
	if err := r.getDB(ctx).Create(rolePermission).Error; err != nil {
		return translateError(err, "role permission")
	}
	return nil
}
