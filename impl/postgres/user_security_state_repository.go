package postgres

import (
	"context"

	"iam-service/entity"
	"iam-service/iam/auth/contract"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userSecurityStateRepository struct {
	baseRepository
}

func NewUserSecurityStateRepository(db *gorm.DB) contract.UserSecurityStateRepository {
	return &userSecurityStateRepository{
		baseRepository: baseRepository{db: db},
	}
}

func (r *userSecurityStateRepository) Create(ctx context.Context, securityState *entity.UserSecurityState) error {
	if err := r.getDB(ctx).Create(securityState).Error; err != nil {
		return translateError(err, "user security state")
	}
	return nil
}

func (r *userSecurityStateRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*entity.UserSecurityState, error) {
	var securityState entity.UserSecurityState
	err := r.getDB(ctx).Where("user_id = ?", userID).First(&securityState).Error
	if err != nil {
		return nil, translateError(err, "user security state")
	}
	return &securityState, nil
}

func (r *userSecurityStateRepository) Update(ctx context.Context, securityState *entity.UserSecurityState) error {
	if err := r.getDB(ctx).Save(securityState).Error; err != nil {
		return translateError(err, "user security state")
	}
	return nil
}
