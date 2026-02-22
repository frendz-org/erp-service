package entity

import (
	"time"

	"github.com/google/uuid"
)

type RegistrationSessionStatus string

const (
	RegistrationSessionStatusPendingVerification RegistrationSessionStatus = "PENDING_VERIFICATION"
	RegistrationSessionStatusVerified            RegistrationSessionStatus = "VERIFIED"
	RegistrationSessionStatusPasswordSet         RegistrationSessionStatus = "PASSWORD_SET"
	RegistrationSessionStatusCompleted           RegistrationSessionStatus = "COMPLETED"
	RegistrationSessionStatusFailed              RegistrationSessionStatus = "FAILED"
	RegistrationSessionStatusExpired             RegistrationSessionStatus = "EXPIRED"
)

type RegistrationSession struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`

	Status RegistrationSessionStatus `json:"status"`

	OTPHash      string    `json:"otp_hash"`
	OTPCreatedAt time.Time `json:"otp_created_at"`
	OTPExpiresAt time.Time `json:"otp_expires_at"`

	Attempts    int `json:"attempts"`
	MaxAttempts int `json:"max_attempts"`

	ResendCount           int        `json:"resend_count"`
	MaxResends            int        `json:"max_resends"`
	LastResentAt          *time.Time `json:"last_resent_at,omitempty"`
	ResendCooldownSeconds int        `json:"resend_cooldown_seconds"`

	VerifiedAt            *time.Time `json:"verified_at,omitempty"`
	RegistrationTokenHash *string    `json:"registration_token_hash,omitempty"`

	PasswordSetAt *time.Time `json:"password_set_at,omitempty"`

	IPAddress string `json:"ip_address"`
	UserAgent string `json:"user_agent"`

	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

func (s *RegistrationSession) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

func (s *RegistrationSession) IsOTPExpired() bool {
	return time.Now().After(s.OTPExpiresAt)
}

func (s *RegistrationSession) IsPendingVerification() bool {
	return s.Status == RegistrationSessionStatusPendingVerification
}

func (s *RegistrationSession) IsVerified() bool {
	return s.Status == RegistrationSessionStatusVerified
}

func (s *RegistrationSession) CanAttemptOTP() bool {
	return s.IsPendingVerification() &&
		!s.IsExpired() &&
		!s.IsOTPExpired() &&
		s.Attempts < s.MaxAttempts
}

func (s *RegistrationSession) CanResendOTP() bool {
	if !s.IsPendingVerification() || s.IsExpired() {
		return false
	}
	if s.ResendCount >= s.MaxResends {
		return false
	}
	if s.LastResentAt != nil {
		cooldownEnd := s.LastResentAt.Add(time.Duration(s.ResendCooldownSeconds) * time.Second)
		if time.Now().Before(cooldownEnd) {
			return false
		}
	}
	return true
}

func (s *RegistrationSession) CooldownRemainingSeconds() int {
	if s.LastResentAt == nil {
		return 0
	}
	cooldownEnd := s.LastResentAt.Add(time.Duration(s.ResendCooldownSeconds) * time.Second)
	remaining := time.Until(cooldownEnd)
	if remaining <= 0 {
		return 0
	}
	return int(remaining.Seconds())
}

func (s *RegistrationSession) RemainingAttempts() int {
	remaining := s.MaxAttempts - s.Attempts
	if remaining < 0 {
		return 0
	}
	return remaining
}

func (s *RegistrationSession) RemainingResends() int {
	remaining := s.MaxResends - s.ResendCount
	if remaining < 0 {
		return 0
	}
	return remaining
}

func (s *RegistrationSession) CanSetPassword() bool {
	return (s.Status == RegistrationSessionStatusVerified || s.Status == RegistrationSessionStatusPasswordSet) && !s.IsExpired()
}

func (s *RegistrationSession) IsPasswordSet() bool {
	return s.Status == RegistrationSessionStatusPasswordSet
}

func (s *RegistrationSession) CanCompleteProfile() bool {
	return s.IsPasswordSet() && !s.IsExpired()
}
