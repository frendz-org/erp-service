package entity

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/google/uuid"
)

type Gender string

const (
	GenderMale   Gender = "GENDER_001"
	GenderFemale Gender = "GENDER_002"
	GenderOther  Gender = "GENDER_003"
)

type MaritalStatus string

const (
	MaritalStatusSingle   MaritalStatus = "single"
	MaritalStatusMarried  MaritalStatus = "married"
	MaritalStatusDivorced MaritalStatus = "divorced"
	MaritalStatusWidowed  MaritalStatus = "widowed"
)

type UserStatus string

const (
	UserStatusPendingVerification UserStatus = "PENDING_VERIFICATION"
	UserStatusActive              UserStatus = "ACTIVE"
	UserStatusInactive            UserStatus = "INACTIVE"
	UserStatusSuspended           UserStatus = "SUSPENDED"
	UserStatusLocked              UserStatus = "LOCKED"
)

type User struct {
	ID                 uuid.UUID    `json:"id" gorm:"column:id;primaryKey;type:uuid;default:uuidv7()" db:"id"`
	Email              string       `json:"email" gorm:"column:email;not null" db:"email"`
	Status             UserStatus   `json:"status" gorm:"column:status;not null;default:PENDING_VERIFICATION" db:"status"`
	StatusChangedAt    *time.Time   `json:"status_changed_at,omitempty" gorm:"column:status_changed_at" db:"status_changed_at"`
	StatusChangedBy    *uuid.UUID   `json:"status_changed_by,omitempty" gorm:"column:status_changed_by" db:"status_changed_by"`
	RegistrationSource string       `json:"registration_source" gorm:"column:registration_source" db:"registration_source"`
	Version            int          `json:"version" gorm:"column:version;not null;default:1" db:"version"`
	CreatedAt          time.Time    `json:"created_at" gorm:"column:created_at" db:"created_at"`
	UpdatedAt          time.Time    `json:"updated_at" gorm:"column:updated_at" db:"updated_at"`
	DeletedAt          sql.NullTime `json:"deleted_at,omitempty" gorm:"column:deleted_at" db:"deleted_at"`
}

func (User) TableName() string {
	return "users"
}

func (u *User) IsActive() bool {
	return u.Status == UserStatusActive
}

type UserProfile struct {
	UserID            uuid.UUID       `json:"user_id" gorm:"column:user_id;primaryKey;type:uuid" db:"user_id"`
	FirstName         string          `json:"first_name" gorm:"column:first_name;not null" db:"first_name"`
	LastName          string          `json:"last_name" gorm:"column:last_name;not null" db:"last_name"`
	PhoneNumber       *string         `json:"phone_number,omitempty" gorm:"column:phone_number" db:"phone_number"`
	DateOfBirth       *time.Time      `json:"date_of_birth,omitempty" gorm:"column:date_of_birth" db:"date_of_birth"`
	Gender            *Gender         `json:"gender,omitempty" gorm:"column:gender" db:"gender"`
	MaritalStatus     *MaritalStatus  `json:"marital_status,omitempty" gorm:"column:marital_status" db:"marital_status"`
	Address           *string         `json:"address,omitempty" gorm:"column:address" db:"address"`
	IDNumber          *string         `json:"id_number,omitempty" gorm:"column:id_number" db:"id_number"`
	ProfilePictureURL *string         `json:"profile_picture_url,omitempty" gorm:"column:profile_picture_url" db:"profile_picture_url"`
	Metadata          json.RawMessage `json:"metadata,omitempty" gorm:"column:metadata;type:jsonb;not null;default:'{}'" db:"metadata"`
	UpdatedAt         time.Time       `json:"updated_at" gorm:"column:updated_at" db:"updated_at"`
}

func (UserProfile) TableName() string {
	return "user_profiles"
}

func (u *UserProfile) FullName() string {
	return u.FirstName + " " + u.LastName
}

type AuthMethodType string

const (
	AuthMethodPassword AuthMethodType = "PASSWORD"
	AuthMethodPIN      AuthMethodType = "PIN"
	AuthMethodSSO      AuthMethodType = "SSO"
)

type PasswordCredentialData struct {
	PasswordHash    string   `json:"password_hash"`
	PasswordHistory []string `json:"password_history,omitempty"`
}

type UserAuthMethod struct {
	ID             uuid.UUID       `json:"id" gorm:"column:id;primaryKey;type:uuid;default:uuidv7()" db:"id"`
	UserID         uuid.UUID       `json:"user_id" gorm:"column:user_id;not null" db:"user_id"`
	MethodType     string          `json:"method_type" gorm:"column:method_type;not null" db:"method_type"`
	CredentialData json.RawMessage `json:"credential_data" gorm:"column:credential_data;type:jsonb" db:"credential_data"`
	IsActive       bool            `json:"is_active" gorm:"column:is_active;not null;default:true" db:"is_active"`
	CreatedAt      time.Time       `json:"created_at" gorm:"column:created_at" db:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at" gorm:"column:updated_at" db:"updated_at"`
}

func (UserAuthMethod) TableName() string {
	return "user_auth_methods"
}

func NewPasswordAuthMethod(userID uuid.UUID, passwordHash string) *UserAuthMethod {
	data := PasswordCredentialData{
		PasswordHash:    passwordHash,
		PasswordHistory: []string{},
	}
	credJSON, err := json.Marshal(data)
	if err != nil {
		panic(fmt.Sprintf("NewPasswordAuthMethod: failed to marshal credential data: %v", err))
	}
	now := time.Now()
	return &UserAuthMethod{
		UserID:         userID,
		MethodType:     string(AuthMethodPassword),
		CredentialData: credJSON,
		IsActive:       true,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

func (m *UserAuthMethod) GetPasswordData() (*PasswordCredentialData, error) {
	var data PasswordCredentialData
	if err := json.Unmarshal(m.CredentialData, &data); err != nil {
		return nil, err
	}
	return &data, nil
}

func (m *UserAuthMethod) GetPasswordHash() string {
	data, err := m.GetPasswordData()
	if err != nil {
		return ""
	}
	return data.PasswordHash
}

type UserSecurityState struct {
	UserID              uuid.UUID  `json:"user_id" gorm:"column:user_id;primaryKey;type:uuid" db:"user_id"`
	FailedLoginAttempts int        `json:"failed_login_attempts" gorm:"column:failed_login_attempts;default:0" db:"failed_login_attempts"`
	FailedPINAttempts   int        `json:"failed_pin_attempts" gorm:"column:failed_pin_attempts;default:0" db:"failed_pin_attempts"`
	LockedUntil         *time.Time `json:"locked_until,omitempty" gorm:"column:locked_until" db:"locked_until"`
	LastLoginAt         *time.Time `json:"last_login_at,omitempty" gorm:"column:last_login_at" db:"last_login_at"`
	LastLoginIP         net.IP     `json:"last_login_ip,omitempty" gorm:"column:last_login_ip" db:"last_login_ip"`
	EmailVerified       bool       `json:"email_verified" gorm:"column:email_verified;default:false" db:"email_verified"`
	EmailVerifiedAt     *time.Time `json:"email_verified_at,omitempty" gorm:"column:email_verified_at" db:"email_verified_at"`
	PINVerified         bool       `json:"pin_verified" gorm:"column:pin_verified;default:false" db:"pin_verified"`
	ForcePasswordChange bool       `json:"force_password_change" gorm:"column:force_password_change;default:false" db:"force_password_change"`
	UpdatedAt           time.Time  `json:"updated_at" gorm:"column:updated_at" db:"updated_at"`
}

func (UserSecurityState) TableName() string {
	return "user_security_states"
}
