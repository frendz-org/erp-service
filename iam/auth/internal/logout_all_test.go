package internal

import (
	"context"
	"testing"
	"time"

	"erp-service/config"
	"erp-service/iam/auth/authdto"
	"erp-service/pkg/errors"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestLogoutAll(t *testing.T) {
	userID := uuid.New()

	tests := []struct {
		name    string
		req     *authdto.LogoutAllRequest
		setup   func(*MockRefreshTokenRepository, *MockUserSessionRepository, *MockInMemoryStore, *MockTransactionManager)
		wantErr bool
	}{
		{
			name: "success - revokes all tokens and sessions and blacklists user",
			req: &authdto.LogoutAllRequest{
				UserID: userID,
			},
			setup: func(mockRefreshTokenRepo *MockRefreshTokenRepository, mockSessionRepo *MockUserSessionRepository, mockBlacklist *MockInMemoryStore, mockTxMgr *MockTransactionManager) {
				mockTxMgr.On("WithTransaction", mock.Anything, mock.Anything).Return(nil)
				mockRefreshTokenRepo.On("RevokeAllByUserID", mock.Anything, userID, "User logout all").Return(nil)
				mockSessionRepo.On("RevokeAllByUserID", mock.Anything, userID).Return(nil)

				mockBlacklist.On("BlacklistUser", mock.Anything, userID, mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "success - no sessions to revoke (idempotent)",
			req: &authdto.LogoutAllRequest{
				UserID: userID,
			},
			setup: func(mockRefreshTokenRepo *MockRefreshTokenRepository, mockSessionRepo *MockUserSessionRepository, mockBlacklist *MockInMemoryStore, mockTxMgr *MockTransactionManager) {
				mockTxMgr.On("WithTransaction", mock.Anything, mock.Anything).Return(nil)
				mockRefreshTokenRepo.On("RevokeAllByUserID", mock.Anything, userID, "User logout all").Return(nil)
				mockSessionRepo.On("RevokeAllByUserID", mock.Anything, userID).Return(nil)

				mockBlacklist.On("BlacklistUser", mock.Anything, userID, mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "error - database failure",
			req: &authdto.LogoutAllRequest{
				UserID: userID,
			},
			setup: func(mockRefreshTokenRepo *MockRefreshTokenRepository, mockSessionRepo *MockUserSessionRepository, mockBlacklist *MockInMemoryStore, mockTxMgr *MockTransactionManager) {
				mockTxMgr.On("WithTransaction", mock.Anything, mock.Anything).Return(nil)
				mockRefreshTokenRepo.On("RevokeAllByUserID", mock.Anything, userID, "User logout all").Return(errors.ErrInternal("database error"))
			},
			wantErr: true,
		},
		{
			name: "success - redis failure on user blacklist (fire-and-forget)",
			req: &authdto.LogoutAllRequest{
				UserID: userID,
			},
			setup: func(mockRefreshTokenRepo *MockRefreshTokenRepository, mockSessionRepo *MockUserSessionRepository, mockBlacklist *MockInMemoryStore, mockTxMgr *MockTransactionManager) {
				mockTxMgr.On("WithTransaction", mock.Anything, mock.Anything).Return(nil)
				mockRefreshTokenRepo.On("RevokeAllByUserID", mock.Anything, userID, "User logout all").Return(nil)
				mockSessionRepo.On("RevokeAllByUserID", mock.Anything, userID).Return(nil)

				mockBlacklist.On("BlacklistUser", mock.Anything, userID, mock.Anything, mock.Anything).Return(errors.ErrInternal("redis down"))
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockTxMgr := NewMockTransactionManager()
			mockRefreshTokenRepo := new(MockRefreshTokenRepository)
			mockSessionRepo := new(MockUserSessionRepository)
			mockBlacklist := new(MockInMemoryStore)

			tt.setup(mockRefreshTokenRepo, mockSessionRepo, mockBlacklist, mockTxMgr)

			uc := &usecase{
				TxManager:        mockTxMgr,
				RefreshTokenRepo: mockRefreshTokenRepo,
				UserSessionRepo:  mockSessionRepo,
				InMemoryStore:    mockBlacklist,
				Config: &config.Config{
					JWT: config.JWTConfig{
						AccessExpiry: 15 * time.Minute,
					},
				},
			}

			err := uc.LogoutAll(context.Background(), tt.req)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRefreshTokenRepo.AssertExpectations(t)
			mockSessionRepo.AssertExpectations(t)
			mockBlacklist.AssertExpectations(t)
		})
	}
}
