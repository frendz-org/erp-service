package postgres

import (
	"context"
	"time"

	"erp-service/entity"
	"erp-service/iam/auth"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type refreshTokenRepository struct {
	baseRepository
}

func NewRefreshTokenRepository(db *gorm.DB) auth.RefreshTokenRepository {
	return &refreshTokenRepository{
		baseRepository: baseRepository{db: db},
	}
}

func (r *refreshTokenRepository) Create(ctx context.Context, token *entity.RefreshToken) error {
	if err := r.getDB(ctx).Create(token).Error; err != nil {
		return translateError(err, "refresh token")
	}
	return nil
}

func (r *refreshTokenRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.RefreshToken, error) {
	var token entity.RefreshToken
	err := r.getDB(ctx).Where("id = ?", id).First(&token).Error
	if err != nil {
		return nil, translateError(err, "refresh token")
	}
	return &token, nil
}

func (r *refreshTokenRepository) GetByTokenHash(ctx context.Context, tokenHash string) (*entity.RefreshToken, error) {
	var token entity.RefreshToken
	err := r.getDB(ctx).Where("token_hash = ?", tokenHash).First(&token).Error
	if err != nil {
		return nil, translateError(err, "refresh token")
	}
	return &token, nil
}

func (r *refreshTokenRepository) SetReplacedBy(ctx context.Context, id uuid.UUID, replacedByID uuid.UUID) error {
	if err := r.getDB(ctx).
		Model(&entity.RefreshToken{}).
		Where("id = ?", id).
		Update("replaced_by_token_id", replacedByID).Error; err != nil {
		return translateError(err, "refresh token")
	}
	return nil
}

func (r *refreshTokenRepository) Revoke(ctx context.Context, id uuid.UUID, reason string) error {
	now := time.Now()
	if err := r.getDB(ctx).
		Model(&entity.RefreshToken{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"revoked_at":     now,
			"revoked_reason": reason,
		}).Error; err != nil {
		return translateError(err, "refresh token")
	}
	return nil
}

func (r *refreshTokenRepository) RevokeAllByUserID(ctx context.Context, userID uuid.UUID, reason string) error {
	now := time.Now()
	if err := r.getDB(ctx).
		Model(&entity.RefreshToken{}).
		Where("user_id = ? AND revoked_at IS NULL", userID).
		Updates(map[string]interface{}{
			"revoked_at":     now,
			"revoked_reason": reason,
		}).Error; err != nil {
		return translateError(err, "refresh token")
	}
	return nil
}

func (r *refreshTokenRepository) RevokeByFamily(ctx context.Context, tokenFamily uuid.UUID, reason string) error {
	now := time.Now()
	if err := r.getDB(ctx).
		Model(&entity.RefreshToken{}).
		Where("token_family = ? AND revoked_at IS NULL", tokenFamily).
		Updates(map[string]interface{}{
			"revoked_at":     now,
			"revoked_reason": reason,
		}).Error; err != nil {
		return translateError(err, "refresh token")
	}
	return nil
}
