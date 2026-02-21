package postgres

import (
	"context"
	"time"

	"erp-service/entity"
	"erp-service/iam/user/contract"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type pinVerificationLogRepository struct {
	baseRepository
}

func NewPINVerificationLogRepository(db *gorm.DB) contract.PINVerificationLogRepository {
	return &pinVerificationLogRepository{
		baseRepository: baseRepository{db: db},
	}
}

func (r *pinVerificationLogRepository) Create(ctx context.Context, log *entity.PINVerificationLog) error {
	if err := r.getDB(ctx).Create(log).Error; err != nil {
		return translateError(err, "PIN verification log")
	}
	return nil
}

func (r *pinVerificationLogRepository) CountRecentFailures(ctx context.Context, userID uuid.UUID, since int) (int, error) {
	var count int64
	sinceTime := time.Now().Add(-time.Duration(since) * time.Minute)
	err := r.getDB(ctx).
		Model(&entity.PINVerificationLog{}).
		Where("user_id = ? AND result = false AND created_at > ?", userID, sinceTime).
		Count(&count).Error
	if err != nil {
		return 0, translateError(err, "PIN verification log")
	}
	return int(count), nil
}
