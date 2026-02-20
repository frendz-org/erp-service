package postgres

import (
	"context"

	"iam-service/entity"
	"iam-service/iam/auth/contract"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userProfileRepository struct {
	baseRepository
}

func NewUserProfileRepository(db *gorm.DB) contract.UserProfileRepository {
	return &userProfileRepository{
		baseRepository: baseRepository{db: db},
	}
}

func (r *userProfileRepository) Create(ctx context.Context, profile *entity.UserProfile) error {
	if err := r.getDB(ctx).Create(profile).Error; err != nil {
		return translateError(err, "user profile")
	}
	return nil
}

func (r *userProfileRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*entity.UserProfile, error) {
	var profile entity.UserProfile
	err := r.getDB(ctx).Where("user_id = ?", userID).First(&profile).Error
	if err != nil {
		return nil, translateError(err, "user profile")
	}
	return &profile, nil
}

func (r *userProfileRepository) Update(ctx context.Context, profile *entity.UserProfile) error {
	if err := r.getDB(ctx).Save(profile).Error; err != nil {
		return translateError(err, "user profile")
	}
	return nil
}
