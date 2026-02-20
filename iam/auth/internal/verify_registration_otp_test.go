package internal

import (
	"context"
	"testing"
	"time"

	"iam-service/config"
	"iam-service/entity"
	"iam-service/iam/auth/authdto"
	"iam-service/pkg/errors"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestVerifyRegistrationOTP(t *testing.T) {
	registrationID := uuid.New()
	email := "test@example.com"
	validOTP := "123456"

	otpHash, _ := bcrypt.GenerateFromPassword([]byte(validOTP), bcrypt.DefaultCost)

	tests := []struct {
		name          string
		req           *authdto.VerifyRegistrationOTPRequest
		setupMocks    func(*MockInMemoryStore)
		expectedError string
		expectedCode  string
	}{
		{
			name: "success - OTP verified",
			req:  &authdto.VerifyRegistrationOTPRequest{Email: email, OTPCode: validOTP},
			setupMocks: func(redis *MockInMemoryStore) {
				session := &entity.RegistrationSession{
					ID:           registrationID,
					Email:        email,
					Status:       entity.RegistrationSessionStatusPendingVerification,
					OTPHash:      string(otpHash),
					OTPExpiresAt: time.Now().Add(10 * time.Minute),
					Attempts:     0,
					MaxAttempts:  5,
					ExpiresAt:    time.Now().Add(10 * time.Minute),
				}
				redis.On("GetRegistrationSession", mock.Anything, registrationID).Return(session, nil)
				redis.On("MarkRegistrationVerified", mock.Anything, registrationID, mock.AnythingOfType("string")).Return(nil)
			},
		},
		{
			name: "error - registration not found",
			req:  &authdto.VerifyRegistrationOTPRequest{Email: email, OTPCode: validOTP},
			setupMocks: func(redis *MockInMemoryStore) {
				redis.On("GetRegistrationSession", mock.Anything, registrationID).Return(nil, errors.ErrNotFound("registration not found"))
			},
			expectedError: "not found",
			expectedCode:  errors.CodeNotFound,
		},
		{
			name: "error - email mismatch",
			req:  &authdto.VerifyRegistrationOTPRequest{Email: "wrong@example.com", OTPCode: validOTP},
			setupMocks: func(redis *MockInMemoryStore) {
				session := &entity.RegistrationSession{
					ID:           registrationID,
					Email:        email,
					Status:       entity.RegistrationSessionStatusPendingVerification,
					OTPHash:      string(otpHash),
					OTPExpiresAt: time.Now().Add(10 * time.Minute),
					ExpiresAt:    time.Now().Add(10 * time.Minute),
				}
				redis.On("GetRegistrationSession", mock.Anything, registrationID).Return(session, nil)
			},
			expectedError: "Email does not match",
			expectedCode:  errors.CodeValidation,
		},
		{
			name: "error - session expired",
			req:  &authdto.VerifyRegistrationOTPRequest{Email: email, OTPCode: validOTP},
			setupMocks: func(redis *MockInMemoryStore) {
				session := &entity.RegistrationSession{
					ID:           registrationID,
					Email:        email,
					Status:       entity.RegistrationSessionStatusPendingVerification,
					OTPHash:      string(otpHash),
					OTPExpiresAt: time.Now().Add(10 * time.Minute),
					ExpiresAt:    time.Now().Add(-1 * time.Minute),
				}
				redis.On("GetRegistrationSession", mock.Anything, registrationID).Return(session, nil)
			},
			expectedError: "expired",
			expectedCode:  "REGISTRATION_EXPIRED",
		},
		{
			name: "error - already verified",
			req:  &authdto.VerifyRegistrationOTPRequest{Email: email, OTPCode: validOTP},
			setupMocks: func(redis *MockInMemoryStore) {
				verifiedAt := time.Now()
				session := &entity.RegistrationSession{
					ID:           registrationID,
					Email:        email,
					Status:       entity.RegistrationSessionStatusVerified,
					OTPHash:      string(otpHash),
					OTPExpiresAt: time.Now().Add(10 * time.Minute),
					VerifiedAt:   &verifiedAt,
					ExpiresAt:    time.Now().Add(10 * time.Minute),
				}
				redis.On("GetRegistrationSession", mock.Anything, registrationID).Return(session, nil)
			},
			expectedError: "already verified",
			expectedCode:  errors.CodeConflict,
		},
		{
			name: "error - OTP expired",
			req:  &authdto.VerifyRegistrationOTPRequest{Email: email, OTPCode: validOTP},
			setupMocks: func(redis *MockInMemoryStore) {
				session := &entity.RegistrationSession{
					ID:           registrationID,
					Email:        email,
					Status:       entity.RegistrationSessionStatusPendingVerification,
					OTPHash:      string(otpHash),
					OTPExpiresAt: time.Now().Add(-1 * time.Minute),
					ExpiresAt:    time.Now().Add(10 * time.Minute),
				}
				redis.On("GetRegistrationSession", mock.Anything, registrationID).Return(session, nil)
			},
			expectedError: "expired",
			expectedCode:  "OTP_EXPIRED",
		},
		{
			name: "error - max attempts exceeded (session failed)",
			req:  &authdto.VerifyRegistrationOTPRequest{Email: email, OTPCode: validOTP},
			setupMocks: func(redis *MockInMemoryStore) {
				session := &entity.RegistrationSession{
					ID:           registrationID,
					Email:        email,
					Status:       entity.RegistrationSessionStatusFailed,
					OTPHash:      string(otpHash),
					OTPExpiresAt: time.Now().Add(10 * time.Minute),
					Attempts:     5,
					MaxAttempts:  5,
					ExpiresAt:    time.Now().Add(10 * time.Minute),
				}
				redis.On("GetRegistrationSession", mock.Anything, registrationID).Return(session, nil)
			},
			expectedError: "Too many failed attempts",
			expectedCode:  errors.CodeTooManyRequests,
		},
		{
			name: "error - invalid OTP",
			req:  &authdto.VerifyRegistrationOTPRequest{Email: email, OTPCode: "000000"},
			setupMocks: func(redis *MockInMemoryStore) {
				session := &entity.RegistrationSession{
					ID:           registrationID,
					Email:        email,
					Status:       entity.RegistrationSessionStatusPendingVerification,
					OTPHash:      string(otpHash),
					OTPExpiresAt: time.Now().Add(10 * time.Minute),
					Attempts:     0,
					MaxAttempts:  5,
					ExpiresAt:    time.Now().Add(10 * time.Minute),
				}
				redis.On("GetRegistrationSession", mock.Anything, registrationID).Return(session, nil)
				redis.On("IncrementRegistrationAttempts", mock.Anything, registrationID).Return(1, nil)
			},
			expectedError: "incorrect",
			expectedCode:  errors.CodeUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			redis := new(MockInMemoryStore)
			tt.setupMocks(redis)

			uc := &usecase{
				Config: &config.Config{
					JWT: config.JWTConfig{
						AccessSecret: "test-secret-key-for-testing-purposes",
					},
				},
				InMemoryStore: redis,
			}

			ctx := context.Background()
			tt.req.RegistrationID = registrationID

			resp, err := uc.VerifyRegistrationOTP(ctx, tt.req)

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
				assert.Equal(t, registrationID.String(), resp.RegistrationID)
				assert.Equal(t, string(entity.RegistrationSessionStatusVerified), resp.Status)
				assert.NotEmpty(t, resp.RegistrationToken)
			}

			redis.AssertExpectations(t)
		})
	}
}
