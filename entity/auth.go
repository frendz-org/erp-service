package entity

import (
	"net"
	"time"

	"github.com/google/uuid"
)

type OTPType string

const (
	OTPTypeEmailVerification OTPType = "email_verification"
	OTPTypeRegistration      OTPType = "registration"
	OTPTypePINReset          OTPType = "pin_reset"
	OTPTypePasswordReset     OTPType = "password_reset"
	OTPTypeAdminInvitation   OTPType = "admin_invitation"
)

type EmailVerification struct {
	ID                  uuid.UUID  `json:"id" gorm:"column:id;primaryKey;type:uuid;default:uuidv7()" db:"id"`
	TenantID            uuid.UUID  `json:"tenant_id" gorm:"column:tenant_id;not null" db:"tenant_id"`
	UserID              uuid.UUID  `json:"user_id" gorm:"column:user_id;not null" db:"user_id"`
	Email               string     `json:"email" gorm:"column:email;not null" db:"email"`
	OTPCode             string     `json:"-" gorm:"column:otp_code;not null" db:"otp_code"`
	OTPHash             string     `json:"-" gorm:"column:otp_hash;not null" db:"otp_hash"`
	OTPType             OTPType    `json:"otp_type" gorm:"column:otp_type;default:email_verification" db:"otp_type"`
	ExpiresAt           time.Time  `json:"expires_at" gorm:"column:expires_at;not null" db:"expires_at"`
	VerifiedAt          *time.Time `json:"verified_at,omitempty" gorm:"column:verified_at" db:"verified_at"`
	IPAddress           net.IP     `json:"ip_address,omitempty" gorm:"column:ip_address" db:"ip_address"`
	UserAgent           string     `json:"user_agent,omitempty" gorm:"column:user_agent" db:"user_agent"`
	CreatedAt           time.Time  `json:"created_at" gorm:"column:created_at" db:"created_at"`
}

func (EmailVerification) TableName() string {
	return "email_verifications"
}

func (e *EmailVerification) IsExpired() bool {
	return time.Now().After(e.ExpiresAt)
}

func (e *EmailVerification) IsVerified() bool {
	return e.VerifiedAt != nil
}

type PasswordResetToken struct {
	ID                   uuid.UUID  `json:"id" gorm:"column:id;primaryKey;type:uuid;default:uuidv7()" db:"id"`
	TenantID             uuid.UUID  `json:"tenant_id" gorm:"column:tenant_id;not null" db:"tenant_id"`
	UserID               uuid.UUID  `json:"user_id" gorm:"column:user_id;not null" db:"user_id"`
	TokenHash            string     `json:"-" gorm:"column:token_hash;uniqueIndex;not null" db:"token_hash"`
	ExpiresAt            time.Time  `json:"expires_at" gorm:"column:expires_at;not null" db:"expires_at"`
	UsedAt               *time.Time `json:"used_at,omitempty" gorm:"column:used_at" db:"used_at"`
	IPAddress            net.IP     `json:"ip_address,omitempty" gorm:"column:ip_address" db:"ip_address"`
	UserAgent            string     `json:"user_agent,omitempty" gorm:"column:user_agent" db:"user_agent"`
	CreatedAt            time.Time  `json:"created_at" gorm:"column:created_at" db:"created_at"`
}

func (PasswordResetToken) TableName() string {
	return "password_reset_tokens"
}

func (p *PasswordResetToken) IsExpired() bool {
	return time.Now().After(p.ExpiresAt)
}

func (p *PasswordResetToken) IsUsed() bool {
	return p.UsedAt != nil
}

type RefreshToken struct {
	ID                uuid.UUID  `json:"id" gorm:"column:id;primaryKey;type:uuid;default:uuidv7()" db:"id"`
	UserID            uuid.UUID  `json:"user_id" gorm:"column:user_id;type:uuid;not null" db:"user_id"`
	TokenHash         string     `json:"-" gorm:"column:token_hash;type:varchar(64);uniqueIndex;not null" db:"token_hash"`
	TokenFamily       uuid.UUID  `json:"token_family" gorm:"column:token_family;type:uuid;default:uuidv7();not null" db:"token_family"`
	ExpiresAt         time.Time  `json:"expires_at" gorm:"column:expires_at;not null" db:"expires_at"`
	RevokedAt         *time.Time `json:"revoked_at,omitempty" gorm:"column:revoked_at" db:"revoked_at"`
	RevokedReason     *string    `json:"revoked_reason,omitempty" gorm:"column:revoked_reason;type:varchar(100)" db:"revoked_reason"`
	ReplacedByTokenID *uuid.UUID `json:"replaced_by_token_id,omitempty" gorm:"column:replaced_by_token_id;type:uuid" db:"replaced_by_token_id"`
	IPAddress         string     `json:"ip_address,omitempty" gorm:"column:ip_address;type:inet" db:"ip_address"`
	UserAgent         string     `json:"user_agent,omitempty" gorm:"column:user_agent;type:text" db:"user_agent"`
	CreatedAt         time.Time  `json:"created_at" gorm:"column:created_at;not null" db:"created_at"`
}

func (RefreshToken) TableName() string {
	return "refresh_tokens"
}

func (r *RefreshToken) IsExpired() bool {
	return time.Now().After(r.ExpiresAt)
}

func (r *RefreshToken) IsRevoked() bool {
	return r.RevokedAt != nil
}

func (r *RefreshToken) IsValid() bool {
	return !r.IsExpired() && !r.IsRevoked()
}

type PINFailureReason string

const (
	PINFailureReasonInvalidPIN    PINFailureReason = "invalid_pin"
	PINFailureReasonRateLimited   PINFailureReason = "rate_limited"
	PINFailureReasonAccountLocked PINFailureReason = "account_locked"
	PINFailureReasonPINExpired    PINFailureReason = "pin_expired"
)

type PINVerificationLog struct {
	ID                   uuid.UUID         `json:"id" gorm:"column:id;primaryKey;type:uuid;default:uuidv7()" db:"id"`
	UserID               uuid.UUID         `json:"user_id" gorm:"column:user_id;not null" db:"user_id"`
	TenantID             uuid.UUID         `json:"tenant_id" gorm:"column:tenant_id;not null" db:"tenant_id"`
	Result               bool              `json:"result" gorm:"column:result;not null" db:"result"`
	FailureReason        *PINFailureReason `json:"failure_reason,omitempty" gorm:"column:failure_reason" db:"failure_reason"`
	IPAddress            net.IP            `json:"ip_address,omitempty" gorm:"column:ip_address" db:"ip_address"`
	UserAgent            string            `json:"user_agent,omitempty" gorm:"column:user_agent" db:"user_agent"`
	Operation            string            `json:"operation,omitempty" gorm:"column:operation" db:"operation"`
	CreatedAt            time.Time         `json:"created_at" gorm:"column:created_at" db:"created_at"`
}

func (PINVerificationLog) TableName() string {
	return "pin_verification_logs"
}

func (p *PINVerificationLog) IsSuccess() bool {
	return p.Result
}

