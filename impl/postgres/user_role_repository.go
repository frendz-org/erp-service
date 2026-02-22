package postgres

import (
	"context"
	"time"

	"erp-service/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userRoleRepository struct {
	baseRepository
}

func NewUserRoleRepository(db *gorm.DB) *userRoleRepository {
	return &userRoleRepository{
		baseRepository: baseRepository{db: db},
	}
}

func (r *userRoleRepository) Create(ctx context.Context, userRole *entity.UserRole) error {
	if err := r.getDB(ctx).Create(userRole).Error; err != nil {
		return translateError(err, "user role")
	}
	return nil
}

func (r *userRoleRepository) ListActiveByUserID(ctx context.Context, userID uuid.UUID, productID *uuid.UUID) ([]entity.UserRole, error) {
	var userRoles []entity.UserRole
	now := time.Now()

	query := r.getDB(ctx).Where("user_id = ? AND deleted_at IS NULL", userID).
		Where("assigned_at <= ?", now).
		Where("expires_at IS NULL OR expires_at > ?", now)

	if productID != nil {
		query = query.Where("product_id = ? OR product_id IS NULL", *productID)
	}

	if err := query.Find(&userRoles).Error; err != nil {
		return nil, translateError(err, "user roles")
	}
	return userRoles, nil
}

func (r *userRoleRepository) SoftDeleteByUserAndProduct(ctx context.Context, userID, productID uuid.UUID) error {
	now := time.Now()
	result := r.getDB(ctx).
		Model(&entity.UserRole{}).
		Where("user_id = ? AND product_id = ? AND deleted_at IS NULL", userID, productID).
		Update("deleted_at", now)

	if result.Error != nil {
		return translateError(result.Error, "user role")
	}
	return nil
}

func (r *userRoleRepository) GetActiveByUserAndProduct(ctx context.Context, userID, productID uuid.UUID) (*entity.UserRole, error) {
	var userRole entity.UserRole
	now := time.Now()
	err := r.getDB(ctx).
		Where("user_id = ? AND product_id = ? AND deleted_at IS NULL", userID, productID).
		Where("assigned_at <= ?", now).
		Where("expires_at IS NULL OR expires_at > ?", now).
		First(&userRole).Error
	if err != nil {
		return nil, translateError(err, "user role")
	}
	return &userRole, nil
}
