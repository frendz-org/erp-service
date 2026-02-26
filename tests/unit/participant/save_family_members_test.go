package participant_test

import (
	"context"
	"testing"

	"erp-service/config"
	"erp-service/entity"
	"erp-service/pkg/errors"
	"erp-service/saving/participant"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func makeFamilyMemberUsecase() (participant.Usecase, *MockTransactionManager, *MockParticipantRepository, *MockParticipantFamilyMemberRepository, *MockFileRepository) {
	txMgr := new(MockTransactionManager)
	participantRepo := new(MockParticipantRepository)
	familyMemberRepo := new(MockParticipantFamilyMemberRepository)
	fileRepo := new(MockFileRepository)

	uc := participant.NewUsecase(
		&config.Config{},
		zap.NewNop(),
		txMgr,
		participantRepo,
		new(MockParticipantIdentityRepository),
		new(MockParticipantAddressRepository),
		new(MockParticipantBankAccountRepository),
		familyMemberRepo,
		new(MockParticipantEmploymentRepository),
		new(MockParticipantPensionRepository),
		new(MockParticipantBeneficiaryRepository),
		new(MockParticipantStatusHistoryRepository),
		new(MockFileStorageAdapter),
		fileRepo,
		nil, nil, nil, nil, nil, nil, nil,
	)
	return uc, txMgr, participantRepo, familyMemberRepo, fileRepo
}

func TestSaveFamilyMembers_EmptySlice_ReturnsBadRequest(t *testing.T) {
	tenantID := uuid.New()
	productID := uuid.New()
	participantID := uuid.New()

	uc, _, _, _, _ := makeFamilyMemberUsecase()

	req := &participant.SaveFamilyMembersRequest{
		TenantID:      tenantID,
		ProductID:     productID,
		ParticipantID: participantID,
		FamilyMembers: []participant.FamilyMemberItem{},
	}

	result, err := uc.SaveFamilyMembers(context.Background(), req)
	require.Error(t, err)
	assert.Nil(t, result)

	var appErr *errors.AppError
	require.True(t, errors.As(err, &appErr))
	assert.Equal(t, errors.KindBadRequest, appErr.Kind)
}

func TestSaveFamilyMembers_Success(t *testing.T) {
	tenantID := uuid.New()
	productID := uuid.New()
	userID := uuid.New()

	uc, txMgr, participantRepo, familyMemberRepo, _ := makeFamilyMemberUsecase()

	p := makeDraftParticipant(tenantID, productID)
	participantID := p.ID

	txMgr.On("WithTransaction", mock.Anything, mock.AnythingOfType("func(context.Context) error")).Return(nil)
	participantRepo.On("GetByID", mock.Anything, participantID).Return(p, nil)
	participantRepo.On("Update", mock.Anything, mock.AnythingOfType("*entity.Participant")).Return(nil)
	familyMemberRepo.On("SoftDeleteAllByParticipantID", mock.Anything, participantID).Return(nil)
	familyMemberRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.ParticipantFamilyMember")).Return(nil)

	req := &participant.SaveFamilyMembersRequest{
		TenantID:      tenantID,
		ProductID:     productID,
		ParticipantID: participantID,
		UserID:        userID,
		FamilyMembers: []participant.FamilyMemberItem{
			{
				FullName:         "John Doe",
				RelationshipType: "SPOUSE",
				IsDependent:      true,
			},
		},
	}

	result, err := uc.SaveFamilyMembers(context.Background(), req)
	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "John Doe", result[0].FullName)

	familyMemberRepo.AssertNotCalled(t, "ListByParticipantID", mock.Anything, mock.Anything)
}

func TestSaveFamilyMembers_WithSupportingDoc_ValidatesOwnershipAndSetsPermanent(t *testing.T) {
	tenantID := uuid.New()
	productID := uuid.New()
	userID := uuid.New()
	fileID := uuid.New()

	uc, txMgr, participantRepo, familyMemberRepo, fileRepo := makeFamilyMemberUsecase()

	p := makeDraftParticipant(tenantID, productID)
	participantID := p.ID

	file := &entity.File{
		ID:        fileID,
		TenantID:  tenantID,
		ProductID: productID,
	}

	txMgr.On("WithTransaction", mock.Anything, mock.AnythingOfType("func(context.Context) error")).Return(nil)
	participantRepo.On("GetByID", mock.Anything, participantID).Return(p, nil)
	participantRepo.On("Update", mock.Anything, mock.AnythingOfType("*entity.Participant")).Return(nil)
	familyMemberRepo.On("SoftDeleteAllByParticipantID", mock.Anything, participantID).Return(nil)
	fileRepo.On("GetByID", mock.Anything, fileID).Return(file, nil)
	fileRepo.On("SetPermanent", mock.Anything, fileID).Return(nil)
	familyMemberRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.ParticipantFamilyMember")).Return(nil)

	req := &participant.SaveFamilyMembersRequest{
		TenantID:      tenantID,
		ProductID:     productID,
		ParticipantID: participantID,
		UserID:        userID,
		FamilyMembers: []participant.FamilyMemberItem{
			{
				FullName:            "Jane",
				RelationshipType:    "PARENT",
				SupportingDocFileID: &fileID,
			},
		},
	}

	result, err := uc.SaveFamilyMembers(context.Background(), req)
	require.NoError(t, err)
	assert.Len(t, result, 1)

	fileRepo.AssertCalled(t, "GetByID", mock.Anything, fileID)
	fileRepo.AssertCalled(t, "SetPermanent", mock.Anything, fileID)

	familyMemberRepo.AssertNotCalled(t, "ListByParticipantID", mock.Anything, mock.Anything)
}

func TestSaveFamilyMembers_FileOwnershipMismatch_ReturnsForbidden(t *testing.T) {
	tenantID := uuid.New()
	productID := uuid.New()
	otherTenantID := uuid.New()
	userID := uuid.New()
	fileID := uuid.New()

	uc, txMgr, participantRepo, familyMemberRepo, fileRepo := makeFamilyMemberUsecase()

	p := makeDraftParticipant(tenantID, productID)
	participantID := p.ID

	file := &entity.File{
		ID:       fileID,
		TenantID: otherTenantID,
	}

	txMgr.On("WithTransaction", mock.Anything, mock.AnythingOfType("func(context.Context) error")).Return(nil)
	participantRepo.On("GetByID", mock.Anything, participantID).Return(p, nil)
	familyMemberRepo.On("SoftDeleteAllByParticipantID", mock.Anything, participantID).Return(nil)
	fileRepo.On("GetByID", mock.Anything, fileID).Return(file, nil)

	req := &participant.SaveFamilyMembersRequest{
		TenantID:      tenantID,
		ProductID:     productID,
		ParticipantID: participantID,
		UserID:        userID,
		FamilyMembers: []participant.FamilyMemberItem{
			{
				FullName:            "Jane",
				RelationshipType:    "PARENT",
				SupportingDocFileID: &fileID,
			},
		},
	}

	result, err := uc.SaveFamilyMembers(context.Background(), req)
	require.Error(t, err)
	assert.Nil(t, result)

	var appErr *errors.AppError
	require.True(t, errors.As(err, &appErr))
	assert.Equal(t, errors.KindForbidden, appErr.Kind)
}

func TestSaveFamilyMembers_ParticipantNotFound(t *testing.T) {
	tenantID := uuid.New()
	productID := uuid.New()
	participantID := uuid.New()

	uc, txMgr, participantRepo, _, _ := makeFamilyMemberUsecase()

	txMgr.On("WithTransaction", mock.Anything, mock.AnythingOfType("func(context.Context) error")).Return(nil)
	participantRepo.On("GetByID", mock.Anything, participantID).Return(nil, errors.ErrNotFound("participant not found"))

	req := &participant.SaveFamilyMembersRequest{
		TenantID:      tenantID,
		ProductID:     productID,
		ParticipantID: participantID,
		FamilyMembers: []participant.FamilyMemberItem{{FullName: "John", RelationshipType: "CHILD"}},
	}

	result, err := uc.SaveFamilyMembers(context.Background(), req)
	require.Error(t, err)
	assert.Nil(t, result)
	assert.True(t, errors.IsNotFound(err))
}

func TestSaveFamilyMembers_StepMarkedComplete(t *testing.T) {
	tenantID := uuid.New()
	productID := uuid.New()
	userID := uuid.New()

	uc, txMgr, participantRepo, familyMemberRepo, _ := makeFamilyMemberUsecase()

	p := makeDraftParticipant(tenantID, productID)
	participantID := p.ID

	var savedParticipant *entity.Participant

	txMgr.On("WithTransaction", mock.Anything, mock.AnythingOfType("func(context.Context) error")).Return(nil)
	participantRepo.On("GetByID", mock.Anything, participantID).Return(p, nil)
	participantRepo.On("Update", mock.Anything, mock.AnythingOfType("*entity.Participant")).
		Run(func(args mock.Arguments) {
			savedParticipant = args.Get(1).(*entity.Participant)
		}).
		Return(nil)
	familyMemberRepo.On("SoftDeleteAllByParticipantID", mock.Anything, participantID).Return(nil)
	familyMemberRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.ParticipantFamilyMember")).Return(nil)

	req := &participant.SaveFamilyMembersRequest{
		TenantID:      tenantID,
		ProductID:     productID,
		ParticipantID: participantID,
		UserID:        userID,
		FamilyMembers: []participant.FamilyMemberItem{{FullName: "Jane", RelationshipType: "PARENT"}},
	}

	result, err := uc.SaveFamilyMembers(context.Background(), req)
	require.NoError(t, err)
	assert.Len(t, result, 1)

	require.NotNil(t, savedParticipant)
	assert.True(t, savedParticipant.StepsCompleted["family_members"], "family_members step should be marked complete")

	familyMemberRepo.AssertNotCalled(t, "ListByParticipantID", mock.Anything, mock.Anything)
}
