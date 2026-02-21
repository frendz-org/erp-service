package contract

import (
	"context"

	"erp-service/entity"

	"github.com/google/uuid"
)

type ParticipantFilter struct {
	TenantID  uuid.UUID
	ProductID uuid.UUID
	Status    *string
	Search    string
	Page      int
	PerPage   int
	SortBy    string
	SortOrder string
}

type ParticipantRepository interface {
	Create(ctx context.Context, participant *entity.Participant) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Participant, error)
	Update(ctx context.Context, participant *entity.Participant) error
	SoftDelete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filter *ParticipantFilter) ([]*entity.Participant, int64, error)
	GetByKTPAndPensionNumber(ctx context.Context, ktpNumber, pensionNumber string, tenantID, productID uuid.UUID) (*entity.Participant, *entity.ParticipantPension, error)
}

type TenantRepository interface {
	GetByCode(ctx context.Context, code string) (*entity.Tenant, error)
}

type ProductRepository interface {
	GetByCodeAndTenant(ctx context.Context, tenantID uuid.UUID, code string) (*entity.Product, error)
}

type ProductRegistrationConfigRepository interface {
	GetByProductAndType(ctx context.Context, productID uuid.UUID, regType string) (*entity.ProductRegistrationConfig, error)
}

type UserTenantRegistrationRepository interface {
	Create(ctx context.Context, reg *entity.UserTenantRegistration) error
	GetByUserAndProduct(ctx context.Context, userID, tenantID, productID uuid.UUID, regType string) (*entity.UserTenantRegistration, error)
}

type UserProfileRepository interface {
	GetByUserID(ctx context.Context, userID uuid.UUID) (*entity.UserProfile, error)
}

type ParticipantIdentityRepository interface {
	Create(ctx context.Context, identity *entity.ParticipantIdentity) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.ParticipantIdentity, error)
	ListByParticipantID(ctx context.Context, participantID uuid.UUID) ([]*entity.ParticipantIdentity, error)
	Update(ctx context.Context, identity *entity.ParticipantIdentity) error
	SoftDelete(ctx context.Context, id uuid.UUID) error
}

type ParticipantAddressRepository interface {
	Create(ctx context.Context, address *entity.ParticipantAddress) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.ParticipantAddress, error)
	ListByParticipantID(ctx context.Context, participantID uuid.UUID) ([]*entity.ParticipantAddress, error)
	Update(ctx context.Context, address *entity.ParticipantAddress) error
	SoftDelete(ctx context.Context, id uuid.UUID) error
}

type ParticipantBankAccountRepository interface {
	Create(ctx context.Context, account *entity.ParticipantBankAccount) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.ParticipantBankAccount, error)
	ListByParticipantID(ctx context.Context, participantID uuid.UUID) ([]*entity.ParticipantBankAccount, error)
	Update(ctx context.Context, account *entity.ParticipantBankAccount) error
	SoftDelete(ctx context.Context, id uuid.UUID) error
	ClearPrimary(ctx context.Context, participantID uuid.UUID) error
}

type ParticipantFamilyMemberRepository interface {
	Create(ctx context.Context, member *entity.ParticipantFamilyMember) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.ParticipantFamilyMember, error)
	ListByParticipantID(ctx context.Context, participantID uuid.UUID) ([]*entity.ParticipantFamilyMember, error)
	Update(ctx context.Context, member *entity.ParticipantFamilyMember) error
	SoftDelete(ctx context.Context, id uuid.UUID) error
}

type ParticipantEmploymentRepository interface {
	Create(ctx context.Context, employment *entity.ParticipantEmployment) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.ParticipantEmployment, error)
	GetByParticipantID(ctx context.Context, participantID uuid.UUID) (*entity.ParticipantEmployment, error)
	Update(ctx context.Context, employment *entity.ParticipantEmployment) error
	SoftDelete(ctx context.Context, id uuid.UUID) error
}

type ParticipantPensionRepository interface {
	Create(ctx context.Context, pension *entity.ParticipantPension) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.ParticipantPension, error)
	GetByParticipantID(ctx context.Context, participantID uuid.UUID) (*entity.ParticipantPension, error)
	Update(ctx context.Context, pension *entity.ParticipantPension) error
	SoftDelete(ctx context.Context, id uuid.UUID) error
}

type ParticipantBeneficiaryRepository interface {
	Create(ctx context.Context, beneficiary *entity.ParticipantBeneficiary) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.ParticipantBeneficiary, error)
	ListByParticipantID(ctx context.Context, participantID uuid.UUID) ([]*entity.ParticipantBeneficiary, error)
	Update(ctx context.Context, beneficiary *entity.ParticipantBeneficiary) error
	SoftDelete(ctx context.Context, id uuid.UUID) error
}

type ParticipantStatusHistoryRepository interface {
	Create(ctx context.Context, history *entity.ParticipantStatusHistory) error
	ListByParticipantID(ctx context.Context, participantID uuid.UUID) ([]*entity.ParticipantStatusHistory, error)
}
