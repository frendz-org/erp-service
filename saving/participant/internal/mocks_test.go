package internal

import (
	"context"
	"io"
	"time"

	"erp-service/entity"
	"erp-service/saving/participant/contract"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockTransactionManager struct {
	mock.Mock
}

func (m *MockTransactionManager) WithTransaction(ctx context.Context, fn func(context.Context) error) error {
	args := m.Called(ctx, fn)
	if args.Get(0) == nil {

		if fn != nil {
			return fn(ctx)
		}
		return nil
	}
	return args.Error(0)
}

type MockParticipantRepository struct {
	mock.Mock
}

func (m *MockParticipantRepository) Create(ctx context.Context, participant *entity.Participant) error {
	args := m.Called(ctx, participant)
	return args.Error(0)
}

func (m *MockParticipantRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Participant, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Participant), args.Error(1)
}

func (m *MockParticipantRepository) Update(ctx context.Context, participant *entity.Participant) error {
	args := m.Called(ctx, participant)
	return args.Error(0)
}

func (m *MockParticipantRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockParticipantRepository) List(ctx context.Context, filter *contract.ParticipantFilter) ([]*entity.Participant, int64, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*entity.Participant), args.Get(1).(int64), args.Error(2)
}

func (m *MockParticipantRepository) GetByKTPAndPensionNumber(ctx context.Context, ktpNumber, pensionNumber string, tenantID, productID uuid.UUID) (*entity.Participant, *entity.ParticipantPension, error) {
	args := m.Called(ctx, ktpNumber, pensionNumber, tenantID, productID)
	var p *entity.Participant
	var pp *entity.ParticipantPension
	if args.Get(0) != nil {
		p = args.Get(0).(*entity.Participant)
	}
	if args.Get(1) != nil {
		pp = args.Get(1).(*entity.ParticipantPension)
	}
	return p, pp, args.Error(2)
}

type MockParticipantIdentityRepository struct {
	mock.Mock
}

func (m *MockParticipantIdentityRepository) Create(ctx context.Context, identity *entity.ParticipantIdentity) error {
	args := m.Called(ctx, identity)
	return args.Error(0)
}

func (m *MockParticipantIdentityRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.ParticipantIdentity, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.ParticipantIdentity), args.Error(1)
}

func (m *MockParticipantIdentityRepository) ListByParticipantID(ctx context.Context, participantID uuid.UUID) ([]*entity.ParticipantIdentity, error) {
	args := m.Called(ctx, participantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.ParticipantIdentity), args.Error(1)
}

func (m *MockParticipantIdentityRepository) Update(ctx context.Context, identity *entity.ParticipantIdentity) error {
	args := m.Called(ctx, identity)
	return args.Error(0)
}

func (m *MockParticipantIdentityRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type MockParticipantAddressRepository struct {
	mock.Mock
}

func (m *MockParticipantAddressRepository) Create(ctx context.Context, address *entity.ParticipantAddress) error {
	args := m.Called(ctx, address)
	return args.Error(0)
}

func (m *MockParticipantAddressRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.ParticipantAddress, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.ParticipantAddress), args.Error(1)
}

func (m *MockParticipantAddressRepository) ListByParticipantID(ctx context.Context, participantID uuid.UUID) ([]*entity.ParticipantAddress, error) {
	args := m.Called(ctx, participantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.ParticipantAddress), args.Error(1)
}

func (m *MockParticipantAddressRepository) Update(ctx context.Context, address *entity.ParticipantAddress) error {
	args := m.Called(ctx, address)
	return args.Error(0)
}

func (m *MockParticipantAddressRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type MockParticipantBankAccountRepository struct {
	mock.Mock
}

func (m *MockParticipantBankAccountRepository) Create(ctx context.Context, account *entity.ParticipantBankAccount) error {
	args := m.Called(ctx, account)
	return args.Error(0)
}

func (m *MockParticipantBankAccountRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.ParticipantBankAccount, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.ParticipantBankAccount), args.Error(1)
}

func (m *MockParticipantBankAccountRepository) ListByParticipantID(ctx context.Context, participantID uuid.UUID) ([]*entity.ParticipantBankAccount, error) {
	args := m.Called(ctx, participantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.ParticipantBankAccount), args.Error(1)
}

func (m *MockParticipantBankAccountRepository) Update(ctx context.Context, account *entity.ParticipantBankAccount) error {
	args := m.Called(ctx, account)
	return args.Error(0)
}

func (m *MockParticipantBankAccountRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockParticipantBankAccountRepository) ClearPrimary(ctx context.Context, participantID uuid.UUID) error {
	args := m.Called(ctx, participantID)
	return args.Error(0)
}

type MockParticipantFamilyMemberRepository struct {
	mock.Mock
}

func (m *MockParticipantFamilyMemberRepository) Create(ctx context.Context, member *entity.ParticipantFamilyMember) error {
	args := m.Called(ctx, member)
	return args.Error(0)
}

func (m *MockParticipantFamilyMemberRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.ParticipantFamilyMember, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.ParticipantFamilyMember), args.Error(1)
}

func (m *MockParticipantFamilyMemberRepository) ListByParticipantID(ctx context.Context, participantID uuid.UUID) ([]*entity.ParticipantFamilyMember, error) {
	args := m.Called(ctx, participantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.ParticipantFamilyMember), args.Error(1)
}

func (m *MockParticipantFamilyMemberRepository) Update(ctx context.Context, member *entity.ParticipantFamilyMember) error {
	args := m.Called(ctx, member)
	return args.Error(0)
}

func (m *MockParticipantFamilyMemberRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type MockParticipantEmploymentRepository struct {
	mock.Mock
}

func (m *MockParticipantEmploymentRepository) Create(ctx context.Context, employment *entity.ParticipantEmployment) error {
	args := m.Called(ctx, employment)
	return args.Error(0)
}

func (m *MockParticipantEmploymentRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.ParticipantEmployment, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.ParticipantEmployment), args.Error(1)
}

func (m *MockParticipantEmploymentRepository) GetByParticipantID(ctx context.Context, participantID uuid.UUID) (*entity.ParticipantEmployment, error) {
	args := m.Called(ctx, participantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.ParticipantEmployment), args.Error(1)
}

func (m *MockParticipantEmploymentRepository) Update(ctx context.Context, employment *entity.ParticipantEmployment) error {
	args := m.Called(ctx, employment)
	return args.Error(0)
}

func (m *MockParticipantEmploymentRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type MockParticipantPensionRepository struct {
	mock.Mock
}

func (m *MockParticipantPensionRepository) Create(ctx context.Context, pension *entity.ParticipantPension) error {
	args := m.Called(ctx, pension)
	return args.Error(0)
}

func (m *MockParticipantPensionRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.ParticipantPension, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.ParticipantPension), args.Error(1)
}

func (m *MockParticipantPensionRepository) GetByParticipantID(ctx context.Context, participantID uuid.UUID) (*entity.ParticipantPension, error) {
	args := m.Called(ctx, participantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.ParticipantPension), args.Error(1)
}

func (m *MockParticipantPensionRepository) Update(ctx context.Context, pension *entity.ParticipantPension) error {
	args := m.Called(ctx, pension)
	return args.Error(0)
}

func (m *MockParticipantPensionRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type MockParticipantBeneficiaryRepository struct {
	mock.Mock
}

func (m *MockParticipantBeneficiaryRepository) Create(ctx context.Context, beneficiary *entity.ParticipantBeneficiary) error {
	args := m.Called(ctx, beneficiary)
	return args.Error(0)
}

func (m *MockParticipantBeneficiaryRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.ParticipantBeneficiary, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.ParticipantBeneficiary), args.Error(1)
}

func (m *MockParticipantBeneficiaryRepository) ListByParticipantID(ctx context.Context, participantID uuid.UUID) ([]*entity.ParticipantBeneficiary, error) {
	args := m.Called(ctx, participantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.ParticipantBeneficiary), args.Error(1)
}

func (m *MockParticipantBeneficiaryRepository) Update(ctx context.Context, beneficiary *entity.ParticipantBeneficiary) error {
	args := m.Called(ctx, beneficiary)
	return args.Error(0)
}

func (m *MockParticipantBeneficiaryRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type MockParticipantStatusHistoryRepository struct {
	mock.Mock
}

func (m *MockParticipantStatusHistoryRepository) Create(ctx context.Context, history *entity.ParticipantStatusHistory) error {
	args := m.Called(ctx, history)
	return args.Error(0)
}

func (m *MockParticipantStatusHistoryRepository) ListByParticipantID(ctx context.Context, participantID uuid.UUID) ([]*entity.ParticipantStatusHistory, error) {
	args := m.Called(ctx, participantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.ParticipantStatusHistory), args.Error(1)
}

type MockFileStorageAdapter struct {
	mock.Mock
}

func (m *MockFileStorageAdapter) UploadFile(ctx context.Context, bucket, objectKey string, data io.Reader, size int64, contentType string) (string, error) {
	args := m.Called(ctx, bucket, objectKey, data, size, contentType)
	return args.String(0), args.Error(1)
}

func (m *MockFileStorageAdapter) DeleteFile(ctx context.Context, bucket, objectKey string) error {
	args := m.Called(ctx, bucket, objectKey)
	return args.Error(0)
}

func (m *MockFileStorageAdapter) GetPresignedURL(ctx context.Context, bucket, objectKey string, expiry time.Duration) (string, error) {
	args := m.Called(ctx, bucket, objectKey, expiry)
	return args.String(0), args.Error(1)
}
