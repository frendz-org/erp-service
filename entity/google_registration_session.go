package entity

import (
	"time"

	"github.com/google/uuid"
)

type GoogleRegistrationSessionStatus string

const (
	GoogleRegistrationSessionStatusPendingProfile GoogleRegistrationSessionStatus = "PENDING_PROFILE"
	GoogleRegistrationSessionStatusCompleted      GoogleRegistrationSessionStatus = "COMPLETED"
	GoogleRegistrationSessionStatusExpired        GoogleRegistrationSessionStatus = "EXPIRED"
)

type GoogleRegistrationSession struct {
	ID                    uuid.UUID                       `json:"id"`
	Email                 string                          `json:"email"`
	GoogleID              string                          `json:"google_id"`
	Name                  string                          `json:"name"`
	Picture               string                          `json:"picture"`
	GoogleEmailVerified   bool                            `json:"google_email_verified"`
	Status                GoogleRegistrationSessionStatus `json:"status"`
	RegistrationTokenHash string                          `json:"registration_token_hash"`
	IPAddress             string                          `json:"ip_address"`
	UserAgent             string                          `json:"user_agent"`
	CreatedAt             time.Time                       `json:"created_at"`
	ExpiresAt             time.Time                       `json:"expires_at"`
}

func (s *GoogleRegistrationSession) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

func (s *GoogleRegistrationSession) IsPendingProfile() bool {
	return s.Status == GoogleRegistrationSessionStatusPendingProfile
}
