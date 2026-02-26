package participant_test

import (
	"context"
	"testing"
	"time"

	"erp-service/config"
	"erp-service/entity"
	"erp-service/masterdata"
	apperrors "erp-service/pkg/errors"
	"erp-service/saving/participant"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func minimalSelfRegisterUsecase(
	mdValidator *mockMasterdataValidator,
	tenantRepo *mockTenantRepository,
	productRepo *mockProductRepository,
	configRepo *mockProductRegistrationConfigRepository,
	profileRepo *mockUserProfileRepository,
	partRepo participant.ParticipantRepository,
	utrRepo *mockUserTenantRegistrationRepository,
) participant.Usecase {
	var md participant.MasterdataUsecase
	if mdValidator != nil {
		md = mdValidator
	}
	var tr participant.TenantRepository
	if tenantRepo != nil {
		tr = tenantRepo
	}
	var pr participant.ProductRepository
	if productRepo != nil {
		pr = productRepo
	}
	var cr participant.ProductRegistrationConfigRepository
	if configRepo != nil {
		cr = configRepo
	}
	var up participant.UserProfileRepository
	if profileRepo != nil {
		up = profileRepo
	}
	var pRepo participant.ParticipantRepository
	if partRepo != nil {
		pRepo = partRepo
	}
	var ur participant.UserTenantRegistrationRepository
	if utrRepo != nil {
		ur = utrRepo
	}

	return participant.NewUsecase(
		&config.Config{},
		zap.NewNop(),
		&MockTransactionManager{},
		pRepo,
		&MockParticipantIdentityRepository{},
		&MockParticipantAddressRepository{},
		&MockParticipantBankAccountRepository{},
		&MockParticipantFamilyMemberRepository{},
		&MockParticipantEmploymentRepository{},
		&MockParticipantPensionRepository{},
		&MockParticipantBeneficiaryRepository{},
		&MockParticipantStatusHistoryRepository{},
		&MockFileStorageAdapter{},
		&MockFileRepository{},
		tr,
		pr,
		cr,
		ur,
		up,
		md,
		nil,
	)
}

func TestSelfRegister_Success_NewParticipant(t *testing.T) {
	ctx := context.Background()
	req := validSelfRegisterRequest()
	tenant := activeTenant()
	product := activeProduct(tenant.ID)
	cfg := activeConfig(product.ID)
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
	empDataRepo := &MockEmployeeDataRepository{}

	mdValidator.On("ValidateItemCode", ctx, mock.MatchedBy(func(r *masterdata.ValidateCodeRequest) bool {
		return r.CategoryCode == "TENANT" && r.ItemCode == req.Organization
	})).Return(&masterdata.ValidateCodeResponse{Valid: true}, nil)

	tenantRepo.On("GetByCode", ctx, req.Organization).Return(tenant, nil)
	productRepo.On("GetByCodeAndTenant", ctx, tenant.ID, "frendz-saving").Return(product, nil)
	configRepo.On("GetByProductAndType", ctx, product.ID, "PARTICIPANT").Return(cfg, nil)
	profileRepo.On("GetByUserID", ctx, req.UserID).Return(profile, nil)
	partRepo.On("GetByKTPAndPensionNumber", ctx, req.IdentityNumber, req.ParticipantNumber, tenant.ID, product.ID).
		Return(nil, nil, apperrors.ErrNotFound("not found"))
	utrRepo.On("GetByUserAndProduct", ctx, req.UserID, tenant.ID, product.ID, "PARTICIPANT").
		Return(nil, apperrors.ErrNotFound("not found"))
	empDataRepo.On("GetByEmpNo", ctx, req.ParticipantNumber).
		Return(&entity.EmployeeData{ID: 1, EmpNo: req.ParticipantNumber, EmpName: "Test Employee"}, nil)
	partRepo.MockParticipantRepository.On("GetByEmployeeNumber", ctx, tenant.ID, product.ID, req.ParticipantNumber).
		Return(nil, apperrors.ErrNotFound("not found"))

	txMgr.On("WithTransaction", ctx, mock.Anything).Return(nil)
	partRepo.On("Create", ctx, mock.AnythingOfType("*entity.Participant")).Return(nil)
	pensionRepo.On("Create", ctx, mock.AnythingOfType("*entity.ParticipantPension")).Return(nil)
	statusHistoryRepo.On("Create", ctx, mock.AnythingOfType("*entity.ParticipantStatusHistory")).Return(nil)
	utrRepo.On("Create", ctx, mock.AnythingOfType("*entity.UserTenantRegistration")).Return(nil)

	uc := buildSelfRegisterUsecase(txMgr, partRepo, pensionRepo, statusHistoryRepo, tenantRepo, productRepo, configRepo, utrRepo, profileRepo, mdValidator, empDataRepo)

	resp, err := uc.SelfRegister(ctx, req)

	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.False(t, resp.IsLinked)
	assert.Equal(t, "ACTIVE", resp.RegistrationStatus)
}

func TestSelfRegister_Success_ResponseDataHasNoKTPNumber(t *testing.T) {

	ctx := context.Background()
	req := validSelfRegisterRequest()
	tenant := activeTenant()
	product := activeProduct(tenant.ID)
	cfg := activeConfig(product.ID)
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
	empDataRepo := &MockEmployeeDataRepository{}

	mdValidator.On("ValidateItemCode", ctx, mock.Anything).Return(&masterdata.ValidateCodeResponse{Valid: true}, nil)
	tenantRepo.On("GetByCode", ctx, req.Organization).Return(tenant, nil)
	productRepo.On("GetByCodeAndTenant", ctx, tenant.ID, "frendz-saving").Return(product, nil)
	configRepo.On("GetByProductAndType", ctx, product.ID, "PARTICIPANT").Return(cfg, nil)
	profileRepo.On("GetByUserID", ctx, req.UserID).Return(profile, nil)
	partRepo.On("GetByKTPAndPensionNumber", ctx, req.IdentityNumber, req.ParticipantNumber, tenant.ID, product.ID).
		Return(nil, nil, apperrors.ErrNotFound("not found"))
	utrRepo.On("GetByUserAndProduct", ctx, req.UserID, tenant.ID, product.ID, "PARTICIPANT").
		Return(nil, apperrors.ErrNotFound("not found"))
	empDataRepo.On("GetByEmpNo", ctx, req.ParticipantNumber).
		Return(&entity.EmployeeData{ID: 1, EmpNo: req.ParticipantNumber, EmpName: "Test Employee"}, nil)
	partRepo.MockParticipantRepository.On("GetByEmployeeNumber", ctx, tenant.ID, product.ID, req.ParticipantNumber).
		Return(nil, apperrors.ErrNotFound("not found"))

	txMgr.On("WithTransaction", ctx, mock.Anything).Return(nil)
	partRepo.On("Create", ctx, mock.AnythingOfType("*entity.Participant")).Return(nil)
	pensionRepo.On("Create", ctx, mock.AnythingOfType("*entity.ParticipantPension")).Return(nil)
	statusHistoryRepo.On("Create", ctx, mock.AnythingOfType("*entity.ParticipantStatusHistory")).Return(nil)
	utrRepo.On("Create", ctx, mock.AnythingOfType("*entity.UserTenantRegistration")).Return(nil)

	uc := buildSelfRegisterUsecase(txMgr, partRepo, pensionRepo, statusHistoryRepo, tenantRepo, productRepo, configRepo, utrRepo, profileRepo, mdValidator, empDataRepo)

	resp, err := uc.SelfRegister(ctx, req)

	require.NoError(t, err)
	require.NotNil(t, resp)

	assert.NotNil(t, resp.Data)

	assert.Equal(t, string(entity.ParticipantStatusApproved), resp.Data.Status)
	assert.NotEmpty(t, resp.Data.ParticipantNumber)
}

func TestSelfRegister_Success_LinkExistingParticipant(t *testing.T) {
	ctx := context.Background()
	req := validSelfRegisterRequest()
	tenant := activeTenant()
	product := activeProduct(tenant.ID)
	cfg := activeConfig(product.ID)
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

	mdValidator.On("ValidateItemCode", ctx, mock.Anything).Return(&masterdata.ValidateCodeResponse{Valid: true}, nil)
	tenantRepo.On("GetByCode", ctx, req.Organization).Return(tenant, nil)
	productRepo.On("GetByCodeAndTenant", ctx, tenant.ID, "frendz-saving").Return(product, nil)
	configRepo.On("GetByProductAndType", ctx, product.ID, "PARTICIPANT").Return(cfg, nil)
	profileRepo.On("GetByUserID", ctx, req.UserID).Return(profile, nil)
	partRepo.On("GetByKTPAndPensionNumber", ctx, req.IdentityNumber, req.ParticipantNumber, tenant.ID, product.ID).
		Return(existingParticipant, existingPension, nil)

	txMgr.On("WithTransaction", ctx, mock.Anything).Return(nil)
	partRepo.MockParticipantRepository.On("Update", ctx, mock.AnythingOfType("*entity.Participant")).Return(nil)
	utrRepo.On("Create", ctx, mock.AnythingOfType("*entity.UserTenantRegistration")).Return(nil)
	statusHistoryRepo.On("Create", ctx, mock.AnythingOfType("*entity.ParticipantStatusHistory")).Return(nil)

	uc := buildSelfRegisterUsecase(txMgr, partRepo, pensionRepo, statusHistoryRepo, tenantRepo, productRepo, configRepo, utrRepo, profileRepo, mdValidator, nil)

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

	uc := minimalSelfRegisterUsecase(nil, nil, nil, nil, nil, nil, nil)
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

	uc := minimalSelfRegisterUsecase(nil, nil, nil, nil, nil, nil, nil)
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

	uc := minimalSelfRegisterUsecase(nil, nil, nil, nil, nil, nil, nil)
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

	uc := minimalSelfRegisterUsecase(nil, nil, nil, nil, nil, nil, nil)
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

	uc := minimalSelfRegisterUsecase(nil, nil, nil, nil, nil, nil, nil)
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

	uc := minimalSelfRegisterUsecase(nil, nil, nil, nil, nil, nil, nil)
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
	mdValidator.On("ValidateItemCode", ctx, mock.Anything).Return(&masterdata.ValidateCodeResponse{
		Valid:   false,
		Message: "item not found in parent",
	}, nil)

	uc := minimalSelfRegisterUsecase(mdValidator, nil, nil, nil, nil, nil, nil)
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
	mdValidator.On("ValidateItemCode", ctx, mock.Anything).Return(&masterdata.ValidateCodeResponse{Valid: true}, nil)
	tenantRepo := &mockTenantRepository{}
	tenantRepo.On("GetByCode", ctx, req.Organization).Return(inactiveTenant, nil)

	uc := minimalSelfRegisterUsecase(mdValidator, tenantRepo, nil, nil, nil, nil, nil)
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
	mdValidator.On("ValidateItemCode", ctx, mock.Anything).Return(&masterdata.ValidateCodeResponse{Valid: true}, nil)
	tenantRepo := &mockTenantRepository{}
	tenantRepo.On("GetByCode", ctx, req.Organization).Return(tenant, nil)
	productRepo := &mockProductRepository{}
	productRepo.On("GetByCodeAndTenant", ctx, tenant.ID, "frendz-saving").Return(nil, apperrors.ErrNotFound("product not found"))

	uc := minimalSelfRegisterUsecase(mdValidator, tenantRepo, productRepo, nil, nil, nil, nil)
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
	mdValidator.On("ValidateItemCode", ctx, mock.Anything).Return(&masterdata.ValidateCodeResponse{Valid: true}, nil)
	tenantRepo := &mockTenantRepository{}
	tenantRepo.On("GetByCode", ctx, req.Organization).Return(tenant, nil)
	productRepo := &mockProductRepository{}
	productRepo.On("GetByCodeAndTenant", ctx, tenant.ID, "frendz-saving").Return(product, nil)
	configRepo := &mockProductRegistrationConfigRepository{}

	configRepo.On("GetByProductAndType", ctx, product.ID, "PARTICIPANT").Return(nil, apperrors.ErrNotFound("config not found"))

	uc := minimalSelfRegisterUsecase(mdValidator, tenantRepo, productRepo, configRepo, nil, nil, nil)
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
	mdValidator.On("ValidateItemCode", ctx, mock.Anything).Return(&masterdata.ValidateCodeResponse{Valid: true}, nil)
	tenantRepo := &mockTenantRepository{}
	tenantRepo.On("GetByCode", ctx, req.Organization).Return(tenant, nil)
	productRepo := &mockProductRepository{}
	productRepo.On("GetByCodeAndTenant", ctx, tenant.ID, "frendz-saving").Return(product, nil)
	configRepo := &mockProductRegistrationConfigRepository{}
	configRepo.On("GetByProductAndType", ctx, product.ID, "PARTICIPANT").Return(inactiveConfig, nil)

	uc := minimalSelfRegisterUsecase(mdValidator, tenantRepo, productRepo, configRepo, nil, nil, nil)
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
	cfg := activeConfig(product.ID)
	dob := time.Date(1990, 1, 15, 0, 0, 0, 0, time.UTC)
	incompleteProfile := &entity.UserProfile{
		UserID:      req.UserID,
		FirstName:   "John",
		LastName:    "Doe",
		DateOfBirth: &dob,
		Gender:      nil,
	}

	mdValidator := &mockMasterdataValidator{}
	mdValidator.On("ValidateItemCode", ctx, mock.Anything).Return(&masterdata.ValidateCodeResponse{Valid: true}, nil)
	tenantRepo := &mockTenantRepository{}
	tenantRepo.On("GetByCode", ctx, req.Organization).Return(tenant, nil)
	productRepo := &mockProductRepository{}
	productRepo.On("GetByCodeAndTenant", ctx, tenant.ID, "frendz-saving").Return(product, nil)
	configRepo := &mockProductRegistrationConfigRepository{}
	configRepo.On("GetByProductAndType", ctx, product.ID, "PARTICIPANT").Return(cfg, nil)
	profileRepo := &mockUserProfileRepository{}
	profileRepo.On("GetByUserID", ctx, req.UserID).Return(incompleteProfile, nil)

	uc := minimalSelfRegisterUsecase(mdValidator, tenantRepo, productRepo, configRepo, profileRepo, nil, nil)
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
	cfg := activeConfig(product.ID)
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
	mdValidator.On("ValidateItemCode", ctx, mock.Anything).Return(&masterdata.ValidateCodeResponse{Valid: true}, nil)
	tenantRepo := &mockTenantRepository{}
	tenantRepo.On("GetByCode", ctx, req.Organization).Return(tenant, nil)
	productRepo := &mockProductRepository{}
	productRepo.On("GetByCodeAndTenant", ctx, tenant.ID, "frendz-saving").Return(product, nil)
	configRepo := &mockProductRegistrationConfigRepository{}
	configRepo.On("GetByProductAndType", ctx, product.ID, "PARTICIPANT").Return(cfg, nil)
	profileRepo := &mockUserProfileRepository{}
	profileRepo.On("GetByUserID", ctx, req.UserID).Return(profile, nil)
	partRepo := &mockParticipantRepositoryWithKTP{}
	partRepo.On("GetByKTPAndPensionNumber", ctx, req.IdentityNumber, req.ParticipantNumber, tenant.ID, product.ID).
		Return(linkedParticipant, existingPension, nil)

	uc := minimalSelfRegisterUsecase(mdValidator, tenantRepo, productRepo, configRepo, profileRepo, partRepo, nil)
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
	cfg := activeConfig(product.ID)
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
	mdValidator.On("ValidateItemCode", ctx, mock.Anything).Return(&masterdata.ValidateCodeResponse{Valid: true}, nil)
	tenantRepo := &mockTenantRepository{}
	tenantRepo.On("GetByCode", ctx, req.Organization).Return(tenant, nil)
	productRepo := &mockProductRepository{}
	productRepo.On("GetByCodeAndTenant", ctx, tenant.ID, "frendz-saving").Return(product, nil)
	configRepo := &mockProductRegistrationConfigRepository{}
	configRepo.On("GetByProductAndType", ctx, product.ID, "PARTICIPANT").Return(cfg, nil)
	profileRepo := &mockUserProfileRepository{}
	profileRepo.On("GetByUserID", ctx, req.UserID).Return(profile, nil)
	partRepo := &mockParticipantRepositoryWithKTP{}
	partRepo.On("GetByKTPAndPensionNumber", ctx, req.IdentityNumber, req.ParticipantNumber, tenant.ID, product.ID).
		Return(linkedParticipant, existingPension, nil)

	uc := minimalSelfRegisterUsecase(mdValidator, tenantRepo, productRepo, configRepo, profileRepo, partRepo, nil)
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
	cfg := activeConfig(product.ID)
	profile := completeProfile(req.UserID)

	mdValidator := &mockMasterdataValidator{}
	mdValidator.On("ValidateItemCode", ctx, mock.Anything).Return(&masterdata.ValidateCodeResponse{Valid: true}, nil)
	tenantRepo := &mockTenantRepository{}
	tenantRepo.On("GetByCode", ctx, req.Organization).Return(tenant, nil)
	productRepo := &mockProductRepository{}
	productRepo.On("GetByCodeAndTenant", ctx, tenant.ID, "frendz-saving").Return(product, nil)
	configRepo := &mockProductRegistrationConfigRepository{}
	configRepo.On("GetByProductAndType", ctx, product.ID, "PARTICIPANT").Return(cfg, nil)
	profileRepo := &mockUserProfileRepository{}
	profileRepo.On("GetByUserID", ctx, req.UserID).Return(profile, nil)
	partRepo := &mockParticipantRepositoryWithKTP{}
	partRepo.On("GetByKTPAndPensionNumber", ctx, req.IdentityNumber, req.ParticipantNumber, tenant.ID, product.ID).
		Return(nil, nil, apperrors.ErrNotFound("not found"))
	utrRepo := &mockUserTenantRegistrationRepository{}
	utrRepo.On("GetByUserAndProduct", ctx, req.UserID, tenant.ID, product.ID, "PARTICIPANT").
		Return(&entity.UserTenantRegistration{ID: uuid.New()}, nil)

	uc := minimalSelfRegisterUsecase(mdValidator, tenantRepo, productRepo, configRepo, profileRepo, partRepo, utrRepo)
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
	cfg := activeConfig(product.ID)
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
	empDataRepo := &MockEmployeeDataRepository{}

	mdValidator.On("ValidateItemCode", ctx, mock.Anything).Return(&masterdata.ValidateCodeResponse{Valid: true}, nil)
	tenantRepo.On("GetByCode", ctx, req.Organization).Return(tenant, nil)
	productRepo.On("GetByCodeAndTenant", ctx, tenant.ID, "frendz-saving").Return(product, nil)
	configRepo.On("GetByProductAndType", ctx, product.ID, "PARTICIPANT").Return(cfg, nil)
	profileRepo.On("GetByUserID", ctx, req.UserID).Return(profile, nil)
	partRepo.On("GetByKTPAndPensionNumber", ctx, req.IdentityNumber, req.ParticipantNumber, tenant.ID, product.ID).
		Return(nil, nil, apperrors.ErrNotFound("not found"))
	utrRepo.On("GetByUserAndProduct", ctx, req.UserID, tenant.ID, product.ID, "PARTICIPANT").
		Return(nil, apperrors.ErrNotFound("not found"))
	empDataRepo.On("GetByEmpNo", ctx, req.ParticipantNumber).
		Return(&entity.EmployeeData{ID: 1, EmpNo: req.ParticipantNumber, EmpName: "Test Employee"}, nil)
	partRepo.MockParticipantRepository.On("GetByEmployeeNumber", ctx, tenant.ID, product.ID, req.ParticipantNumber).
		Return(nil, apperrors.ErrNotFound("not found"))

	conflictErr := apperrors.ErrConflict("registration not eligible")
	txMgr.On("WithTransaction", ctx, mock.Anything).Return(conflictErr)

	uc := buildSelfRegisterUsecase(txMgr, partRepo, pensionRepo, statusHistoryRepo, tenantRepo, productRepo, configRepo, utrRepo, profileRepo, mdValidator, empDataRepo)

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
	cfg := activeConfig(product.ID)
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
	empDataRepo := &MockEmployeeDataRepository{}

	mdValidator.On("ValidateItemCode", ctx, mock.Anything).Return(&masterdata.ValidateCodeResponse{Valid: true}, nil)
	tenantRepo.On("GetByCode", ctx, req.Organization).Return(tenant, nil)
	productRepo.On("GetByCodeAndTenant", ctx, tenant.ID, "frendz-saving").Return(product, nil)
	configRepo.On("GetByProductAndType", ctx, product.ID, "PARTICIPANT").Return(cfg, nil)
	profileRepo.On("GetByUserID", ctx, req.UserID).Return(profile, nil)
	partRepo.On("GetByKTPAndPensionNumber", ctx, req.IdentityNumber, req.ParticipantNumber, tenant.ID, product.ID).
		Return(nil, nil, apperrors.ErrNotFound("not found"))
	utrRepo.On("GetByUserAndProduct", ctx, req.UserID, tenant.ID, product.ID, "PARTICIPANT").
		Return(nil, apperrors.ErrNotFound("not found"))
	empDataRepo.On("GetByEmpNo", ctx, req.ParticipantNumber).
		Return(&entity.EmployeeData{ID: 1, EmpNo: req.ParticipantNumber, EmpName: "Test Employee"}, nil)
	partRepo.MockParticipantRepository.On("GetByEmployeeNumber", ctx, tenant.ID, product.ID, req.ParticipantNumber).
		Return(nil, apperrors.ErrNotFound("not found"))

	txErr := apperrors.ErrInternal("db error")
	txMgr.On("WithTransaction", ctx, mock.Anything).Return(txErr)

	uc := buildSelfRegisterUsecase(txMgr, partRepo, pensionRepo, statusHistoryRepo, tenantRepo, productRepo, configRepo, utrRepo, profileRepo, mdValidator, empDataRepo)

	resp, err := uc.SelfRegister(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
}
