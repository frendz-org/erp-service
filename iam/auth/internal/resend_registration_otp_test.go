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
)

func TestResendRegistrationOTP(t *testing.T) {
	registrationID := uuid.New()
	email := "test@example.com"

	tests := []struct {
		name          string
		req           *authdto.ResendRegistrationOTPRequest
		setupMocks    func(*MockInMemoryStore, *MockEmailService)
		expectedError string
		expectedCode  string
	}{
		{
			name: "success - OTP resent",
			req:  &authdto.ResendRegistrationOTPRequest{Email: email},
			setupMocks: func(redis *MockInMemoryStore, emailSvc *MockEmailService) {
				session := &entity.RegistrationSession{
					ID:                    registrationID,
					Email:                 email,
					Status:                entity.RegistrationSessionStatusPendingVerification,
					ResendCount:           0,
					MaxResends:            3,
					ResendCooldownSeconds: 60,
					LastResentAt:          nil,
					ExpiresAt:             time.Now().Add(10 * time.Minute),
				}
				redis.On("GetRegistrationSession", mock.Anything, registrationID).Return(session, nil)
				redis.On("UpdateRegistrationOTP", mock.Anything, registrationID, mock.AnythingOfType("string"), mock.AnythingOfType("time.Time")).Return(nil)
				emailSvc.On("SendRegistrationOTP", mock.Anything, email, mock.AnythingOfType("string"), RegistrationOTPExpiryMinutes).Return(nil)
			},
		},
		{
			name: "error - registration not found",
			req:  &authdto.ResendRegistrationOTPRequest{Email: email},
			setupMocks: func(redis *MockInMemoryStore, emailSvc *MockEmailService) {
				redis.On("GetRegistrationSession", mock.Anything, registrationID).Return(nil, errors.ErrNotFound("registration not found"))
			},
			expectedError: "not found",
			expectedCode:  errors.CodeNotFound,
		},
		{
			name: "error - email mismatch",
			req:  &authdto.ResendRegistrationOTPRequest{Email: "wrong@example.com"},
			setupMocks: func(redis *MockInMemoryStore, emailSvc *MockEmailService) {
				session := &entity.RegistrationSession{
					ID:        registrationID,
					Email:     email,
					Status:    entity.RegistrationSessionStatusPendingVerification,
					ExpiresAt: time.Now().Add(10 * time.Minute),
				}
				redis.On("GetRegistrationSession", mock.Anything, registrationID).Return(session, nil)
			},
			expectedError: "Email does not match",
			expectedCode:  errors.CodeValidation,
		},
		{
			name: "error - session expired",
			req:  &authdto.ResendRegistrationOTPRequest{Email: email},
			setupMocks: func(redis *MockInMemoryStore, emailSvc *MockEmailService) {
				session := &entity.RegistrationSession{
					ID:        registrationID,
					Email:     email,
					Status:    entity.RegistrationSessionStatusPendingVerification,
					ExpiresAt: time.Now().Add(-1 * time.Minute),
				}
				redis.On("GetRegistrationSession", mock.Anything, registrationID).Return(session, nil)
			},
			expectedError: "expired",
			expectedCode:  "REGISTRATION_EXPIRED",
		},
		{
			name: "error - already verified",
			req:  &authdto.ResendRegistrationOTPRequest{Email: email},
			setupMocks: func(redis *MockInMemoryStore, emailSvc *MockEmailService) {
				verifiedAt := time.Now()
				session := &entity.RegistrationSession{
					ID:         registrationID,
					Email:      email,
					Status:     entity.RegistrationSessionStatusVerified,
					VerifiedAt: &verifiedAt,
					ExpiresAt:  time.Now().Add(10 * time.Minute),
				}
				redis.On("GetRegistrationSession", mock.Anything, registrationID).Return(session, nil)
			},
			expectedError: "already verified",
			expectedCode:  errors.CodeConflict,
		},
		{
			name: "error - max resends exceeded",
			req:  &authdto.ResendRegistrationOTPRequest{Email: email},
			setupMocks: func(redis *MockInMemoryStore, emailSvc *MockEmailService) {
				session := &entity.RegistrationSession{
					ID:          registrationID,
					Email:       email,
					Status:      entity.RegistrationSessionStatusPendingVerification,
					ResendCount: 3,
					MaxResends:  3,
					ExpiresAt:   time.Now().Add(10 * time.Minute),
				}
				redis.On("GetRegistrationSession", mock.Anything, registrationID).Return(session, nil)
			},
			expectedError: "Maximum number of resends",
			expectedCode:  errors.CodeTooManyRequests,
		},
		{
			name: "error - resend cooldown",
			req:  &authdto.ResendRegistrationOTPRequest{Email: email},
			setupMocks: func(redis *MockInMemoryStore, emailSvc *MockEmailService) {
				lastResent := time.Now().Add(-30 * time.Second)
				session := &entity.RegistrationSession{
					ID:                    registrationID,
					Email:                 email,
					Status:                entity.RegistrationSessionStatusPendingVerification,
					ResendCount:           1,
					MaxResends:            3,
					ResendCooldownSeconds: 60,
					LastResentAt:          &lastResent,
					ExpiresAt:             time.Now().Add(10 * time.Minute),
				}
				redis.On("GetRegistrationSession", mock.Anything, registrationID).Return(session, nil)
			},
			expectedError: "Please wait before requesting",
			expectedCode:  errors.CodeTooManyRequests,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			redis := new(MockInMemoryStore)
			emailSvc := new(MockEmailService)
			tt.setupMocks(redis, emailSvc)

			uc := &usecase{
				Config:       &config.Config{},
				InMemoryStore:        redis,
				EmailService: emailSvc,
			}

			ctx := context.Background()
			tt.req.RegistrationID = registrationID

			resp, err := uc.ResendRegistrationOTP(ctx, tt.req)

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
				assert.True(t, resp.ResendsRemaining >= 0)
			}

			redis.AssertExpectations(t)
		})
	}
}
