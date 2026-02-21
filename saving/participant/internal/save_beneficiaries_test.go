package internal

import (
	"context"
	"testing"

	"erp-service/entity"
	"erp-service/pkg/errors"
	"erp-service/saving/participant/participantdto"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func makeBeneficiaryUsecase() (*usecase, *MockTransactionManager, *MockParticipantRepository, *MockParticipantBeneficiaryRepository, *MockParticipantFamilyMemberRepository, *MockFileRepository) {
	txMgr := new(MockTransactionManager)
	participantRepo := new(MockParticipantRepository)
	beneficiaryRepo := new(MockParticipantBeneficiaryRepository)
	familyMemberRepo := new(MockParticipantFamilyMemberRepository)
	fileRepo := new(MockFileRepository)

	uc := &usecase{
		txManager:         txMgr,
		participantRepo:   participantRepo,
		beneficiaryRepo:   beneficiaryRepo,
		familyMemberRepo:  familyMemberRepo,
		identityRepo:      new(MockParticipantIdentityRepository),
		addressRepo:       new(MockParticipantAddressRepository),
		bankAccountRepo:   new(MockParticipantBankAccountRepository),
		employmentRepo:    new(MockParticipantEmploymentRepository),
		pensionRepo:       new(MockParticipantPensionRepository),
		statusHistoryRepo: new(MockParticipantStatusHistoryRepository),
		fileStorage:       new(MockFileStorageAdapter),
		fileRepo:          fileRepo,
	}
	return uc, txMgr, participantRepo, beneficiaryRepo, familyMemberRepo, fileRepo
}

func TestSaveBeneficiaries_EmptySlice_ReturnsBadRequest(t *testing.T) {
	tenantID := uuid.New()
	productID := uuid.New()
	participantID := uuid.New()

	uc, _, _, _, _, _ := makeBeneficiaryUsecase()

	req := &participantdto.SaveBeneficiariesRequest{
		TenantID:      tenantID,
		ProductID:     productID,
		ParticipantID: participantID,
		Beneficiaries: []participantdto.BeneficiaryItem{},
	}

	result, err := uc.SaveBeneficiaries(context.Background(), req)
	require.Error(t, err)
	assert.Nil(t, result)

	var appErr *errors.AppError
	require.True(t, errors.As(err, &appErr))
	assert.Equal(t, errors.KindBadRequest, appErr.Kind)
}

func TestSaveBeneficiaries_Success(t *testing.T) {
	tenantID := uuid.New()
	productID := uuid.New()
	userID := uuid.New()
	familyMemberID := uuid.New()

	uc, txMgr, participantRepo, beneficiaryRepo, familyMemberRepo, _ := makeBeneficiaryUsecase()

	participant := makeDraftParticipant(tenantID, productID)
	participantID := participant.ID

	familyMember := &entity.ParticipantFamilyMember{
		ID:            familyMemberID,
		ParticipantID: participantID,
		FullName:      "Jane Doe",
	}

	txMgr.On("WithTransaction", mock.Anything, mock.AnythingOfType("func(context.Context) error")).Return(nil)
	participantRepo.On("GetByID", mock.Anything, participantID).Return(participant, nil)
	participantRepo.On("Update", mock.Anything, mock.AnythingOfType("*entity.Participant")).Return(nil)
	beneficiaryRepo.On("SoftDeleteAllByParticipantID", mock.Anything, participantID).Return(nil)
	familyMemberRepo.On("GetByID", mock.Anything, familyMemberID).Return(familyMember, nil)
	beneficiaryRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.ParticipantBeneficiary")).Return(nil)

	req := &participantdto.SaveBeneficiariesRequest{
		TenantID:      tenantID,
		ProductID:     productID,
		ParticipantID: participantID,
		UserID:        userID,
		Beneficiaries: []participantdto.BeneficiaryItem{
			{
				FamilyMemberID: familyMemberID,
			},
		},
	}

	result, err := uc.SaveBeneficiaries(context.Background(), req)
	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, familyMemberID, result[0].FamilyMemberID)

	beneficiaryRepo.AssertNotCalled(t, "ListByParticipantID", mock.Anything, mock.Anything)
}

func TestSaveBeneficiaries_FamilyMemberBelongsToOtherParticipant_ReturnsForbidden(t *testing.T) {
	tenantID := uuid.New()
	productID := uuid.New()
	userID := uuid.New()
	familyMemberID := uuid.New()
	otherParticipantID := uuid.New()

	uc, txMgr, participantRepo, beneficiaryRepo, familyMemberRepo, _ := makeBeneficiaryUsecase()

	participant := makeDraftParticipant(tenantID, productID)
	participantID := participant.ID

	familyMember := &entity.ParticipantFamilyMember{
		ID:            familyMemberID,
		ParticipantID: otherParticipantID,
		FullName:      "Other Person",
	}

	txMgr.On("WithTransaction", mock.Anything, mock.AnythingOfType("func(context.Context) error")).Return(nil)
	participantRepo.On("GetByID", mock.Anything, participantID).Return(participant, nil)
	beneficiaryRepo.On("SoftDeleteAllByParticipantID", mock.Anything, participantID).Return(nil)
	familyMemberRepo.On("GetByID", mock.Anything, familyMemberID).Return(familyMember, nil)

	req := &participantdto.SaveBeneficiariesRequest{
		TenantID:      tenantID,
		ProductID:     productID,
		ParticipantID: participantID,
		UserID:        userID,
		Beneficiaries: []participantdto.BeneficiaryItem{
			{FamilyMemberID: familyMemberID},
		},
	}

	result, err := uc.SaveBeneficiaries(context.Background(), req)
	require.Error(t, err)
	assert.Nil(t, result)

	var appErr *errors.AppError
	require.True(t, errors.As(err, &appErr))
	assert.Equal(t, errors.KindForbidden, appErr.Kind)
}

func TestSaveBeneficiaries_FileOwnershipMismatch_ReturnsForbidden(t *testing.T) {
	tenantID := uuid.New()
	productID := uuid.New()
	otherTenantID := uuid.New()
	userID := uuid.New()
	familyMemberID := uuid.New()
	fileID := uuid.New()

	uc, txMgr, participantRepo, beneficiaryRepo, familyMemberRepo, fileRepo := makeBeneficiaryUsecase()

	participant := makeDraftParticipant(tenantID, productID)
	participantID := participant.ID

	familyMember := &entity.ParticipantFamilyMember{
		ID:            familyMemberID,
		ParticipantID: participantID,
		FullName:      "Jane Doe",
	}

	file := &entity.File{
		ID:       fileID,
		TenantID: otherTenantID,
	}

	txMgr.On("WithTransaction", mock.Anything, mock.AnythingOfType("func(context.Context) error")).Return(nil)
	participantRepo.On("GetByID", mock.Anything, participantID).Return(participant, nil)
	beneficiaryRepo.On("SoftDeleteAllByParticipantID", mock.Anything, participantID).Return(nil)
	familyMemberRepo.On("GetByID", mock.Anything, familyMemberID).Return(familyMember, nil)
	fileRepo.On("GetByID", mock.Anything, fileID).Return(file, nil)

	req := &participantdto.SaveBeneficiariesRequest{
		TenantID:      tenantID,
		ProductID:     productID,
		ParticipantID: participantID,
		UserID:        userID,
		Beneficiaries: []participantdto.BeneficiaryItem{
			{
				FamilyMemberID:      familyMemberID,
				IdentityPhotoFileID: &fileID,
			},
		},
	}

	result, err := uc.SaveBeneficiaries(context.Background(), req)
	require.Error(t, err)
	assert.Nil(t, result)

	var appErr *errors.AppError
	require.True(t, errors.As(err, &appErr))
	assert.Equal(t, errors.KindForbidden, appErr.Kind)
}

func TestSaveBeneficiaries_WithFileIDs_SetsPermanent(t *testing.T) {
	tenantID := uuid.New()
	productID := uuid.New()
	userID := uuid.New()
	familyMemberID := uuid.New()
	identityPhotoFileID := uuid.New()

	uc, txMgr, participantRepo, beneficiaryRepo, familyMemberRepo, fileRepo := makeBeneficiaryUsecase()

	participant := makeDraftParticipant(tenantID, productID)
	participantID := participant.ID

	familyMember := &entity.ParticipantFamilyMember{
		ID:            familyMemberID,
		ParticipantID: participantID,
		FullName:      "Jane",
	}

	file := &entity.File{
		ID:        identityPhotoFileID,
		TenantID:  tenantID,
		ProductID: productID,
	}

	txMgr.On("WithTransaction", mock.Anything, mock.AnythingOfType("func(context.Context) error")).Return(nil)
	participantRepo.On("GetByID", mock.Anything, participantID).Return(participant, nil)
	participantRepo.On("Update", mock.Anything, mock.AnythingOfType("*entity.Participant")).Return(nil)
	beneficiaryRepo.On("SoftDeleteAllByParticipantID", mock.Anything, participantID).Return(nil)
	familyMemberRepo.On("GetByID", mock.Anything, familyMemberID).Return(familyMember, nil)
	fileRepo.On("GetByID", mock.Anything, identityPhotoFileID).Return(file, nil)
	fileRepo.On("SetPermanent", mock.Anything, identityPhotoFileID).Return(nil)
	beneficiaryRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.ParticipantBeneficiary")).Return(nil)

	req := &participantdto.SaveBeneficiariesRequest{
		TenantID:      tenantID,
		ProductID:     productID,
		ParticipantID: participantID,
		UserID:        userID,
		Beneficiaries: []participantdto.BeneficiaryItem{
			{
				FamilyMemberID:      familyMemberID,
				IdentityPhotoFileID: &identityPhotoFileID,
			},
		},
	}

	result, err := uc.SaveBeneficiaries(context.Background(), req)
	require.NoError(t, err)
	assert.Len(t, result, 1)

	fileRepo.AssertCalled(t, "SetPermanent", mock.Anything, identityPhotoFileID)

	beneficiaryRepo.AssertNotCalled(t, "ListByParticipantID", mock.Anything, mock.Anything)
}

func TestSaveBeneficiaries_ParticipantNotFound(t *testing.T) {
	tenantID := uuid.New()
	productID := uuid.New()
	participantID := uuid.New()

	uc, txMgr, participantRepo, _, _, _ := makeBeneficiaryUsecase()

	txMgr.On("WithTransaction", mock.Anything, mock.AnythingOfType("func(context.Context) error")).Return(nil)
	participantRepo.On("GetByID", mock.Anything, participantID).Return(nil, errors.ErrNotFound("participant not found"))

	req := &participantdto.SaveBeneficiariesRequest{
		TenantID:      tenantID,
		ProductID:     productID,
		ParticipantID: participantID,
		Beneficiaries: []participantdto.BeneficiaryItem{{FamilyMemberID: uuid.New()}},
	}

	result, err := uc.SaveBeneficiaries(context.Background(), req)
	require.Error(t, err)
	assert.Nil(t, result)
	assert.True(t, errors.IsNotFound(err))
}

func TestSaveBeneficiaries_StepMarkedComplete(t *testing.T) {
	tenantID := uuid.New()
	productID := uuid.New()
	userID := uuid.New()
	familyMemberID := uuid.New()

	uc, txMgr, participantRepo, beneficiaryRepo, familyMemberRepo, _ := makeBeneficiaryUsecase()

	participant := makeDraftParticipant(tenantID, productID)
	participantID := participant.ID

	familyMember := &entity.ParticipantFamilyMember{
		ID:            familyMemberID,
		ParticipantID: participantID,
	}

	var savedParticipant *entity.Participant

	txMgr.On("WithTransaction", mock.Anything, mock.AnythingOfType("func(context.Context) error")).Return(nil)
	participantRepo.On("GetByID", mock.Anything, participantID).Return(participant, nil)
	participantRepo.On("Update", mock.Anything, mock.AnythingOfType("*entity.Participant")).
		Run(func(args mock.Arguments) {
			savedParticipant = args.Get(1).(*entity.Participant)
		}).
		Return(nil)
	beneficiaryRepo.On("SoftDeleteAllByParticipantID", mock.Anything, participantID).Return(nil)
	familyMemberRepo.On("GetByID", mock.Anything, familyMemberID).Return(familyMember, nil)
	beneficiaryRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.ParticipantBeneficiary")).Return(nil)

	req := &participantdto.SaveBeneficiariesRequest{
		TenantID:      tenantID,
		ProductID:     productID,
		ParticipantID: participantID,
		UserID:        userID,
		Beneficiaries: []participantdto.BeneficiaryItem{{FamilyMemberID: familyMemberID}},
	}

	result, err := uc.SaveBeneficiaries(context.Background(), req)
	require.NoError(t, err)
	assert.Len(t, result, 1)

	require.NotNil(t, savedParticipant)
	assert.True(t, savedParticipant.StepsCompleted["beneficiaries"], "beneficiaries step should be marked complete")

	beneficiaryRepo.AssertNotCalled(t, "ListByParticipantID", mock.Anything, mock.Anything)
}

func TestSaveBeneficiaries_SoftDeleteFailure(t *testing.T) {
	tenantID := uuid.New()
	productID := uuid.New()
	participantID := uuid.New()

	uc, txMgr, participantRepo, beneficiaryRepo, _, _ := makeBeneficiaryUsecase()

	participant := makeDraftParticipant(tenantID, productID)
	participant.ID = participantID

	txMgr.On("WithTransaction", mock.Anything, mock.AnythingOfType("func(context.Context) error")).Return(nil)
	participantRepo.On("GetByID", mock.Anything, participantID).Return(participant, nil)
	beneficiaryRepo.On("SoftDeleteAllByParticipantID", mock.Anything, participantID).Return(errors.ErrInternal("db error"))

	req := &participantdto.SaveBeneficiariesRequest{
		TenantID:      tenantID,
		ProductID:     productID,
		ParticipantID: participantID,
		Beneficiaries: []participantdto.BeneficiaryItem{{FamilyMemberID: uuid.New()}},
	}

	result, err := uc.SaveBeneficiaries(context.Background(), req)
	require.Error(t, err)
	assert.Nil(t, result)
}
