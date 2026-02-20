package entity

import (
	"time"

	"github.com/google/uuid"
)

type LoginSessionStatus string

const (
	LoginSessionStatusPendingVerification LoginSessionStatus = "PENDING_VERIFICATION"
	LoginSessionStatusVerified            LoginSessionStatus = "VERIFIED"
	LoginSessionStatusFailed              LoginSessionStatus = "FAILED"
	LoginSessionStatusExpired             LoginSessionStatus = "EXPIRED"
)

type LoginSession struct {
	ID     uuid.UUID `json:"id"`
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`

	Status LoginSessionStatus `json:"status"`

	OTPHash      string    `json:"otp_hash"`
	OTPCreatedAt time.Time `json:"otp_created_at"`
	OTPExpiresAt time.Time `json:"otp_expires_at"`

	Attempts    int `json:"attempts"`
	MaxAttempts int `json:"max_attempts"`

	ResendCount           int        `json:"resend_count"`
	MaxResends            int        `json:"max_resends"`
	LastResentAt          *time.Time `json:"last_resent_at,omitempty"`
	ResendCooldownSeconds int        `json:"resend_cooldown_seconds"`

	IPAddress string `json:"ip_address"`
	UserAgent string `json:"user_agent"`

	CreatedAt  time.Time  `json:"created_at"`
	ExpiresAt  time.Time  `json:"expires_at"`
	VerifiedAt *time.Time `json:"verified_at,omitempty"`
}

func (s *LoginSession) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

func (s *LoginSession) IsOTPExpired() bool {
	return time.Now().After(s.OTPExpiresAt)
}

func (s *LoginSession) IsPendingVerification() bool {
	return s.Status == LoginSessionStatusPendingVerification
}

func (s *LoginSession) IsVerified() bool {
	return s.Status == LoginSessionStatusVerified
}

func (s *LoginSession) IsLocked() bool {
	return s.Attempts >= s.MaxAttempts
}

func (s *LoginSession) CanAttemptOTP() bool {
	return s.IsPendingVerification() &&
		!s.IsExpired() &&
		!s.IsOTPExpired() &&
		!s.IsLocked()
}

func (s *LoginSession) CanResend() bool {
	if !s.IsPendingVerification() || s.IsExpired() {
		return false
	}
	return s.ResendCount < s.MaxResends
}

func (s *LoginSession) IsInCooldown() bool {
	if s.LastResentAt == nil {
		return false
	}
	cooldownEnd := s.LastResentAt.Add(time.Duration(s.ResendCooldownSeconds) * time.Second)
	return time.Now().Before(cooldownEnd)
}

func (s *LoginSession) CooldownRemainingSeconds() int {
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

func (s *LoginSession) RemainingAttempts() int {
	remaining := s.MaxAttempts - s.Attempts
	if remaining < 0 {
		return 0
	}
	return remaining
}

func (s *LoginSession) RemainingResends() int {
	remaining := s.MaxResends - s.ResendCount
	if remaining < 0 {
		return 0
	}
	return remaining
}
