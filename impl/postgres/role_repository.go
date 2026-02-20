package postgres

import (
	"context"

	"iam-service/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type roleRepository struct {
	baseRepository
}

func NewRoleRepository(db *gorm.DB) *roleRepository {
	return &roleRepository{
		baseRepository: baseRepository{db: db},
	}
}

func (r *roleRepository) Create(ctx context.Context, role *entity.Role) error {
	if err := r.getDB(ctx).Create(role).Error; err != nil {
		return translateError(err, "role")
	}
	return nil
}

func (r *roleRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Role, error) {
	var role entity.Role
	err := r.getDB(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&role).Error
	if err != nil {
		return nil, translateError(err, "role")
	}
	return &role, nil
}

func (r *roleRepository) GetByName(ctx context.Context, productID uuid.UUID, name string) (*entity.Role, error) {
	var role entity.Role
	err := r.getDB(ctx).Where("product_id = ? AND name = ? AND deleted_at IS NULL", productID, name).First(&role).Error
	if err != nil {
		return nil, translateError(err, "role")
	}
	return &role, nil
}

func (r *roleRepository) GetByCode(ctx context.Context, productID uuid.UUID, code string) (*entity.Role, error) {
	var role entity.Role
	err := r.getDB(ctx).Where("product_id = ? AND code = ? AND deleted_at IS NULL", productID, code).First(&role).Error
	if err != nil {
		return nil, translateError(err, "role")
	}
	return &role, nil
}

func (r *roleRepository) Update(ctx context.Context, role *entity.Role) error {
	if err := r.getDB(ctx).Save(role).Error; err != nil {
		return translateError(err, "role")
	}
	return nil
}

func (r *roleRepository) GetByIDs(ctx context.Context, ids []uuid.UUID) ([]*entity.Role, error) {
	var roles []*entity.Role
	err := r.getDB(ctx).Where("id IN ? AND status = 'ACTIVE' AND deleted_at IS NULL", ids).Find(&roles).Error
	if err != nil {
		return nil, translateError(err, "roles")
	}
	return roles, nil
}

func (r *roleRepository) GetByCodeAndProduct(ctx context.Context, productID uuid.UUID, code string) (*entity.Role, error) {
	var role entity.Role
	err := r.getDB(ctx).Table("roles").
		Where("product_id = ? AND code = ? AND status = ? AND deleted_at IS NULL", productID, code, "ACTIVE").
		First(&role).Error
	if err != nil {
		return nil, translateError(err, "role")
	}
	return &role, nil
}
