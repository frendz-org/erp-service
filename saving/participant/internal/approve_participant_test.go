package internal

import (
	"context"
	"testing"

	"iam-service/entity"
	"iam-service/saving/participant/participantdto"
	"iam-service/pkg/errors"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUsecase_ApproveParticipant(t *testing.T) {
	tenantID := uuid.New()
	productID := uuid.New()
	userID := uuid.New()
	approverID := uuid.New()
	otherTenantID := uuid.New()

	tests := []struct {
		name    string
		req     *participantdto.ApproveParticipantRequest
		setup   func(*MockTransactionManager, *MockParticipantRepository, *MockParticipantStatusHistoryRepository, *MockParticipantIdentityRepository, *MockParticipantAddressRepository, *MockParticipantBankAccountRepository, *MockParticipantFamilyMemberRepository, *MockParticipantEmploymentRepository, *MockParticipantPensionRepository, *MockParticipantBeneficiaryRepository)
		wantErr bool
		errKind errors.Kind
	}{
		{
			name: "success - approves PENDING_APPROVAL participant",
			req: &participantdto.ApproveParticipantRequest{
				ParticipantID: uuid.New(),
				TenantID:      tenantID,
				ProductID: productID,
				UserID:        approverID,
			},
			setup: func(txMgr *MockTransactionManager, partRepo *MockParticipantRepository, histRepo *MockParticipantStatusHistoryRepository, identRepo *MockParticipantIdentityRepository, addrRepo *MockParticipantAddressRepository, bankRepo *MockParticipantBankAccountRepository, famRepo *MockParticipantFamilyMemberRepository, empRepo *MockParticipantEmploymentRepository, penRepo *MockParticipantPensionRepository, benRepo *MockParticipantBeneficiaryRepository) {
				txMgr.On("WithTransaction", mock.Anything, mock.Anything).Return(nil)
				participant := createMockParticipant(entity.ParticipantStatusPendingApproval, tenantID, productID, userID)
				partRepo.On("GetByID", mock.Anything, mock.Anything).Return(participant, nil)
				partRepo.On("Update", mock.Anything, mock.MatchedBy(func(p *entity.Participant) bool {
					return p.Status == entity.ParticipantStatusApproved && p.ApprovedBy != nil && *p.ApprovedBy == approverID
				})).Return(nil)
				histRepo.On("Create", mock.Anything, mock.MatchedBy(func(h *entity.ParticipantStatusHistory) bool {
					return h.ToStatus == string(entity.ParticipantStatusApproved) && *h.FromStatus == string(entity.ParticipantStatusPendingApproval)
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
			req: &participantdto.ApproveParticipantRequest{
				ParticipantID: uuid.New(),
				TenantID:      tenantID,
				ProductID: productID,
				UserID:        approverID,
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
			req: &participantdto.ApproveParticipantRequest{
				ParticipantID: uuid.New(),
				TenantID:      otherTenantID,
				ProductID: productID,
				UserID:        approverID,
			},
			setup: func(txMgr *MockTransactionManager, partRepo *MockParticipantRepository, histRepo *MockParticipantStatusHistoryRepository, identRepo *MockParticipantIdentityRepository, addrRepo *MockParticipantAddressRepository, bankRepo *MockParticipantBankAccountRepository, famRepo *MockParticipantFamilyMemberRepository, empRepo *MockParticipantEmploymentRepository, penRepo *MockParticipantPensionRepository, benRepo *MockParticipantBeneficiaryRepository) {
				txMgr.On("WithTransaction", mock.Anything, mock.Anything).Return(nil)
				participant := createMockParticipant(entity.ParticipantStatusPendingApproval, tenantID, productID, userID)
				partRepo.On("GetByID", mock.Anything, mock.Anything).Return(participant, nil)
			},
			wantErr: true,
			errKind: errors.KindForbidden,
		},
		{
			name: "error - cannot approve DRAFT participant",
			req: &participantdto.ApproveParticipantRequest{
				ParticipantID: uuid.New(),
				TenantID:      tenantID,
				ProductID: productID,
				UserID:        approverID,
			},
			setup: func(txMgr *MockTransactionManager, partRepo *MockParticipantRepository, histRepo *MockParticipantStatusHistoryRepository, identRepo *MockParticipantIdentityRepository, addrRepo *MockParticipantAddressRepository, bankRepo *MockParticipantBankAccountRepository, famRepo *MockParticipantFamilyMemberRepository, empRepo *MockParticipantEmploymentRepository, penRepo *MockParticipantPensionRepository, benRepo *MockParticipantBeneficiaryRepository) {
				txMgr.On("WithTransaction", mock.Anything, mock.Anything).Return(nil)
				participant := createMockParticipant(entity.ParticipantStatusDraft, tenantID, productID, userID)
				partRepo.On("GetByID", mock.Anything, mock.Anything).Return(participant, nil)
			},
			wantErr: true,
			errKind: errors.KindBadRequest,
		},
		{
			name: "error - cannot approve APPROVED participant",
			req: &participantdto.ApproveParticipantRequest{
				ParticipantID: uuid.New(),
				TenantID:      tenantID,
				ProductID: productID,
				UserID:        approverID,
			},
			setup: func(txMgr *MockTransactionManager, partRepo *MockParticipantRepository, histRepo *MockParticipantStatusHistoryRepository, identRepo *MockParticipantIdentityRepository, addrRepo *MockParticipantAddressRepository, bankRepo *MockParticipantBankAccountRepository, famRepo *MockParticipantFamilyMemberRepository, empRepo *MockParticipantEmploymentRepository, penRepo *MockParticipantPensionRepository, benRepo *MockParticipantBeneficiaryRepository) {
				txMgr.On("WithTransaction", mock.Anything, mock.Anything).Return(nil)
				participant := createMockParticipant(entity.ParticipantStatusApproved, tenantID, productID, userID)
				partRepo.On("GetByID", mock.Anything, mock.Anything).Return(participant, nil)
			},
			wantErr: true,
			errKind: errors.KindBadRequest,
		},
		{
			name: "error - cannot approve REJECTED participant",
			req: &participantdto.ApproveParticipantRequest{
				ParticipantID: uuid.New(),
				TenantID:      tenantID,
				ProductID: productID,
				UserID:        approverID,
			},
			setup: func(txMgr *MockTransactionManager, partRepo *MockParticipantRepository, histRepo *MockParticipantStatusHistoryRepository, identRepo *MockParticipantIdentityRepository, addrRepo *MockParticipantAddressRepository, bankRepo *MockParticipantBankAccountRepository, famRepo *MockParticipantFamilyMemberRepository, empRepo *MockParticipantEmploymentRepository, penRepo *MockParticipantPensionRepository, benRepo *MockParticipantBeneficiaryRepository) {
				txMgr.On("WithTransaction", mock.Anything, mock.Anything).Return(nil)
				participant := createMockParticipant(entity.ParticipantStatusRejected, tenantID, productID, userID)
				partRepo.On("GetByID", mock.Anything, mock.Anything).Return(participant, nil)
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

			resp, err := uc.ApproveParticipant(context.Background(), tt.req)

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
				assert.Equal(t, string(entity.ParticipantStatusApproved), resp.Status)
				assert.NotNil(t, resp.ApprovedBy)
				assert.NotNil(t, resp.ApprovedAt)
			}

			txMgr.AssertExpectations(t)
			partRepo.AssertExpectations(t)
		})
	}
}
