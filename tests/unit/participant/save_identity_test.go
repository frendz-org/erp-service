package participant_test

import (
	"context"
	"testing"

	"erp-service/entity"
	"erp-service/pkg/errors"
	"erp-service/saving/participant"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUsecase_SaveIdentity(t *testing.T) {
	tenantID := uuid.New()
	productID := uuid.New()
	userID := uuid.New()
	participantID := uuid.New()
	identityID := uuid.New()
	otherParticipantID := uuid.New()
	otherTenantID := uuid.New()

	tests := []struct {
		name     string
		req      *participant.SaveIdentityRequest
		setup    func(*MockTransactionManager, *MockParticipantRepository, *MockParticipantIdentityRepository)
		wantErr  bool
		errKind  errors.Kind
		isUpdate bool
	}{
		{
			name: "success - creates new identity",
			req: &participant.SaveIdentityRequest{
				ParticipantID:  participantID,
				TenantID:       tenantID,
				ProductID:      productID,
				IdentityType:   "KTP",
				IdentityNumber: "1234567890123456",
			},
			setup: func(txMgr *MockTransactionManager, partRepo *MockParticipantRepository, identRepo *MockParticipantIdentityRepository) {
				txMgr.On("WithTransaction", mock.Anything, mock.Anything).Return(nil)
				p := createMockParticipant(entity.ParticipantStatusDraft, tenantID, productID, userID)
				p.ID = participantID
				partRepo.On("GetByID", mock.Anything, participantID).Return(p, nil)
				identRepo.On("Create", mock.Anything, mock.MatchedBy(func(i *entity.ParticipantIdentity) bool {
					return i.ParticipantID == participantID &&
						i.IdentityType == "KTP" &&
						i.IdentityNumber == "1234567890123456"
				})).Return(nil)
			},
			wantErr:  false,
			isUpdate: false,
		},
		{
			name: "success - updates existing identity",
			req: &participant.SaveIdentityRequest{
				ID:             &identityID,
				ParticipantID:  participantID,
				TenantID:       tenantID,
				ProductID:      productID,
				IdentityType:   "PASSPORT",
				IdentityNumber: "A1234567",
			},
			setup: func(txMgr *MockTransactionManager, partRepo *MockParticipantRepository, identRepo *MockParticipantIdentityRepository) {
				txMgr.On("WithTransaction", mock.Anything, mock.Anything).Return(nil)
				p := createMockParticipant(entity.ParticipantStatusDraft, tenantID, productID, userID)
				p.ID = participantID
				partRepo.On("GetByID", mock.Anything, participantID).Return(p, nil)

				identity := createMockIdentity(participantID)
				identity.ID = identityID
				identRepo.On("GetByID", mock.Anything, identityID).Return(identity, nil)
				identRepo.On("Update", mock.Anything, mock.MatchedBy(func(i *entity.ParticipantIdentity) bool {
					return i.ID == identityID &&
						i.IdentityType == "PASSPORT" &&
						i.IdentityNumber == "A1234567"
				})).Return(nil)
			},
			wantErr:  false,
			isUpdate: true,
		},
		{
			name: "success - saves to REJECTED participant",
			req: &participant.SaveIdentityRequest{
				ParticipantID:  participantID,
				TenantID:       tenantID,
				ProductID:      productID,
				IdentityType:   "KTP",
				IdentityNumber: "1234567890123456",
			},
			setup: func(txMgr *MockTransactionManager, partRepo *MockParticipantRepository, identRepo *MockParticipantIdentityRepository) {
				txMgr.On("WithTransaction", mock.Anything, mock.Anything).Return(nil)
				p := createMockParticipant(entity.ParticipantStatusRejected, tenantID, productID, userID)
				p.ID = participantID
				partRepo.On("GetByID", mock.Anything, participantID).Return(p, nil)
				identRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
			},
			wantErr:  false,
			isUpdate: false,
		},
		{
			name: "error - participant not found",
			req: &participant.SaveIdentityRequest{
				ParticipantID:  participantID,
				TenantID:       tenantID,
				ProductID:      productID,
				IdentityType:   "KTP",
				IdentityNumber: "1234567890123456",
			},
			setup: func(txMgr *MockTransactionManager, partRepo *MockParticipantRepository, identRepo *MockParticipantIdentityRepository) {
				txMgr.On("WithTransaction", mock.Anything, mock.Anything).Return(nil)
				partRepo.On("GetByID", mock.Anything, participantID).Return(nil, errors.ErrNotFound("participant not found"))
			},
			wantErr: true,
			errKind: errors.KindNotFound,
		},
		{
			name: "error - BOLA: wrong tenant",
			req: &participant.SaveIdentityRequest{
				ParticipantID:  participantID,
				TenantID:       otherTenantID,
				ProductID:      productID,
				IdentityType:   "KTP",
				IdentityNumber: "1234567890123456",
			},
			setup: func(txMgr *MockTransactionManager, partRepo *MockParticipantRepository, identRepo *MockParticipantIdentityRepository) {
				txMgr.On("WithTransaction", mock.Anything, mock.Anything).Return(nil)
				p := createMockParticipant(entity.ParticipantStatusDraft, tenantID, productID, userID)
				p.ID = participantID
				partRepo.On("GetByID", mock.Anything, participantID).Return(p, nil)
			},
			wantErr: true,
			errKind: errors.KindForbidden,
		},
		{
			name: "error - cannot save to PENDING_APPROVAL participant",
			req: &participant.SaveIdentityRequest{
				ParticipantID:  participantID,
				TenantID:       tenantID,
				ProductID:      productID,
				IdentityType:   "KTP",
				IdentityNumber: "1234567890123456",
			},
			setup: func(txMgr *MockTransactionManager, partRepo *MockParticipantRepository, identRepo *MockParticipantIdentityRepository) {
				txMgr.On("WithTransaction", mock.Anything, mock.Anything).Return(nil)
				p := createMockParticipant(entity.ParticipantStatusPendingApproval, tenantID, productID, userID)
				p.ID = participantID
				partRepo.On("GetByID", mock.Anything, participantID).Return(p, nil)
			},
			wantErr: true,
			errKind: errors.KindBadRequest,
		},
		{
			name: "error - cannot save to APPROVED participant",
			req: &participant.SaveIdentityRequest{
				ParticipantID:  participantID,
				TenantID:       tenantID,
				ProductID:      productID,
				IdentityType:   "KTP",
				IdentityNumber: "1234567890123456",
			},
			setup: func(txMgr *MockTransactionManager, partRepo *MockParticipantRepository, identRepo *MockParticipantIdentityRepository) {
				txMgr.On("WithTransaction", mock.Anything, mock.Anything).Return(nil)
				p := createMockParticipant(entity.ParticipantStatusApproved, tenantID, productID, userID)
				p.ID = participantID
				partRepo.On("GetByID", mock.Anything, participantID).Return(p, nil)
			},
			wantErr: true,
			errKind: errors.KindBadRequest,
		},
		{
			name: "error - identity not found for update",
			req: &participant.SaveIdentityRequest{
				ID:             &identityID,
				ParticipantID:  participantID,
				TenantID:       tenantID,
				ProductID:      productID,
				IdentityType:   "PASSPORT",
				IdentityNumber: "A1234567",
			},
			setup: func(txMgr *MockTransactionManager, partRepo *MockParticipantRepository, identRepo *MockParticipantIdentityRepository) {
				txMgr.On("WithTransaction", mock.Anything, mock.Anything).Return(nil)
				p := createMockParticipant(entity.ParticipantStatusDraft, tenantID, productID, userID)
				p.ID = participantID
				partRepo.On("GetByID", mock.Anything, participantID).Return(p, nil)
				identRepo.On("GetByID", mock.Anything, identityID).Return(nil, errors.ErrNotFound("identity not found"))
			},
			wantErr:  true,
			errKind:  errors.KindNotFound,
			isUpdate: true,
		},
		{
			name: "error - BOLA: identity belongs to different participant",
			req: &participant.SaveIdentityRequest{
				ID:             &identityID,
				ParticipantID:  participantID,
				TenantID:       tenantID,
				ProductID:      productID,
				IdentityType:   "PASSPORT",
				IdentityNumber: "A1234567",
			},
			setup: func(txMgr *MockTransactionManager, partRepo *MockParticipantRepository, identRepo *MockParticipantIdentityRepository) {
				txMgr.On("WithTransaction", mock.Anything, mock.Anything).Return(nil)
				p := createMockParticipant(entity.ParticipantStatusDraft, tenantID, productID, userID)
				p.ID = participantID
				partRepo.On("GetByID", mock.Anything, participantID).Return(p, nil)

				identity := createMockIdentity(otherParticipantID)
				identity.ID = identityID
				identRepo.On("GetByID", mock.Anything, identityID).Return(identity, nil)
			},
			wantErr:  true,
			errKind:  errors.KindForbidden,
			isUpdate: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			txMgr := new(MockTransactionManager)
			partRepo := new(MockParticipantRepository)
			identRepo := new(MockParticipantIdentityRepository)
			addrRepo := new(MockParticipantAddressRepository)
			bankRepo := new(MockParticipantBankAccountRepository)
			famRepo := new(MockParticipantFamilyMemberRepository)
			empRepo := new(MockParticipantEmploymentRepository)
			penRepo := new(MockParticipantPensionRepository)
			benRepo := new(MockParticipantBeneficiaryRepository)
			histRepo := new(MockParticipantStatusHistoryRepository)
			fileStorage := new(MockFileStorageAdapter)

			tt.setup(txMgr, partRepo, identRepo)

			uc := newTestUsecase(txMgr, partRepo, identRepo, addrRepo, bankRepo, famRepo, empRepo, penRepo, benRepo, histRepo, fileStorage)

			resp, err := uc.SaveIdentity(context.Background(), tt.req)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, resp)
				if tt.errKind != 0 {
					var appErr *errors.AppError
					require.True(t, errors.As(err, &appErr))
					assert.Equal(t, tt.errKind, appErr.Kind)
				}
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
				assert.Equal(t, tt.req.IdentityType, resp.IdentityType)
				assert.Equal(t, tt.req.IdentityNumber, resp.IdentityNumber)
			}

			txMgr.AssertExpectations(t)
			partRepo.AssertExpectations(t)
			identRepo.AssertExpectations(t)
		})
	}
}
