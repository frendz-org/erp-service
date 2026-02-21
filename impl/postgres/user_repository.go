package postgres

import (
	"context"
	"fmt"
	"time"

	"erp-service/entity"
	"erp-service/iam/user/contract"
	"erp-service/pkg/errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userRepository struct {
	baseRepository
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{
		baseRepository: baseRepository{db: db},
	}
}

func (r *userRepository) Create(ctx context.Context, user *entity.User) error {
	if err := r.getDB(ctx).Create(user).Error; err != nil {
		return translateError(err, "user")
	}
	return nil
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	var user entity.User
	err := r.getDB(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&user).Error
	if err != nil {
		return nil, translateError(err, "user")
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	err := r.getDB(ctx).Where("email = ? AND deleted_at IS NULL", email).First(&user).Error
	if err != nil {
		return nil, translateError(err, "user")
	}
	return &user, nil
}

func (r *userRepository) Update(ctx context.Context, user *entity.User) error {
	if err := r.getDB(ctx).Save(user).Error; err != nil {
		return translateError(err, "user")
	}
	return nil
}

func (r *userRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	var count int64
	err := r.getDB(ctx).Model(&entity.User{}).
		Where("email = ? AND deleted_at IS NULL", email).
		Count(&count).Error
	if err != nil {
		return false, translateError(err, "user")
	}
	return count > 0, nil
}

func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	now := time.Now()
	result := r.getDB(ctx).Model(&entity.User{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(map[string]interface{}{
			"deleted_at":        now,
			"status":            string(entity.UserStatusLocked),
			"status_changed_at": now,
		})
	if result.Error != nil {
		return translateError(result.Error, "user")
	}
	if result.RowsAffected == 0 {
		return errors.ErrNotFound("user not found")
	}
	return nil
}

func (r *userRepository) List(ctx context.Context, filter *contract.UserListFilter) ([]*entity.User, int64, error) {
	var users []*entity.User
	var total int64

	query := r.getDB(ctx).Model(&entity.User{}).Where("deleted_at IS NULL")

	if filter.Status != nil {
		query = query.Where("status = ?", string(*filter.Status))
	}

	if filter.Search != "" {
		searchTerm := "%" + filter.Search + "%"
		query = query.Where(
			"email ILIKE ? OR id IN (SELECT user_id FROM user_profiles WHERE first_name ILIKE ? OR last_name ILIKE ?)",
			searchTerm, searchTerm, searchTerm,
		)
	}

	if filter.RoleID != nil {
		query = query.Where(
			"id IN (SELECT user_id FROM user_roles WHERE role_id = ? AND deleted_at IS NULL AND (effective_to IS NULL OR effective_to > NOW()))",
			*filter.RoleID,
		)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, translateError(err, "user")
	}

	validSortColumns := map[string]string{
		"created_at": "created_at",
		"updated_at": "updated_at",
		"email":      "email",
	}

	sortColumn := "created_at"
	if col, ok := validSortColumns[filter.SortBy]; ok {
		sortColumn = col
	}

	sortOrder := "DESC"
	if filter.SortOrder == "asc" {
		sortOrder = "ASC"
	}

	offset := (filter.Page - 1) * filter.PerPage
	err := query.
		Order(fmt.Sprintf("%s %s", sortColumn, sortOrder)).
		Offset(offset).
		Limit(filter.PerPage).
		Find(&users).Error
	if err != nil {
		return nil, 0, translateError(err, "user")
	}

	return users, total, nil
}
