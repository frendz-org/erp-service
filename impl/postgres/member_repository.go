package postgres

import (
	"context"

	"erp-service/entity"
	"erp-service/saving/member"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type memberRepository struct {
	baseRepository
}

func NewMemberRepository(db *gorm.DB) member.MemberRepository {
	return &memberRepository{
		baseRepository: baseRepository{db: db},
	}
}

func (r *memberRepository) Create(ctx context.Context, m *entity.Member) error {
	if err := r.getDB(ctx).Create(m).Error; err != nil {
		return translateError(err, "member")
	}
	return nil
}

func (r *memberRepository) GetByRegistrationID(ctx context.Context, registrationID uuid.UUID) (*entity.Member, error) {
	var m entity.Member
	err := r.getDB(ctx).
		Where("user_tenant_registration_id = ? AND deleted_at IS NULL", registrationID).
		First(&m).Error
	if err != nil {
		return nil, translateError(err, "member")
	}
	return &m, nil
}

func (r *memberRepository) GetByUserTenantProduct(ctx context.Context, userID, tenantID, productID uuid.UUID) (*entity.Member, error) {
	var m entity.Member
	err := r.getDB(ctx).
		Where("user_id = ? AND tenant_id = ? AND product_id = ? AND deleted_at IS NULL",
			userID, tenantID, productID).
		First(&m).Error
	if err != nil {
		return nil, translateError(err, "member")
	}
	return &m, nil
}

func (r *memberRepository) GetByParticipantNumber(ctx context.Context, tenantID, productID uuid.UUID, participantNumber string) (*entity.Member, error) {
	var m entity.Member
	err := r.getDB(ctx).
		Where("tenant_id = ? AND product_id = ? AND participant_number = ? AND deleted_at IS NULL",
			tenantID, productID, participantNumber).
		First(&m).Error
	if err != nil {
		return nil, translateError(err, "member")
	}
	return &m, nil
}
