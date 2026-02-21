package postgres

import (
	"context"
	"time"

	"erp-service/entity"
	"erp-service/iam/user/contract"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type emailVerificationRepository struct {
	baseRepository
}

func NewEmailVerificationRepository(db *gorm.DB) contract.EmailVerificationRepository {
	return &emailVerificationRepository{
		baseRepository: baseRepository{db: db},
	}
}

func (r *emailVerificationRepository) Create(ctx context.Context, verification *entity.EmailVerification) error {
	if err := r.getDB(ctx).Create(verification).Error; err != nil {
		return translateError(err, "email verification")
	}
	return nil
}

func (r *emailVerificationRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.EmailVerification, error) {
	var verification entity.EmailVerification
	err := r.getDB(ctx).Where("id = ?", id).First(&verification).Error
	if err != nil {
		return nil, translateError(err, "email verification")
	}
	return &verification, nil
}

func (r *emailVerificationRepository) GetLatestByEmail(ctx context.Context, email string, otpType entity.OTPType) (*entity.EmailVerification, error) {
	var verification entity.EmailVerification
	err := r.getDB(ctx).
		Where("email = ? AND otp_type = ? AND verified_at IS NULL AND expires_at > ?", email, otpType, time.Now()).
		Order("created_at DESC").
		First(&verification).Error
	if err != nil {
		return nil, translateError(err, "email verification")
	}
	return &verification, nil
}

func (r *emailVerificationRepository) GetLatestByUserID(ctx context.Context, userID uuid.UUID, otpType entity.OTPType) (*entity.EmailVerification, error) {
	var verification entity.EmailVerification
	err := r.getDB(ctx).
		Where("user_id = ? AND otp_type = ? AND verified_at IS NULL AND expires_at > ?", userID, otpType, time.Now()).
		Order("created_at DESC").
		First(&verification).Error
	if err != nil {
		return nil, translateError(err, "email verification")
	}
	return &verification, nil
}

func (r *emailVerificationRepository) MarkAsVerified(ctx context.Context, id uuid.UUID) error {
	now := time.Now()
	if err := r.getDB(ctx).
		Model(&entity.EmailVerification{}).
		Where("id = ?", id).
		Update("verified_at", now).Error; err != nil {
		return translateError(err, "email verification")
	}
	return nil
}

func (r *emailVerificationRepository) CountActiveOTPsByEmail(ctx context.Context, email string, otpType entity.OTPType) (int, error) {
	var count int64
	err := r.getDB(ctx).
		Model(&entity.EmailVerification{}).
		Where("email = ? AND otp_type = ? AND verified_at IS NULL AND expires_at > ?", email, otpType, time.Now()).
		Count(&count).Error
	if err != nil {
		return 0, translateError(err, "email verification")
	}
	return int(count), nil
}

func (r *emailVerificationRepository) DeleteExpiredByEmail(ctx context.Context, email string) error {
	if err := r.getDB(ctx).
		Where("email = ? AND expires_at < ?", email, time.Now()).
		Delete(&entity.EmailVerification{}).Error; err != nil {
		return translateError(err, "email verification")
	}
	return nil
}
