package postgres

import (
	"context"

	"iam-service/entity"
	"iam-service/saving/participant/contract"
	"iam-service/pkg/errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type participantIdentityRepository struct {
	baseRepository
}

func NewParticipantIdentityRepository(db *gorm.DB) contract.ParticipantIdentityRepository {
	return &participantIdentityRepository{
		baseRepository: baseRepository{db: db},
	}
}

func (r *participantIdentityRepository) Create(ctx context.Context, identity *entity.ParticipantIdentity) error {
	if err := r.getDB(ctx).Create(identity).Error; err != nil {
		return translateError(err, "participant identity")
	}
	return nil
}

func (r *participantIdentityRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.ParticipantIdentity, error) {
	var identity entity.ParticipantIdentity
	err := r.getDB(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&identity).Error
	if err != nil {
		return nil, translateError(err, "participant identity")
	}
	return &identity, nil
}

func (r *participantIdentityRepository) ListByParticipantID(ctx context.Context, participantID uuid.UUID) ([]*entity.ParticipantIdentity, error) {
	var identities []*entity.ParticipantIdentity
	err := r.getDB(ctx).
		Where("participant_id = ? AND deleted_at IS NULL", participantID).
		Order("created_at ASC").
		Find(&identities).Error
	if err != nil {
		return nil, translateError(err, "participant identity")
	}
	return identities, nil
}

func (r *participantIdentityRepository) Update(ctx context.Context, identity *entity.ParticipantIdentity) error {
	oldVersion := identity.Version
	identity.Version = oldVersion + 1

	result := r.getDB(ctx).Where("version = ? AND deleted_at IS NULL", oldVersion).Save(identity)
	if result.Error != nil {
		identity.Version = oldVersion
		return translateError(result.Error, "participant identity")
	}
	if result.RowsAffected == 0 {
		identity.Version = oldVersion
		return errors.ErrConflict("participant identity was modified by another request")
	}
	return nil
}

func (r *participantIdentityRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	err := r.getDB(ctx).Model(&entity.ParticipantIdentity{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Update("deleted_at", gorm.Expr("NOW()")).Error
	if err != nil {
		return translateError(err, "participant identity")
	}
	return nil
}
