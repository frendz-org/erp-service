package auth_test

import (
	"context"
	"testing"
	"time"

	"erp-service/config"
	"erp-service/entity"
	"erp-service/iam/auth"
	"erp-service/pkg/errors"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestInitiateRegistration(t *testing.T) {
	email := "test@example.com"
	ipAddress := "192.168.1.1"
	userAgent := "Mozilla/5.0"

	tests := []struct {
		name          string
		req           *auth.InitiateRegistrationRequest
		setupMocks    func(*MockUserRepository, *MockInMemoryStore, *MockEmailService)
		expectedError string
		expectedCode  string
	}{
		{
			name: "success - registration initiated",
			req:  &auth.InitiateRegistrationRequest{Email: email, IPAddress: ipAddress, UserAgent: userAgent},
			setupMocks: func(userRepo *MockUserRepository, redis *MockInMemoryStore, emailSvc *MockEmailService) {
				userRepo.On("EmailExists", mock.Anything, email).Return(false, nil)
				redis.On("IncrementRegistrationRateLimit", mock.Anything, email, mock.Anything).Return(int64(1), nil)
				redis.On("IsRegistrationEmailLocked", mock.Anything, email).Return(false, nil)
				redis.On("LockRegistrationEmail", mock.Anything, email, mock.Anything).Return(true, nil)
				redis.On("CreateRegistrationSession", mock.Anything, mock.AnythingOfType("*entity.RegistrationSession"), mock.Anything).Return(nil)
				emailSvc.On("SendRegistrationOTP", mock.Anything, email, mock.AnythingOfType("string"), auth.RegistrationOTPExpiryMinutes).Return(nil)
			},
		},
		{
			name: "success - email already exists (returns fake success to prevent enumeration)",
			req:  &auth.InitiateRegistrationRequest{Email: email, IPAddress: ipAddress, UserAgent: userAgent},
			setupMocks: func(userRepo *MockUserRepository, redis *MockInMemoryStore, emailSvc *MockEmailService) {
				redis.On("IncrementRegistrationRateLimit", mock.Anything, email, mock.Anything).Return(int64(1), nil)
				userRepo.On("EmailExists", mock.Anything, email).Return(true, nil)
			},
		},
		{
			name: "error - rate limit exceeded",
			req:  &auth.InitiateRegistrationRequest{Email: email, IPAddress: ipAddress, UserAgent: userAgent},
			setupMocks: func(userRepo *MockUserRepository, redis *MockInMemoryStore, emailSvc *MockEmailService) {
				redis.On("IncrementRegistrationRateLimit", mock.Anything, email, mock.Anything).Return(int64(auth.RegistrationRateLimitPerHour+1), nil)
			},
			expectedError: "Too many registration attempts",
			expectedCode:  errors.CodeTooManyRequests,
		},
		{
			name: "error - email already locked (registration in progress)",
			req:  &auth.InitiateRegistrationRequest{Email: email, IPAddress: ipAddress, UserAgent: userAgent},
			setupMocks: func(userRepo *MockUserRepository, redis *MockInMemoryStore, emailSvc *MockEmailService) {
				userRepo.On("EmailExists", mock.Anything, email).Return(false, nil)
				redis.On("IncrementRegistrationRateLimit", mock.Anything, email, mock.Anything).Return(int64(1), nil)
				redis.On("IsRegistrationEmailLocked", mock.Anything, email).Return(true, nil)
			},
			expectedError: "An active registration already exists",
			expectedCode:  errors.CodeConflict,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := new(MockUserRepository)
			redis := new(MockInMemoryStore)
			emailSvc := new(MockEmailService)

			tt.setupMocks(userRepo, redis, emailSvc)

			uc := auth.NewUsecase(nil, &config.Config{}, userRepo, nil, nil, nil, nil, nil, nil, nil, nil, nil, emailSvc, redis, nil, nil, nil, nil, nil)

			ctx := context.Background()
			resp, err := uc.InitiateRegistration(ctx, tt.req)

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
				assert.Equal(t, email, resp.Email)
				assert.Equal(t, string(entity.RegistrationSessionStatusPendingVerification), resp.Status)
				assert.NotEmpty(t, resp.RegistrationID)
				assert.True(t, resp.ExpiresAt.After(time.Now()))
			}

			userRepo.AssertExpectations(t)
			redis.AssertExpectations(t)
		})
	}
}
