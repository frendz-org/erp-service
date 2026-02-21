package auth_test

import (
	"context"
	"testing"
	"time"

	"erp-service/config"
	"erp-service/entity"
	"erp-service/iam/auth"
	"erp-service/pkg/errors"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestLogout(t *testing.T) {
	userID := uuid.New()
	refreshTokenID := uuid.New()
	sessionID := uuid.New()
	tokenHash := "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08"
	jti := uuid.New().String()

	tests := []struct {
		name    string
		req     *auth.LogoutRequest
		setup   func(*MockRefreshTokenRepository, *MockUserSessionRepository, *MockInMemoryStore, *MockTransactionManager)
		wantErr bool
		errType string
	}{
		{
			name: "success - revokes token and session and blacklists JTI",
			req: &auth.LogoutRequest{
				RefreshToken:   "test-refresh-token",
				AccessTokenJTI: jti,
				AccessTokenExp: time.Now().Add(15 * time.Minute),
				UserID:         userID,
			},
			setup: func(mockRefreshTokenRepo *MockRefreshTokenRepository, mockSessionRepo *MockUserSessionRepository, mockBlacklist *MockInMemoryStore, mockTxMgr *MockTransactionManager) {
				refreshToken := &entity.RefreshToken{
					ID:        refreshTokenID,
					UserID:    userID,
					TokenHash: tokenHash,
					ExpiresAt: time.Now().Add(30 * time.Minute),
					RevokedAt: nil,
				}
				mockRefreshTokenRepo.On("GetByTokenHash", mock.Anything, mock.Anything).Return(refreshToken, nil)

				session := &entity.UserSession{
					ID:     sessionID,
					UserID: userID,
					Status: entity.UserSessionStatusActive,
				}
				mockSessionRepo.On("GetByRefreshTokenID", mock.Anything, refreshTokenID).Return(session, nil)

				mockTxMgr.On("WithTransaction", mock.Anything, mock.Anything).Return(nil)
				mockRefreshTokenRepo.On("Revoke", mock.Anything, refreshTokenID, "User logout").Return(nil)
				mockSessionRepo.On("Revoke", mock.Anything, sessionID).Return(nil)

				mockBlacklist.On("BlacklistToken", mock.Anything, jti, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "success - token not found (idempotent)",
			req: &auth.LogoutRequest{
				RefreshToken:   "nonexistent-token",
				AccessTokenJTI: jti,
				UserID:         userID,
			},
			setup: func(mockRefreshTokenRepo *MockRefreshTokenRepository, mockSessionRepo *MockUserSessionRepository, mockBlacklist *MockInMemoryStore, mockTxMgr *MockTransactionManager) {
				mockRefreshTokenRepo.On("GetByTokenHash", mock.Anything, mock.Anything).Return(nil, errors.ErrNotFound("token not found"))
			},
			wantErr: false,
		},
		{
			name: "success - already revoked (idempotent)",
			req: &auth.LogoutRequest{
				RefreshToken: "test-refresh-token",
				UserID:       userID,
			},
			setup: func(mockRefreshTokenRepo *MockRefreshTokenRepository, mockSessionRepo *MockUserSessionRepository, mockBlacklist *MockInMemoryStore, mockTxMgr *MockTransactionManager) {
				now := time.Now()
				refreshToken := &entity.RefreshToken{
					ID:        refreshTokenID,
					UserID:    userID,
					TokenHash: tokenHash,
					RevokedAt: &now,
				}
				mockRefreshTokenRepo.On("GetByTokenHash", mock.Anything, mock.Anything).Return(refreshToken, nil)
			},
			wantErr: false,
		},
		{
			name: "success - expired (idempotent)",
			req: &auth.LogoutRequest{
				RefreshToken: "test-refresh-token",
				UserID:       userID,
			},
			setup: func(mockRefreshTokenRepo *MockRefreshTokenRepository, mockSessionRepo *MockUserSessionRepository, mockBlacklist *MockInMemoryStore, mockTxMgr *MockTransactionManager) {
				refreshToken := &entity.RefreshToken{
					ID:        refreshTokenID,
					UserID:    userID,
					TokenHash: tokenHash,
					ExpiresAt: time.Now().Add(-1 * time.Hour),
					RevokedAt: nil,
				}
				mockRefreshTokenRepo.On("GetByTokenHash", mock.Anything, mock.Anything).Return(refreshToken, nil)
			},
			wantErr: false,
		},
		{
			name: "success - BOLA violation (different user, idempotent)",
			req: &auth.LogoutRequest{
				RefreshToken: "test-refresh-token",
				UserID:       uuid.New(),
			},
			setup: func(mockRefreshTokenRepo *MockRefreshTokenRepository, mockSessionRepo *MockUserSessionRepository, mockBlacklist *MockInMemoryStore, mockTxMgr *MockTransactionManager) {
				refreshToken := &entity.RefreshToken{
					ID:        refreshTokenID,
					UserID:    userID,
					TokenHash: tokenHash,
					ExpiresAt: time.Now().Add(30 * time.Minute),
					RevokedAt: nil,
				}
				mockRefreshTokenRepo.On("GetByTokenHash", mock.Anything, mock.Anything).Return(refreshToken, nil)
			},
			wantErr: false,
		},
		{
			name: "success - no linked session",
			req: &auth.LogoutRequest{
				RefreshToken:   "test-refresh-token",
				AccessTokenJTI: jti,
				AccessTokenExp: time.Now().Add(15 * time.Minute),
				UserID:         userID,
			},
			setup: func(mockRefreshTokenRepo *MockRefreshTokenRepository, mockSessionRepo *MockUserSessionRepository, mockBlacklist *MockInMemoryStore, mockTxMgr *MockTransactionManager) {
				refreshToken := &entity.RefreshToken{
					ID:        refreshTokenID,
					UserID:    userID,
					TokenHash: tokenHash,
					ExpiresAt: time.Now().Add(30 * time.Minute),
					RevokedAt: nil,
				}
				mockRefreshTokenRepo.On("GetByTokenHash", mock.Anything, mock.Anything).Return(refreshToken, nil)
				mockSessionRepo.On("GetByRefreshTokenID", mock.Anything, refreshTokenID).Return(nil, errors.ErrNotFound("session not found"))

				mockTxMgr.On("WithTransaction", mock.Anything, mock.Anything).Return(nil)
				mockRefreshTokenRepo.On("Revoke", mock.Anything, refreshTokenID, "User logout").Return(nil)

				mockBlacklist.On("BlacklistToken", mock.Anything, jti, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "error - empty refresh token",
			req: &auth.LogoutRequest{
				RefreshToken: "",
				UserID:       userID,
			},
			setup: func(mockRefreshTokenRepo *MockRefreshTokenRepository, mockSessionRepo *MockUserSessionRepository, mockBlacklist *MockInMemoryStore, mockTxMgr *MockTransactionManager) {

			},
			wantErr: true,
			errType: "bad_request",
		},
		{
			name: "error - database failure",
			req: &auth.LogoutRequest{
				RefreshToken: "test-refresh-token",
				UserID:       userID,
			},
			setup: func(mockRefreshTokenRepo *MockRefreshTokenRepository, mockSessionRepo *MockUserSessionRepository, mockBlacklist *MockInMemoryStore, mockTxMgr *MockTransactionManager) {
				mockRefreshTokenRepo.On("GetByTokenHash", mock.Anything, mock.Anything).Return(nil, errors.ErrInternal("database error"))
			},
			wantErr: true,
			errType: "internal",
		},
		{
			name: "success - redis failure on blacklist (fire-and-forget)",
			req: &auth.LogoutRequest{
				RefreshToken:   "test-refresh-token",
				AccessTokenJTI: jti,
				AccessTokenExp: time.Now().Add(15 * time.Minute),
				UserID:         userID,
			},
			setup: func(mockRefreshTokenRepo *MockRefreshTokenRepository, mockSessionRepo *MockUserSessionRepository, mockBlacklist *MockInMemoryStore, mockTxMgr *MockTransactionManager) {
				refreshToken := &entity.RefreshToken{
					ID:        refreshTokenID,
					UserID:    userID,
					TokenHash: tokenHash,
					ExpiresAt: time.Now().Add(30 * time.Minute),
					RevokedAt: nil,
				}
				mockRefreshTokenRepo.On("GetByTokenHash", mock.Anything, mock.Anything).Return(refreshToken, nil)
				mockSessionRepo.On("GetByRefreshTokenID", mock.Anything, refreshTokenID).Return(nil, errors.ErrNotFound("session not found"))

				mockTxMgr.On("WithTransaction", mock.Anything, mock.Anything).Return(nil)
				mockRefreshTokenRepo.On("Revoke", mock.Anything, refreshTokenID, "User logout").Return(nil)

				mockBlacklist.On("BlacklistToken", mock.Anything, jti, mock.Anything).Return(errors.ErrInternal("redis down"))
			},
			wantErr: false,
		},
		{
			name: "success - skip blacklist if token already expired",
			req: &auth.LogoutRequest{
				RefreshToken:   "test-refresh-token",
				AccessTokenJTI: jti,
				AccessTokenExp: time.Now().Add(-5 * time.Minute),
				UserID:         userID,
			},
			setup: func(mockRefreshTokenRepo *MockRefreshTokenRepository, mockSessionRepo *MockUserSessionRepository, mockBlacklist *MockInMemoryStore, mockTxMgr *MockTransactionManager) {
				refreshToken := &entity.RefreshToken{
					ID:        refreshTokenID,
					UserID:    userID,
					TokenHash: tokenHash,
					ExpiresAt: time.Now().Add(30 * time.Minute),
					RevokedAt: nil,
				}
				mockRefreshTokenRepo.On("GetByTokenHash", mock.Anything, mock.Anything).Return(refreshToken, nil)
				mockSessionRepo.On("GetByRefreshTokenID", mock.Anything, refreshTokenID).Return(nil, errors.ErrNotFound("session not found"))

				mockTxMgr.On("WithTransaction", mock.Anything, mock.Anything).Return(nil)
				mockRefreshTokenRepo.On("Revoke", mock.Anything, refreshTokenID, "User logout").Return(nil)

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

			uc := auth.NewUsecase(mockTxMgr, &config.Config{}, nil, nil, nil, nil, nil, nil, mockRefreshTokenRepo, nil, nil, nil, nil, mockBlacklist, mockSessionRepo, nil, nil, nil, nil)

			err := uc.Logout(context.Background(), tt.req)

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
