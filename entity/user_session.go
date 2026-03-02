package entity

import (
	"time"

	"github.com/google/uuid"
)

type UserSessionStatus string

const (
	UserSessionStatusActive  UserSessionStatus = "ACTIVE"
	UserSessionStatusRevoked UserSessionStatus = "REVOKED"
	UserSessionStatusExpired UserSessionStatus = "EXPIRED"
)

type UserSessionLoginMethod string

const (
	UserSessionLoginMethodEmailOTP      UserSessionLoginMethod = "EMAIL_OTP"
	UserSessionLoginMethodGoogleOAuth   UserSessionLoginMethod = "GOOGLE_OAUTH"
	UserSessionLoginMethodTransferToken UserSessionLoginMethod = "TRANSFER_TOKEN"
)

type UserSession struct {
	ID                uuid.UUID              `json:"id" gorm:"column:id;primaryKey;type:uuid;default:uuidv7()" db:"id"`
	UserID            uuid.UUID              `json:"user_id" gorm:"column:user_id;type:uuid;not null" db:"user_id"`
	ParentSessionID   *uuid.UUID             `json:"parent_session_id,omitempty" gorm:"column:parent_session_id;type:uuid" db:"parent_session_id"`
	RefreshTokenID    *uuid.UUID             `json:"refresh_token_id,omitempty" gorm:"column:refresh_token_id;type:uuid" db:"refresh_token_id"`
	IPAddress         string                 `json:"ip_address" gorm:"column:ip_address;type:inet;not null" db:"ip_address"`
	UserAgent         *string                `json:"user_agent,omitempty" gorm:"column:user_agent;type:text" db:"user_agent"`
	DeviceFingerprint *string                `json:"device_fingerprint,omitempty" gorm:"column:device_fingerprint;type:varchar(255)" db:"device_fingerprint"`
	LoginMethod       UserSessionLoginMethod `json:"login_method" gorm:"column:login_method;type:varchar(20);not null" db:"login_method"`
	Status            UserSessionStatus      `json:"status" gorm:"column:status;type:varchar(20);not null;default:ACTIVE" db:"status"`
	LastActiveAt      time.Time              `json:"last_active_at" gorm:"column:last_active_at;not null" db:"last_active_at"`
	ExpiresAt         time.Time              `json:"expires_at" gorm:"column:expires_at;not null" db:"expires_at"`
	RevokedAt         *time.Time             `json:"revoked_at,omitempty" gorm:"column:revoked_at" db:"revoked_at"`
	CreatedAt         time.Time              `json:"created_at" gorm:"column:created_at;not null" db:"created_at"`
	UpdatedAt         time.Time              `json:"updated_at" gorm:"column:updated_at;not null" db:"updated_at"`
}

func (UserSession) TableName() string {
	return "user_sessions"
}

func (s *UserSession) IsActive() bool {
	return s.Status == UserSessionStatusActive
}

func (s *UserSession) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}
