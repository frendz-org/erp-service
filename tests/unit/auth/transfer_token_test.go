package auth_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"erp-service/config"
	"erp-service/entity"
	"erp-service/iam/auth"
	"erp-service/pkg/errors"
	"erp-service/pkg/logger"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func newTestUsecase(
	txMgr *MockTransactionManager,
	userRepo *MockUserRepository,
	profileRepo *MockUserProfileRepository,
	refreshTokenRepo *MockRefreshTokenRepository,
	sessionRepo *MockUserSessionRepository,
	store *MockInMemoryStore,
	regRepo *MockUserTenantRegistrationRepository,
	prodByTenantRepo *MockProductsByTenantRepository,
	userRoleRepo *MockUserRoleRepository,
	roleRepo *MockRoleRepository,
	permRepo *MockPermissionRepository,
) auth.Usecase {
	return auth.NewUsecase(
		txMgr,
		&config.Config{
			JWT: config.JWTConfig{
				AccessExpiry:  15 * time.Minute,
				RefreshExpiry: 7 * 24 * time.Hour,
				AccessSecret:  "test-access-secret-key-at-least-32-bytes-long!!",
				RefreshSecret: "test-refresh-secret-key-at-least-32-bytes!!!!!",
				Issuer:        "test-issuer",
				Audience:      []string{"test-audience"},
				SigningMethod: "HS256",
			},
		},
		userRepo,
		profileRepo,
		nil,
		nil,
		nil,
		roleRepo,
		refreshTokenRepo,
		userRoleRepo,
		nil,
		permRepo,
		nil,
		store,
		sessionRepo,
		regRepo,
		prodByTenantRepo,
		logger.NewNoopAuditLogger(),
		nil,
	)
}

func TestCreateTransferToken(t *testing.T) {
	userID := uuid.New()
	sessionID := uuid.New()
	tenantID := uuid.New()
	productID := uuid.New()

	activeSession := &entity.UserSession{
		ID:     sessionID,
		UserID: userID,
		Status: entity.UserSessionStatusActive,
	}

	activeRegs := []entity.UserTenantRegistration{
		{
			ID:        uuid.New(),
			UserID:    userID,
			TenantID:  tenantID,
			ProductID: &productID,
			Status:    "ACTIVE",
		},
	}

	activeProducts := []entity.Product{
		{
			ID:       productID,
			TenantID: tenantID,
			Code:     "PRODUCT_B",
			Name:     "Product B",
			Status:   "ACTIVE",
		},
	}

	tests := []struct {
		name    string
		req     *auth.CreateTransferTokenRequest
		setup   func(*MockInMemoryStore, *MockUserSessionRepository, *MockUserTenantRegistrationRepository, *MockProductsByTenantRepository)
		wantErr bool
		errCode string
	}{
		{
			name: "success - creates transfer token",
			req: &auth.CreateTransferTokenRequest{
				ProductCode: "PRODUCT_B",
				UserID:      userID,
				SessionID:   sessionID,
			},
			setup: func(store *MockInMemoryStore, sessRepo *MockUserSessionRepository, regRepo *MockUserTenantRegistrationRepository, prodRepo *MockProductsByTenantRepository) {
				store.On("IncrementTransferTokenRateLimit", mock.Anything, userID, mock.Anything).Return(int64(1), nil)
				regRepo.On("ListActiveByUserID", mock.Anything, userID).Return(activeRegs, nil)
				prodRepo.On("ListActiveByTenantID", mock.Anything, tenantID).Return(activeProducts, nil)
				sessRepo.On("GetByID", mock.Anything, sessionID).Return(activeSession, nil)
				store.On("StoreTransferToken", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "error - rate limit exceeded",
			req: &auth.CreateTransferTokenRequest{
				ProductCode: "PRODUCT_B",
				UserID:      userID,
				SessionID:   sessionID,
			},
			setup: func(store *MockInMemoryStore, sessRepo *MockUserSessionRepository, regRepo *MockUserTenantRegistrationRepository, prodRepo *MockProductsByTenantRepository) {
				store.On("IncrementTransferTokenRateLimit", mock.Anything, userID, mock.Anything).Return(int64(6), nil)
			},
			wantErr: true,
			errCode: "RATE_LIMIT_EXCEEDED",
		},
		{
			name: "error - no product access",
			req: &auth.CreateTransferTokenRequest{
				ProductCode: "NONEXISTENT",
				UserID:      userID,
				SessionID:   sessionID,
			},
			setup: func(store *MockInMemoryStore, sessRepo *MockUserSessionRepository, regRepo *MockUserTenantRegistrationRepository, prodRepo *MockProductsByTenantRepository) {
				store.On("IncrementTransferTokenRateLimit", mock.Anything, userID, mock.Anything).Return(int64(1), nil)
				regRepo.On("ListActiveByUserID", mock.Anything, userID).Return(activeRegs, nil)
				prodRepo.On("ListActiveByTenantID", mock.Anything, tenantID).Return(activeProducts, nil)
			},
			wantErr: true,
			errCode: "NO_PRODUCT_ACCESS",
		},
		{
			name: "error - source session revoked",
			req: &auth.CreateTransferTokenRequest{
				ProductCode: "PRODUCT_B",
				UserID:      userID,
				SessionID:   sessionID,
			},
			setup: func(store *MockInMemoryStore, sessRepo *MockUserSessionRepository, regRepo *MockUserTenantRegistrationRepository, prodRepo *MockProductsByTenantRepository) {
				store.On("IncrementTransferTokenRateLimit", mock.Anything, userID, mock.Anything).Return(int64(1), nil)
				regRepo.On("ListActiveByUserID", mock.Anything, userID).Return(activeRegs, nil)
				prodRepo.On("ListActiveByTenantID", mock.Anything, tenantID).Return(activeProducts, nil)
				revokedSession := &entity.UserSession{
					ID:     sessionID,
					UserID: userID,
					Status: entity.UserSessionStatusRevoked,
				}
				sessRepo.On("GetByID", mock.Anything, sessionID).Return(revokedSession, nil)
			},
			wantErr: true,
			errCode: "SESSION_INVALID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := new(MockInMemoryStore)
			sessRepo := new(MockUserSessionRepository)
			regRepo := new(MockUserTenantRegistrationRepository)
			prodRepo := new(MockProductsByTenantRepository)

			tt.setup(store, sessRepo, regRepo, prodRepo)

			uc := newTestUsecase(
				NewMockTransactionManager(),
				new(MockUserRepository),
				new(MockUserProfileRepository),
				new(MockRefreshTokenRepository),
				sessRepo,
				store,
				regRepo,
				prodRepo,
				new(MockUserRoleRepository),
				new(MockRoleRepository),
				new(MockPermissionRepository),
			)

			resp, err := uc.CreateTransferToken(context.Background(), tt.req)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errCode)
				assert.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
				assert.Len(t, resp.Code, 64)
				assert.Equal(t, 30, resp.ExpiresIn)
			}
		})
	}
}

func TestExchangeTransferToken(t *testing.T) {
	userID := uuid.New()
	sourceSessionID := uuid.New()

	activeUser := &entity.User{
		ID:     userID,
		Email:  "test@example.com",
		Status: entity.UserStatusActive,
	}

	activeSession := &entity.UserSession{
		ID:     sourceSessionID,
		UserID: userID,
		Status: entity.UserSessionStatusActive,
	}

	transferData := struct {
		UserID          uuid.UUID `json:"user_id"`
		SourceSessionID uuid.UUID `json:"source_session_id"`
		ProductCode     string    `json:"product_code"`
		CreatedAt       time.Time `json:"created_at"`
	}{
		UserID:          userID,
		SourceSessionID: sourceSessionID,
		ProductCode:     "PRODUCT_B",
		CreatedAt:       time.Now(),
	}
	transferJSON, _ := json.Marshal(transferData)

	tests := []struct {
		name    string
		req     *auth.ExchangeTransferTokenRequest
		setup   func(*MockInMemoryStore, *MockUserSessionRepository, *MockUserRepository, *MockUserProfileRepository, *MockRefreshTokenRepository, *MockUserTenantRegistrationRepository, *MockProductsByTenantRepository, *MockUserRoleRepository, *MockRoleRepository, *MockPermissionRepository)
		wantErr bool
		errCode string
	}{
		{
			name: "success - exchanges token for new session",
			req: &auth.ExchangeTransferTokenRequest{
				Code:      "a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2",
				IPAddress: "127.0.0.1",
				UserAgent: "TestAgent",
			},
			setup: func(store *MockInMemoryStore, sessRepo *MockUserSessionRepository, userRepo *MockUserRepository, profileRepo *MockUserProfileRepository, rtRepo *MockRefreshTokenRepository, regRepo *MockUserTenantRegistrationRepository, prodRepo *MockProductsByTenantRepository, urRepo *MockUserRoleRepository, roleRepo *MockRoleRepository, permRepo *MockPermissionRepository) {
				store.On("GetAndDeleteTransferToken", mock.Anything, mock.Anything).Return(transferJSON, nil)
				sessRepo.On("GetByID", mock.Anything, sourceSessionID).Return(activeSession, nil)
				userRepo.On("GetByID", mock.Anything, userID).Return(activeUser, nil)
				store.On("GetUserBlacklistTimestamp", mock.Anything, userID).Return(nil, nil)
				regRepo.On("ListByUserIDForClaims", mock.Anything, userID).Return([]entity.UserTenantRegistration{}, nil)
				urRepo.On("ListActiveByUserID", mock.Anything, userID, mock.Anything).Return([]entity.UserRole{}, nil)
				rtRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
				sessRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
				profileRepo.On("GetByUserID", mock.Anything, userID).Return(&entity.UserProfile{
					FirstName: "Test",
					LastName:  "User",
				}, nil)
			},
			wantErr: false,
		},
		{
			name: "error - code not found / expired",
			req: &auth.ExchangeTransferTokenRequest{
				Code: "deadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef",
			},
			setup: func(store *MockInMemoryStore, sessRepo *MockUserSessionRepository, userRepo *MockUserRepository, profileRepo *MockUserProfileRepository, rtRepo *MockRefreshTokenRepository, regRepo *MockUserTenantRegistrationRepository, prodRepo *MockProductsByTenantRepository, urRepo *MockUserRoleRepository, roleRepo *MockRoleRepository, permRepo *MockPermissionRepository) {
				store.On("GetAndDeleteTransferToken", mock.Anything, mock.Anything).Return(nil, nil)
			},
			wantErr: true,
			errCode: "INVALID_CODE",
		},
		{
			name: "error - parent session revoked",
			req: &auth.ExchangeTransferTokenRequest{
				Code: "a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2",
			},
			setup: func(store *MockInMemoryStore, sessRepo *MockUserSessionRepository, userRepo *MockUserRepository, profileRepo *MockUserProfileRepository, rtRepo *MockRefreshTokenRepository, regRepo *MockUserTenantRegistrationRepository, prodRepo *MockProductsByTenantRepository, urRepo *MockUserRoleRepository, roleRepo *MockRoleRepository, permRepo *MockPermissionRepository) {
				store.On("GetAndDeleteTransferToken", mock.Anything, mock.Anything).Return(transferJSON, nil)
				revokedSession := &entity.UserSession{
					ID:     sourceSessionID,
					UserID: userID,
					Status: entity.UserSessionStatusRevoked,
				}
				sessRepo.On("GetByID", mock.Anything, sourceSessionID).Return(revokedSession, nil)
			},
			wantErr: true,
			errCode: "SESSION_REVOKED",
		},
		{
			name: "error - user inactive",
			req: &auth.ExchangeTransferTokenRequest{
				Code: "a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2",
			},
			setup: func(store *MockInMemoryStore, sessRepo *MockUserSessionRepository, userRepo *MockUserRepository, profileRepo *MockUserProfileRepository, rtRepo *MockRefreshTokenRepository, regRepo *MockUserTenantRegistrationRepository, prodRepo *MockProductsByTenantRepository, urRepo *MockUserRoleRepository, roleRepo *MockRoleRepository, permRepo *MockPermissionRepository) {
				store.On("GetAndDeleteTransferToken", mock.Anything, mock.Anything).Return(transferJSON, nil)
				sessRepo.On("GetByID", mock.Anything, sourceSessionID).Return(activeSession, nil)
				inactiveUser := &entity.User{
					ID:     userID,
					Email:  "test@example.com",
					Status: entity.UserStatusInactive,
				}
				userRepo.On("GetByID", mock.Anything, userID).Return(inactiveUser, nil)
			},
			wantErr: true,
			errCode: "USER_INACTIVE",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := new(MockInMemoryStore)
			sessRepo := new(MockUserSessionRepository)
			userRepo := new(MockUserRepository)
			profileRepo := new(MockUserProfileRepository)
			rtRepo := new(MockRefreshTokenRepository)
			regRepo := new(MockUserTenantRegistrationRepository)
			prodRepo := new(MockProductsByTenantRepository)
			urRepo := new(MockUserRoleRepository)
			roleRepo := new(MockRoleRepository)
			permRepo := new(MockPermissionRepository)

			tt.setup(store, sessRepo, userRepo, profileRepo, rtRepo, regRepo, prodRepo, urRepo, roleRepo, permRepo)

			uc := newTestUsecase(
				NewMockTransactionManager(),
				userRepo,
				profileRepo,
				rtRepo,
				sessRepo,
				store,
				regRepo,
				prodRepo,
				urRepo,
				roleRepo,
				permRepo,
			)

			resp, err := uc.ExchangeTransferToken(context.Background(), tt.req)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errCode)
				assert.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
				assert.NotEmpty(t, resp.AccessToken)
				assert.NotEmpty(t, resp.RefreshToken)
				assert.Equal(t, "Bearer", resp.TokenType)
				assert.Equal(t, userID, resp.User.ID)
				assert.Equal(t, "test@example.com", resp.User.Email)
			}
		})
	}
}

func TestLogoutTree(t *testing.T) {
	userID := uuid.New()
	rootSessionID := uuid.New()
	childSessionID := uuid.New()
	refreshTokenID := uuid.New()
	childRefreshTokenID := uuid.New()

	now := time.Now()

	tests := []struct {
		name    string
		req     *auth.LogoutTreeRequest
		setup   func(*MockRefreshTokenRepository, *MockUserSessionRepository, *MockInMemoryStore)
		wantErr bool
	}{
		{
			name: "success - revokes entire session tree",
			req: &auth.LogoutTreeRequest{
				RefreshToken:   "valid-refresh-token",
				UserID:         userID,
				AccessTokenJTI: "test-jti",
				AccessTokenExp: now.Add(15 * time.Minute),
			},
			setup: func(rtRepo *MockRefreshTokenRepository, sessRepo *MockUserSessionRepository, store *MockInMemoryStore) {
				rt := &entity.RefreshToken{
					ID:        refreshTokenID,
					UserID:    userID,
					ExpiresAt: now.Add(24 * time.Hour),
				}
				rtRepo.On("GetByTokenHash", mock.Anything, mock.Anything).Return(rt, nil)

				rootSession := &entity.UserSession{
					ID:             rootSessionID,
					UserID:         userID,
					RefreshTokenID: &refreshTokenID,
					Status:         entity.UserSessionStatusActive,
				}
				sessRepo.On("GetByRefreshTokenID", mock.Anything, refreshTokenID).Return(rootSession, nil)
				sessRepo.On("GetByID", mock.Anything, rootSessionID).Return(rootSession, nil)

				descendantIDs := []uuid.UUID{rootSessionID, childSessionID}
				sessRepo.On("GetDescendantSessionIDs", mock.Anything, rootSessionID).Return(descendantIDs, nil)

				childSession := &entity.UserSession{
					ID:              childSessionID,
					UserID:          userID,
					ParentSessionID: &rootSessionID,
					RefreshTokenID:  &childRefreshTokenID,
					Status:          entity.UserSessionStatusActive,
				}
				sessRepo.On("GetByID", mock.Anything, childSessionID).Return(childSession, nil)

				sessRepo.On("RevokeByIDs", mock.Anything, descendantIDs).Return(nil)
				rtRepo.On("RevokeByIDs", mock.Anything, []uuid.UUID{refreshTokenID, childRefreshTokenID}, "Session tree logout").Return(nil)
				store.On("BlacklistToken", mock.Anything, "test-jti", mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "idempotent - refresh token not found",
			req: &auth.LogoutTreeRequest{
				RefreshToken: "nonexistent-token",
				UserID:       userID,
			},
			setup: func(rtRepo *MockRefreshTokenRepository, sessRepo *MockUserSessionRepository, store *MockInMemoryStore) {
				rtRepo.On("GetByTokenHash", mock.Anything, mock.Anything).Return(nil, errors.ErrNotFound("refresh token not found"))
			},
			wantErr: false,
		},
		{
			name: "idempotent - different user (BOLA prevention)",
			req: &auth.LogoutTreeRequest{
				RefreshToken: "other-user-token",
				UserID:       userID,
			},
			setup: func(rtRepo *MockRefreshTokenRepository, sessRepo *MockUserSessionRepository, store *MockInMemoryStore) {
				otherUserID := uuid.New()
				rt := &entity.RefreshToken{
					ID:        uuid.New(),
					UserID:    otherUserID,
					ExpiresAt: now.Add(24 * time.Hour),
				}
				rtRepo.On("GetByTokenHash", mock.Anything, mock.Anything).Return(rt, nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rtRepo := new(MockRefreshTokenRepository)
			sessRepo := new(MockUserSessionRepository)
			store := new(MockInMemoryStore)

			tt.setup(rtRepo, sessRepo, store)

			uc := newTestUsecase(
				NewMockTransactionManager(),
				new(MockUserRepository),
				new(MockUserProfileRepository),
				rtRepo,
				sessRepo,
				store,
				new(MockUserTenantRegistrationRepository),
				new(MockProductsByTenantRepository),
				new(MockUserRoleRepository),
				new(MockRoleRepository),
				new(MockPermissionRepository),
			)

			err := uc.LogoutTree(context.Background(), tt.req)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
