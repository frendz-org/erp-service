package internal

import (
	"context"
	"testing"
	"time"

	"iam-service/entity"
	"iam-service/saving/participant/participantdto"
	"iam-service/pkg/errors"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUsecase_SavePension(t *testing.T) {
	tenantID := uuid.New()
	productID := uuid.New()
	userID := uuid.New()
	participantID := uuid.New()
	pensionID := uuid.New()
	otherParticipantID := uuid.New()
	otherTenantID := uuid.New()

	tests := []struct {
		name    string
		req     *participantdto.SavePensionRequest
		setup   func(*MockTransactionManager, *MockParticipantRepository, *MockParticipantPensionRepository)
		wantErr bool
		errKind errors.Kind
	}{
		{
			name: "success - creates new pension (no existing)",
			req: &participantdto.SavePensionRequest{
				ParticipantID:   participantID,
				TenantID:        tenantID,
				ProductID:   productID,
				ParticipantNumber: strPtr("PEN-001"),
				PensionCategory:   strPtr("PARTICIPANT_CATEGORY_001"),
				PensionStatus:     strPtr("PARTICIPANT_STATUS_001"),
				EffectiveDate:     timePtr(time.Now()),
			},
			setup: func(txMgr *MockTransactionManager, partRepo *MockParticipantRepository, penRepo *MockParticipantPensionRepository) {
				txMgr.On("WithTransaction", mock.Anything, mock.Anything).Return(nil)
				participant := createMockParticipant(entity.ParticipantStatusDraft, tenantID, productID, userID)
				participant.ID = participantID
				partRepo.On("GetByID", mock.Anything, participantID).Return(participant, nil)
				penRepo.On("GetByParticipantID", mock.Anything, participantID).Return(nil, errors.ErrNotFound("not found"))
				penRepo.On("Create", mock.Anything, mock.MatchedBy(func(p *entity.ParticipantPension) bool {
					return p.ParticipantID == participantID &&
						*p.ParticipantNumber == "PEN-001" &&
						*p.PensionCategory == "PARTICIPANT_CATEGORY_001"
				})).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "success - updates existing pension (found by participant ID)",
			req: &participantdto.SavePensionRequest{
				ParticipantID:   participantID,
				TenantID:        tenantID,
				ProductID:   productID,
				ParticipantNumber: strPtr("PEN-002"),
				PensionStatus:     strPtr("PARTICIPANT_STATUS_002"),
			},
			setup: func(txMgr *MockTransactionManager, partRepo *MockParticipantRepository, penRepo *MockParticipantPensionRepository) {
				txMgr.On("WithTransaction", mock.Anything, mock.Anything).Return(nil)
				participant := createMockParticipant(entity.ParticipantStatusDraft, tenantID, productID, userID)
				participant.ID = participantID
				partRepo.On("GetByID", mock.Anything, participantID).Return(participant, nil)

				existing := createMockPension(participantID)
				penRepo.On("GetByParticipantID", mock.Anything, participantID).Return(existing, nil)
				penRepo.On("Update", mock.Anything, mock.MatchedBy(func(p *entity.ParticipantPension) bool {
					return p.ParticipantID == participantID &&
						*p.ParticipantNumber == "PEN-002"
				})).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "success - updates existing pension by ID",
			req: &participantdto.SavePensionRequest{
				ID:              &pensionID,
				ParticipantID:   participantID,
				TenantID:        tenantID,
				ProductID:   productID,
				ParticipantNumber: strPtr("PEN-003"),
			},
			setup: func(txMgr *MockTransactionManager, partRepo *MockParticipantRepository, penRepo *MockParticipantPensionRepository) {
				txMgr.On("WithTransaction", mock.Anything, mock.Anything).Return(nil)
				participant := createMockParticipant(entity.ParticipantStatusDraft, tenantID, productID, userID)
				participant.ID = participantID
				partRepo.On("GetByID", mock.Anything, participantID).Return(participant, nil)

				existing := createMockPension(participantID)
				existing.ID = pensionID
				penRepo.On("GetByID", mock.Anything, pensionID).Return(existing, nil)
				penRepo.On("Update", mock.Anything, mock.MatchedBy(func(p *entity.ParticipantPension) bool {
					return p.ID == pensionID && *p.ParticipantNumber == "PEN-003"
				})).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "success - saves to REJECTED participant",
			req: &participantdto.SavePensionRequest{
				ParticipantID:   participantID,
				TenantID:        tenantID,
				ProductID:   productID,
				ParticipantNumber: strPtr("PEN-004"),
			},
			setup: func(txMgr *MockTransactionManager, partRepo *MockParticipantRepository, penRepo *MockParticipantPensionRepository) {
				txMgr.On("WithTransaction", mock.Anything, mock.Anything).Return(nil)
				participant := createMockParticipant(entity.ParticipantStatusRejected, tenantID, productID, userID)
				participant.ID = participantID
				partRepo.On("GetByID", mock.Anything, participantID).Return(participant, nil)
				penRepo.On("GetByParticipantID", mock.Anything, participantID).Return(nil, errors.ErrNotFound("not found"))
				penRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "error - participant not found",
			req: &participantdto.SavePensionRequest{
				ParticipantID:   participantID,
				TenantID:        tenantID,
				ProductID:   productID,
				ParticipantNumber: strPtr("PEN-005"),
			},
			setup: func(txMgr *MockTransactionManager, partRepo *MockParticipantRepository, penRepo *MockParticipantPensionRepository) {
				txMgr.On("WithTransaction", mock.Anything, mock.Anything).Return(nil)
				partRepo.On("GetByID", mock.Anything, participantID).Return(nil, errors.ErrNotFound("participant not found"))
			},
			wantErr: true,
			errKind: errors.KindNotFound,
		},
		{
			name: "error - BOLA: wrong tenant",
			req: &participantdto.SavePensionRequest{
				ParticipantID:   participantID,
				TenantID:        otherTenantID,
				ProductID:   productID,
				ParticipantNumber: strPtr("PEN-006"),
			},
			setup: func(txMgr *MockTransactionManager, partRepo *MockParticipantRepository, penRepo *MockParticipantPensionRepository) {
				txMgr.On("WithTransaction", mock.Anything, mock.Anything).Return(nil)
				participant := createMockParticipant(entity.ParticipantStatusDraft, tenantID, productID, userID)
				participant.ID = participantID
				partRepo.On("GetByID", mock.Anything, participantID).Return(participant, nil)
			},
			wantErr: true,
			errKind: errors.KindForbidden,
		},
		{
			name: "error - cannot save to PENDING_APPROVAL participant",
			req: &participantdto.SavePensionRequest{
				ParticipantID:   participantID,
				TenantID:        tenantID,
				ProductID:   productID,
				ParticipantNumber: strPtr("PEN-007"),
			},
			setup: func(txMgr *MockTransactionManager, partRepo *MockParticipantRepository, penRepo *MockParticipantPensionRepository) {
				txMgr.On("WithTransaction", mock.Anything, mock.Anything).Return(nil)
				participant := createMockParticipant(entity.ParticipantStatusPendingApproval, tenantID, productID, userID)
				participant.ID = participantID
				partRepo.On("GetByID", mock.Anything, participantID).Return(participant, nil)
			},
			wantErr: true,
			errKind: errors.KindBadRequest,
		},
		{
			name: "error - cannot save to APPROVED participant",
			req: &participantdto.SavePensionRequest{
				ParticipantID:   participantID,
				TenantID:        tenantID,
				ProductID:   productID,
				ParticipantNumber: strPtr("PEN-008"),
			},
			setup: func(txMgr *MockTransactionManager, partRepo *MockParticipantRepository, penRepo *MockParticipantPensionRepository) {
				txMgr.On("WithTransaction", mock.Anything, mock.Anything).Return(nil)
				participant := createMockParticipant(entity.ParticipantStatusApproved, tenantID, productID, userID)
				participant.ID = participantID
				partRepo.On("GetByID", mock.Anything, participantID).Return(participant, nil)
			},
			wantErr: true,
			errKind: errors.KindBadRequest,
		},
		{
			name: "error - pension not found for update by ID",
			req: &participantdto.SavePensionRequest{
				ID:              &pensionID,
				ParticipantID:   participantID,
				TenantID:        tenantID,
				ProductID:   productID,
				ParticipantNumber: strPtr("PEN-009"),
			},
			setup: func(txMgr *MockTransactionManager, partRepo *MockParticipantRepository, penRepo *MockParticipantPensionRepository) {
				txMgr.On("WithTransaction", mock.Anything, mock.Anything).Return(nil)
				participant := createMockParticipant(entity.ParticipantStatusDraft, tenantID, productID, userID)
				participant.ID = participantID
				partRepo.On("GetByID", mock.Anything, participantID).Return(participant, nil)
				penRepo.On("GetByID", mock.Anything, pensionID).Return(nil, errors.ErrNotFound("pension not found"))
			},
			wantErr: true,
			errKind: errors.KindNotFound,
		},
		{
			name: "error - BOLA: pension belongs to different participant",
			req: &participantdto.SavePensionRequest{
				ID:              &pensionID,
				ParticipantID:   participantID,
				TenantID:        tenantID,
				ProductID:   productID,
				ParticipantNumber: strPtr("PEN-010"),
			},
			setup: func(txMgr *MockTransactionManager, partRepo *MockParticipantRepository, penRepo *MockParticipantPensionRepository) {
				txMgr.On("WithTransaction", mock.Anything, mock.Anything).Return(nil)
				participant := createMockParticipant(entity.ParticipantStatusDraft, tenantID, productID, userID)
				participant.ID = participantID
				partRepo.On("GetByID", mock.Anything, participantID).Return(participant, nil)

				pension := createMockPension(otherParticipantID)
				pension.ID = pensionID
				penRepo.On("GetByID", mock.Anything, pensionID).Return(pension, nil)
			},
			wantErr: true,
			errKind: errors.KindForbidden,
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

			tt.setup(txMgr, partRepo, penRepo)

			uc := newTestUsecase(txMgr, partRepo, identRepo, addrRepo, bankRepo, famRepo, empRepo, penRepo, benRepo, histRepo, fileStorage)

			resp, err := uc.SavePension(context.Background(), tt.req)

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
			}

			txMgr.AssertExpectations(t)
			partRepo.AssertExpectations(t)
			penRepo.AssertExpectations(t)
		})
	}
}
