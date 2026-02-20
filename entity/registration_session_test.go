package entity

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestRegistrationSession_CanSetPassword(t *testing.T) {
	tests := []struct {
		name     string
		session  *RegistrationSession
		expected bool
	}{
		{
			name: "can set password when status is VERIFIED",
			session: &RegistrationSession{
				Status:    RegistrationSessionStatusVerified,
				ExpiresAt: time.Now().Add(10 * time.Minute),
			},
			expected: true,
		},
		{
			name: "can set password when status is PASSWORD_SET (idempotent)",
			session: &RegistrationSession{
				Status:    RegistrationSessionStatusPasswordSet,
				ExpiresAt: time.Now().Add(10 * time.Minute),
			},
			expected: true,
		},
		{
			name: "cannot set password when status is PENDING_VERIFICATION",
			session: &RegistrationSession{
				Status:    RegistrationSessionStatusPendingVerification,
				ExpiresAt: time.Now().Add(10 * time.Minute),
			},
			expected: false,
		},
		{
			name: "cannot set password when status is COMPLETED",
			session: &RegistrationSession{
				Status:    RegistrationSessionStatusCompleted,
				ExpiresAt: time.Now().Add(10 * time.Minute),
			},
			expected: false,
		},
		{
			name: "cannot set password when session is expired",
			session: &RegistrationSession{
				Status:    RegistrationSessionStatusVerified,
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
		session  *RegistrationSession
		expected bool
	}{
		{
			name: "true when status is PASSWORD_SET",
			session: &RegistrationSession{
				Status: RegistrationSessionStatusPasswordSet,
			},
			expected: true,
		},
		{
			name: "false when status is VERIFIED",
			session: &RegistrationSession{
				Status: RegistrationSessionStatusVerified,
			},
			expected: false,
		},
		{
			name: "false when status is PENDING_VERIFICATION",
			session: &RegistrationSession{
				Status: RegistrationSessionStatusPendingVerification,
			},
			expected: false,
		},
		{
			name: "false when status is COMPLETED",
			session: &RegistrationSession{
				Status: RegistrationSessionStatusCompleted,
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
		session  *RegistrationSession
		expected bool
	}{
		{
			name: "can complete profile when password is set and not expired",
			session: &RegistrationSession{
				Status:    RegistrationSessionStatusPasswordSet,
				ExpiresAt: time.Now().Add(10 * time.Minute),
			},
			expected: true,
		},
		{
			name: "cannot complete profile when status is VERIFIED (password not set)",
			session: &RegistrationSession{
				Status:    RegistrationSessionStatusVerified,
				ExpiresAt: time.Now().Add(10 * time.Minute),
			},
			expected: false,
		},
		{
			name: "cannot complete profile when status is PENDING_VERIFICATION",
			session: &RegistrationSession{
				Status:    RegistrationSessionStatusPendingVerification,
				ExpiresAt: time.Now().Add(10 * time.Minute),
			},
			expected: false,
		},
		{
			name: "cannot complete profile when status is COMPLETED",
			session: &RegistrationSession{
				Status:    RegistrationSessionStatusCompleted,
				ExpiresAt: time.Now().Add(10 * time.Minute),
			},
			expected: false,
		},
		{
			name: "cannot complete profile when session is expired",
			session: &RegistrationSession{
				Status:    RegistrationSessionStatusPasswordSet,
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
		session  *RegistrationSession
		expected bool
	}{
		{
			name: "can complete with legacy flow (VERIFIED status)",
			session: &RegistrationSession{
				Status:    RegistrationSessionStatusVerified,
				ExpiresAt: time.Now().Add(10 * time.Minute),
			},
			expected: true,
		},
		{
			name: "can complete with new flow (PASSWORD_SET status)",
			session: &RegistrationSession{
				Status:    RegistrationSessionStatusPasswordSet,
				ExpiresAt: time.Now().Add(10 * time.Minute),
			},
			expected: true,
		},
		{
			name: "cannot complete when status is PENDING_VERIFICATION",
			session: &RegistrationSession{
				Status:    RegistrationSessionStatusPendingVerification,
				ExpiresAt: time.Now().Add(10 * time.Minute),
			},
			expected: false,
		},
		{
			name: "cannot complete when status is COMPLETED",
			session: &RegistrationSession{
				Status:    RegistrationSessionStatusCompleted,
				ExpiresAt: time.Now().Add(10 * time.Minute),
			},
			expected: false,
		},
		{
			name: "cannot complete when session is expired (VERIFIED)",
			session: &RegistrationSession{
				Status:    RegistrationSessionStatusVerified,
				ExpiresAt: time.Now().Add(-1 * time.Minute),
			},
			expected: false,
		},
		{
			name: "cannot complete when session is expired (PASSWORD_SET)",
			session: &RegistrationSession{
				Status:    RegistrationSessionStatusPasswordSet,
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

		session := &RegistrationSession{
			ID:            uuid.New(),
			Email:         "test@example.com",
			Status:        RegistrationSessionStatusPasswordSet,
			PasswordSetAt: &passwordSetAt,
			ExpiresAt:     time.Now().Add(10 * time.Minute),
		}

		assert.NotNil(t, session.PasswordSetAt)
		assert.WithinDuration(t, passwordSetAt, *session.PasswordSetAt, time.Second)
	})

	t.Run("password set timestamp can be nil initially", func(t *testing.T) {
		session := &RegistrationSession{
			ID:        uuid.New(),
			Email:     "test@example.com",
			Status:    RegistrationSessionStatusVerified,
			ExpiresAt: time.Now().Add(10 * time.Minute),
		}

		assert.Nil(t, session.PasswordSetAt)
	})
}
