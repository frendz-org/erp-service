package postgres

import (
	"context"

	"erp-service/entity"
	"erp-service/saving/participant"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type participantStatusHistoryRepository struct {
	baseRepository
}

func NewParticipantStatusHistoryRepository(db *gorm.DB) participant.ParticipantStatusHistoryRepository {
	return &participantStatusHistoryRepository{
		baseRepository: baseRepository{db: db},
	}
}

func (r *participantStatusHistoryRepository) Create(ctx context.Context, history *entity.ParticipantStatusHistory) error {
	if err := r.getDB(ctx).Create(history).Error; err != nil {
		return translateError(err, "participant status history")
	}
	return nil
}

func (r *participantStatusHistoryRepository) ListByParticipantID(ctx context.Context, participantID uuid.UUID) ([]*entity.ParticipantStatusHistory, error) {
	var histories []*entity.ParticipantStatusHistory
	err := r.getDB(ctx).
		Where("participant_id = ?", participantID).
		Order("changed_at DESC").
		Find(&histories).Error
	if err != nil {
		return nil, translateError(err, "participant status history")
	}
	return histories, nil
}
