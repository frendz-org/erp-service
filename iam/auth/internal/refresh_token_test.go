package internal

import (
	"context"
	"testing"
	"time"

	"iam-service/config"
	"iam-service/entity"
	"iam-service/iam/auth/authdto"
	"iam-service/pkg/errors"
	jwtpkg "iam-service/pkg/jwt"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func newTestJWTConfig() *config.JWTConfig {
	return &config.JWTConfig{
		SigningMethod: "HS256",
		AccessSecret:  "test-access-secret-must-be-32chars!",
		RefreshSecret: "test-refresh-secret-must-be-32chars",
		AccessExpiry:  15 * time.Minute,
		RefreshExpiry: 7 * 24 * time.Hour,
		Issuer:        "iam-service",
		Audience:      []string{"iam-service"},
	}
}

func generateTestRefreshToken(userID, sessionID uuid.UUID, jwtCfg *config.JWTConfig) string {
	tokenConfig := &jwtpkg.TokenConfig{
		SigningMethod: jwtCfg.SigningMethod,
		RefreshSecret: jwtCfg.RefreshSecret,
		RefreshExpiry: jwtCfg.RefreshExpiry,
		Issuer:        jwtCfg.Issuer,
		Audience:      jwtCfg.Audience,
	}
	token, _ := jwtpkg.GenerateRefreshToken(userID, sessionID, tokenConfig)
	return token
}

func TestRefreshToken(t *testing.T) {
	userID := uuid.New()
	sessionID := uuid.New()
	refreshTokenID := uuid.New()
	sessionRecordID := uuid.New()
	tokenFamily := uuid.New()

	jwtCfg := newTestJWTConfig()
	validRefreshToken := generateTestRefreshToken(userID, sessionID, jwtCfg)
	tokenHash := hashToken(validRefreshToken)

	tests := []struct {
		name    string
		req     *authdto.RefreshTokenRequest
		setup   func(*MockRefreshTokenRepository, *MockUserSessionRepository, *MockInMemoryStore, *MockTransactionManager, *MockUserRepository, *MockUserProfileRepository, *MockUserTenantRegistrationRepository, *MockProductsByTenantRepository, *MockUserRoleRepository, *MockRoleRepository, *MockPermissionRepository)
		wantErr bool
		errCode string
	}{
		{
			name: "success - rotates token and returns new pair",
			req: &authdto.RefreshTokenRequest{
				RefreshToken: validRefreshToken,
				IPAddress:    "192.168.1.1",
				UserAgent:    "TestBrowser/1.0",
			},
			setup: func(mockRefresh *MockRefreshTokenRepository, mockSession *MockUserSessionRepository, mockStore *MockInMemoryStore, mockTx *MockTransactionManager, mockUser *MockUserRepository, mockProfile *MockUserProfileRepository, mockTenantReg *MockUserTenantRegistrationRepository, mockProdByTenant *MockProductsByTenantRepository, mockUserRole *MockUserRoleRepository, mockRole *MockRoleRepository, mockPerm *MockPermissionRepository) {
				oldToken := &entity.RefreshToken{
					ID:          refreshTokenID,
					UserID:      userID,
					TokenHash:   tokenHash,
					TokenFamily: tokenFamily,
					ExpiresAt:   time.Now().Add(7 * 24 * time.Hour),
					CreatedAt:   time.Now(),
				}
				mockRefresh.On("GetByTokenHash", mock.Anything, mock.Anything).Return(oldToken, nil)

				mockStore.On("GetUserBlacklistTimestamp", mock.Anything, userID).Return(nil, nil)

				mockUser.On("GetByID", mock.Anything, userID).Return(&entity.User{
					ID:     userID,
					Email:  "test@example.com",
					Status: entity.UserStatusActive,
				}, nil)

				mockTenantReg.On("ListActiveByUserID", mock.Anything, userID).Return([]entity.UserTenantRegistration{}, nil)
				mockUserRole.On("ListActiveByUserID", mock.Anything, userID, (*uuid.UUID)(nil)).Return([]entity.UserRole{}, nil)

				session := &entity.UserSession{
					ID:     sessionRecordID,
					UserID: userID,
					Status: entity.UserSessionStatusActive,
				}
				mockSession.On("GetByRefreshTokenID", mock.Anything, refreshTokenID).Return(session, nil)

				mockRefresh.On("Create", mock.Anything, mock.Anything).Return(nil)
				mockRefresh.On("Revoke", mock.Anything, refreshTokenID, "Token rotation").Return(nil)
				mockRefresh.On("SetReplacedBy", mock.Anything, refreshTokenID, mock.Anything).Return(nil)
				mockSession.On("UpdateRefreshTokenID", mock.Anything, sessionRecordID, mock.Anything).Return(nil)

				mockProfile.On("GetByUserID", mock.Anything, userID).Return(&entity.UserProfile{
					FirstName: "John",
					LastName:  "Doe",
				}, nil)
			},
			wantErr: false,
		},
		{
			name: "error - empty refresh token",
			req: &authdto.RefreshTokenRequest{
				RefreshToken: "invalid-not-a-jwt-token",
				IPAddress:    "192.168.1.1",
			},
			setup: func(mockRefresh *MockRefreshTokenRepository, mockSession *MockUserSessionRepository, mockStore *MockInMemoryStore, mockTx *MockTransactionManager, mockUser *MockUserRepository, mockProfile *MockUserProfileRepository, mockTenantReg *MockUserTenantRegistrationRepository, mockProdByTenant *MockProductsByTenantRepository, mockUserRole *MockUserRoleRepository, mockRole *MockRoleRepository, mockPerm *MockPermissionRepository) {
			},
			wantErr: true,
			errCode: "INVALID_TOKEN",
		},
		{
			name: "error - token not found in database",
			req: &authdto.RefreshTokenRequest{
				RefreshToken: validRefreshToken,
				IPAddress:    "192.168.1.1",
			},
			setup: func(mockRefresh *MockRefreshTokenRepository, mockSession *MockUserSessionRepository, mockStore *MockInMemoryStore, mockTx *MockTransactionManager, mockUser *MockUserRepository, mockProfile *MockUserProfileRepository, mockTenantReg *MockUserTenantRegistrationRepository, mockProdByTenant *MockProductsByTenantRepository, mockUserRole *MockUserRoleRepository, mockRole *MockRoleRepository, mockPerm *MockPermissionRepository) {
				mockRefresh.On("GetByTokenHash", mock.Anything, mock.Anything).Return(nil, errors.ErrNotFound("token not found"))
			},
			wantErr: true,
			errCode: "INVALID_TOKEN",
		},
		{
			name: "error - token reuse detected (already revoked) triggers family revocation",
			req: &authdto.RefreshTokenRequest{
				RefreshToken: validRefreshToken,
				IPAddress:    "192.168.1.1",
			},
			setup: func(mockRefresh *MockRefreshTokenRepository, mockSession *MockUserSessionRepository, mockStore *MockInMemoryStore, mockTx *MockTransactionManager, mockUser *MockUserRepository, mockProfile *MockUserProfileRepository, mockTenantReg *MockUserTenantRegistrationRepository, mockProdByTenant *MockProductsByTenantRepository, mockUserRole *MockUserRoleRepository, mockRole *MockRoleRepository, mockPerm *MockPermissionRepository) {
				now := time.Now()
				oldToken := &entity.RefreshToken{
					ID:          refreshTokenID,
					UserID:      userID,
					TokenHash:   tokenHash,
					TokenFamily: tokenFamily,
					ExpiresAt:   time.Now().Add(7 * 24 * time.Hour),
					RevokedAt:   &now,
					CreatedAt:   time.Now(),
				}
				mockRefresh.On("GetByTokenHash", mock.Anything, mock.Anything).Return(oldToken, nil)
				mockRefresh.On("RevokeByFamily", mock.Anything, tokenFamily, "Token reuse detected").Return(nil)
			},
			wantErr: true,
			errCode: "TOKEN_REUSED",
		},
		{
			name: "error - token expired",
			req: &authdto.RefreshTokenRequest{
				RefreshToken: validRefreshToken,
				IPAddress:    "192.168.1.1",
			},
			setup: func(mockRefresh *MockRefreshTokenRepository, mockSession *MockUserSessionRepository, mockStore *MockInMemoryStore, mockTx *MockTransactionManager, mockUser *MockUserRepository, mockProfile *MockUserProfileRepository, mockTenantReg *MockUserTenantRegistrationRepository, mockProdByTenant *MockProductsByTenantRepository, mockUserRole *MockUserRoleRepository, mockRole *MockRoleRepository, mockPerm *MockPermissionRepository) {
				oldToken := &entity.RefreshToken{
					ID:          refreshTokenID,
					UserID:      userID,
					TokenHash:   tokenHash,
					TokenFamily: tokenFamily,
					ExpiresAt:   time.Now().Add(-1 * time.Hour),
					CreatedAt:   time.Now().Add(-8 * 24 * time.Hour),
				}
				mockRefresh.On("GetByTokenHash", mock.Anything, mock.Anything).Return(oldToken, nil)
			},
			wantErr: true,
			errCode: "INVALID_TOKEN",
		},
		{
			name: "error - user blacklisted",
			req: &authdto.RefreshTokenRequest{
				RefreshToken: validRefreshToken,
				IPAddress:    "192.168.1.1",
			},
			setup: func(mockRefresh *MockRefreshTokenRepository, mockSession *MockUserSessionRepository, mockStore *MockInMemoryStore, mockTx *MockTransactionManager, mockUser *MockUserRepository, mockProfile *MockUserProfileRepository, mockTenantReg *MockUserTenantRegistrationRepository, mockProdByTenant *MockProductsByTenantRepository, mockUserRole *MockUserRoleRepository, mockRole *MockRoleRepository, mockPerm *MockPermissionRepository) {
				createdAt := time.Now().Add(-1 * time.Hour)
				oldToken := &entity.RefreshToken{
					ID:          refreshTokenID,
					UserID:      userID,
					TokenHash:   tokenHash,
					TokenFamily: tokenFamily,
					ExpiresAt:   time.Now().Add(7 * 24 * time.Hour),
					CreatedAt:   createdAt,
				}
				mockRefresh.On("GetByTokenHash", mock.Anything, mock.Anything).Return(oldToken, nil)

				blacklistTime := time.Now()
				mockStore.On("GetUserBlacklistTimestamp", mock.Anything, userID).Return(&blacklistTime, nil)
			},
			wantErr: true,
			errCode: "USER_BLACKLISTED",
		},
		{
			name: "success - user blacklist check fails (fail-open)",
			req: &authdto.RefreshTokenRequest{
				RefreshToken: validRefreshToken,
				IPAddress:    "192.168.1.1",
			},
			setup: func(mockRefresh *MockRefreshTokenRepository, mockSession *MockUserSessionRepository, mockStore *MockInMemoryStore, mockTx *MockTransactionManager, mockUser *MockUserRepository, mockProfile *MockUserProfileRepository, mockTenantReg *MockUserTenantRegistrationRepository, mockProdByTenant *MockProductsByTenantRepository, mockUserRole *MockUserRoleRepository, mockRole *MockRoleRepository, mockPerm *MockPermissionRepository) {
				oldToken := &entity.RefreshToken{
					ID:          refreshTokenID,
					UserID:      userID,
					TokenHash:   tokenHash,
					TokenFamily: tokenFamily,
					ExpiresAt:   time.Now().Add(7 * 24 * time.Hour),
					CreatedAt:   time.Now(),
				}
				mockRefresh.On("GetByTokenHash", mock.Anything, mock.Anything).Return(oldToken, nil)

				mockStore.On("GetUserBlacklistTimestamp", mock.Anything, userID).Return(nil, errors.ErrInternal("redis down"))

				mockUser.On("GetByID", mock.Anything, userID).Return(&entity.User{
					ID:     userID,
					Email:  "test@example.com",
					Status: entity.UserStatusActive,
				}, nil)

				mockTenantReg.On("ListActiveByUserID", mock.Anything, userID).Return([]entity.UserTenantRegistration{}, nil)
				mockUserRole.On("ListActiveByUserID", mock.Anything, userID, (*uuid.UUID)(nil)).Return([]entity.UserRole{}, nil)

				session := &entity.UserSession{
					ID:     sessionRecordID,
					UserID: userID,
					Status: entity.UserSessionStatusActive,
				}
				mockSession.On("GetByRefreshTokenID", mock.Anything, refreshTokenID).Return(session, nil)

				mockRefresh.On("Create", mock.Anything, mock.Anything).Return(nil)
				mockRefresh.On("Revoke", mock.Anything, refreshTokenID, "Token rotation").Return(nil)
				mockRefresh.On("SetReplacedBy", mock.Anything, refreshTokenID, mock.Anything).Return(nil)
				mockSession.On("UpdateRefreshTokenID", mock.Anything, sessionRecordID, mock.Anything).Return(nil)

				mockProfile.On("GetByUserID", mock.Anything, userID).Return(nil, errors.ErrNotFound("profile not found"))
			},
			wantErr: false,
		},
		{
			name: "error - user not found",
			req: &authdto.RefreshTokenRequest{
				RefreshToken: validRefreshToken,
				IPAddress:    "192.168.1.1",
			},
			setup: func(mockRefresh *MockRefreshTokenRepository, mockSession *MockUserSessionRepository, mockStore *MockInMemoryStore, mockTx *MockTransactionManager, mockUser *MockUserRepository, mockProfile *MockUserProfileRepository, mockTenantReg *MockUserTenantRegistrationRepository, mockProdByTenant *MockProductsByTenantRepository, mockUserRole *MockUserRoleRepository, mockRole *MockRoleRepository, mockPerm *MockPermissionRepository) {
				oldToken := &entity.RefreshToken{
					ID:          refreshTokenID,
					UserID:      userID,
					TokenHash:   tokenHash,
					TokenFamily: tokenFamily,
					ExpiresAt:   time.Now().Add(7 * 24 * time.Hour),
					CreatedAt:   time.Now(),
				}
				mockRefresh.On("GetByTokenHash", mock.Anything, mock.Anything).Return(oldToken, nil)
				mockStore.On("GetUserBlacklistTimestamp", mock.Anything, userID).Return(nil, nil)
				mockUser.On("GetByID", mock.Anything, userID).Return(nil, errors.ErrNotFound("user not found"))
			},
			wantErr: true,
			errCode: "INVALID_TOKEN",
		},
		{
			name: "error - user suspended",
			req: &authdto.RefreshTokenRequest{
				RefreshToken: validRefreshToken,
				IPAddress:    "192.168.1.1",
			},
			setup: func(mockRefresh *MockRefreshTokenRepository, mockSession *MockUserSessionRepository, mockStore *MockInMemoryStore, mockTx *MockTransactionManager, mockUser *MockUserRepository, mockProfile *MockUserProfileRepository, mockTenantReg *MockUserTenantRegistrationRepository, mockProdByTenant *MockProductsByTenantRepository, mockUserRole *MockUserRoleRepository, mockRole *MockRoleRepository, mockPerm *MockPermissionRepository) {
				oldToken := &entity.RefreshToken{
					ID:          refreshTokenID,
					UserID:      userID,
					TokenHash:   tokenHash,
					TokenFamily: tokenFamily,
					ExpiresAt:   time.Now().Add(7 * 24 * time.Hour),
					CreatedAt:   time.Now(),
				}
				mockRefresh.On("GetByTokenHash", mock.Anything, mock.Anything).Return(oldToken, nil)
				mockStore.On("GetUserBlacklistTimestamp", mock.Anything, userID).Return(nil, nil)
				mockUser.On("GetByID", mock.Anything, userID).Return(&entity.User{
					ID:     userID,
					Email:  "test@example.com",
					Status: entity.UserStatusSuspended,
				}, nil)
			},
			wantErr: true,
			errCode: "USER_INACTIVE",
		},
		{
			name: "error - BOLA violation (token belongs to different user)",
			req: &authdto.RefreshTokenRequest{
				RefreshToken: validRefreshToken,
				IPAddress:    "192.168.1.1",
			},
			setup: func(mockRefresh *MockRefreshTokenRepository, mockSession *MockUserSessionRepository, mockStore *MockInMemoryStore, mockTx *MockTransactionManager, mockUser *MockUserRepository, mockProfile *MockUserProfileRepository, mockTenantReg *MockUserTenantRegistrationRepository, mockProdByTenant *MockProductsByTenantRepository, mockUserRole *MockUserRoleRepository, mockRole *MockRoleRepository, mockPerm *MockPermissionRepository) {
				differentUserID := uuid.New()
				oldToken := &entity.RefreshToken{
					ID:          refreshTokenID,
					UserID:      differentUserID,
					TokenHash:   tokenHash,
					TokenFamily: tokenFamily,
					ExpiresAt:   time.Now().Add(7 * 24 * time.Hour),
					CreatedAt:   time.Now(),
				}
				mockRefresh.On("GetByTokenHash", mock.Anything, mock.Anything).Return(oldToken, nil)
			},
			wantErr: true,
			errCode: "INVALID_TOKEN",
		},
		{
			name: "success - no linked session (orphaned token)",
			req: &authdto.RefreshTokenRequest{
				RefreshToken: validRefreshToken,
				IPAddress:    "192.168.1.1",
			},
			setup: func(mockRefresh *MockRefreshTokenRepository, mockSession *MockUserSessionRepository, mockStore *MockInMemoryStore, mockTx *MockTransactionManager, mockUser *MockUserRepository, mockProfile *MockUserProfileRepository, mockTenantReg *MockUserTenantRegistrationRepository, mockProdByTenant *MockProductsByTenantRepository, mockUserRole *MockUserRoleRepository, mockRole *MockRoleRepository, mockPerm *MockPermissionRepository) {
				oldToken := &entity.RefreshToken{
					ID:          refreshTokenID,
					UserID:      userID,
					TokenHash:   tokenHash,
					TokenFamily: tokenFamily,
					ExpiresAt:   time.Now().Add(7 * 24 * time.Hour),
					CreatedAt:   time.Now(),
				}
				mockRefresh.On("GetByTokenHash", mock.Anything, mock.Anything).Return(oldToken, nil)
				mockStore.On("GetUserBlacklistTimestamp", mock.Anything, userID).Return(nil, nil)
				mockUser.On("GetByID", mock.Anything, userID).Return(&entity.User{
					ID:     userID,
					Email:  "test@example.com",
					Status: entity.UserStatusActive,
				}, nil)
				mockTenantReg.On("ListActiveByUserID", mock.Anything, userID).Return([]entity.UserTenantRegistration{}, nil)
				mockUserRole.On("ListActiveByUserID", mock.Anything, userID, (*uuid.UUID)(nil)).Return([]entity.UserRole{}, nil)

				mockSession.On("GetByRefreshTokenID", mock.Anything, refreshTokenID).Return(nil, errors.ErrNotFound("no session"))

				mockRefresh.On("Create", mock.Anything, mock.Anything).Return(nil)
				mockRefresh.On("Revoke", mock.Anything, refreshTokenID, "Token rotation").Return(nil)
				mockRefresh.On("SetReplacedBy", mock.Anything, refreshTokenID, mock.Anything).Return(nil)

				mockProfile.On("GetByUserID", mock.Anything, userID).Return(&entity.UserProfile{
					FirstName: "Jane",
				}, nil)
			},
			wantErr: false,
		},
		{
			name: "error - database failure on token lookup",
			req: &authdto.RefreshTokenRequest{
				RefreshToken: validRefreshToken,
				IPAddress:    "192.168.1.1",
			},
			setup: func(mockRefresh *MockRefreshTokenRepository, mockSession *MockUserSessionRepository, mockStore *MockInMemoryStore, mockTx *MockTransactionManager, mockUser *MockUserRepository, mockProfile *MockUserProfileRepository, mockTenantReg *MockUserTenantRegistrationRepository, mockProdByTenant *MockProductsByTenantRepository, mockUserRole *MockUserRoleRepository, mockRole *MockRoleRepository, mockPerm *MockPermissionRepository) {
				mockRefresh.On("GetByTokenHash", mock.Anything, mock.Anything).Return(nil, errors.ErrInternal("database error"))
			},
			wantErr: true,
		},
		{
			name: "error - transaction failure",
			req: &authdto.RefreshTokenRequest{
				RefreshToken: validRefreshToken,
				IPAddress:    "192.168.1.1",
			},
			setup: func(mockRefresh *MockRefreshTokenRepository, mockSession *MockUserSessionRepository, mockStore *MockInMemoryStore, mockTx *MockTransactionManager, mockUser *MockUserRepository, mockProfile *MockUserProfileRepository, mockTenantReg *MockUserTenantRegistrationRepository, mockProdByTenant *MockProductsByTenantRepository, mockUserRole *MockUserRoleRepository, mockRole *MockRoleRepository, mockPerm *MockPermissionRepository) {
				oldToken := &entity.RefreshToken{
					ID:          refreshTokenID,
					UserID:      userID,
					TokenHash:   tokenHash,
					TokenFamily: tokenFamily,
					ExpiresAt:   time.Now().Add(7 * 24 * time.Hour),
					CreatedAt:   time.Now(),
				}
				mockRefresh.On("GetByTokenHash", mock.Anything, mock.Anything).Return(oldToken, nil)
				mockStore.On("GetUserBlacklistTimestamp", mock.Anything, userID).Return(nil, nil)
				mockUser.On("GetByID", mock.Anything, userID).Return(&entity.User{
					ID:     userID,
					Email:  "test@example.com",
					Status: entity.UserStatusActive,
				}, nil)
				mockTenantReg.On("ListActiveByUserID", mock.Anything, userID).Return([]entity.UserTenantRegistration{}, nil)
				mockUserRole.On("ListActiveByUserID", mock.Anything, userID, (*uuid.UUID)(nil)).Return([]entity.UserRole{}, nil)

				mockSession.On("GetByRefreshTokenID", mock.Anything, refreshTokenID).Return(nil, errors.ErrNotFound("no session"))

				mockRefresh.On("Create", mock.Anything, mock.Anything).Return(errors.ErrInternal("db write failed"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockTxMgr := NewMockTransactionManager()
			mockRefreshTokenRepo := new(MockRefreshTokenRepository)
			mockSessionRepo := new(MockUserSessionRepository)
			mockInMemory := new(MockInMemoryStore)
			mockUserRepo := new(MockUserRepository)
			mockProfileRepo := new(MockUserProfileRepository)
			mockTenantRegRepo := new(MockUserTenantRegistrationRepository)
			mockProdByTenantRepo := new(MockProductsByTenantRepository)
			mockUserRoleRepo := new(MockUserRoleRepository)
			mockRoleRepo := new(MockRoleRepository)
			mockPermRepo := new(MockPermissionRepository)

			tt.setup(mockRefreshTokenRepo, mockSessionRepo, mockInMemory, mockTxMgr, mockUserRepo, mockProfileRepo, mockTenantRegRepo, mockProdByTenantRepo, mockUserRoleRepo, mockRoleRepo, mockPermRepo)

			uc := &usecase{
				TxManager:            mockTxMgr,
				RefreshTokenRepo:     mockRefreshTokenRepo,
				UserSessionRepo:      mockSessionRepo,
				InMemoryStore:        mockInMemory,
				UserRepo:             mockUserRepo,
				UserProfileRepo:      mockProfileRepo,
				UserTenantRegRepo:    mockTenantRegRepo,
				ProductsByTenantRepo: mockProdByTenantRepo,
				UserRoleRepo:         mockUserRoleRepo,
				RoleRepo:             mockRoleRepo,
				PermissionRepo:       mockPermRepo,
				Config: &config.Config{
					JWT: *jwtCfg,
				},
			}

			resp, err := uc.RefreshToken(context.Background(), tt.req)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errCode != "" {
					appErr, ok := err.(*errors.AppError)
					if ok {
						assert.Equal(t, tt.errCode, appErr.Code)
					}
				}
				assert.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
				assert.NotEmpty(t, resp.AccessToken)
				assert.NotEmpty(t, resp.RefreshToken)
				assert.Equal(t, "Bearer", resp.TokenType)
				assert.Greater(t, resp.ExpiresIn, 0)
				assert.Equal(t, userID, resp.User.ID)
			}

			mockRefreshTokenRepo.AssertExpectations(t)
			mockSessionRepo.AssertExpectations(t)
			mockInMemory.AssertExpectations(t)
		})
	}
}
