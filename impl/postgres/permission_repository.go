package postgres

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type permissionRepository struct {
	baseRepository
}

func NewPermissionRepository(db *gorm.DB) *permissionRepository {
	return &permissionRepository{
		baseRepository: baseRepository{db: db},
	}
}

func (r *permissionRepository) GetCodesByRoleIDs(ctx context.Context, roleIDs []uuid.UUID) ([]string, error) {
	if len(roleIDs) == 0 {
		return []string{}, nil
	}

	type permissionResult struct {
		Code string
	}

	var results []permissionResult
	err := r.getDB(ctx).Raw(`
		SELECT DISTINCT p.code
		FROM role_permissions rp
		INNER JOIN permissions p ON p.id = rp.permission_id
		WHERE rp.role_id IN ? AND p.deleted_at IS NULL
		ORDER BY p.code
	`, roleIDs).Scan(&results).Error

	if err != nil {
		return nil, translateError(err, "permissions")
	}

	permissions := make([]string, len(results))
	for i, r := range results {
		permissions[i] = r.Code
	}

	return permissions, nil
}
