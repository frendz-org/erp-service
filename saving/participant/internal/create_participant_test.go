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

func TestUsecase_CreateParticipant(t *testing.T) {
	tests := []struct {
		name    string
		req     *participantdto.CreateParticipantRequest
		setup   func(*MockTransactionManager, *MockParticipantRepository, *MockParticipantStatusHistoryRepository, *MockParticipantIdentityRepository, *MockParticipantAddressRepository, *MockParticipantBankAccountRepository, *MockParticipantFamilyMemberRepository, *MockParticipantEmploymentRepository, *MockParticipantPensionRepository, *MockParticipantBeneficiaryRepository)
		wantErr bool
		errKind errors.Kind
	}{
		{
			name: "success - creates participant in DRAFT status",
			req: &participantdto.CreateParticipantRequest{
				TenantID:      uuid.New(),
				ProductID: uuid.New(),
				UserID:        uuid.New(),
				FullName:      "John Doe",
			},
			setup: func(txMgr *MockTransactionManager, partRepo *MockParticipantRepository, histRepo *MockParticipantStatusHistoryRepository, identRepo *MockParticipantIdentityRepository, addrRepo *MockParticipantAddressRepository, bankRepo *MockParticipantBankAccountRepository, famRepo *MockParticipantFamilyMemberRepository, empRepo *MockParticipantEmploymentRepository, penRepo *MockParticipantPensionRepository, benRepo *MockParticipantBeneficiaryRepository) {
				txMgr.On("WithTransaction", mock.Anything, mock.Anything).Return(nil)
				partRepo.On("Create", mock.Anything, mock.MatchedBy(func(p *entity.Participant) bool {
					return p.FullName == "John Doe" && p.Status == entity.ParticipantStatusDraft
				})).Return(nil)
				histRepo.On("Create", mock.Anything, mock.MatchedBy(func(h *entity.ParticipantStatusHistory) bool {
					return h.ToStatus == string(entity.ParticipantStatusDraft) && h.FromStatus == nil
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
			name: "error - repository create fails",
			req: &participantdto.CreateParticipantRequest{
				TenantID:      uuid.New(),
				ProductID: uuid.New(),
				UserID:        uuid.New(),
				FullName:      "John Doe",
			},
			setup: func(txMgr *MockTransactionManager, partRepo *MockParticipantRepository, histRepo *MockParticipantStatusHistoryRepository, identRepo *MockParticipantIdentityRepository, addrRepo *MockParticipantAddressRepository, bankRepo *MockParticipantBankAccountRepository, famRepo *MockParticipantFamilyMemberRepository, empRepo *MockParticipantEmploymentRepository, penRepo *MockParticipantPensionRepository, benRepo *MockParticipantBeneficiaryRepository) {
				txMgr.On("WithTransaction", mock.Anything, mock.Anything).Return(nil)
				partRepo.On("Create", mock.Anything, mock.Anything).Return(errors.ErrInternal("database error"))
			},
			wantErr: true,
		},
		{
			name: "error - status history creation fails",
			req: &participantdto.CreateParticipantRequest{
				TenantID:      uuid.New(),
				ProductID: uuid.New(),
				UserID:        uuid.New(),
				FullName:      "John Doe",
			},
			setup: func(txMgr *MockTransactionManager, partRepo *MockParticipantRepository, histRepo *MockParticipantStatusHistoryRepository, identRepo *MockParticipantIdentityRepository, addrRepo *MockParticipantAddressRepository, bankRepo *MockParticipantBankAccountRepository, famRepo *MockParticipantFamilyMemberRepository, empRepo *MockParticipantEmploymentRepository, penRepo *MockParticipantPensionRepository, benRepo *MockParticipantBeneficiaryRepository) {
				txMgr.On("WithTransaction", mock.Anything, mock.Anything).Return(nil)
				partRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
				histRepo.On("Create", mock.Anything, mock.Anything).Return(errors.ErrInternal("database error"))
			},
			wantErr: true,
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

			resp, err := uc.CreateParticipant(context.Background(), tt.req)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
				assert.Equal(t, tt.req.FullName, resp.FullName)
				assert.Equal(t, string(entity.ParticipantStatusDraft), resp.Status)
				assert.Equal(t, tt.req.TenantID, resp.TenantID)
				assert.Equal(t, tt.req.ProductID, resp.ProductID)
			}

			txMgr.AssertExpectations(t)
			partRepo.AssertExpectations(t)
			histRepo.AssertExpectations(t)
		})
	}
}
