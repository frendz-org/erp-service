package postgres

import (
	"context"
	"time"

	"iam-service/entity"
	"iam-service/iam/auth/contract"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userSessionRepository struct {
	baseRepository
}

func NewUserSessionRepository(db *gorm.DB) contract.UserSessionRepository {
	return &userSessionRepository{
		baseRepository: baseRepository{db: db},
	}
}

func (r *userSessionRepository) Create(ctx context.Context, session *entity.UserSession) error {
	if err := r.getDB(ctx).Create(session).Error; err != nil {
		return translateError(err, "user session")
	}
	return nil
}

func (r *userSessionRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.UserSession, error) {
	var session entity.UserSession
	err := r.getDB(ctx).Where("id = ?", id).First(&session).Error
	if err != nil {
		return nil, translateError(err, "user session")
	}
	return &session, nil
}

func (r *userSessionRepository) UpdateLastActive(ctx context.Context, id uuid.UUID) error {
	if err := r.getDB(ctx).
		Model(&entity.UserSession{}).
		Where("id = ?", id).
		Update("last_active_at", time.Now()).Error; err != nil {
		return translateError(err, "user session")
	}
	return nil
}

func (r *userSessionRepository) GetByRefreshTokenID(ctx context.Context, refreshTokenID uuid.UUID) (*entity.UserSession, error) {
	var session entity.UserSession
	err := r.getDB(ctx).
		Where("refresh_token_id = ? AND status = ?", refreshTokenID, entity.UserSessionStatusActive).
		First(&session).Error
	if err != nil {
		return nil, translateError(err, "user session")
	}
	return &session, nil
}

func (r *userSessionRepository) UpdateRefreshTokenID(ctx context.Context, sessionID uuid.UUID, refreshTokenID uuid.UUID) error {
	if err := r.getDB(ctx).
		Model(&entity.UserSession{}).
		Where("id = ?", sessionID).
		Update("refresh_token_id", refreshTokenID).Error; err != nil {
		return translateError(err, "user session")
	}
	return nil
}

func (r *userSessionRepository) Revoke(ctx context.Context, id uuid.UUID) error {
	now := time.Now()
	if err := r.getDB(ctx).
		Model(&entity.UserSession{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":     entity.UserSessionStatusRevoked,
			"revoked_at": now,
		}).Error; err != nil {
		return translateError(err, "user session")
	}
	return nil
}

func (r *userSessionRepository) RevokeAllByUserID(ctx context.Context, userID uuid.UUID) error {
	now := time.Now()
	if err := r.getDB(ctx).
		Model(&entity.UserSession{}).
		Where("user_id = ? AND status = ?", userID, entity.UserSessionStatusActive).
		Updates(map[string]interface{}{
			"status":     entity.UserSessionStatusRevoked,
			"revoked_at": now,
		}).Error; err != nil {
		return translateError(err, "user session")
	}
	return nil
}
