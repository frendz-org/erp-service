package auth_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"erp-service/config"
	"erp-service/entity"
	"erp-service/iam/auth"
	"erp-service/pkg/logger"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestInitiateLogin_BcryptHashUpgrade(t *testing.T) {
	userID := uuid.New()
	email := "test@example.com"
	password := "TestPassword123!"

	lowCostHash, err := bcrypt.GenerateFromPassword([]byte(password), 4)
	require.NoError(t, err)

	targetCostHash, err := bcrypt.GenerateFromPassword([]byte(password), auth.BcryptTargetCost)
	require.NoError(t, err)

	activeUser := &entity.User{
		ID:     userID,
		Email:  email,
		Status: entity.UserStatusActive,
	}

	makeAuthMethod := func(hash []byte) *entity.UserAuthMethod {
		credData := entity.PasswordCredentialData{
			PasswordHash:    string(hash),
			PasswordHistory: []string{},
		}
		credJSON, _ := json.Marshal(credData)
		return &entity.UserAuthMethod{
			ID:             uuid.New(),
			UserID:         userID,
			MethodType:     string(entity.AuthMethodPassword),
			CredentialData: credJSON,
			IsActive:       true,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}
	}

	tests := []struct {
		name          string
		hash          []byte
		setup         func(*MockUserAuthMethodRepository, *MockUserRepository, *MockInMemoryStore, *MockEmailService)
		expectUpgrade bool
	}{
		{
			name: "upgrades hash when cost below target",
			hash: lowCostHash,
			setup: func(mockAuthRepo *MockUserAuthMethodRepository, mockUserRepo *MockUserRepository, mockStore *MockInMemoryStore, mockEmail *MockEmailService) {
				mockAuthRepo.On("Update", mock.Anything, mock.MatchedBy(func(am *entity.UserAuthMethod) bool {
					newHash := am.GetPasswordHash()
					newCost, err := bcrypt.Cost([]byte(newHash))
					return err == nil && newCost == auth.BcryptTargetCost
				})).Return(nil).Once()
			},
			expectUpgrade: true,
		},
		{
			name: "skips upgrade when cost already at target",
			hash: targetCostHash,
			setup: func(mockAuthRepo *MockUserAuthMethodRepository, mockUserRepo *MockUserRepository, mockStore *MockInMemoryStore, mockEmail *MockEmailService) {

			},
			expectUpgrade: false,
		},
		{
			name: "continues login when upgrade save fails",
			hash: lowCostHash,
			setup: func(mockAuthRepo *MockUserAuthMethodRepository, mockUserRepo *MockUserRepository, mockStore *MockInMemoryStore, mockEmail *MockEmailService) {
				mockAuthRepo.On("Update", mock.Anything, mock.Anything).Return(assert.AnError).Once()
			},
			expectUpgrade: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockTxMgr := NewMockTransactionManager()
			mockUserRepo := new(MockUserRepository)
			mockAuthRepo := new(MockUserAuthMethodRepository)
			mockStore := new(MockInMemoryStore)
			mockEmail := new(MockEmailService)
			auditLog := logger.NewNoopAuditLogger()

			authMethod := makeAuthMethod(tt.hash)

			mockUserRepo.On("GetByEmail", mock.Anything, email).Return(activeUser, nil)
			mockAuthRepo.On("GetByUserID", mock.Anything, userID).Return(authMethod, nil)
			mockStore.On("IncrementLoginRateLimit", mock.Anything, email, mock.Anything).Return(int64(1), nil)
			mockStore.On("CreateLoginSession", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			mockEmail.On("SendLoginOTP", mock.Anything, email, mock.Anything, mock.Anything).Return(nil)

			tt.setup(mockAuthRepo, mockUserRepo, mockStore, mockEmail)

			uc := auth.NewUsecase(
				mockTxMgr,
				&config.Config{},
				mockUserRepo,
				nil,
				mockAuthRepo,
				nil,
				nil,
				nil,
				nil,
				nil,
				nil,
				nil,
				mockEmail,
				mockStore,
				nil,
				nil,
				nil,
				auditLog,
				nil,
			)

			resp, err := uc.InitiateLogin(context.Background(), &auth.InitiateLoginRequest{
				Email:    email,
				Password: password,
			})

			require.NoError(t, err)
			require.NotNil(t, resp)

			if tt.expectUpgrade {
				newHash := authMethod.GetPasswordHash()
				newCost, costErr := bcrypt.Cost([]byte(newHash))
				require.NoError(t, costErr)
				assert.Equal(t, auth.BcryptTargetCost, newCost)
			}

			mockAuthRepo.AssertExpectations(t)
			mockUserRepo.AssertExpectations(t)
			mockStore.AssertExpectations(t)
		})
	}
}

func TestSetPasswordHash(t *testing.T) {
	password := "TestPassword123!"
	originalHash, err := bcrypt.GenerateFromPassword([]byte(password), 4)
	require.NoError(t, err)

	credData := entity.PasswordCredentialData{
		PasswordHash:    string(originalHash),
		PasswordHistory: []string{"old_hash_1"},
	}
	credJSON, err := json.Marshal(credData)
	require.NoError(t, err)

	authMethod := &entity.UserAuthMethod{
		ID:             uuid.New(),
		MethodType:     string(entity.AuthMethodPassword),
		CredentialData: credJSON,
	}

	newHash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	require.NoError(t, err)

	err = authMethod.SetPasswordHash(string(newHash))
	require.NoError(t, err)

	assert.Equal(t, string(newHash), authMethod.GetPasswordHash())

	data, err := authMethod.GetPasswordData()
	require.NoError(t, err)
	assert.Equal(t, []string{"old_hash_1"}, data.PasswordHistory)
}
