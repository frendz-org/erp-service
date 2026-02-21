package internal

import (
	"context"
	"testing"
	"time"

	"erp-service/entity"
	"erp-service/masterdata/masterdatadto"
	apperrors "erp-service/pkg/errors"
	"erp-service/saving/participant/participantdto"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
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

func (m *mockMasterdataValidator) ValidateItemCode(ctx context.Context, req *masterdatadto.ValidateCodeRequest) (*masterdatadto.ValidateCodeResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*masterdatadto.ValidateCodeResponse), args.Error(1)
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

func validSelfRegisterRequest() *participantdto.SelfRegisterRequest {
	return &participantdto.SelfRegisterRequest{
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
) *usecase {
	return &usecase{
		txManager:         txMgr,
		participantRepo:   partRepo,
		identityRepo:      &MockParticipantIdentityRepository{},
		addressRepo:       &MockParticipantAddressRepository{},
		bankAccountRepo:   &MockParticipantBankAccountRepository{},
		familyMemberRepo:  &MockParticipantFamilyMemberRepository{},
		employmentRepo:    &MockParticipantEmploymentRepository{},
		pensionRepo:       pensionRepo,
		beneficiaryRepo:   &MockParticipantBeneficiaryRepository{},
		statusHistoryRepo: statusHistoryRepo,
		tenantRepo:        tenantRepo,
		productRepo:       productRepo,
		configRepo:        configRepo,
		utrRepo:           utrRepo,
		userProfileRepo:   profileRepo,
		masterdataUsecase: mdValidator,
	}
}

func TestSelfRegister_Success_NewParticipant(t *testing.T) {
	ctx := context.Background()
	req := validSelfRegisterRequest()
	tenant := activeTenant()
	product := activeProduct(tenant.ID)
	config := activeConfig(product.ID)
	profile := completeProfile(req.UserID)

	txMgr := &MockTransactionManager{}
	partRepo := &mockParticipantRepositoryWithKTP{}
	pensionRepo := &MockParticipantPensionRepository{}
	statusHistoryRepo := &MockParticipantStatusHistoryRepository{}
	tenantRepo := &mockTenantRepository{}
	productRepo := &mockProductRepository{}
	configRepo := &mockProductRegistrationConfigRepository{}
	utrRepo := &mockUserTenantRegistrationRepository{}
	profileRepo := &mockUserProfileRepository{}
	mdValidator := &mockMasterdataValidator{}

	mdValidator.On("ValidateItemCode", ctx, mock.MatchedBy(func(r *masterdatadto.ValidateCodeRequest) bool {
		return r.CategoryCode == "TENANT_TYPE" && r.ItemCode == req.Organization
	})).Return(&masterdatadto.ValidateCodeResponse{Valid: true}, nil)

	tenantRepo.On("GetByCode", ctx, req.Organization).Return(tenant, nil)
	productRepo.On("GetByCodeAndTenant", ctx, tenant.ID, "frendz-saving").Return(product, nil)
	configRepo.On("GetByProductAndType", ctx, product.ID, "PARTICIPANT").Return(config, nil)
	profileRepo.On("GetByUserID", ctx, req.UserID).Return(profile, nil)
	partRepo.On("GetByKTPAndPensionNumber", ctx, req.IdentityNumber, req.ParticipantNumber, tenant.ID, product.ID).
		Return(nil, nil, apperrors.ErrNotFound("not found"))
	utrRepo.On("GetByUserAndProduct", ctx, req.UserID, tenant.ID, product.ID, "PARTICIPANT").
		Return(nil, apperrors.ErrNotFound("not found"))

	txMgr.On("WithTransaction", ctx, mock.Anything).Return(nil)
	partRepo.On("Create", ctx, mock.AnythingOfType("*entity.Participant")).Return(nil)
	pensionRepo.On("Create", ctx, mock.AnythingOfType("*entity.ParticipantPension")).Return(nil)
	statusHistoryRepo.On("Create", ctx, mock.AnythingOfType("*entity.ParticipantStatusHistory")).Return(nil)
	utrRepo.On("Create", ctx, mock.AnythingOfType("*entity.UserTenantRegistration")).Return(nil)

	partRepo.MockParticipantRepository.On("GetByID", ctx, mock.Anything).Return(&entity.Participant{
		ID:        uuid.New(),
		TenantID:  tenant.ID,
		ProductID: product.ID,
		UserID:    &req.UserID,
		FullName:  profile.FullName(),
		Status:    entity.ParticipantStatusDraft,
		CreatedBy: req.UserID,
	}, nil).Maybe()

	uc := buildSelfRegisterUsecase(txMgr, partRepo, pensionRepo, statusHistoryRepo, tenantRepo, productRepo, configRepo, utrRepo, profileRepo, mdValidator)

	resp, err := uc.SelfRegister(ctx, req)

	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.False(t, resp.IsLinked)
	assert.Equal(t, "PENDING_APPROVAL", resp.RegistrationStatus)
}

func TestSelfRegister_Success_ResponseDataHasNoKTPNumber(t *testing.T) {

	ctx := context.Background()
	req := validSelfRegisterRequest()
	tenant := activeTenant()
	product := activeProduct(tenant.ID)
	config := activeConfig(product.ID)
	profile := completeProfile(req.UserID)

	txMgr := &MockTransactionManager{}
	partRepo := &mockParticipantRepositoryWithKTP{}
	pensionRepo := &MockParticipantPensionRepository{}
	statusHistoryRepo := &MockParticipantStatusHistoryRepository{}
	tenantRepo := &mockTenantRepository{}
	productRepo := &mockProductRepository{}
	configRepo := &mockProductRegistrationConfigRepository{}
	utrRepo := &mockUserTenantRegistrationRepository{}
	profileRepo := &mockUserProfileRepository{}
	mdValidator := &mockMasterdataValidator{}

	mdValidator.On("ValidateItemCode", ctx, mock.Anything).Return(&masterdatadto.ValidateCodeResponse{Valid: true}, nil)
	tenantRepo.On("GetByCode", ctx, req.Organization).Return(tenant, nil)
	productRepo.On("GetByCodeAndTenant", ctx, tenant.ID, "frendz-saving").Return(product, nil)
	configRepo.On("GetByProductAndType", ctx, product.ID, "PARTICIPANT").Return(config, nil)
	profileRepo.On("GetByUserID", ctx, req.UserID).Return(profile, nil)
	partRepo.On("GetByKTPAndPensionNumber", ctx, req.IdentityNumber, req.ParticipantNumber, tenant.ID, product.ID).
		Return(nil, nil, apperrors.ErrNotFound("not found"))
	utrRepo.On("GetByUserAndProduct", ctx, req.UserID, tenant.ID, product.ID, "PARTICIPANT").
		Return(nil, apperrors.ErrNotFound("not found"))

	txMgr.On("WithTransaction", ctx, mock.Anything).Return(nil)
	partRepo.On("Create", ctx, mock.AnythingOfType("*entity.Participant")).Return(nil)
	pensionRepo.On("Create", ctx, mock.AnythingOfType("*entity.ParticipantPension")).Return(nil)
	statusHistoryRepo.On("Create", ctx, mock.AnythingOfType("*entity.ParticipantStatusHistory")).Return(nil)
	utrRepo.On("Create", ctx, mock.AnythingOfType("*entity.UserTenantRegistration")).Return(nil)

	uc := buildSelfRegisterUsecase(txMgr, partRepo, pensionRepo, statusHistoryRepo, tenantRepo, productRepo, configRepo, utrRepo, profileRepo, mdValidator)

	resp, err := uc.SelfRegister(ctx, req)

	require.NoError(t, err)
	require.NotNil(t, resp)

	assert.NotNil(t, resp.Data)

	assert.Equal(t, string(entity.ParticipantStatusDraft), resp.Data.Status)
	assert.NotEmpty(t, resp.Data.ParticipantNumber)
}

func TestSelfRegister_Success_LinkExistingParticipant(t *testing.T) {
	ctx := context.Background()
	req := validSelfRegisterRequest()
	tenant := activeTenant()
	product := activeProduct(tenant.ID)
	config := activeConfig(product.ID)
	profile := completeProfile(req.UserID)

	existingParticipant := &entity.Participant{
		ID:        uuid.New(),
		TenantID:  tenant.ID,
		ProductID: product.ID,
		UserID:    nil,
		FullName:  "Existing Person",
		Status:    entity.ParticipantStatusDraft,
		CreatedBy: uuid.New(),
	}
	existingPension := &entity.ParticipantPension{
		ID:            uuid.New(),
		ParticipantID: existingParticipant.ID,
	}

	txMgr := &MockTransactionManager{}
	partRepo := &mockParticipantRepositoryWithKTP{}
	pensionRepo := &MockParticipantPensionRepository{}
	statusHistoryRepo := &MockParticipantStatusHistoryRepository{}
	tenantRepo := &mockTenantRepository{}
	productRepo := &mockProductRepository{}
	configRepo := &mockProductRegistrationConfigRepository{}
	utrRepo := &mockUserTenantRegistrationRepository{}
	profileRepo := &mockUserProfileRepository{}
	mdValidator := &mockMasterdataValidator{}

	mdValidator.On("ValidateItemCode", ctx, mock.Anything).Return(&masterdatadto.ValidateCodeResponse{Valid: true}, nil)
	tenantRepo.On("GetByCode", ctx, req.Organization).Return(tenant, nil)
	productRepo.On("GetByCodeAndTenant", ctx, tenant.ID, "frendz-saving").Return(product, nil)
	configRepo.On("GetByProductAndType", ctx, product.ID, "PARTICIPANT").Return(config, nil)
	profileRepo.On("GetByUserID", ctx, req.UserID).Return(profile, nil)
	partRepo.On("GetByKTPAndPensionNumber", ctx, req.IdentityNumber, req.ParticipantNumber, tenant.ID, product.ID).
		Return(existingParticipant, existingPension, nil)

	txMgr.On("WithTransaction", ctx, mock.Anything).Return(nil)
	partRepo.MockParticipantRepository.On("Update", ctx, mock.AnythingOfType("*entity.Participant")).Return(nil)
	utrRepo.On("Create", ctx, mock.AnythingOfType("*entity.UserTenantRegistration")).Return(nil)
	statusHistoryRepo.On("Create", ctx, mock.AnythingOfType("*entity.ParticipantStatusHistory")).Return(nil)

	uc := buildSelfRegisterUsecase(txMgr, partRepo, pensionRepo, statusHistoryRepo, tenantRepo, productRepo, configRepo, utrRepo, profileRepo, mdValidator)

	resp, err := uc.SelfRegister(ctx, req)

	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.IsLinked)
	assert.Equal(t, "PENDING_APPROVAL", resp.RegistrationStatus)
}

func TestSelfRegister_InvalidOrganizationFormat(t *testing.T) {
	ctx := context.Background()
	req := validSelfRegisterRequest()
	req.Organization = "INVALID"

	uc := &usecase{}
	resp, err := uc.SelfRegister(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	appErr := apperrors.GetAppError(err)
	require.NotNil(t, appErr)
	assert.Equal(t, apperrors.CodeValidation, appErr.Code)
}

func TestSelfRegister_InvalidNIK_NotSixteenDigits(t *testing.T) {
	ctx := context.Background()
	req := validSelfRegisterRequest()
	req.IdentityNumber = "123456789"

	uc := &usecase{}
	resp, err := uc.SelfRegister(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	appErr := apperrors.GetAppError(err)
	require.NotNil(t, appErr)
	assert.Equal(t, apperrors.CodeValidation, appErr.Code)
}

func TestSelfRegister_InvalidNIK_BadDay_Zero(t *testing.T) {
	ctx := context.Background()
	req := validSelfRegisterRequest()

	req.IdentityNumber = "3201010001900001"

	uc := &usecase{}
	resp, err := uc.SelfRegister(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	appErr := apperrors.GetAppError(err)
	require.NotNil(t, appErr)
	assert.Equal(t, apperrors.CodeValidation, appErr.Code)
}

func TestSelfRegister_InvalidNIK_BadMonth_Thirteen(t *testing.T) {
	ctx := context.Background()
	req := validSelfRegisterRequest()

	req.IdentityNumber = "3201011513990001"

	uc := &usecase{}
	resp, err := uc.SelfRegister(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	appErr := apperrors.GetAppError(err)
	require.NotNil(t, appErr)
	assert.Equal(t, apperrors.CodeValidation, appErr.Code)
}

func TestSelfRegister_InvalidParticipantNumber(t *testing.T) {
	ctx := context.Background()
	req := validSelfRegisterRequest()
	req.ParticipantNumber = "invalid-format-123"

	uc := &usecase{}
	resp, err := uc.SelfRegister(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	appErr := apperrors.GetAppError(err)
	require.NotNil(t, appErr)
	assert.Equal(t, apperrors.CodeValidation, appErr.Code)
}

func TestSelfRegister_InvalidPhoneNumber(t *testing.T) {
	ctx := context.Background()
	req := validSelfRegisterRequest()
	req.PhoneNumber = "12345"

	uc := &usecase{}
	resp, err := uc.SelfRegister(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	appErr := apperrors.GetAppError(err)
	require.NotNil(t, appErr)
	assert.Equal(t, apperrors.CodeValidation, appErr.Code)
}

func TestSelfRegister_OrganizationNotInTenantType002(t *testing.T) {
	ctx := context.Background()
	req := validSelfRegisterRequest()

	mdValidator := &mockMasterdataValidator{}
	mdValidator.On("ValidateItemCode", ctx, mock.Anything).Return(&masterdatadto.ValidateCodeResponse{
		Valid:   false,
		Message: "item not found in parent",
	}, nil)

	uc := &usecase{
		masterdataUsecase: mdValidator,
	}
	resp, err := uc.SelfRegister(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	appErr := apperrors.GetAppError(err)
	require.NotNil(t, appErr)
	assert.Equal(t, apperrors.CodeUnprocessable, appErr.Code)
}

func TestSelfRegister_TenantInactive(t *testing.T) {
	ctx := context.Background()
	req := validSelfRegisterRequest()

	inactiveTenant := &entity.Tenant{
		ID:     uuid.New(),
		Code:   req.Organization,
		Status: entity.TenantStatusInactive,
	}

	mdValidator := &mockMasterdataValidator{}
	mdValidator.On("ValidateItemCode", ctx, mock.Anything).Return(&masterdatadto.ValidateCodeResponse{Valid: true}, nil)
	tenantRepo := &mockTenantRepository{}
	tenantRepo.On("GetByCode", ctx, req.Organization).Return(inactiveTenant, nil)

	uc := &usecase{
		masterdataUsecase: mdValidator,
		tenantRepo:        tenantRepo,
	}
	resp, err := uc.SelfRegister(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	appErr := apperrors.GetAppError(err)
	require.NotNil(t, appErr)
	assert.Equal(t, apperrors.CodeUnprocessable, appErr.Code)
}

func TestSelfRegister_ProductNotFound(t *testing.T) {
	ctx := context.Background()
	req := validSelfRegisterRequest()
	tenant := activeTenant()

	mdValidator := &mockMasterdataValidator{}
	mdValidator.On("ValidateItemCode", ctx, mock.Anything).Return(&masterdatadto.ValidateCodeResponse{Valid: true}, nil)
	tenantRepo := &mockTenantRepository{}
	tenantRepo.On("GetByCode", ctx, req.Organization).Return(tenant, nil)
	productRepo := &mockProductRepository{}
	productRepo.On("GetByCodeAndTenant", ctx, tenant.ID, "frendz-saving").Return(nil, apperrors.ErrNotFound("product not found"))

	uc := &usecase{
		masterdataUsecase: mdValidator,
		tenantRepo:        tenantRepo,
		productRepo:       productRepo,
	}
	resp, err := uc.SelfRegister(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	appErr := apperrors.GetAppError(err)
	require.NotNil(t, appErr)
	assert.Equal(t, apperrors.CodeNotFound, appErr.Code)
}

func TestSelfRegister_ConfigNotFound(t *testing.T) {
	ctx := context.Background()
	req := validSelfRegisterRequest()
	tenant := activeTenant()
	product := activeProduct(tenant.ID)

	mdValidator := &mockMasterdataValidator{}
	mdValidator.On("ValidateItemCode", ctx, mock.Anything).Return(&masterdatadto.ValidateCodeResponse{Valid: true}, nil)
	tenantRepo := &mockTenantRepository{}
	tenantRepo.On("GetByCode", ctx, req.Organization).Return(tenant, nil)
	productRepo := &mockProductRepository{}
	productRepo.On("GetByCodeAndTenant", ctx, tenant.ID, "frendz-saving").Return(product, nil)
	configRepo := &mockProductRegistrationConfigRepository{}

	configRepo.On("GetByProductAndType", ctx, product.ID, "PARTICIPANT").Return(nil, apperrors.ErrNotFound("config not found"))

	uc := &usecase{
		masterdataUsecase: mdValidator,
		tenantRepo:        tenantRepo,
		productRepo:       productRepo,
		configRepo:        configRepo,
	}
	resp, err := uc.SelfRegister(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	appErr := apperrors.GetAppError(err)
	require.NotNil(t, appErr)
	assert.Equal(t, apperrors.CodeUnprocessable, appErr.Code)
}

func TestSelfRegister_ConfigNotActive(t *testing.T) {
	ctx := context.Background()
	req := validSelfRegisterRequest()
	tenant := activeTenant()
	product := activeProduct(tenant.ID)
	inactiveConfig := &entity.ProductRegistrationConfig{
		ID:               uuid.New(),
		ProductID:        product.ID,
		RegistrationType: "PARTICIPANT",
		IsActive:         false,
	}

	mdValidator := &mockMasterdataValidator{}
	mdValidator.On("ValidateItemCode", ctx, mock.Anything).Return(&masterdatadto.ValidateCodeResponse{Valid: true}, nil)
	tenantRepo := &mockTenantRepository{}
	tenantRepo.On("GetByCode", ctx, req.Organization).Return(tenant, nil)
	productRepo := &mockProductRepository{}
	productRepo.On("GetByCodeAndTenant", ctx, tenant.ID, "frendz-saving").Return(product, nil)
	configRepo := &mockProductRegistrationConfigRepository{}
	configRepo.On("GetByProductAndType", ctx, product.ID, "PARTICIPANT").Return(inactiveConfig, nil)

	uc := &usecase{
		masterdataUsecase: mdValidator,
		tenantRepo:        tenantRepo,
		productRepo:       productRepo,
		configRepo:        configRepo,
	}
	resp, err := uc.SelfRegister(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	appErr := apperrors.GetAppError(err)
	require.NotNil(t, appErr)
	assert.Equal(t, apperrors.CodeUnprocessable, appErr.Code)
}

func TestSelfRegister_IncompleteProfile_MissingGender(t *testing.T) {
	ctx := context.Background()
	req := validSelfRegisterRequest()
	tenant := activeTenant()
	product := activeProduct(tenant.ID)
	config := activeConfig(product.ID)
	dob := time.Date(1990, 1, 15, 0, 0, 0, 0, time.UTC)
	incompleteProfile := &entity.UserProfile{
		UserID:      req.UserID,
		FirstName:   "John",
		LastName:    "Doe",
		DateOfBirth: &dob,
		Gender:      nil,
	}

	mdValidator := &mockMasterdataValidator{}
	mdValidator.On("ValidateItemCode", ctx, mock.Anything).Return(&masterdatadto.ValidateCodeResponse{Valid: true}, nil)
	tenantRepo := &mockTenantRepository{}
	tenantRepo.On("GetByCode", ctx, req.Organization).Return(tenant, nil)
	productRepo := &mockProductRepository{}
	productRepo.On("GetByCodeAndTenant", ctx, tenant.ID, "frendz-saving").Return(product, nil)
	configRepo := &mockProductRegistrationConfigRepository{}
	configRepo.On("GetByProductAndType", ctx, product.ID, "PARTICIPANT").Return(config, nil)
	profileRepo := &mockUserProfileRepository{}
	profileRepo.On("GetByUserID", ctx, req.UserID).Return(incompleteProfile, nil)

	uc := &usecase{
		masterdataUsecase: mdValidator,
		tenantRepo:        tenantRepo,
		productRepo:       productRepo,
		configRepo:        configRepo,
		userProfileRepo:   profileRepo,
	}
	resp, err := uc.SelfRegister(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	appErr := apperrors.GetAppError(err)
	require.NotNil(t, appErr)
	assert.Equal(t, apperrors.CodeUnprocessable, appErr.Code)
}

func TestSelfRegister_AlreadyLinked_DifferentUser(t *testing.T) {
	ctx := context.Background()
	req := validSelfRegisterRequest()
	tenant := activeTenant()
	product := activeProduct(tenant.ID)
	config := activeConfig(product.ID)
	profile := completeProfile(req.UserID)
	existingUserID := uuid.New()
	linkedParticipant := &entity.Participant{
		ID:        uuid.New(),
		TenantID:  tenant.ID,
		ProductID: product.ID,
		UserID:    &existingUserID,
		FullName:  "Someone Else",
		Status:    entity.ParticipantStatusApproved,
	}
	existingPension := &entity.ParticipantPension{
		ID:            uuid.New(),
		ParticipantID: linkedParticipant.ID,
	}

	mdValidator := &mockMasterdataValidator{}
	mdValidator.On("ValidateItemCode", ctx, mock.Anything).Return(&masterdatadto.ValidateCodeResponse{Valid: true}, nil)
	tenantRepo := &mockTenantRepository{}
	tenantRepo.On("GetByCode", ctx, req.Organization).Return(tenant, nil)
	productRepo := &mockProductRepository{}
	productRepo.On("GetByCodeAndTenant", ctx, tenant.ID, "frendz-saving").Return(product, nil)
	configRepo := &mockProductRegistrationConfigRepository{}
	configRepo.On("GetByProductAndType", ctx, product.ID, "PARTICIPANT").Return(config, nil)
	profileRepo := &mockUserProfileRepository{}
	profileRepo.On("GetByUserID", ctx, req.UserID).Return(profile, nil)
	partRepo := &mockParticipantRepositoryWithKTP{}
	partRepo.On("GetByKTPAndPensionNumber", ctx, req.IdentityNumber, req.ParticipantNumber, tenant.ID, product.ID).
		Return(linkedParticipant, existingPension, nil)

	uc := &usecase{
		masterdataUsecase: mdValidator,
		tenantRepo:        tenantRepo,
		productRepo:       productRepo,
		configRepo:        configRepo,
		userProfileRepo:   profileRepo,
		participantRepo:   partRepo,
	}
	resp, err := uc.SelfRegister(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	appErr := apperrors.GetAppError(err)
	require.NotNil(t, appErr)
	assert.Equal(t, apperrors.CodeConflict, appErr.Code)

	assert.Equal(t, "registration not eligible", appErr.Message)
}

func TestSelfRegister_AlreadyLinked_SameUser(t *testing.T) {
	ctx := context.Background()
	req := validSelfRegisterRequest()
	tenant := activeTenant()
	product := activeProduct(tenant.ID)
	config := activeConfig(product.ID)
	profile := completeProfile(req.UserID)

	sameUserID := req.UserID
	linkedParticipant := &entity.Participant{
		ID:        uuid.New(),
		TenantID:  tenant.ID,
		ProductID: product.ID,
		UserID:    &sameUserID,
		FullName:  "Same User",
		Status:    entity.ParticipantStatusApproved,
	}
	existingPension := &entity.ParticipantPension{
		ID:            uuid.New(),
		ParticipantID: linkedParticipant.ID,
	}

	mdValidator := &mockMasterdataValidator{}
	mdValidator.On("ValidateItemCode", ctx, mock.Anything).Return(&masterdatadto.ValidateCodeResponse{Valid: true}, nil)
	tenantRepo := &mockTenantRepository{}
	tenantRepo.On("GetByCode", ctx, req.Organization).Return(tenant, nil)
	productRepo := &mockProductRepository{}
	productRepo.On("GetByCodeAndTenant", ctx, tenant.ID, "frendz-saving").Return(product, nil)
	configRepo := &mockProductRegistrationConfigRepository{}
	configRepo.On("GetByProductAndType", ctx, product.ID, "PARTICIPANT").Return(config, nil)
	profileRepo := &mockUserProfileRepository{}
	profileRepo.On("GetByUserID", ctx, req.UserID).Return(profile, nil)
	partRepo := &mockParticipantRepositoryWithKTP{}
	partRepo.On("GetByKTPAndPensionNumber", ctx, req.IdentityNumber, req.ParticipantNumber, tenant.ID, product.ID).
		Return(linkedParticipant, existingPension, nil)

	uc := &usecase{
		masterdataUsecase: mdValidator,
		tenantRepo:        tenantRepo,
		productRepo:       productRepo,
		configRepo:        configRepo,
		userProfileRepo:   profileRepo,
		participantRepo:   partRepo,
	}
	resp, err := uc.SelfRegister(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	appErr := apperrors.GetAppError(err)
	require.NotNil(t, appErr)
	assert.Equal(t, apperrors.CodeConflict, appErr.Code)

	assert.Equal(t, "registration not eligible", appErr.Message)
}

func TestSelfRegister_AlreadyRegistered_UTRExists(t *testing.T) {
	ctx := context.Background()
	req := validSelfRegisterRequest()
	tenant := activeTenant()
	product := activeProduct(tenant.ID)
	config := activeConfig(product.ID)
	profile := completeProfile(req.UserID)

	mdValidator := &mockMasterdataValidator{}
	mdValidator.On("ValidateItemCode", ctx, mock.Anything).Return(&masterdatadto.ValidateCodeResponse{Valid: true}, nil)
	tenantRepo := &mockTenantRepository{}
	tenantRepo.On("GetByCode", ctx, req.Organization).Return(tenant, nil)
	productRepo := &mockProductRepository{}
	productRepo.On("GetByCodeAndTenant", ctx, tenant.ID, "frendz-saving").Return(product, nil)
	configRepo := &mockProductRegistrationConfigRepository{}
	configRepo.On("GetByProductAndType", ctx, product.ID, "PARTICIPANT").Return(config, nil)
	profileRepo := &mockUserProfileRepository{}
	profileRepo.On("GetByUserID", ctx, req.UserID).Return(profile, nil)
	partRepo := &mockParticipantRepositoryWithKTP{}
	partRepo.On("GetByKTPAndPensionNumber", ctx, req.IdentityNumber, req.ParticipantNumber, tenant.ID, product.ID).
		Return(nil, nil, apperrors.ErrNotFound("not found"))
	utrRepo := &mockUserTenantRegistrationRepository{}
	utrRepo.On("GetByUserAndProduct", ctx, req.UserID, tenant.ID, product.ID, "PARTICIPANT").
		Return(&entity.UserTenantRegistration{ID: uuid.New()}, nil)

	uc := &usecase{
		masterdataUsecase: mdValidator,
		tenantRepo:        tenantRepo,
		productRepo:       productRepo,
		configRepo:        configRepo,
		userProfileRepo:   profileRepo,
		participantRepo:   partRepo,
		utrRepo:           utrRepo,
	}
	resp, err := uc.SelfRegister(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	appErr := apperrors.GetAppError(err)
	require.NotNil(t, appErr)
	assert.Equal(t, apperrors.CodeConflict, appErr.Code)

	assert.Equal(t, "registration not eligible", appErr.Message)
}

func TestSelfRegister_UniqueConstraintViolation_OnCreate(t *testing.T) {
	ctx := context.Background()
	req := validSelfRegisterRequest()
	tenant := activeTenant()
	product := activeProduct(tenant.ID)
	config := activeConfig(product.ID)
	profile := completeProfile(req.UserID)

	txMgr := &MockTransactionManager{}
	partRepo := &mockParticipantRepositoryWithKTP{}
	pensionRepo := &MockParticipantPensionRepository{}
	statusHistoryRepo := &MockParticipantStatusHistoryRepository{}
	tenantRepo := &mockTenantRepository{}
	productRepo := &mockProductRepository{}
	configRepo := &mockProductRegistrationConfigRepository{}
	utrRepo := &mockUserTenantRegistrationRepository{}
	profileRepo := &mockUserProfileRepository{}
	mdValidator := &mockMasterdataValidator{}

	mdValidator.On("ValidateItemCode", ctx, mock.Anything).Return(&masterdatadto.ValidateCodeResponse{Valid: true}, nil)
	tenantRepo.On("GetByCode", ctx, req.Organization).Return(tenant, nil)
	productRepo.On("GetByCodeAndTenant", ctx, tenant.ID, "frendz-saving").Return(product, nil)
	configRepo.On("GetByProductAndType", ctx, product.ID, "PARTICIPANT").Return(config, nil)
	profileRepo.On("GetByUserID", ctx, req.UserID).Return(profile, nil)
	partRepo.On("GetByKTPAndPensionNumber", ctx, req.IdentityNumber, req.ParticipantNumber, tenant.ID, product.ID).
		Return(nil, nil, apperrors.ErrNotFound("not found"))
	utrRepo.On("GetByUserAndProduct", ctx, req.UserID, tenant.ID, product.ID, "PARTICIPANT").
		Return(nil, apperrors.ErrNotFound("not found"))

	conflictErr := apperrors.ErrConflict("registration not eligible")
	txMgr.On("WithTransaction", ctx, mock.Anything).Return(conflictErr)

	uc := buildSelfRegisterUsecase(txMgr, partRepo, pensionRepo, statusHistoryRepo, tenantRepo, productRepo, configRepo, utrRepo, profileRepo, mdValidator)

	resp, err := uc.SelfRegister(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	appErr := apperrors.GetAppError(err)
	require.NotNil(t, appErr)
	assert.Equal(t, apperrors.CodeConflict, appErr.Code)
	assert.Equal(t, "registration not eligible", appErr.Message)
}

func TestSelfRegister_TransactionRollback_OnUTRCreateFailure(t *testing.T) {
	ctx := context.Background()
	req := validSelfRegisterRequest()
	tenant := activeTenant()
	product := activeProduct(tenant.ID)
	config := activeConfig(product.ID)
	profile := completeProfile(req.UserID)

	txMgr := &MockTransactionManager{}
	partRepo := &mockParticipantRepositoryWithKTP{}
	pensionRepo := &MockParticipantPensionRepository{}
	statusHistoryRepo := &MockParticipantStatusHistoryRepository{}
	tenantRepo := &mockTenantRepository{}
	productRepo := &mockProductRepository{}
	configRepo := &mockProductRegistrationConfigRepository{}
	utrRepo := &mockUserTenantRegistrationRepository{}
	profileRepo := &mockUserProfileRepository{}
	mdValidator := &mockMasterdataValidator{}

	mdValidator.On("ValidateItemCode", ctx, mock.Anything).Return(&masterdatadto.ValidateCodeResponse{Valid: true}, nil)
	tenantRepo.On("GetByCode", ctx, req.Organization).Return(tenant, nil)
	productRepo.On("GetByCodeAndTenant", ctx, tenant.ID, "frendz-saving").Return(product, nil)
	configRepo.On("GetByProductAndType", ctx, product.ID, "PARTICIPANT").Return(config, nil)
	profileRepo.On("GetByUserID", ctx, req.UserID).Return(profile, nil)
	partRepo.On("GetByKTPAndPensionNumber", ctx, req.IdentityNumber, req.ParticipantNumber, tenant.ID, product.ID).
		Return(nil, nil, apperrors.ErrNotFound("not found"))
	utrRepo.On("GetByUserAndProduct", ctx, req.UserID, tenant.ID, product.ID, "PARTICIPANT").
		Return(nil, apperrors.ErrNotFound("not found"))

	txErr := apperrors.ErrInternal("db error")
	txMgr.On("WithTransaction", ctx, mock.Anything).Return(txErr)

	uc := buildSelfRegisterUsecase(txMgr, partRepo, pensionRepo, statusHistoryRepo, tenantRepo, productRepo, configRepo, utrRepo, profileRepo, mdValidator)

	resp, err := uc.SelfRegister(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
}
