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

func makeTestUsecase() (*usecase, *MockTransactionManager, *MockParticipantRepository, *MockParticipantAddressRepository, *MockParticipantStatusHistoryRepository) {
	txMgr := new(MockTransactionManager)
	participantRepo := new(MockParticipantRepository)
	addressRepo := new(MockParticipantAddressRepository)
	statusHistoryRepo := new(MockParticipantStatusHistoryRepository)

	uc := &usecase{
		txManager:         txMgr,
		participantRepo:   participantRepo,
		addressRepo:       addressRepo,
		statusHistoryRepo: statusHistoryRepo,
		identityRepo:      new(MockParticipantIdentityRepository),
		bankAccountRepo:   new(MockParticipantBankAccountRepository),
		familyMemberRepo:  new(MockParticipantFamilyMemberRepository),
		employmentRepo:    new(MockParticipantEmploymentRepository),
		pensionRepo:       new(MockParticipantPensionRepository),
		beneficiaryRepo:   new(MockParticipantBeneficiaryRepository),
		fileStorage:       new(MockFileStorageAdapter),
		fileRepo:          new(MockFileRepository),
	}
	return uc, txMgr, participantRepo, addressRepo, statusHistoryRepo
}

func makeDraftParticipant(tenantID, productID uuid.UUID) *entity.Participant {
	return &entity.Participant{
		ID:        uuid.New(),
		TenantID:  tenantID,
		ProductID: productID,
		Status:    entity.ParticipantStatusDraft,
		FullName:  "Test Participant",
		Version:   1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func TestSaveAddresses_EmptySlice_ReturnsBadRequest(t *testing.T) {
	tenantID := uuid.New()
	productID := uuid.New()
	participantID := uuid.New()

	uc, _, _, _, _ := makeTestUsecase()

	req := &participantdto.SaveAddressesRequest{
		TenantID:      tenantID,
		ProductID:     productID,
		ParticipantID: participantID,
		Addresses:     []participantdto.AddressItem{},
	}

	result, err := uc.SaveAddresses(context.Background(), req)
	require.Error(t, err)
	assert.Nil(t, result)

	var appErr *errors.AppError
	require.True(t, errors.As(err, &appErr))
	assert.Equal(t, errors.KindBadRequest, appErr.Kind)
}

func TestSaveAddresses_Success(t *testing.T) {
	tenantID := uuid.New()
	productID := uuid.New()
	userID := uuid.New()

	uc, txMgr, participantRepo, addressRepo, _ := makeTestUsecase()

	participant := makeDraftParticipant(tenantID, productID)
	participantID := participant.ID

	txMgr.On("WithTransaction", mock.Anything, mock.AnythingOfType("func(context.Context) error")).Return(nil)
	participantRepo.On("GetByID", mock.Anything, participantID).Return(participant, nil)
	participantRepo.On("Update", mock.Anything, mock.AnythingOfType("*entity.Participant")).Return(nil)
	addressRepo.On("SoftDeleteAllByParticipantID", mock.Anything, participantID).Return(nil)
	addressRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.ParticipantAddress")).Return(nil)

	req := &participantdto.SaveAddressesRequest{
		TenantID:      tenantID,
		ProductID:     productID,
		ParticipantID: participantID,
		UserID:        userID,
		Addresses: []participantdto.AddressItem{
			{
				AddressType: "KTP",
				IsPrimary:   true,
			},
		},
	}

	result, err := uc.SaveAddresses(context.Background(), req)
	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "KTP", result[0].AddressType)

	addressRepo.AssertNotCalled(t, "ListByParticipantID", mock.Anything, mock.Anything)
}

func TestSaveAddresses_ParticipantNotFound(t *testing.T) {
	tenantID := uuid.New()
	productID := uuid.New()
	participantID := uuid.New()

	uc, txMgr, participantRepo, _, _ := makeTestUsecase()

	txMgr.On("WithTransaction", mock.Anything, mock.AnythingOfType("func(context.Context) error")).Return(nil)
	participantRepo.On("GetByID", mock.Anything, participantID).Return(nil, errors.ErrNotFound("participant not found"))

	req := &participantdto.SaveAddressesRequest{
		TenantID:      tenantID,
		ProductID:     productID,
		ParticipantID: participantID,
		Addresses:     []participantdto.AddressItem{{AddressType: "KTP"}},
	}

	result, err := uc.SaveAddresses(context.Background(), req)
	require.Error(t, err)
	assert.Nil(t, result)
	assert.True(t, errors.IsNotFound(err))
}

func TestSaveAddresses_WrongTenant(t *testing.T) {
	tenantID := uuid.New()
	otherTenantID := uuid.New()
	productID := uuid.New()
	participantID := uuid.New()

	uc, txMgr, participantRepo, _, _ := makeTestUsecase()

	participant := &entity.Participant{
		ID:        participantID,
		TenantID:  otherTenantID,
		ProductID: productID,
		Status:    entity.ParticipantStatusDraft,
	}

	txMgr.On("WithTransaction", mock.Anything, mock.AnythingOfType("func(context.Context) error")).Return(nil)
	participantRepo.On("GetByID", mock.Anything, participantID).Return(participant, nil)

	req := &participantdto.SaveAddressesRequest{
		TenantID:      tenantID,
		ProductID:     productID,
		ParticipantID: participantID,
		Addresses:     []participantdto.AddressItem{{AddressType: "KTP"}},
	}

	result, err := uc.SaveAddresses(context.Background(), req)
	require.Error(t, err)
	assert.Nil(t, result)
}

func TestSaveAddresses_NonEditableStatus(t *testing.T) {
	tenantID := uuid.New()
	productID := uuid.New()
	participantID := uuid.New()

	uc, txMgr, participantRepo, _, _ := makeTestUsecase()

	participant := &entity.Participant{
		ID:        participantID,
		TenantID:  tenantID,
		ProductID: productID,
		Status:    entity.ParticipantStatusApproved,
	}

	txMgr.On("WithTransaction", mock.Anything, mock.AnythingOfType("func(context.Context) error")).Return(nil)
	participantRepo.On("GetByID", mock.Anything, participantID).Return(participant, nil)

	req := &participantdto.SaveAddressesRequest{
		TenantID:      tenantID,
		ProductID:     productID,
		ParticipantID: participantID,
		Addresses:     []participantdto.AddressItem{{AddressType: "KTP"}},
	}

	result, err := uc.SaveAddresses(context.Background(), req)
	require.Error(t, err)
	assert.Nil(t, result)
}

func TestSaveAddresses_SoftDeleteFailure(t *testing.T) {
	tenantID := uuid.New()
	productID := uuid.New()
	participantID := uuid.New()

	uc, txMgr, participantRepo, addressRepo, _ := makeTestUsecase()

	participant := makeDraftParticipant(tenantID, productID)
	participant.ID = participantID

	txMgr.On("WithTransaction", mock.Anything, mock.AnythingOfType("func(context.Context) error")).Return(nil)
	participantRepo.On("GetByID", mock.Anything, participantID).Return(participant, nil)
	addressRepo.On("SoftDeleteAllByParticipantID", mock.Anything, participantID).Return(errors.ErrInternal("db error"))

	req := &participantdto.SaveAddressesRequest{
		TenantID:      tenantID,
		ProductID:     productID,
		ParticipantID: participantID,
		Addresses:     []participantdto.AddressItem{{AddressType: "KTP"}},
	}

	result, err := uc.SaveAddresses(context.Background(), req)
	require.Error(t, err)
	assert.Nil(t, result)
}

func TestSaveAddresses_StepMarkedComplete(t *testing.T) {
	tenantID := uuid.New()
	productID := uuid.New()
	userID := uuid.New()

	uc, txMgr, participantRepo, addressRepo, _ := makeTestUsecase()

	participant := makeDraftParticipant(tenantID, productID)
	participantID := participant.ID

	var savedParticipant *entity.Participant

	txMgr.On("WithTransaction", mock.Anything, mock.AnythingOfType("func(context.Context) error")).Return(nil)
	participantRepo.On("GetByID", mock.Anything, participantID).Return(participant, nil)
	participantRepo.On("Update", mock.Anything, mock.AnythingOfType("*entity.Participant")).
		Run(func(args mock.Arguments) {
			savedParticipant = args.Get(1).(*entity.Participant)
		}).
		Return(nil)
	addressRepo.On("SoftDeleteAllByParticipantID", mock.Anything, participantID).Return(nil)
	addressRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.ParticipantAddress")).Return(nil)

	req := &participantdto.SaveAddressesRequest{
		TenantID:      tenantID,
		ProductID:     productID,
		ParticipantID: participantID,
		UserID:        userID,
		Addresses:     []participantdto.AddressItem{{AddressType: "KTP"}},
	}

	result, err := uc.SaveAddresses(context.Background(), req)
	require.NoError(t, err)
	assert.Len(t, result, 1)

	require.NotNil(t, savedParticipant)
	assert.True(t, savedParticipant.StepsCompleted["address"], "address step should be marked complete")

	addressRepo.AssertNotCalled(t, "ListByParticipantID", mock.Anything, mock.Anything)
}
