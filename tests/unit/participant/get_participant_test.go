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

func TestUsecase_GetParticipant(t *testing.T) {
	tenantID := uuid.New()
	productID := uuid.New()
	userID := uuid.New()
	participantID := uuid.New()
	otherTenantID := uuid.New()

	tests := []struct {
		name    string
		req     *participant.GetParticipantRequest
		setup   func(*MockParticipantRepository, *MockParticipantIdentityRepository, *MockParticipantAddressRepository, *MockParticipantBankAccountRepository, *MockParticipantFamilyMemberRepository, *MockParticipantEmploymentRepository, *MockParticipantPensionRepository, *MockParticipantBeneficiaryRepository)
		wantErr bool
		errKind errors.Kind
	}{
		{
			name: "success - retrieves participant with all child entities",
			req: &participant.GetParticipantRequest{
				ParticipantID: participantID,
				TenantID:      tenantID,
				ProductID:     productID,
			},
			setup: func(partRepo *MockParticipantRepository, identRepo *MockParticipantIdentityRepository, addrRepo *MockParticipantAddressRepository, bankRepo *MockParticipantBankAccountRepository, famRepo *MockParticipantFamilyMemberRepository, empRepo *MockParticipantEmploymentRepository, penRepo *MockParticipantPensionRepository, benRepo *MockParticipantBeneficiaryRepository) {
				p := createMockParticipant(entity.ParticipantStatusDraft, tenantID, productID, userID)
				p.ID = participantID
				partRepo.On("GetByID", mock.Anything, participantID).Return(p, nil)

				identity := createMockIdentity(participantID)
				identRepo.On("ListByParticipantID", mock.Anything, participantID).Return([]*entity.ParticipantIdentity{identity}, nil)

				address := createMockAddress(participantID)
				addrRepo.On("ListByParticipantID", mock.Anything, participantID).Return([]*entity.ParticipantAddress{address}, nil)

				bankAccount := createMockBankAccount(participantID)
				bankRepo.On("ListByParticipantID", mock.Anything, participantID).Return([]*entity.ParticipantBankAccount{bankAccount}, nil)

				familyMember := createMockFamilyMember(participantID)
				famRepo.On("ListByParticipantID", mock.Anything, participantID).Return([]*entity.ParticipantFamilyMember{familyMember}, nil)

				employment := createMockEmployment(participantID)
				empRepo.On("GetByParticipantID", mock.Anything, participantID).Return(employment, nil)

				penRepo.On("GetByParticipantID", mock.Anything, participantID).Return(nil, errors.ErrNotFound("not found"))

				beneficiary := createMockBeneficiary(participantID, uuid.New())
				benRepo.On("ListByParticipantID", mock.Anything, participantID).Return([]*entity.ParticipantBeneficiary{beneficiary}, nil)
			},
			wantErr: false,
		},
		{
			name: "success - retrieves participant without employment",
			req: &participant.GetParticipantRequest{
				ParticipantID: participantID,
				TenantID:      tenantID,
				ProductID:     productID,
			},
			setup: func(partRepo *MockParticipantRepository, identRepo *MockParticipantIdentityRepository, addrRepo *MockParticipantAddressRepository, bankRepo *MockParticipantBankAccountRepository, famRepo *MockParticipantFamilyMemberRepository, empRepo *MockParticipantEmploymentRepository, penRepo *MockParticipantPensionRepository, benRepo *MockParticipantBeneficiaryRepository) {
				p := createMockParticipant(entity.ParticipantStatusDraft, tenantID, productID, userID)
				p.ID = participantID
				partRepo.On("GetByID", mock.Anything, participantID).Return(p, nil)

				identRepo.On("ListByParticipantID", mock.Anything, participantID).Return([]*entity.ParticipantIdentity{}, nil)
				addrRepo.On("ListByParticipantID", mock.Anything, participantID).Return([]*entity.ParticipantAddress{}, nil)
				bankRepo.On("ListByParticipantID", mock.Anything, participantID).Return([]*entity.ParticipantBankAccount{}, nil)
				famRepo.On("ListByParticipantID", mock.Anything, participantID).Return([]*entity.ParticipantFamilyMember{}, nil)
				empRepo.On("GetByParticipantID", mock.Anything, participantID).Return(nil, errors.ErrNotFound("not found"))
				penRepo.On("GetByParticipantID", mock.Anything, participantID).Return(nil, errors.ErrNotFound("not found"))
				benRepo.On("ListByParticipantID", mock.Anything, participantID).Return([]*entity.ParticipantBeneficiary{}, nil)
			},
			wantErr: false,
		},
		{
			name: "error - participant not found",
			req: &participant.GetParticipantRequest{
				ParticipantID: participantID,
				TenantID:      tenantID,
				ProductID:     productID,
			},
			setup: func(partRepo *MockParticipantRepository, identRepo *MockParticipantIdentityRepository, addrRepo *MockParticipantAddressRepository, bankRepo *MockParticipantBankAccountRepository, famRepo *MockParticipantFamilyMemberRepository, empRepo *MockParticipantEmploymentRepository, penRepo *MockParticipantPensionRepository, benRepo *MockParticipantBeneficiaryRepository) {
				partRepo.On("GetByID", mock.Anything, participantID).Return(nil, errors.ErrNotFound("participant not found"))
			},
			wantErr: true,
			errKind: errors.KindNotFound,
		},
		{
			name: "error - BOLA: wrong tenant",
			req: &participant.GetParticipantRequest{
				ParticipantID: participantID,
				TenantID:      otherTenantID,
				ProductID:     productID,
			},
			setup: func(partRepo *MockParticipantRepository, identRepo *MockParticipantIdentityRepository, addrRepo *MockParticipantAddressRepository, bankRepo *MockParticipantBankAccountRepository, famRepo *MockParticipantFamilyMemberRepository, empRepo *MockParticipantEmploymentRepository, penRepo *MockParticipantPensionRepository, benRepo *MockParticipantBeneficiaryRepository) {
				p := createMockParticipant(entity.ParticipantStatusDraft, tenantID, productID, userID)
				p.ID = participantID
				partRepo.On("GetByID", mock.Anything, participantID).Return(p, nil)
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

			tt.setup(partRepo, identRepo, addrRepo, bankRepo, famRepo, empRepo, penRepo, benRepo)

			uc := newTestUsecase(txMgr, partRepo, identRepo, addrRepo, bankRepo, famRepo, empRepo, penRepo, benRepo, histRepo, fileStorage)

			resp, err := uc.GetParticipant(context.Background(), tt.req)

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
				assert.Equal(t, participantID, resp.ID)
				assert.Equal(t, tenantID, resp.TenantID)
			}

			partRepo.AssertExpectations(t)
			if !tt.wantErr {
				identRepo.AssertExpectations(t)
				addrRepo.AssertExpectations(t)
				bankRepo.AssertExpectations(t)
				famRepo.AssertExpectations(t)
				empRepo.AssertExpectations(t)
				benRepo.AssertExpectations(t)
			}
		})
	}
}
