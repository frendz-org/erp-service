package postgres

import (
	"context"

	"erp-service/entity"
	"erp-service/iam/auth"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userAuthMethodRepository struct {
	baseRepository
}

func NewUserAuthMethodRepository(db *gorm.DB) auth.UserAuthMethodRepository {
	return &userAuthMethodRepository{
		baseRepository: baseRepository{db: db},
	}
}

func (r *userAuthMethodRepository) Create(ctx context.Context, authMethod *entity.UserAuthMethod) error {
	if err := r.getDB(ctx).Create(authMethod).Error; err != nil {
		return translateError(err, "user auth method")
	}
	return nil
}

func (r *userAuthMethodRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*entity.UserAuthMethod, error) {
	var authMethod entity.UserAuthMethod
	err := r.getDB(ctx).Where("user_id = ? AND is_active = true", userID).First(&authMethod).Error
	if err != nil {
		return nil, translateError(err, "user auth method")
	}
	return &authMethod, nil
}

func (r *userAuthMethodRepository) Update(ctx context.Context, authMethod *entity.UserAuthMethod) error {
	if err := r.getDB(ctx).Save(authMethod).Error; err != nil {
		return translateError(err, "user auth method")
	}
	return nil
}
