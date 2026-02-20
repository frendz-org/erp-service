package entity

import (
	"encoding/json"
	"net"
	"time"

	"github.com/google/uuid"
)

type RegistrationStatus string

const (
	RegistrationStatusPendingVerification RegistrationStatus = "pending_verification"
	RegistrationStatusVerified            RegistrationStatus = "verified"
	RegistrationStatusCompleted           RegistrationStatus = "completed"
	RegistrationStatusExpired             RegistrationStatus = "expired"
	RegistrationStatusCancelled           RegistrationStatus = "cancelled"
)

type Registration struct {
	ID             uuid.UUID  `json:"id" gorm:"column:id;primaryKey" db:"id"`
	TenantID       uuid.UUID  `json:"tenant_id" gorm:"column:tenant_id;not null" db:"tenant_id"`
	BranchID       *uuid.UUID `json:"branch_id,omitempty" gorm:"column:branch_id" db:"branch_id"`

	Email        string `json:"email" gorm:"column:email;not null" db:"email"`
	PasswordHash string `json:"-" gorm:"column:password_hash;not null" db:"password_hash"`

	UserAgent *string         `json:"user_agent,omitempty" gorm:"column:user_agent" db:"user_agent"`
	IPAddress *net.IP         `json:"ip_address,omitempty" gorm:"column:ip_address;type:inet" db:"ip_address"`
	Referrer  *string         `json:"referrer,omitempty" gorm:"column:referrer" db:"referrer"`
	Metadata  json.RawMessage `json:"metadata,omitempty" gorm:"column:metadata;type:jsonb;default:'{}'" db:"metadata"`

	Status RegistrationStatus `json:"status" gorm:"column:status;not null;default:'pending_verification'" db:"status"`

	CreatedAt   time.Time  `json:"created_at" gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" db:"created_at"`
	ExpiresAt   time.Time  `json:"expires_at" gorm:"column:expires_at;not null" db:"expires_at"`
	VerifiedAt  *time.Time `json:"verified_at,omitempty" gorm:"column:verified_at" db:"verified_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty" gorm:"column:completed_at" db:"completed_at"`

	UserID             *uuid.UUID `json:"user_id,omitempty" gorm:"column:user_id" db:"user_id"`
	CancellationReason *string    `json:"cancellation_reason,omitempty" gorm:"column:cancellation_reason" db:"cancellation_reason"`
}

func (Registration) TableName() string {
	return "registrations"
}

func (r *Registration) IsExpired() bool {
	return time.Now().After(r.ExpiresAt)
}

func (r *Registration) IsPending() bool {
	return r.Status == RegistrationStatusPendingVerification
}

func (r *Registration) IsVerified() bool {
	return r.Status == RegistrationStatusVerified
}

func (r *Registration) IsCompleted() bool {
	return r.Status == RegistrationStatusCompleted
}

func (r *Registration) CanBeVerified() bool {
	return r.Status == RegistrationStatusPendingVerification && !r.IsExpired()
}

func (r *Registration) CanBeCompleted() bool {
	return r.Status == RegistrationStatusVerified && !r.IsExpired()
}
