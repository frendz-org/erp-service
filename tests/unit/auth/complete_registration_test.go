package auth_test

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"testing"
	"time"

	"erp-service/config"
	"erp-service/entity"
	"erp-service/iam/auth"
	"erp-service/pkg/errors"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCompleteRegistration(t *testing.T) {
	registrationID := uuid.New()
	email := "test@example.com"
	jwtSecret := "test-secret-key-for-testing-purposes"

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

	validReq := &auth.CompleteRegistrationRequest{
		RegistrationID:       registrationID,
		Password:             "SecureP@ssw0rd!",
		PasswordConfirmation: "SecureP@ssw0rd!",
		FirstName:            "John",
		LastName:             "Doe",
	}

	tests := []struct {
		name          string
		req           *auth.CompleteRegistrationRequest
		setupToken    func() string
		setupMocks    func(*MockTransactionManager, *MockInMemoryStore, *MockUserRepository, *MockUserProfileRepository, *MockUserAuthMethodRepository, *MockUserSecurityStateRepository, *MockEmailService, string)
		expectedError string
		expectedCode  string
	}{
		{
			name: "success - registration completed",
			req:  validReq,
			setupToken: func() string {
				tokenString, _ := generateValidToken()
				return tokenString
			},
			setupMocks: func(txManager *MockTransactionManager, redis *MockInMemoryStore, userRepo *MockUserRepository, profileRepo *MockUserProfileRepository, authMethodRepo *MockUserAuthMethodRepository, securityStateRepo *MockUserSecurityStateRepository, emailSvc *MockEmailService, tokenHash string) {
				session := &entity.RegistrationSession{
					ID:                    registrationID,
					Email:                 email,
					Status:                entity.RegistrationSessionStatusVerified,
					RegistrationTokenHash: &tokenHash,
					ExpiresAt:             time.Now().Add(10 * time.Minute),
				}
				redis.On("GetRegistrationSession", mock.Anything, registrationID).Return(session, nil)
				userRepo.On("EmailExists", mock.Anything, email).Return(false, nil)
				userRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.User")).Return(nil)
				authMethodRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.UserAuthMethod")).Return(nil)
				profileRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.UserProfile")).Return(nil)
				securityStateRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.UserSecurityState")).Return(nil)
				redis.On("DeleteRegistrationSession", mock.Anything, registrationID).Return(nil)
				redis.On("UnlockRegistrationEmail", mock.Anything, email).Return(nil)
				emailSvc.On("SendWelcome", mock.Anything, email, "John").Return(nil)
			},
		},
		{
			name: "error - invalid token format",
			req:  validReq,
			setupToken: func() string {
				return "invalid-token"
			},
			setupMocks: func(txManager *MockTransactionManager, redis *MockInMemoryStore, userRepo *MockUserRepository, profileRepo *MockUserProfileRepository, authMethodRepo *MockUserAuthMethodRepository, securityStateRepo *MockUserSecurityStateRepository, emailSvc *MockEmailService, tokenHash string) {

			},
			expectedError: "invalid",
			expectedCode:  errors.CodeUnauthorized,
		},
		{
			name: "error - expired token",
			req:  validReq,
			setupToken: func() string {
				claims := jwt.MapClaims{
					"registration_id": registrationID.String(),
					"email":           email,
					"purpose":         auth.RegistrationCompleteTokenPurpose,
					"exp":             time.Now().Add(-1 * time.Minute).Unix(),
					"iat":             time.Now().Add(-16 * time.Minute).Unix(),
					"jti":             uuid.New().String(),
				}
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				tokenString, _ := token.SignedString([]byte(jwtSecret))
				return tokenString
			},
			setupMocks: func(txManager *MockTransactionManager, redis *MockInMemoryStore, userRepo *MockUserRepository, profileRepo *MockUserProfileRepository, authMethodRepo *MockUserAuthMethodRepository, securityStateRepo *MockUserSecurityStateRepository, emailSvc *MockEmailService, tokenHash string) {

			},
			expectedError: "invalid or expired",
			expectedCode:  errors.CodeUnauthorized,
		},
		{
			name: "error - token purpose mismatch",
			req:  validReq,
			setupToken: func() string {
				claims := jwt.MapClaims{
					"registration_id": registrationID.String(),
					"email":           email,
					"purpose":         "wrong_purpose",
					"exp":             time.Now().Add(15 * time.Minute).Unix(),
					"iat":             time.Now().Unix(),
					"jti":             uuid.New().String(),
				}
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				tokenString, _ := token.SignedString([]byte(jwtSecret))
				return tokenString
			},
			setupMocks: func(txManager *MockTransactionManager, redis *MockInMemoryStore, userRepo *MockUserRepository, profileRepo *MockUserProfileRepository, authMethodRepo *MockUserAuthMethodRepository, securityStateRepo *MockUserSecurityStateRepository, emailSvc *MockEmailService, tokenHash string) {

			},
			expectedError: "not a registration completion token",
			expectedCode:  errors.CodeUnauthorized,
		},
		{
			name: "error - registration not found",
			req:  validReq,
			setupToken: func() string {
				tokenString, _ := generateValidToken()
				return tokenString
			},
			setupMocks: func(txManager *MockTransactionManager, redis *MockInMemoryStore, userRepo *MockUserRepository, profileRepo *MockUserProfileRepository, authMethodRepo *MockUserAuthMethodRepository, securityStateRepo *MockUserSecurityStateRepository, emailSvc *MockEmailService, tokenHash string) {
				redis.On("GetRegistrationSession", mock.Anything, registrationID).Return(nil, errors.ErrNotFound("registration not found"))
			},
			expectedError: "not found",
			expectedCode:  errors.CodeNotFound,
		},
		{
			name: "error - session expired",
			req:  validReq,
			setupToken: func() string {
				tokenString, _ := generateValidToken()
				return tokenString
			},
			setupMocks: func(txManager *MockTransactionManager, redis *MockInMemoryStore, userRepo *MockUserRepository, profileRepo *MockUserProfileRepository, authMethodRepo *MockUserAuthMethodRepository, securityStateRepo *MockUserSecurityStateRepository, emailSvc *MockEmailService, tokenHash string) {
				session := &entity.RegistrationSession{
					ID:        registrationID,
					Email:     email,
					Status:    entity.RegistrationSessionStatusVerified,
					ExpiresAt: time.Now().Add(-1 * time.Minute),
				}
				redis.On("GetRegistrationSession", mock.Anything, registrationID).Return(session, nil)
			},
			expectedError: "expired",
			expectedCode:  "REGISTRATION_EXPIRED",
		},
		{
			name: "error - email not verified",
			req:  validReq,
			setupToken: func() string {
				tokenString, _ := generateValidToken()
				return tokenString
			},
			setupMocks: func(txManager *MockTransactionManager, redis *MockInMemoryStore, userRepo *MockUserRepository, profileRepo *MockUserProfileRepository, authMethodRepo *MockUserAuthMethodRepository, securityStateRepo *MockUserSecurityStateRepository, emailSvc *MockEmailService, tokenHash string) {
				session := &entity.RegistrationSession{
					ID:        registrationID,
					Email:     email,
					Status:    entity.RegistrationSessionStatusPendingVerification,
					ExpiresAt: time.Now().Add(10 * time.Minute),
				}
				redis.On("GetRegistrationSession", mock.Anything, registrationID).Return(session, nil)
			},
			expectedError: "not been verified",
			expectedCode:  errors.CodeForbidden,
		},
		{
			name: "error - token already used (hash mismatch)",
			req:  validReq,
			setupToken: func() string {
				tokenString, _ := generateValidToken()
				return tokenString
			},
			setupMocks: func(txManager *MockTransactionManager, redis *MockInMemoryStore, userRepo *MockUserRepository, profileRepo *MockUserProfileRepository, authMethodRepo *MockUserAuthMethodRepository, securityStateRepo *MockUserSecurityStateRepository, emailSvc *MockEmailService, tokenHash string) {
				wrongHash := "different-hash-value"
				session := &entity.RegistrationSession{
					ID:                    registrationID,
					Email:                 email,
					Status:                entity.RegistrationSessionStatusVerified,
					RegistrationTokenHash: &wrongHash,
					ExpiresAt:             time.Now().Add(10 * time.Minute),
				}
				redis.On("GetRegistrationSession", mock.Anything, registrationID).Return(session, nil)
			},
			expectedError: "already been used",
			expectedCode:  errors.CodeUnauthorized,
		},
		{
			name: "error - weak password",
			req: &auth.CompleteRegistrationRequest{
				Password:             "weak",
				PasswordConfirmation: "weak",
				FirstName:            "John",
				LastName:             "Doe",
			},
			setupToken: func() string {
				tokenString, tokenHash := generateValidToken()

				_ = tokenHash
				return tokenString
			},
			setupMocks: func(txManager *MockTransactionManager, redis *MockInMemoryStore, userRepo *MockUserRepository, profileRepo *MockUserProfileRepository, authMethodRepo *MockUserAuthMethodRepository, securityStateRepo *MockUserSecurityStateRepository, emailSvc *MockEmailService, tokenHash string) {
				session := &entity.RegistrationSession{
					ID:                    registrationID,
					Email:                 email,
					Status:                entity.RegistrationSessionStatusVerified,
					RegistrationTokenHash: &tokenHash,
					ExpiresAt:             time.Now().Add(10 * time.Minute),
				}
				redis.On("GetRegistrationSession", mock.Anything, registrationID).Return(session, nil)
			},
			expectedError: "Password",
			expectedCode:  errors.CodeValidation,
		},
		{
			name: "error - email already registered (race condition)",
			req:  validReq,
			setupToken: func() string {
				tokenString, _ := generateValidToken()
				return tokenString
			},
			setupMocks: func(txManager *MockTransactionManager, redis *MockInMemoryStore, userRepo *MockUserRepository, profileRepo *MockUserProfileRepository, authMethodRepo *MockUserAuthMethodRepository, securityStateRepo *MockUserSecurityStateRepository, emailSvc *MockEmailService, tokenHash string) {
				session := &entity.RegistrationSession{
					ID:                    registrationID,
					Email:                 email,
					Status:                entity.RegistrationSessionStatusVerified,
					RegistrationTokenHash: &tokenHash,
					ExpiresAt:             time.Now().Add(10 * time.Minute),
				}
				redis.On("GetRegistrationSession", mock.Anything, registrationID).Return(session, nil)
				userRepo.On("EmailExists", mock.Anything, email).Return(true, nil)
			},
			expectedError: "already been registered",
			expectedCode:  errors.CodeConflict,
		},
		{
			name: "error - transaction rollback on user creation failure",
			req:  validReq,
			setupToken: func() string {
				tokenString, _ := generateValidToken()
				return tokenString
			},
			setupMocks: func(txManager *MockTransactionManager, redis *MockInMemoryStore, userRepo *MockUserRepository, profileRepo *MockUserProfileRepository, authMethodRepo *MockUserAuthMethodRepository, securityStateRepo *MockUserSecurityStateRepository, emailSvc *MockEmailService, tokenHash string) {
				session := &entity.RegistrationSession{
					ID:                    registrationID,
					Email:                 email,
					Status:                entity.RegistrationSessionStatusVerified,
					RegistrationTokenHash: &tokenHash,
					ExpiresAt:             time.Now().Add(10 * time.Minute),
				}
				redis.On("GetRegistrationSession", mock.Anything, registrationID).Return(session, nil)
				userRepo.On("EmailExists", mock.Anything, email).Return(false, nil)
				userRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.User")).Return(errors.ErrInternal("database error"))
			},
			expectedError: "failed to create user",
			expectedCode:  errors.CodeInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			txManager := NewMockTransactionManager()
			redis := new(MockInMemoryStore)
			userRepo := new(MockUserRepository)
			profileRepo := new(MockUserProfileRepository)
			authMethodRepo := new(MockUserAuthMethodRepository)
			securityStateRepo := new(MockUserSecurityStateRepository)
			emailSvc := new(MockEmailService)
			refreshTokenRepo := new(MockRefreshTokenRepository)

			token := tt.setupToken()
			hash := sha256.Sum256([]byte(token))
			tokenHash := hex.EncodeToString(hash[:])

			tt.setupMocks(txManager, redis, userRepo, profileRepo, authMethodRepo, securityStateRepo, emailSvc, tokenHash)

			refreshTokenRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.RefreshToken")).Return(nil)

			tt.req.RegistrationID = registrationID
			tt.req.RegistrationToken = token

			cfg := &config.Config{
				JWT: config.JWTConfig{
					AccessSecret:  jwtSecret,
					RefreshSecret: "refresh-secret",
					SigningMethod: "HS256",
					AccessExpiry:  3600 * time.Second,
					RefreshExpiry: 86400 * time.Second,
					Issuer:        "iam-service",
					Audience:      []string{"iam-api"},
				},
			}

			uc := auth.NewUsecase(txManager, cfg, userRepo, profileRepo, authMethodRepo, securityStateRepo, nil, nil, refreshTokenRepo, nil, nil, nil, emailSvc, redis, nil, nil, nil, nil, nil)

			ctx := context.Background()
			resp, err := uc.CompleteRegistration(ctx, tt.req)

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				if tt.expectedCode != "" {
					appErr := errors.GetAppError(err)
					require.NotNil(t, appErr, "Expected AppError but got: %v", err)
					assert.Equal(t, tt.expectedCode, appErr.Code)
				}
				assert.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
				assert.NotEqual(t, uuid.Nil, resp.UserID)
				assert.Equal(t, email, resp.Email)
			}

			redis.AssertExpectations(t)
			userRepo.AssertExpectations(t)
			profileRepo.AssertExpectations(t)
			authMethodRepo.AssertExpectations(t)
			securityStateRepo.AssertExpectations(t)
		})
	}
}
