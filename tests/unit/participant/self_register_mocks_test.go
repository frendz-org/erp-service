package participant_test

import (
	"context"
	"time"

	"erp-service/config"
	"erp-service/entity"
	"erp-service/masterdata"
	"erp-service/saving/participant"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

type mockTenantRepository struct {
	mock.Mock
}

func (m *mockTenantRepository) GetByCode(ctx context.Context, code string) (*entity.Tenant, error) {
	args := m.Called(ctx, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Tenant), args.Error(1)
}

type mockProductRepository struct {
	mock.Mock
}

func (m *mockProductRepository) GetByCodeAndTenant(ctx context.Context, tenantID uuid.UUID, code string) (*entity.Product, error) {
	args := m.Called(ctx, tenantID, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Product), args.Error(1)
}

type mockProductRegistrationConfigRepository struct {
	mock.Mock
}

func (m *mockProductRegistrationConfigRepository) GetByProductAndType(ctx context.Context, productID uuid.UUID, regType string) (*entity.ProductRegistrationConfig, error) {
	args := m.Called(ctx, productID, regType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.ProductRegistrationConfig), args.Error(1)
}

type mockUserTenantRegistrationRepository struct {
	mock.Mock
}

func (m *mockUserTenantRegistrationRepository) Create(ctx context.Context, reg *entity.UserTenantRegistration) error {
	args := m.Called(ctx, reg)
	return args.Error(0)
}

func (m *mockUserTenantRegistrationRepository) GetByUserAndProduct(ctx context.Context, userID, tenantID, productID uuid.UUID, regType string) (*entity.UserTenantRegistration, error) {
	args := m.Called(ctx, userID, tenantID, productID, regType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.UserTenantRegistration), args.Error(1)
}

type mockUserProfileRepository struct {
	mock.Mock
}

func (m *mockUserProfileRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*entity.UserProfile, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.UserProfile), args.Error(1)
}

type mockMasterdataValidator struct {
	mock.Mock
}

func (m *mockMasterdataValidator) ValidateItemCode(ctx context.Context, req *masterdata.ValidateCodeRequest) (*masterdata.ValidateCodeResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*masterdata.ValidateCodeResponse), args.Error(1)
}

type mockParticipantRepositoryWithKTP struct {
	MockParticipantRepository
}

func (m *mockParticipantRepositoryWithKTP) GetByKTPAndPensionNumber(ctx context.Context, ktpNumber, pensionNumber string, tenantID, productID uuid.UUID) (*entity.Participant, *entity.ParticipantPension, error) {
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

func validSelfRegisterRequest() *participant.SelfRegisterRequest {
	return &participant.SelfRegisterRequest{
		UserID:            uuid.New(),
		Organization:      "TENANT_001",
		IdentityNumber:    "3201011501900001",
		ParticipantNumber: "ABC12345",
		PhoneNumber:       "+6281234567890",
	}
}

func activeTenant() *entity.Tenant {
	return &entity.Tenant{
		ID:     uuid.New(),
		Code:   "TENANT_001",
		Status: entity.TenantStatusActive,
	}
}

func activeProduct(tenantID uuid.UUID) *entity.Product {
	return &entity.Product{
		ID:       uuid.New(),
		TenantID: tenantID,
		Code:     "frendz-saving",
	}
}

func activeConfig(productID uuid.UUID) *entity.ProductRegistrationConfig {
	return &entity.ProductRegistrationConfig{
		ID:               uuid.New(),
		ProductID:        productID,
		RegistrationType: "PARTICIPANT",
		IsActive:         true,
	}
}

func completeProfile(userID uuid.UUID) *entity.UserProfile {
	dob := time.Date(1990, 1, 15, 0, 0, 0, 0, time.UTC)
	gender := entity.GenderMale
	return &entity.UserProfile{
		UserID:      userID,
		FirstName:   "John",
		LastName:    "Doe",
		DateOfBirth: &dob,
		Gender:      &gender,
	}
}

func buildSelfRegisterUsecase(
	txMgr *MockTransactionManager,
	partRepo *mockParticipantRepositoryWithKTP,
	pensionRepo *MockParticipantPensionRepository,
	statusHistoryRepo *MockParticipantStatusHistoryRepository,
	tenantRepo *mockTenantRepository,
	productRepo *mockProductRepository,
	configRepo *mockProductRegistrationConfigRepository,
	utrRepo *mockUserTenantRegistrationRepository,
	profileRepo *mockUserProfileRepository,
	mdValidator *mockMasterdataValidator,
	empDataRepo *MockEmployeeDataRepository,
) participant.Usecase {
	var edr participant.EmployeeDataRepository
	if empDataRepo != nil {
		edr = empDataRepo
	}
	return participant.NewUsecase(
		&config.Config{},
		zap.NewNop(),
		txMgr,
		partRepo,
		&MockParticipantIdentityRepository{},
		&MockParticipantAddressRepository{},
		&MockParticipantBankAccountRepository{},
		&MockParticipantFamilyMemberRepository{},
		&MockParticipantEmploymentRepository{},
		pensionRepo,
		&MockParticipantBeneficiaryRepository{},
		statusHistoryRepo,
		&MockFileStorageAdapter{},
		&MockFileRepository{},
		tenantRepo,
		productRepo,
		configRepo,
		utrRepo,
		profileRepo,
		mdValidator,
		edr,
	)
}
