package postgres

import (
	"context"

	"iam-service/entity"
	"iam-service/saving/participant/contract"
	"iam-service/pkg/errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type participantAddressRepository struct {
	baseRepository
}

func NewParticipantAddressRepository(db *gorm.DB) contract.ParticipantAddressRepository {
	return &participantAddressRepository{
		baseRepository: baseRepository{db: db},
	}
}

func (r *participantAddressRepository) Create(ctx context.Context, address *entity.ParticipantAddress) error {
	if err := r.getDB(ctx).Create(address).Error; err != nil {
		return translateError(err, "participant address")
	}
	return nil
}

func (r *participantAddressRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.ParticipantAddress, error) {
	var address entity.ParticipantAddress
	err := r.getDB(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&address).Error
	if err != nil {
		return nil, translateError(err, "participant address")
	}
	return &address, nil
}

func (r *participantAddressRepository) ListByParticipantID(ctx context.Context, participantID uuid.UUID) ([]*entity.ParticipantAddress, error) {
	var addresses []*entity.ParticipantAddress
	err := r.getDB(ctx).
		Where("participant_id = ? AND deleted_at IS NULL", participantID).
		Order("is_primary DESC, created_at ASC").
		Find(&addresses).Error
	if err != nil {
		return nil, translateError(err, "participant address")
	}
	return addresses, nil
}

func (r *participantAddressRepository) Update(ctx context.Context, address *entity.ParticipantAddress) error {
	oldVersion := address.Version
	address.Version = oldVersion + 1

	result := r.getDB(ctx).Where("version = ? AND deleted_at IS NULL", oldVersion).Save(address)
	if result.Error != nil {
		address.Version = oldVersion
		return translateError(result.Error, "participant address")
	}
	if result.RowsAffected == 0 {
		address.Version = oldVersion
		return errors.ErrConflict("participant address was modified by another request")
	}
	return nil
}

func (r *participantAddressRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	err := r.getDB(ctx).Model(&entity.ParticipantAddress{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Update("deleted_at", gorm.Expr("NOW()")).Error
	if err != nil {
		return translateError(err, "participant address")
	}
	return nil
}
