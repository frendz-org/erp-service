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

func TestUsecase_RejectParticipant(t *testing.T) {
	tenantID := uuid.New()
	productID := uuid.New()
	userID := uuid.New()
	rejecterID := uuid.New()
	otherTenantID := uuid.New()

	tests := []struct {
		name    string
		req     *participant.RejectParticipantRequest
		setup   func(*MockTransactionManager, *MockParticipantRepository, *MockParticipantStatusHistoryRepository, *MockParticipantIdentityRepository, *MockParticipantAddressRepository, *MockParticipantBankAccountRepository, *MockParticipantFamilyMemberRepository, *MockParticipantEmploymentRepository, *MockParticipantPensionRepository, *MockParticipantBeneficiaryRepository)
		wantErr bool
		errKind errors.Kind
	}{
		{
			name: "success - rejects PENDING_APPROVAL participant with reason",
			req: &participant.RejectParticipantRequest{
				ParticipantID: uuid.New(),
				TenantID:      tenantID,
				ProductID:     productID,
				UserID:        rejecterID,
				Reason:        "Incomplete documents",
			},
			setup: func(txMgr *MockTransactionManager, partRepo *MockParticipantRepository, histRepo *MockParticipantStatusHistoryRepository, identRepo *MockParticipantIdentityRepository, addrRepo *MockParticipantAddressRepository, bankRepo *MockParticipantBankAccountRepository, famRepo *MockParticipantFamilyMemberRepository, empRepo *MockParticipantEmploymentRepository, penRepo *MockParticipantPensionRepository, benRepo *MockParticipantBeneficiaryRepository) {
				txMgr.On("WithTransaction", mock.Anything, mock.Anything).Return(nil)
				p := createMockParticipant(entity.ParticipantStatusPendingApproval, tenantID, productID, userID)
				partRepo.On("GetByID", mock.Anything, mock.Anything).Return(p, nil)
				partRepo.On("Update", mock.Anything, mock.MatchedBy(func(p *entity.Participant) bool {
					return p.Status == entity.ParticipantStatusRejected &&
						p.RejectedBy != nil &&
						*p.RejectedBy == rejecterID
				})).Return(nil)
				histRepo.On("Create", mock.Anything, mock.MatchedBy(func(h *entity.ParticipantStatusHistory) bool {
					return h.ToStatus == string(entity.ParticipantStatusRejected) &&
						*h.FromStatus == string(entity.ParticipantStatusPendingApproval)
				})).Return(nil)
				identRepo.On("ListByParticipantID", mock.Anything, mock.Anything).Return([]*entity.ParticipantIdentity{}, nil)
				addrRepo.On("ListByParticipantID", mock.Anything, mock.Anything).Return([]*entity.ParticipantAddress{}, nil)
				bankRepo.On("ListByParticipantID", mock.Anything, mock.Anything).Return([]*entity.ParticipantBankAccount{}, nil)
				famRepo.On("ListByParticipantID", mock.Anything, mock.Anything).Return([]*entity.ParticipantFamilyMember{}, nil)
				empRepo.On("GetByParticipantID", mock.Anything, mock.Anything).Return(nil, errors.ErrNotFound("not found"))
				penRepo.On("GetByParticipantID", mock.Anything, mock.Anything).Return(nil, errors.ErrNotFound("not found"))
				benRepo.On("ListByParticipantID", mock.Anything, mock.Anything).Return([]*entity.ParticipantBeneficiary{}, nil)
			},
			wantErr: false,
		},
		{
			name: "error - participant not found",
			req: &participant.RejectParticipantRequest{
				ParticipantID: uuid.New(),
				TenantID:      tenantID,
				ProductID:     productID,
				UserID:        rejecterID,
			},
			setup: func(txMgr *MockTransactionManager, partRepo *MockParticipantRepository, histRepo *MockParticipantStatusHistoryRepository, identRepo *MockParticipantIdentityRepository, addrRepo *MockParticipantAddressRepository, bankRepo *MockParticipantBankAccountRepository, famRepo *MockParticipantFamilyMemberRepository, empRepo *MockParticipantEmploymentRepository, penRepo *MockParticipantPensionRepository, benRepo *MockParticipantBeneficiaryRepository) {
				txMgr.On("WithTransaction", mock.Anything, mock.Anything).Return(nil)
				partRepo.On("GetByID", mock.Anything, mock.Anything).Return(nil, errors.ErrNotFound("participant not found"))
			},
			wantErr: true,
			errKind: errors.KindNotFound,
		},
		{
			name: "error - BOLA: wrong tenant",
			req: &participant.RejectParticipantRequest{
				ParticipantID: uuid.New(),
				TenantID:      otherTenantID,
				ProductID:     productID,
				UserID:        rejecterID,
			},
			setup: func(txMgr *MockTransactionManager, partRepo *MockParticipantRepository, histRepo *MockParticipantStatusHistoryRepository, identRepo *MockParticipantIdentityRepository, addrRepo *MockParticipantAddressRepository, bankRepo *MockParticipantBankAccountRepository, famRepo *MockParticipantFamilyMemberRepository, empRepo *MockParticipantEmploymentRepository, penRepo *MockParticipantPensionRepository, benRepo *MockParticipantBeneficiaryRepository) {
				txMgr.On("WithTransaction", mock.Anything, mock.Anything).Return(nil)
				p := createMockParticipant(entity.ParticipantStatusPendingApproval, tenantID, productID, userID)
				partRepo.On("GetByID", mock.Anything, mock.Anything).Return(p, nil)
			},
			wantErr: true,
			errKind: errors.KindForbidden,
		},
		{
			name: "error - cannot reject DRAFT participant",
			req: &participant.RejectParticipantRequest{
				ParticipantID: uuid.New(),
				TenantID:      tenantID,
				ProductID:     productID,
				UserID:        rejecterID,
			},
			setup: func(txMgr *MockTransactionManager, partRepo *MockParticipantRepository, histRepo *MockParticipantStatusHistoryRepository, identRepo *MockParticipantIdentityRepository, addrRepo *MockParticipantAddressRepository, bankRepo *MockParticipantBankAccountRepository, famRepo *MockParticipantFamilyMemberRepository, empRepo *MockParticipantEmploymentRepository, penRepo *MockParticipantPensionRepository, benRepo *MockParticipantBeneficiaryRepository) {
				txMgr.On("WithTransaction", mock.Anything, mock.Anything).Return(nil)
				p := createMockParticipant(entity.ParticipantStatusDraft, tenantID, productID, userID)
				partRepo.On("GetByID", mock.Anything, mock.Anything).Return(p, nil)
			},
			wantErr: true,
			errKind: errors.KindBadRequest,
		},
		{
			name: "error - cannot reject APPROVED participant",
			req: &participant.RejectParticipantRequest{
				ParticipantID: uuid.New(),
				TenantID:      tenantID,
				ProductID:     productID,
				UserID:        rejecterID,
			},
			setup: func(txMgr *MockTransactionManager, partRepo *MockParticipantRepository, histRepo *MockParticipantStatusHistoryRepository, identRepo *MockParticipantIdentityRepository, addrRepo *MockParticipantAddressRepository, bankRepo *MockParticipantBankAccountRepository, famRepo *MockParticipantFamilyMemberRepository, empRepo *MockParticipantEmploymentRepository, penRepo *MockParticipantPensionRepository, benRepo *MockParticipantBeneficiaryRepository) {
				txMgr.On("WithTransaction", mock.Anything, mock.Anything).Return(nil)
				p := createMockParticipant(entity.ParticipantStatusApproved, tenantID, productID, userID)
				partRepo.On("GetByID", mock.Anything, mock.Anything).Return(p, nil)
			},
			wantErr: true,
			errKind: errors.KindBadRequest,
		},
		{
			name: "error - cannot reject REJECTED participant",
			req: &participant.RejectParticipantRequest{
				ParticipantID: uuid.New(),
				TenantID:      tenantID,
				ProductID:     productID,
				UserID:        rejecterID,
			},
			setup: func(txMgr *MockTransactionManager, partRepo *MockParticipantRepository, histRepo *MockParticipantStatusHistoryRepository, identRepo *MockParticipantIdentityRepository, addrRepo *MockParticipantAddressRepository, bankRepo *MockParticipantBankAccountRepository, famRepo *MockParticipantFamilyMemberRepository, empRepo *MockParticipantEmploymentRepository, penRepo *MockParticipantPensionRepository, benRepo *MockParticipantBeneficiaryRepository) {
				txMgr.On("WithTransaction", mock.Anything, mock.Anything).Return(nil)
				p := createMockParticipant(entity.ParticipantStatusRejected, tenantID, productID, userID)
				partRepo.On("GetByID", mock.Anything, mock.Anything).Return(p, nil)
			},
			wantErr: true,
			errKind: errors.KindBadRequest,
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

			tt.setup(txMgr, partRepo, histRepo, identRepo, addrRepo, bankRepo, famRepo, empRepo, penRepo, benRepo)

			uc := newTestUsecase(txMgr, partRepo, identRepo, addrRepo, bankRepo, famRepo, empRepo, penRepo, benRepo, histRepo, fileStorage)

			resp, err := uc.RejectParticipant(context.Background(), tt.req)

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
				assert.Equal(t, string(entity.ParticipantStatusRejected), resp.Status)
				assert.NotNil(t, resp.RejectedBy)
				assert.NotNil(t, resp.RejectedAt)
			}

			txMgr.AssertExpectations(t)
			partRepo.AssertExpectations(t)
		})
	}
}
