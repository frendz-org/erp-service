package postgres

import (
	"context"

	"erp-service/entity"
	"erp-service/pkg/errors"
	"erp-service/saving/participant/contract"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type participantBankAccountRepository struct {
	baseRepository
}

func NewParticipantBankAccountRepository(db *gorm.DB) contract.ParticipantBankAccountRepository {
	return &participantBankAccountRepository{
		baseRepository: baseRepository{db: db},
	}
}

func (r *participantBankAccountRepository) Create(ctx context.Context, account *entity.ParticipantBankAccount) error {
	if err := r.getDB(ctx).Create(account).Error; err != nil {
		return translateError(err, "participant bank account")
	}
	return nil
}

func (r *participantBankAccountRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.ParticipantBankAccount, error) {
	var account entity.ParticipantBankAccount
	err := r.getDB(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&account).Error
	if err != nil {
		return nil, translateError(err, "participant bank account")
	}
	return &account, nil
}

func (r *participantBankAccountRepository) ListByParticipantID(ctx context.Context, participantID uuid.UUID) ([]*entity.ParticipantBankAccount, error) {
	var accounts []*entity.ParticipantBankAccount
	err := r.getDB(ctx).
		Where("participant_id = ? AND deleted_at IS NULL", participantID).
		Order("is_primary DESC, created_at ASC").
		Find(&accounts).Error
	if err != nil {
		return nil, translateError(err, "participant bank account")
	}
	return accounts, nil
}

func (r *participantBankAccountRepository) Update(ctx context.Context, account *entity.ParticipantBankAccount) error {
	oldVersion := account.Version
	account.Version = oldVersion + 1

	result := r.getDB(ctx).Where("version = ? AND deleted_at IS NULL", oldVersion).Save(account)
	if result.Error != nil {
		account.Version = oldVersion
		return translateError(result.Error, "participant bank account")
	}
	if result.RowsAffected == 0 {
		account.Version = oldVersion
		return errors.ErrConflict("participant bank account was modified by another request")
	}
	return nil
}

func (r *participantBankAccountRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	err := r.getDB(ctx).Model(&entity.ParticipantBankAccount{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Update("deleted_at", gorm.Expr("NOW()")).Error
	if err != nil {
		return translateError(err, "participant bank account")
	}
	return nil
}

func (r *participantBankAccountRepository) ClearPrimary(ctx context.Context, participantID uuid.UUID) error {
	err := r.getDB(ctx).Model(&entity.ParticipantBankAccount{}).
		Where("participant_id = ? AND deleted_at IS NULL", participantID).
		Update("is_primary", false).Error
	if err != nil {
		return translateError(err, "participant bank account")
	}
	return nil
}
