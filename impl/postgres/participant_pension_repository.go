package postgres

import (
	"context"

	"erp-service/entity"
	"erp-service/pkg/errors"
	"erp-service/saving/participant/contract"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type participantPensionRepository struct {
	baseRepository
}

func NewParticipantPensionRepository(db *gorm.DB) contract.ParticipantPensionRepository {
	return &participantPensionRepository{
		baseRepository: baseRepository{db: db},
	}
}

func (r *participantPensionRepository) Create(ctx context.Context, pension *entity.ParticipantPension) error {
	if err := r.getDB(ctx).Create(pension).Error; err != nil {
		return translateError(err, "participant pension")
	}
	return nil
}

func (r *participantPensionRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.ParticipantPension, error) {
	var pension entity.ParticipantPension
	err := r.getDB(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&pension).Error
	if err != nil {
		return nil, translateError(err, "participant pension")
	}
	return &pension, nil
}

func (r *participantPensionRepository) GetByParticipantID(ctx context.Context, participantID uuid.UUID) (*entity.ParticipantPension, error) {
	var pension entity.ParticipantPension
	err := r.getDB(ctx).Where("participant_id = ? AND deleted_at IS NULL", participantID).First(&pension).Error
	if err != nil {
		return nil, translateError(err, "participant pension")
	}
	return &pension, nil
}

func (r *participantPensionRepository) Update(ctx context.Context, pension *entity.ParticipantPension) error {
	oldVersion := pension.Version
	pension.Version = oldVersion + 1

	result := r.getDB(ctx).Where("version = ? AND deleted_at IS NULL", oldVersion).Save(pension)
	if result.Error != nil {
		pension.Version = oldVersion
		return translateError(result.Error, "participant pension")
	}
	if result.RowsAffected == 0 {
		pension.Version = oldVersion
		return errors.ErrConflict("participant pension was modified by another request")
	}
	return nil
}

func (r *participantPensionRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	result := r.getDB(ctx).Model(&entity.ParticipantPension{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Update("deleted_at", gorm.Expr("NOW()"))
	if result.Error != nil {
		return translateError(result.Error, "participant pension")
	}
	if result.RowsAffected == 0 {
		return errors.ErrNotFound("participant pension not found")
	}
	return nil
}
