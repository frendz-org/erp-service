package entity_test

import (
	"testing"
	"time"

	"erp-service/entity"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestRegistrationSession_CanSetPassword(t *testing.T) {
	tests := []struct {
		name     string
		session  *entity.RegistrationSession
		expected bool
	}{
		{
			name: "can set password when status is VERIFIED",
			session: &entity.RegistrationSession{
				Status:    entity.RegistrationSessionStatusVerified,
				ExpiresAt: time.Now().Add(10 * time.Minute),
			},
			expected: true,
		},
		{
			name: "can set password when status is PASSWORD_SET (idempotent)",
			session: &entity.RegistrationSession{
				Status:    entity.RegistrationSessionStatusPasswordSet,
				ExpiresAt: time.Now().Add(10 * time.Minute),
			},
			expected: true,
		},
		{
			name: "cannot set password when status is PENDING_VERIFICATION",
			session: &entity.RegistrationSession{
				Status:    entity.RegistrationSessionStatusPendingVerification,
				ExpiresAt: time.Now().Add(10 * time.Minute),
			},
			expected: false,
		},
		{
			name: "cannot set password when status is COMPLETED",
			session: &entity.RegistrationSession{
				Status:    entity.RegistrationSessionStatusCompleted,
				ExpiresAt: time.Now().Add(10 * time.Minute),
			},
			expected: false,
		},
		{
			name: "cannot set password when session is expired",
			session: &entity.RegistrationSession{
				Status:    entity.RegistrationSessionStatusVerified,
				ExpiresAt: time.Now().Add(-1 * time.Minute),
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.session.CanSetPassword()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestRegistrationSession_IsPasswordSet(t *testing.T) {
	tests := []struct {
		name     string
		session  *entity.RegistrationSession
		expected bool
	}{
		{
			name: "true when status is PASSWORD_SET",
			session: &entity.RegistrationSession{
				Status: entity.RegistrationSessionStatusPasswordSet,
			},
			expected: true,
		},
		{
			name: "false when status is VERIFIED",
			session: &entity.RegistrationSession{
				Status: entity.RegistrationSessionStatusVerified,
			},
			expected: false,
		},
		{
			name: "false when status is PENDING_VERIFICATION",
			session: &entity.RegistrationSession{
				Status: entity.RegistrationSessionStatusPendingVerification,
			},
			expected: false,
		},
		{
			name: "false when status is COMPLETED",
			session: &entity.RegistrationSession{
				Status: entity.RegistrationSessionStatusCompleted,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.session.IsPasswordSet()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestRegistrationSession_CanCompleteProfile(t *testing.T) {
	tests := []struct {
		name     string
		session  *entity.RegistrationSession
		expected bool
	}{
		{
			name: "can complete profile when password is set and not expired",
			session: &entity.RegistrationSession{
				Status:    entity.RegistrationSessionStatusPasswordSet,
				ExpiresAt: time.Now().Add(10 * time.Minute),
			},
			expected: true,
		},
		{
			name: "cannot complete profile when status is VERIFIED (password not set)",
			session: &entity.RegistrationSession{
				Status:    entity.RegistrationSessionStatusVerified,
				ExpiresAt: time.Now().Add(10 * time.Minute),
			},
			expected: false,
		},
		{
			name: "cannot complete profile when status is PENDING_VERIFICATION",
			session: &entity.RegistrationSession{
				Status:    entity.RegistrationSessionStatusPendingVerification,
				ExpiresAt: time.Now().Add(10 * time.Minute),
			},
			expected: false,
		},
		{
			name: "cannot complete profile when status is COMPLETED",
			session: &entity.RegistrationSession{
				Status:    entity.RegistrationSessionStatusCompleted,
				ExpiresAt: time.Now().Add(10 * time.Minute),
			},
			expected: false,
		},
		{
			name: "cannot complete profile when session is expired",
			session: &entity.RegistrationSession{
				Status:    entity.RegistrationSessionStatusPasswordSet,
				ExpiresAt: time.Now().Add(-1 * time.Minute),
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.session.CanCompleteProfile()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestRegistrationSession_CanComplete_UpdatedForNewFlow(t *testing.T) {
	tests := []struct {
		name     string
		session  *entity.RegistrationSession
		expected bool
	}{
		{
			name: "can complete with legacy flow (VERIFIED status)",
			session: &entity.RegistrationSession{
				Status:    entity.RegistrationSessionStatusVerified,
				ExpiresAt: time.Now().Add(10 * time.Minute),
			},
			expected: true,
		},
		{
			name: "can complete with new flow (PASSWORD_SET status)",
			session: &entity.RegistrationSession{
				Status:    entity.RegistrationSessionStatusPasswordSet,
				ExpiresAt: time.Now().Add(10 * time.Minute),
			},
			expected: true,
		},
		{
			name: "cannot complete when status is PENDING_VERIFICATION",
			session: &entity.RegistrationSession{
				Status:    entity.RegistrationSessionStatusPendingVerification,
				ExpiresAt: time.Now().Add(10 * time.Minute),
			},
			expected: false,
		},
		{
			name: "cannot complete when status is COMPLETED",
			session: &entity.RegistrationSession{
				Status:    entity.RegistrationSessionStatusCompleted,
				ExpiresAt: time.Now().Add(10 * time.Minute),
			},
			expected: false,
		},
		{
			name: "cannot complete when session is expired (VERIFIED)",
			session: &entity.RegistrationSession{
				Status:    entity.RegistrationSessionStatusVerified,
				ExpiresAt: time.Now().Add(-1 * time.Minute),
			},
			expected: false,
		},
		{
			name: "cannot complete when session is expired (PASSWORD_SET)",
			session: &entity.RegistrationSession{
				Status:    entity.RegistrationSessionStatusPasswordSet,
				ExpiresAt: time.Now().Add(-1 * time.Minute),
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.session.CanComplete()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestRegistrationSession_PasswordFields(t *testing.T) {
	t.Run("password set timestamp is set correctly", func(t *testing.T) {
		passwordSetAt := time.Now()

		session := &entity.RegistrationSession{
			ID:            uuid.New(),
			Email:         "test@example.com",
			Status:        entity.RegistrationSessionStatusPasswordSet,
			PasswordSetAt: &passwordSetAt,
			ExpiresAt:     time.Now().Add(10 * time.Minute),
		}

		assert.NotNil(t, session.PasswordSetAt)
		assert.WithinDuration(t, passwordSetAt, *session.PasswordSetAt, time.Second)
	})

	t.Run("password set timestamp can be nil initially", func(t *testing.T) {
		session := &entity.RegistrationSession{
			ID:        uuid.New(),
			Email:     "test@example.com",
			Status:    entity.RegistrationSessionStatusVerified,
			ExpiresAt: time.Now().Add(10 * time.Minute),
		}

		assert.Nil(t, session.PasswordSetAt)
	})
}
