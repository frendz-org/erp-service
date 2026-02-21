package postgres

import (
	"context"

	"erp-service/entity"
	"erp-service/pkg/errors"
	"erp-service/saving/participant/contract"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type participantEmploymentRepository struct {
	baseRepository
}

func NewParticipantEmploymentRepository(db *gorm.DB) contract.ParticipantEmploymentRepository {
	return &participantEmploymentRepository{
		baseRepository: baseRepository{db: db},
	}
}

func (r *participantEmploymentRepository) Create(ctx context.Context, employment *entity.ParticipantEmployment) error {
	if err := r.getDB(ctx).Create(employment).Error; err != nil {
		return translateError(err, "participant employment")
	}
	return nil
}

func (r *participantEmploymentRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.ParticipantEmployment, error) {
	var employment entity.ParticipantEmployment
	err := r.getDB(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&employment).Error
	if err != nil {
		return nil, translateError(err, "participant employment")
	}
	return &employment, nil
}

func (r *participantEmploymentRepository) GetByParticipantID(ctx context.Context, participantID uuid.UUID) (*entity.ParticipantEmployment, error) {
	var employment entity.ParticipantEmployment
	err := r.getDB(ctx).Where("participant_id = ? AND deleted_at IS NULL", participantID).First(&employment).Error
	if err != nil {
		return nil, translateError(err, "participant employment")
	}
	return &employment, nil
}

func (r *participantEmploymentRepository) Update(ctx context.Context, employment *entity.ParticipantEmployment) error {
	oldVersion := employment.Version
	employment.Version = oldVersion + 1

	result := r.getDB(ctx).Where("version = ? AND deleted_at IS NULL", oldVersion).Save(employment)
	if result.Error != nil {
		employment.Version = oldVersion
		return translateError(result.Error, "participant employment")
	}
	if result.RowsAffected == 0 {
		employment.Version = oldVersion
		return errors.ErrConflict("participant employment was modified by another request")
	}
	return nil
}

func (r *participantEmploymentRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	err := r.getDB(ctx).Model(&entity.ParticipantEmployment{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Update("deleted_at", gorm.Expr("NOW()")).Error
	if err != nil {
		return translateError(err, "participant employment")
	}
	return nil
}
