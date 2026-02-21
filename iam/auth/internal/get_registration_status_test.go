package internal

import (
	"context"
	"testing"
	"time"

	"erp-service/config"
	"erp-service/entity"
	"erp-service/pkg/errors"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetRegistrationStatus(t *testing.T) {
	registrationID := uuid.New()
	email := "test@example.com"

	tests := []struct {
		name          string
		email         string
		setupMocks    func(*MockInMemoryStore)
		expectedError string
		expectedCode  string
		checkResponse func(*testing.T, *entity.RegistrationSession)
	}{
		{
			name:  "success - status retrieved",
			email: email,
			setupMocks: func(redis *MockInMemoryStore) {
				session := &entity.RegistrationSession{
					ID:          registrationID,
					Email:       email,
					Status:      entity.RegistrationSessionStatusPendingVerification,
					Attempts:    2,
					MaxAttempts: 5,
					ResendCount: 1,
					MaxResends:  3,
					ExpiresAt:   time.Now().Add(10 * time.Minute),
				}
				redis.On("GetRegistrationSession", mock.Anything, registrationID).Return(session, nil)
			},
			checkResponse: func(t *testing.T, session *entity.RegistrationSession) {},
		},
		{
			name:  "success - case insensitive email match",
			email: "TEST@EXAMPLE.COM",
			setupMocks: func(redis *MockInMemoryStore) {
				session := &entity.RegistrationSession{
					ID:          registrationID,
					Email:       email,
					Status:      entity.RegistrationSessionStatusVerified,
					Attempts:    0,
					MaxAttempts: 5,
					ResendCount: 0,
					MaxResends:  3,
					ExpiresAt:   time.Now().Add(10 * time.Minute),
				}
				redis.On("GetRegistrationSession", mock.Anything, registrationID).Return(session, nil)
			},
		},
		{
			name:  "error - registration not found",
			email: email,
			setupMocks: func(redis *MockInMemoryStore) {
				redis.On("GetRegistrationSession", mock.Anything, registrationID).Return(nil, errors.ErrNotFound("registration not found"))
			},
			expectedError: "not found",
			expectedCode:  errors.CodeNotFound,
		},
		{
			name:  "error - email mismatch (returns not found for security)",
			email: "wrong@example.com",
			setupMocks: func(redis *MockInMemoryStore) {
				session := &entity.RegistrationSession{
					ID:        registrationID,
					Email:     email,
					Status:    entity.RegistrationSessionStatusPendingVerification,
					ExpiresAt: time.Now().Add(10 * time.Minute),
				}
				redis.On("GetRegistrationSession", mock.Anything, registrationID).Return(session, nil)
			},
			expectedError: "not found",
			expectedCode:  errors.CodeNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			redis := new(MockInMemoryStore)
			tt.setupMocks(redis)

			uc := &usecase{
				Config:        &config.Config{},
				InMemoryStore: redis,
			}

			ctx := context.Background()
			resp, err := uc.GetRegistrationStatus(ctx, registrationID, tt.email)

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
				assert.NotEmpty(t, resp.Status)

				assert.Contains(t, resp.Email, "***@")
				assert.True(t, resp.OTPAttemptsRemaining >= 0)
				assert.True(t, resp.ResendsRemaining >= 0)
			}

			redis.AssertExpectations(t)
		})
	}
}

func TestMaskEmailForRegistration(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"user@example.com", "u***@example.com"},
		{"a@example.com", "a***@example.com"},
		{"ab@example.com", "a***@example.com"},
		{"invalid", "***"},
		{"@example.com", "***@example.com"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := maskEmail(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
