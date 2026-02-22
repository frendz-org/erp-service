package auth_test

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"testing"
	"time"

	"erp-service/config"
	"erp-service/entity"
	"erp-service/iam/auth"
	"erp-service/masterdata"
	pkgerrors "erp-service/pkg/errors"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCompleteProfileRegistration(t *testing.T) {
	registrationID := uuid.New()
	email := "test@example.com"
	jwtSecret := "test-secret-key-for-testing-purposes"
	passwordHash := "$2a$10$hashedpassword"

	generateValidToken := func() (string, string) {
		claims := jwt.MapClaims{
			"registration_id": registrationID.String(),
			"email":           email,
			"purpose":         auth.RegistrationCompleteTokenPurpose,
			"exp":             time.Now().Add(15 * time.Minute).Unix(),
			"iat":             time.Now().Unix(),
			"jti":             uuid.New().String(),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, _ := token.SignedString([]byte(jwtSecret))
		hash := sha256.Sum256([]byte(tokenString))
		tokenHash := hex.EncodeToString(hash[:])
		return tokenString, tokenHash
	}

	validReq := &auth.CompleteProfileRegistrationRequest{
		FullName:    "John Michael Smith",
		DateOfBirth: "1990-01-15",
		Gender:      "GENDER_001",
	}

	type setupMockFn func(
		txManager *MockTransactionManager,
		redis *MockInMemoryStore,
		userRepo *MockUserRepository,
		profileRepo *MockUserProfileRepository,
		authMethodRepo *MockUserAuthMethodRepository,
		securityStateRepo *MockUserSecurityStateRepository,
		emailSvc *MockEmailService,
		refreshTokenRepo *MockRefreshTokenRepository,
		masterdataUsecase *MockMasterdataUsecase,
		tokenHash string,
	)

	tests := []struct {
		name          string
		req           *auth.CompleteProfileRegistrationRequest
		setupToken    func() string
		setupMocks    setupMockFn
		expectedError string
		expectedCode  string
		validateResp  func(*testing.T, *auth.CompleteProfileRegistrationResponse)
	}{
		{
			name: "success - new flow with PASSWORD_SET status",
			req:  validReq,
			setupToken: func() string {
				tokenString, _ := generateValidToken()
				return tokenString
			},
			setupMocks: func(txManager *MockTransactionManager, redis *MockInMemoryStore, userRepo *MockUserRepository, profileRepo *MockUserProfileRepository, authMethodRepo *MockUserAuthMethodRepository, securityStateRepo *MockUserSecurityStateRepository, emailSvc *MockEmailService, refreshTokenRepo *MockRefreshTokenRepository, mdValidator *MockMasterdataUsecase, tokenHash string) {
				session := &entity.RegistrationSession{
					ID:                    registrationID,
					Email:                 email,
					Status:                entity.RegistrationSessionStatusPasswordSet,
					RegistrationTokenHash: &tokenHash,
					ExpiresAt:             time.Now().Add(10 * time.Minute),
				}
				redis.On("GetRegistrationSession", mock.Anything, registrationID).Return(session, nil)
				mdValidator.On("ValidateItemCode", mock.Anything, mock.Anything).Return(&masterdata.ValidateCodeResponse{Valid: true}, nil)
				redis.On("GetRegistrationPasswordHash", mock.Anything, registrationID).Return(passwordHash, nil)
				userRepo.On("EmailExists", mock.Anything, email).Return(false, nil)
				userRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.User")).Return(nil)
				authMethodRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.UserAuthMethod")).Return(nil)
				profileRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.UserProfile")).Return(nil)
				securityStateRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.UserSecurityState")).Return(nil)
				refreshTokenRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.RefreshToken")).Return(nil)
				redis.On("DeleteRegistrationSession", mock.Anything, registrationID).Return(nil)
				redis.On("UnlockRegistrationEmail", mock.Anything, email).Return(nil)
				emailSvc.On("SendWelcome", mock.Anything, email, "John Michael").Return(nil)
			},
			validateResp: func(t *testing.T, resp *auth.CompleteProfileRegistrationResponse) {
				assert.Equal(t, email, resp.Email)
				assert.NotEmpty(t, resp.AccessToken)
				assert.NotEmpty(t, resp.RefreshToken)
				assert.Equal(t, "Bearer", resp.TokenType)
				assert.Greater(t, resp.ExpiresIn, 0)
				assert.Equal(t, "John Michael", resp.Profile.FirstName)
				assert.Equal(t, "Smith", resp.Profile.LastName)
			},
		},
		{
			name: "error - invalid token",
			req:  validReq,
			setupToken: func() string {
				return "invalid-token"
			},
			setupMocks: func(txManager *MockTransactionManager, redis *MockInMemoryStore, userRepo *MockUserRepository, profileRepo *MockUserProfileRepository, authMethodRepo *MockUserAuthMethodRepository, securityStateRepo *MockUserSecurityStateRepository, emailSvc *MockEmailService, refreshTokenRepo *MockRefreshTokenRepository, mdValidator *MockMasterdataUsecase, tokenHash string) {
			},
			expectedError: "invalid",
			expectedCode:  pkgerrors.CodeUnauthorized,
		},
		{
			name: "error - expired session",
			req:  validReq,
			setupToken: func() string {
				tokenString, _ := generateValidToken()
				return tokenString
			},
			setupMocks: func(txManager *MockTransactionManager, redis *MockInMemoryStore, userRepo *MockUserRepository, profileRepo *MockUserProfileRepository, authMethodRepo *MockUserAuthMethodRepository, securityStateRepo *MockUserSecurityStateRepository, emailSvc *MockEmailService, refreshTokenRepo *MockRefreshTokenRepository, mdValidator *MockMasterdataUsecase, tokenHash string) {
				session := &entity.RegistrationSession{
					ID:                    registrationID,
					Email:                 email,
					Status:                entity.RegistrationSessionStatusPasswordSet,
					RegistrationTokenHash: &tokenHash,
					ExpiresAt:             time.Now().Add(-1 * time.Minute),
				}
				redis.On("GetRegistrationSession", mock.Anything, registrationID).Return(session, nil)
			},
			expectedError: "expired",
			expectedCode:  "REGISTRATION_EXPIRED",
		},
		{
			name: "error - wrong status PENDING_VERIFICATION",
			req:  validReq,
			setupToken: func() string {
				tokenString, _ := generateValidToken()
				return tokenString
			},
			setupMocks: func(txManager *MockTransactionManager, redis *MockInMemoryStore, userRepo *MockUserRepository, profileRepo *MockUserProfileRepository, authMethodRepo *MockUserAuthMethodRepository, securityStateRepo *MockUserSecurityStateRepository, emailSvc *MockEmailService, refreshTokenRepo *MockRefreshTokenRepository, mdValidator *MockMasterdataUsecase, tokenHash string) {
				session := &entity.RegistrationSession{
					ID:                    registrationID,
					Email:                 email,
					Status:                entity.RegistrationSessionStatusPendingVerification,
					RegistrationTokenHash: &tokenHash,
					ExpiresAt:             time.Now().Add(10 * time.Minute),
				}
				redis.On("GetRegistrationSession", mock.Anything, registrationID).Return(session, nil)
			},
			expectedError: "not ready",
			expectedCode:  pkgerrors.CodeForbidden,
		},
		{
			name: "error - missing full_name",
			req: &auth.CompleteProfileRegistrationRequest{
				FullName:    "",
				DateOfBirth: "1990-01-15",
				Gender:      "GENDER_001",
			},
			setupToken: func() string {
				tokenString, _ := generateValidToken()
				return tokenString
			},
			setupMocks: func(txManager *MockTransactionManager, redis *MockInMemoryStore, userRepo *MockUserRepository, profileRepo *MockUserProfileRepository, authMethodRepo *MockUserAuthMethodRepository, securityStateRepo *MockUserSecurityStateRepository, emailSvc *MockEmailService, refreshTokenRepo *MockRefreshTokenRepository, mdValidator *MockMasterdataUsecase, tokenHash string) {
				session := &entity.RegistrationSession{
					ID:                    registrationID,
					Email:                 email,
					Status:                entity.RegistrationSessionStatusPasswordSet,
					RegistrationTokenHash: &tokenHash,
					ExpiresAt:             time.Now().Add(10 * time.Minute),
				}
				redis.On("GetRegistrationSession", mock.Anything, registrationID).Return(session, nil)
			},
			expectedError: "full_name",
			expectedCode:  pkgerrors.CodeValidation,
		},
		{
			name: "error - age under 18",
			req: &auth.CompleteProfileRegistrationRequest{
				FullName:    "Young Person",
				DateOfBirth: time.Now().AddDate(-17, 0, 0).Format("2006-01-02"),
				Gender:      "GENDER_001",
			},
			setupToken: func() string {
				tokenString, _ := generateValidToken()
				return tokenString
			},
			setupMocks: func(txManager *MockTransactionManager, redis *MockInMemoryStore, userRepo *MockUserRepository, profileRepo *MockUserProfileRepository, authMethodRepo *MockUserAuthMethodRepository, securityStateRepo *MockUserSecurityStateRepository, emailSvc *MockEmailService, refreshTokenRepo *MockRefreshTokenRepository, mdValidator *MockMasterdataUsecase, tokenHash string) {
				session := &entity.RegistrationSession{
					ID:                    registrationID,
					Email:                 email,
					Status:                entity.RegistrationSessionStatusPasswordSet,
					RegistrationTokenHash: &tokenHash,
					ExpiresAt:             time.Now().Add(10 * time.Minute),
				}
				redis.On("GetRegistrationSession", mock.Anything, registrationID).Return(session, nil)
				mdValidator.On("ValidateItemCode", mock.Anything, mock.Anything).Return(&masterdata.ValidateCodeResponse{Valid: true}, nil)
			},
			expectedError: "18 years",
			expectedCode:  pkgerrors.CodeValidation,
		},
		{
			name: "error - Phase A fails - invalid gender format",
			req: &auth.CompleteProfileRegistrationRequest{
				FullName:    "John Doe",
				DateOfBirth: "1990-01-15",
				Gender:      "male",
			},
			setupToken: func() string {
				tokenString, _ := generateValidToken()
				return tokenString
			},
			setupMocks: func(txManager *MockTransactionManager, redis *MockInMemoryStore, userRepo *MockUserRepository, profileRepo *MockUserProfileRepository, authMethodRepo *MockUserAuthMethodRepository, securityStateRepo *MockUserSecurityStateRepository, emailSvc *MockEmailService, refreshTokenRepo *MockRefreshTokenRepository, mdValidator *MockMasterdataUsecase, tokenHash string) {
				session := &entity.RegistrationSession{
					ID:                    registrationID,
					Email:                 email,
					Status:                entity.RegistrationSessionStatusPasswordSet,
					RegistrationTokenHash: &tokenHash,
					ExpiresAt:             time.Now().Add(10 * time.Minute),
				}
				redis.On("GetRegistrationSession", mock.Anything, registrationID).Return(session, nil)

			},
			expectedError: "gender must be in format GENDER_NNN",
			expectedCode:  pkgerrors.CodeValidation,
		},
		{
			name: "error - Phase A fails - empty gender",
			req: &auth.CompleteProfileRegistrationRequest{
				FullName:    "John Doe",
				DateOfBirth: "1990-01-15",
				Gender:      "",
			},
			setupToken: func() string {
				tokenString, _ := generateValidToken()
				return tokenString
			},
			setupMocks: func(txManager *MockTransactionManager, redis *MockInMemoryStore, userRepo *MockUserRepository, profileRepo *MockUserProfileRepository, authMethodRepo *MockUserAuthMethodRepository, securityStateRepo *MockUserSecurityStateRepository, emailSvc *MockEmailService, refreshTokenRepo *MockRefreshTokenRepository, mdValidator *MockMasterdataUsecase, tokenHash string) {
				session := &entity.RegistrationSession{
					ID:                    registrationID,
					Email:                 email,
					Status:                entity.RegistrationSessionStatusPasswordSet,
					RegistrationTokenHash: &tokenHash,
					ExpiresAt:             time.Now().Add(10 * time.Minute),
				}
				redis.On("GetRegistrationSession", mock.Anything, registrationID).Return(session, nil)
			},
			expectedError: "gender must be in format GENDER_NNN",
			expectedCode:  pkgerrors.CodeValidation,
		},
		{
			name: "error - Phase B fails - gender not in masterdata",
			req: &auth.CompleteProfileRegistrationRequest{
				FullName:    "John Doe",
				DateOfBirth: "1990-01-15",
				Gender:      "GENDER_999",
			},
			setupToken: func() string {
				tokenString, _ := generateValidToken()
				return tokenString
			},
			setupMocks: func(txManager *MockTransactionManager, redis *MockInMemoryStore, userRepo *MockUserRepository, profileRepo *MockUserProfileRepository, authMethodRepo *MockUserAuthMethodRepository, securityStateRepo *MockUserSecurityStateRepository, emailSvc *MockEmailService, refreshTokenRepo *MockRefreshTokenRepository, mdValidator *MockMasterdataUsecase, tokenHash string) {
				session := &entity.RegistrationSession{
					ID:                    registrationID,
					Email:                 email,
					Status:                entity.RegistrationSessionStatusPasswordSet,
					RegistrationTokenHash: &tokenHash,
					ExpiresAt:             time.Now().Add(10 * time.Minute),
				}
				redis.On("GetRegistrationSession", mock.Anything, registrationID).Return(session, nil)
				mdValidator.On("ValidateItemCode", mock.Anything, mock.Anything).Return(&masterdata.ValidateCodeResponse{Valid: false}, nil)
			},
			expectedError: "gender",
			expectedCode:  pkgerrors.CodeValidation,
		},
		{
			name: "error - Phase B fails - masterdata service error",
			req: &auth.CompleteProfileRegistrationRequest{
				FullName:    "John Doe",
				DateOfBirth: "1990-01-15",
				Gender:      "GENDER_001",
			},
			setupToken: func() string {
				tokenString, _ := generateValidToken()
				return tokenString
			},
			setupMocks: func(txManager *MockTransactionManager, redis *MockInMemoryStore, userRepo *MockUserRepository, profileRepo *MockUserProfileRepository, authMethodRepo *MockUserAuthMethodRepository, securityStateRepo *MockUserSecurityStateRepository, emailSvc *MockEmailService, refreshTokenRepo *MockRefreshTokenRepository, mdValidator *MockMasterdataUsecase, tokenHash string) {
				session := &entity.RegistrationSession{
					ID:                    registrationID,
					Email:                 email,
					Status:                entity.RegistrationSessionStatusPasswordSet,
					RegistrationTokenHash: &tokenHash,
					ExpiresAt:             time.Now().Add(10 * time.Minute),
				}
				redis.On("GetRegistrationSession", mock.Anything, registrationID).Return(session, nil)
				mdValidator.On("ValidateItemCode", mock.Anything, mock.Anything).Return((*masterdata.ValidateCodeResponse)(nil), errors.New("masterdata service unavailable"))
			},
			expectedError: "failed to validate gender",
			expectedCode:  pkgerrors.CodeInternal,
		},
		{
			name: "error - password hash not found in store",
			req:  validReq,
			setupToken: func() string {
				tokenString, _ := generateValidToken()
				return tokenString
			},
			setupMocks: func(txManager *MockTransactionManager, redis *MockInMemoryStore, userRepo *MockUserRepository, profileRepo *MockUserProfileRepository, authMethodRepo *MockUserAuthMethodRepository, securityStateRepo *MockUserSecurityStateRepository, emailSvc *MockEmailService, refreshTokenRepo *MockRefreshTokenRepository, mdValidator *MockMasterdataUsecase, tokenHash string) {
				session := &entity.RegistrationSession{
					ID:                    registrationID,
					Email:                 email,
					Status:                entity.RegistrationSessionStatusPasswordSet,
					RegistrationTokenHash: &tokenHash,
					ExpiresAt:             time.Now().Add(10 * time.Minute),
				}
				redis.On("GetRegistrationSession", mock.Anything, registrationID).Return(session, nil)
				mdValidator.On("ValidateItemCode", mock.Anything, mock.Anything).Return(&masterdata.ValidateCodeResponse{Valid: true}, nil)
				userRepo.On("EmailExists", mock.Anything, email).Return(false, nil)
				redis.On("GetRegistrationPasswordHash", mock.Anything, registrationID).Return("", errors.New("key not found"))
			},
			expectedError: "Password has not been set",
			expectedCode:  pkgerrors.CodeForbidden,
		},
		{
			name: "error - email already registered (race condition)",
			req:  validReq,
			setupToken: func() string {
				tokenString, _ := generateValidToken()
				return tokenString
			},
			setupMocks: func(txManager *MockTransactionManager, redis *MockInMemoryStore, userRepo *MockUserRepository, profileRepo *MockUserProfileRepository, authMethodRepo *MockUserAuthMethodRepository, securityStateRepo *MockUserSecurityStateRepository, emailSvc *MockEmailService, refreshTokenRepo *MockRefreshTokenRepository, mdValidator *MockMasterdataUsecase, tokenHash string) {
				session := &entity.RegistrationSession{
					ID:                    registrationID,
					Email:                 email,
					Status:                entity.RegistrationSessionStatusPasswordSet,
					RegistrationTokenHash: &tokenHash,
					ExpiresAt:             time.Now().Add(10 * time.Minute),
				}
				redis.On("GetRegistrationSession", mock.Anything, registrationID).Return(session, nil)
				mdValidator.On("ValidateItemCode", mock.Anything, mock.Anything).Return(&masterdata.ValidateCodeResponse{Valid: true}, nil)
				userRepo.On("EmailExists", mock.Anything, email).Return(true, nil)
			},
			expectedError: "already been registered",
			expectedCode:  pkgerrors.CodeConflict,
		},
		{
			name: "success - name splitting: John Smith",
			req: &auth.CompleteProfileRegistrationRequest{
				FullName:    "John Smith",
				DateOfBirth: "1990-01-15",
				Gender:      "GENDER_001",
			},
			setupToken: func() string {
				tokenString, _ := generateValidToken()
				return tokenString
			},
			setupMocks: func(txManager *MockTransactionManager, redis *MockInMemoryStore, userRepo *MockUserRepository, profileRepo *MockUserProfileRepository, authMethodRepo *MockUserAuthMethodRepository, securityStateRepo *MockUserSecurityStateRepository, emailSvc *MockEmailService, refreshTokenRepo *MockRefreshTokenRepository, mdValidator *MockMasterdataUsecase, tokenHash string) {
				session := &entity.RegistrationSession{
					ID:                    registrationID,
					Email:                 email,
					Status:                entity.RegistrationSessionStatusPasswordSet,
					RegistrationTokenHash: &tokenHash,
					ExpiresAt:             time.Now().Add(10 * time.Minute),
				}
				redis.On("GetRegistrationSession", mock.Anything, registrationID).Return(session, nil)
				mdValidator.On("ValidateItemCode", mock.Anything, mock.Anything).Return(&masterdata.ValidateCodeResponse{Valid: true}, nil)
				redis.On("GetRegistrationPasswordHash", mock.Anything, registrationID).Return(passwordHash, nil)
				userRepo.On("EmailExists", mock.Anything, email).Return(false, nil)
				userRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.User")).Return(nil)
				authMethodRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.UserAuthMethod")).Return(nil)
				profileRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.UserProfile")).Return(nil)
				securityStateRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.UserSecurityState")).Return(nil)
				refreshTokenRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.RefreshToken")).Return(nil)
				redis.On("DeleteRegistrationSession", mock.Anything, registrationID).Return(nil)
				redis.On("UnlockRegistrationEmail", mock.Anything, email).Return(nil)
				emailSvc.On("SendWelcome", mock.Anything, email, "John").Return(nil)
			},
			validateResp: func(t *testing.T, resp *auth.CompleteProfileRegistrationResponse) {
				assert.Equal(t, "John", resp.Profile.FirstName)
				assert.Equal(t, "Smith", resp.Profile.LastName)
			},
		},
		{
			name: "success - name splitting: Madonna (single name)",
			req: &auth.CompleteProfileRegistrationRequest{
				FullName:    "Madonna",
				DateOfBirth: "1990-01-15",
				Gender:      "GENDER_002",
			},
			setupToken: func() string {
				tokenString, _ := generateValidToken()
				return tokenString
			},
			setupMocks: func(txManager *MockTransactionManager, redis *MockInMemoryStore, userRepo *MockUserRepository, profileRepo *MockUserProfileRepository, authMethodRepo *MockUserAuthMethodRepository, securityStateRepo *MockUserSecurityStateRepository, emailSvc *MockEmailService, refreshTokenRepo *MockRefreshTokenRepository, mdValidator *MockMasterdataUsecase, tokenHash string) {
				session := &entity.RegistrationSession{
					ID:                    registrationID,
					Email:                 email,
					Status:                entity.RegistrationSessionStatusPasswordSet,
					RegistrationTokenHash: &tokenHash,
					ExpiresAt:             time.Now().Add(10 * time.Minute),
				}
				redis.On("GetRegistrationSession", mock.Anything, registrationID).Return(session, nil)
				mdValidator.On("ValidateItemCode", mock.Anything, mock.Anything).Return(&masterdata.ValidateCodeResponse{Valid: true}, nil)
				redis.On("GetRegistrationPasswordHash", mock.Anything, registrationID).Return(passwordHash, nil)
				userRepo.On("EmailExists", mock.Anything, email).Return(false, nil)
				userRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.User")).Return(nil)
				authMethodRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.UserAuthMethod")).Return(nil)
				profileRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.UserProfile")).Return(nil)
				securityStateRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.UserSecurityState")).Return(nil)
				refreshTokenRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.RefreshToken")).Return(nil)
				redis.On("DeleteRegistrationSession", mock.Anything, registrationID).Return(nil)
				redis.On("UnlockRegistrationEmail", mock.Anything, email).Return(nil)
				emailSvc.On("SendWelcome", mock.Anything, email, "Madonna").Return(nil)
			},
			validateResp: func(t *testing.T, resp *auth.CompleteProfileRegistrationResponse) {
				assert.Equal(t, "Madonna", resp.Profile.FirstName)
				assert.Equal(t, "", resp.Profile.LastName)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			txManager := NewMockTransactionManager()
			redis := &MockInMemoryStore{}
			userRepo := &MockUserRepository{}
			profileRepo := &MockUserProfileRepository{}
			authMethodRepo := &MockUserAuthMethodRepository{}
			securityStateRepo := &MockUserSecurityStateRepository{}
			emailSvc := &MockEmailService{}
			refreshTokenRepo := &MockRefreshTokenRepository{}
			mdValidator := &MockMasterdataUsecase{}

			tokenString := tt.setupToken()
			_, tokenHash := generateValidToken()
			if tokenString != "invalid-token" {
				hash := sha256.Sum256([]byte(tokenString))
				tokenHash = hex.EncodeToString(hash[:])
			}
			tt.setupMocks(txManager, redis, userRepo, profileRepo, authMethodRepo, securityStateRepo, emailSvc, refreshTokenRepo, mdValidator, tokenHash)

			cfg := &config.Config{
				JWT: config.JWTConfig{
					AccessSecret:  jwtSecret,
					RefreshSecret: "refresh-secret",
					SigningMethod: "HS256",
					AccessExpiry:  3600 * time.Second,
					RefreshExpiry: 86400 * time.Second,
					Issuer:        "erp-service",
					Audience:      []string{"erp-api"},
				},
			}

			uc := auth.NewUsecase(txManager, cfg, userRepo, profileRepo, authMethodRepo, securityStateRepo, nil, nil, refreshTokenRepo, nil, nil, nil, emailSvc, redis, nil, nil, nil, nil, mdValidator)

			tt.req.RegistrationID = registrationID
			tt.req.RegistrationToken = tokenString

			response, err := uc.CompleteProfileRegistration(
				context.Background(),
				tt.req,
			)

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				appErr, ok := err.(*pkgerrors.AppError)
				require.True(t, ok, "Error should be AppError")
				assert.Equal(t, tt.expectedCode, appErr.Code)
				return
			}

			require.NoError(t, err)
			assert.NotNil(t, response)

			if tt.validateResp != nil {
				tt.validateResp(t, response)
			}
		})
	}
}
