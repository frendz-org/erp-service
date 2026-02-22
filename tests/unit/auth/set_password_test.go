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

func TestSetPassword(t *testing.T) {
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

	validReq := &auth.SetPasswordRequest{
		Password:             "SecureP@ssw0rd!",
		ConfirmationPassword: "SecureP@ssw0rd!",
	}

	tests := []struct {
		name          string
		req           *auth.SetPasswordRequest
		setupToken    func() string
		setupMocks    func(*MockInMemoryStore, string)
		expectedError string
		expectedCode  string
	}{
		{
			name: "success - password set for verified session",
			req:  validReq,
			setupToken: func() string {
				tokenString, _ := generateValidToken()
				return tokenString
			},
			setupMocks: func(redis *MockInMemoryStore, tokenHash string) {
				session := &entity.RegistrationSession{
					ID:                    registrationID,
					Email:                 email,
					Status:                entity.RegistrationSessionStatusVerified,
					RegistrationTokenHash: &tokenHash,
					ExpiresAt:             time.Now().Add(10 * time.Minute),
				}
				redis.On("GetRegistrationSession", mock.Anything, registrationID).Return(session, nil)
				redis.On("MarkRegistrationPasswordSet", mock.Anything, registrationID, mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil)
			},
		},
		{
			name: "success - password re-set (idempotent)",
			req:  validReq,
			setupToken: func() string {
				tokenString, _ := generateValidToken()
				return tokenString
			},
			setupMocks: func(redis *MockInMemoryStore, tokenHash string) {
				session := &entity.RegistrationSession{
					ID:                    registrationID,
					Email:                 email,
					Status:                entity.RegistrationSessionStatusPasswordSet,
					RegistrationTokenHash: &tokenHash,
					ExpiresAt:             time.Now().Add(10 * time.Minute),
				}
				redis.On("GetRegistrationSession", mock.Anything, registrationID).Return(session, nil)
				redis.On("MarkRegistrationPasswordSet", mock.Anything, registrationID, mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil)
			},
		},
		{
			name: "error - invalid token",
			req:  validReq,
			setupToken: func() string {
				return "invalid-token"
			},
			setupMocks:    func(redis *MockInMemoryStore, tokenHash string) {},
			expectedError: "invalid",
			expectedCode:  errors.CodeUnauthorized,
		},
		{
			name: "error - expired session",
			req:  validReq,
			setupToken: func() string {
				tokenString, _ := generateValidToken()
				return tokenString
			},
			setupMocks: func(redis *MockInMemoryStore, tokenHash string) {
				session := &entity.RegistrationSession{
					ID:                    registrationID,
					Email:                 email,
					Status:                entity.RegistrationSessionStatusVerified,
					RegistrationTokenHash: &tokenHash,
					ExpiresAt:             time.Now().Add(-1 * time.Minute),
				}
				redis.On("GetRegistrationSession", mock.Anything, registrationID).Return(session, nil)
			},
			expectedError: "expired",
			expectedCode:  "REGISTRATION_EXPIRED",
		},
		{
			name: "error - session not verified",
			req:  validReq,
			setupToken: func() string {
				tokenString, _ := generateValidToken()
				return tokenString
			},
			setupMocks: func(redis *MockInMemoryStore, tokenHash string) {
				session := &entity.RegistrationSession{
					ID:                    registrationID,
					Email:                 email,
					Status:                entity.RegistrationSessionStatusPendingVerification,
					RegistrationTokenHash: &tokenHash,
					ExpiresAt:             time.Now().Add(10 * time.Minute),
				}
				redis.On("GetRegistrationSession", mock.Anything, registrationID).Return(session, nil)
			},
			expectedError: "not been verified",
			expectedCode:  errors.CodeForbidden,
		},
		{
			name: "error - weak password (no uppercase)",
			req: &auth.SetPasswordRequest{
				Password:             "weakpassword1!",
				ConfirmationPassword: "weakpassword1!",
			},
			setupToken: func() string {
				tokenString, _ := generateValidToken()
				return tokenString
			},
			setupMocks: func(redis *MockInMemoryStore, tokenHash string) {
				session := &entity.RegistrationSession{
					ID:                    registrationID,
					Email:                 email,
					Status:                entity.RegistrationSessionStatusVerified,
					RegistrationTokenHash: &tokenHash,
					ExpiresAt:             time.Now().Add(10 * time.Minute),
				}
				redis.On("GetRegistrationSession", mock.Anything, registrationID).Return(session, nil)
			},
			expectedError: "uppercase",
			expectedCode:  errors.CodeValidation,
		},
		{
			name: "error - password mismatch",
			req: &auth.SetPasswordRequest{
				Password:             "SecureP@ssw0rd!",
				ConfirmationPassword: "DifferentP@ssw0rd!",
			},
			setupToken: func() string {
				tokenString, _ := generateValidToken()
				return tokenString
			},
			setupMocks: func(redis *MockInMemoryStore, tokenHash string) {
				session := &entity.RegistrationSession{
					ID:                    registrationID,
					Email:                 email,
					Status:                entity.RegistrationSessionStatusVerified,
					RegistrationTokenHash: &tokenHash,
					ExpiresAt:             time.Now().Add(10 * time.Minute),
				}
				redis.On("GetRegistrationSession", mock.Anything, registrationID).Return(session, nil)
			},
			expectedError: "do not match",
			expectedCode:  errors.CodeValidation,
		},
		{
			name: "error - session already completed",
			req:  validReq,
			setupToken: func() string {
				tokenString, _ := generateValidToken()
				return tokenString
			},
			setupMocks: func(redis *MockInMemoryStore, tokenHash string) {
				session := &entity.RegistrationSession{
					ID:                    registrationID,
					Email:                 email,
					Status:                entity.RegistrationSessionStatusCompleted,
					RegistrationTokenHash: &tokenHash,
					ExpiresAt:             time.Now().Add(10 * time.Minute),
				}
				redis.On("GetRegistrationSession", mock.Anything, registrationID).Return(session, nil)
			},
			expectedError: "not been verified",
			expectedCode:  errors.CodeForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			redis := &MockInMemoryStore{}
			tokenString := tt.setupToken()
			_, tokenHash := generateValidToken()
			if tokenString != "invalid-token" {
				hash := sha256.Sum256([]byte(tokenString))
				tokenHash = hex.EncodeToString(hash[:])
			}
			tt.setupMocks(redis, tokenHash)

			cfg := &config.Config{
				JWT: config.JWTConfig{
					AccessSecret: jwtSecret,
				},
			}

			uc := auth.NewUsecase(nil, cfg, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, redis, nil, nil, nil, nil, nil)

			tt.req.RegistrationID = registrationID
			tt.req.RegistrationToken = tokenString

			response, err := uc.SetPassword(
				context.Background(),
				tt.req,
			)

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				appErr, ok := err.(*errors.AppError)
				require.True(t, ok, "Error should be AppError")
				assert.Equal(t, tt.expectedCode, appErr.Code)
				redis.AssertExpectations(t)
				return
			}

			require.NoError(t, err)
			assert.NotNil(t, response)
			assert.Equal(t, registrationID.String(), response.RegistrationID)
			assert.Equal(t, string(entity.RegistrationSessionStatusPasswordSet), response.Status)
			assert.Contains(t, response.Message, "Password set successfully")
			assert.NotEmpty(t, response.RegistrationToken)
			assert.Equal(t, "set-profile", response.NextStep.Action)
			assert.Contains(t, response.NextStep.Endpoint, "/complete-profile")
			assert.Equal(t, []string{"full_name", "gender", "date_of_birth"}, response.NextStep.RequiredFields)
			redis.AssertExpectations(t)
		})
	}
}
