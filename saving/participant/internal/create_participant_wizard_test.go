package internal

import (
	"context"
	"testing"
	"time"

	"erp-service/entity"
	"erp-service/pkg/errors"
	"erp-service/saving/participant/participantdto"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func makeCreateUsecase() (*usecase, *MockTransactionManager, *MockParticipantRepository, *MockParticipantStatusHistoryRepository) {
	txMgr := new(MockTransactionManager)
	participantRepo := new(MockParticipantRepository)
	statusHistoryRepo := new(MockParticipantStatusHistoryRepository)

	uc := &usecase{
		txManager:         txMgr,
		participantRepo:   participantRepo,
		statusHistoryRepo: statusHistoryRepo,
		identityRepo:      new(MockParticipantIdentityRepository),
		addressRepo:       new(MockParticipantAddressRepository),
		bankAccountRepo:   new(MockParticipantBankAccountRepository),
		familyMemberRepo:  new(MockParticipantFamilyMemberRepository),
		employmentRepo:    new(MockParticipantEmploymentRepository),
		pensionRepo:       new(MockParticipantPensionRepository),
		beneficiaryRepo:   new(MockParticipantBeneficiaryRepository),
		fileStorage:       new(MockFileStorageAdapter),
		fileRepo:          new(MockFileRepository),
	}
	return uc, txMgr, participantRepo, statusHistoryRepo
}

func TestCreateParticipant_DraftExists_ReturnsDraftExistsError(t *testing.T) {
	tenantID := uuid.New()
	productID := uuid.New()
	userID := uuid.New()
	existingID := uuid.New()

	uc, txMgr, participantRepo, _ := makeCreateUsecase()

	existingDraft := &entity.Participant{
		ID:             existingID,
		TenantID:       tenantID,
		ProductID:      productID,
		Status:         entity.ParticipantStatusDraft,
		FullName:       "Existing Draft",
		KTPNumber:      strPtr("1234567890123456"),
		EmployeeNumber: strPtr("EMP001"),
		Version:        1,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	txMgr.On("WithTransaction", mock.Anything, mock.AnythingOfType("func(context.Context) error")).Return(nil)

	participantRepo.On("GetByKTPNumber", mock.Anything, tenantID, productID, "1234567890123456").
		Return(existingDraft, nil)

	req := &participantdto.CreateParticipantRequest{
		TenantID:       tenantID,
		ProductID:      productID,
		UserID:         userID,
		FullName:       "New Participant",
		KTPNumber:      "1234567890123456",
		EmployeeNumber: "EMP001",
	}

	result, err := uc.CreateParticipant(context.Background(), req)
	require.Error(t, err)
	assert.Nil(t, result)

	appErr := errors.GetAppError(err)
	require.NotNil(t, appErr)
	assert.Equal(t, "PARTICIPANT_DRAFT_EXISTS", appErr.Code)

	assert.NotNil(t, appErr.Details)
	assert.Equal(t, existingID, appErr.Details["participant_id"])
}

func TestCreateParticipant_AlreadyRegistered_ReturnsAlreadyRegisteredError(t *testing.T) {
	tenantID := uuid.New()
	productID := uuid.New()
	userID := uuid.New()
	existingID := uuid.New()

	uc, txMgr, participantRepo, _ := makeCreateUsecase()

	existingApproved := &entity.Participant{
		ID:        existingID,
		TenantID:  tenantID,
		ProductID: productID,
		Status:    entity.ParticipantStatusApproved,
		FullName:  "Approved Participant",
		Version:   1,
	}

	txMgr.On("WithTransaction", mock.Anything, mock.AnythingOfType("func(context.Context) error")).Return(nil)
	participantRepo.On("GetByKTPNumber", mock.Anything, tenantID, productID, "1234567890123456").
		Return(existingApproved, nil)

	req := &participantdto.CreateParticipantRequest{
		TenantID:       tenantID,
		ProductID:      productID,
		UserID:         userID,
		FullName:       "New Participant",
		KTPNumber:      "1234567890123456",
		EmployeeNumber: "EMP002",
	}

	result, err := uc.CreateParticipant(context.Background(), req)
	require.Error(t, err)
	assert.Nil(t, result)

	appErr := errors.GetAppError(err)
	require.NotNil(t, appErr)
	assert.Equal(t, "PARTICIPANT_ALREADY_REGISTERED", appErr.Code)
}

func TestCreateParticipant_EmployeeNumberConflict_AlreadyRegistered(t *testing.T) {
	tenantID := uuid.New()
	productID := uuid.New()
	userID := uuid.New()
	existingID := uuid.New()

	uc, txMgr, participantRepo, _ := makeCreateUsecase()

	existingApproved := &entity.Participant{
		ID:        existingID,
		TenantID:  tenantID,
		ProductID: productID,
		Status:    entity.ParticipantStatusApproved,
		FullName:  "Approved Participant",
		Version:   1,
	}

	txMgr.On("WithTransaction", mock.Anything, mock.AnythingOfType("func(context.Context) error")).Return(nil)

	participantRepo.On("GetByKTPNumber", mock.Anything, tenantID, productID, "9999999999999999").
		Return(nil, errors.ErrNotFound("participant not found"))
	participantRepo.On("GetByEmployeeNumber", mock.Anything, tenantID, productID, "EMP-TAKEN").
		Return(existingApproved, nil)

	req := &participantdto.CreateParticipantRequest{
		TenantID:       tenantID,
		ProductID:      productID,
		UserID:         userID,
		FullName:       "New Participant",
		KTPNumber:      "9999999999999999",
		EmployeeNumber: "EMP-TAKEN",
	}

	result, err := uc.CreateParticipant(context.Background(), req)
	require.Error(t, err)
	assert.Nil(t, result)

	appErr := errors.GetAppError(err)
	require.NotNil(t, appErr)
	assert.Equal(t, "PARTICIPANT_ALREADY_REGISTERED", appErr.Code)
}

func TestCreateParticipant_Success_StoresKTPAndEmployee(t *testing.T) {
	tenantID := uuid.New()
	productID := uuid.New()
	userID := uuid.New()

	uc, txMgr, participantRepo, statusHistoryRepo := makeCreateUsecase()

	txMgr.On("WithTransaction", mock.Anything, mock.AnythingOfType("func(context.Context) error")).Return(nil)
	participantRepo.On("GetByKTPNumber", mock.Anything, tenantID, productID, "1234567890123456").
		Return(nil, errors.ErrNotFound("not found"))
	participantRepo.On("GetByEmployeeNumber", mock.Anything, tenantID, productID, "EMP001").
		Return(nil, errors.ErrNotFound("not found"))

	var createdParticipant *entity.Participant
	participantRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.Participant")).
		Run(func(args mock.Arguments) {
			createdParticipant = args.Get(1).(*entity.Participant)
			createdParticipant.ID = uuid.New()
		}).
		Return(nil)
	statusHistoryRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.ParticipantStatusHistory")).
		Return(nil)

	uc.identityRepo.(*MockParticipantIdentityRepository).On("ListByParticipantID", mock.Anything, mock.AnythingOfType("uuid.UUID")).
		Return([]*entity.ParticipantIdentity{}, nil)
	uc.addressRepo.(*MockParticipantAddressRepository).On("ListByParticipantID", mock.Anything, mock.AnythingOfType("uuid.UUID")).
		Return([]*entity.ParticipantAddress{}, nil)
	uc.bankAccountRepo.(*MockParticipantBankAccountRepository).On("ListByParticipantID", mock.Anything, mock.AnythingOfType("uuid.UUID")).
		Return([]*entity.ParticipantBankAccount{}, nil)
	uc.familyMemberRepo.(*MockParticipantFamilyMemberRepository).On("ListByParticipantID", mock.Anything, mock.AnythingOfType("uuid.UUID")).
		Return([]*entity.ParticipantFamilyMember{}, nil)
	uc.employmentRepo.(*MockParticipantEmploymentRepository).On("GetByParticipantID", mock.Anything, mock.AnythingOfType("uuid.UUID")).
		Return(nil, errors.ErrNotFound("not found"))
	uc.pensionRepo.(*MockParticipantPensionRepository).On("GetByParticipantID", mock.Anything, mock.AnythingOfType("uuid.UUID")).
		Return(nil, errors.ErrNotFound("not found"))
	uc.beneficiaryRepo.(*MockParticipantBeneficiaryRepository).On("ListByParticipantID", mock.Anything, mock.AnythingOfType("uuid.UUID")).
		Return([]*entity.ParticipantBeneficiary{}, nil)

	req := &participantdto.CreateParticipantRequest{
		TenantID:       tenantID,
		ProductID:      productID,
		UserID:         userID,
		FullName:       "New Participant",
		KTPNumber:      "1234567890123456",
		EmployeeNumber: "EMP001",
	}

	result, err := uc.CreateParticipant(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, result)

	require.NotNil(t, createdParticipant)
	require.NotNil(t, createdParticipant.KTPNumber)
	assert.Equal(t, "1234567890123456", *createdParticipant.KTPNumber)
	require.NotNil(t, createdParticipant.EmployeeNumber)
	assert.Equal(t, "EMP001", *createdParticipant.EmployeeNumber)
}
