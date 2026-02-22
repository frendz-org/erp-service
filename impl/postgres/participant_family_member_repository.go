package postgres

import (
	"context"

	"erp-service/entity"
	"erp-service/pkg/errors"
	"erp-service/saving/participant"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type participantFamilyMemberRepository struct {
	baseRepository
}

func NewParticipantFamilyMemberRepository(db *gorm.DB) participant.ParticipantFamilyMemberRepository {
	return &participantFamilyMemberRepository{
		baseRepository: baseRepository{db: db},
	}
}

func (r *participantFamilyMemberRepository) Create(ctx context.Context, member *entity.ParticipantFamilyMember) error {
	if err := r.getDB(ctx).Create(member).Error; err != nil {
		return translateError(err, "participant family member")
	}
	return nil
}

func (r *participantFamilyMemberRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.ParticipantFamilyMember, error) {
	var member entity.ParticipantFamilyMember
	err := r.getDB(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&member).Error
	if err != nil {
		return nil, translateError(err, "participant family member")
	}
	return &member, nil
}

func (r *participantFamilyMemberRepository) ListByParticipantID(ctx context.Context, participantID uuid.UUID) ([]*entity.ParticipantFamilyMember, error) {
	var members []*entity.ParticipantFamilyMember
	err := r.getDB(ctx).
		Where("participant_id = ? AND deleted_at IS NULL", participantID).
		Order("created_at ASC").
		Find(&members).Error
	if err != nil {
		return nil, translateError(err, "participant family member")
	}
	return members, nil
}

func (r *participantFamilyMemberRepository) Update(ctx context.Context, member *entity.ParticipantFamilyMember) error {
	oldVersion := member.Version
	member.Version = oldVersion + 1

	result := r.getDB(ctx).Where("version = ? AND deleted_at IS NULL", oldVersion).Save(member)
	if result.Error != nil {
		member.Version = oldVersion
		return translateError(result.Error, "participant family member")
	}
	if result.RowsAffected == 0 {
		member.Version = oldVersion
		return errors.ErrConflict("participant family member was modified by another request")
	}
	return nil
}

func (r *participantFamilyMemberRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	err := r.getDB(ctx).Model(&entity.ParticipantFamilyMember{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Update("deleted_at", gorm.Expr("NOW()")).Error
	if err != nil {
		return translateError(err, "participant family member")
	}
	return nil
}

func (r *participantFamilyMemberRepository) SoftDeleteAllByParticipantID(ctx context.Context, participantID uuid.UUID) error {
	err := r.getDB(ctx).Model(&entity.ParticipantFamilyMember{}).
		Where("participant_id = ? AND deleted_at IS NULL", participantID).
		Update("deleted_at", gorm.Expr("NOW()")).Error
	if err != nil {
		return translateError(err, "participant family member")
	}
	return nil
}
