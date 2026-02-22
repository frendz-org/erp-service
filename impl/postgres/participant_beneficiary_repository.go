package postgres

import (
	"context"

	"erp-service/entity"
	"erp-service/pkg/errors"
	"erp-service/saving/participant"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type participantBeneficiaryRepository struct {
	baseRepository
}

func NewParticipantBeneficiaryRepository(db *gorm.DB) participant.ParticipantBeneficiaryRepository {
	return &participantBeneficiaryRepository{
		baseRepository: baseRepository{db: db},
	}
}

func (r *participantBeneficiaryRepository) Create(ctx context.Context, beneficiary *entity.ParticipantBeneficiary) error {
	if err := r.getDB(ctx).Create(beneficiary).Error; err != nil {
		return translateError(err, "participant beneficiary")
	}
	return nil
}

func (r *participantBeneficiaryRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.ParticipantBeneficiary, error) {
	var beneficiary entity.ParticipantBeneficiary
	err := r.getDB(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&beneficiary).Error
	if err != nil {
		return nil, translateError(err, "participant beneficiary")
	}
	return &beneficiary, nil
}

func (r *participantBeneficiaryRepository) ListByParticipantID(ctx context.Context, participantID uuid.UUID) ([]*entity.ParticipantBeneficiary, error) {
	var beneficiaries []*entity.ParticipantBeneficiary
	err := r.getDB(ctx).
		Where("participant_id = ? AND deleted_at IS NULL", participantID).
		Order("created_at ASC").
		Find(&beneficiaries).Error
	if err != nil {
		return nil, translateError(err, "participant beneficiary")
	}
	return beneficiaries, nil
}

func (r *participantBeneficiaryRepository) Update(ctx context.Context, beneficiary *entity.ParticipantBeneficiary) error {
	oldVersion := beneficiary.Version
	beneficiary.Version = oldVersion + 1

	result := r.getDB(ctx).Where("version = ? AND deleted_at IS NULL", oldVersion).Save(beneficiary)
	if result.Error != nil {
		beneficiary.Version = oldVersion
		return translateError(result.Error, "participant beneficiary")
	}
	if result.RowsAffected == 0 {
		beneficiary.Version = oldVersion
		return errors.ErrConflict("participant beneficiary was modified by another request")
	}
	return nil
}

func (r *participantBeneficiaryRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	err := r.getDB(ctx).Model(&entity.ParticipantBeneficiary{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Update("deleted_at", gorm.Expr("NOW()")).Error
	if err != nil {
		return translateError(err, "participant beneficiary")
	}
	return nil
}

func (r *participantBeneficiaryRepository) SoftDeleteAllByParticipantID(ctx context.Context, participantID uuid.UUID) error {
	err := r.getDB(ctx).Model(&entity.ParticipantBeneficiary{}).
		Where("participant_id = ? AND deleted_at IS NULL", participantID).
		Update("deleted_at", gorm.Expr("NOW()")).Error
	if err != nil {
		return translateError(err, "participant beneficiary")
	}
	return nil
}
