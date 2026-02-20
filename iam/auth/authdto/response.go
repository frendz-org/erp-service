package authdto

import (
	"time"

	"github.com/google/uuid"
)

type UserResponse struct {
	ID         uuid.UUID  `json:"id"`
	Email      string     `json:"email"`
	FullName   string     `json:"full_name"`
	TenantID   *uuid.UUID `json:"tenant_id,omitempty"`
	ProductID  *uuid.UUID `json:"product_id,omitempty"`
	BranchID   *uuid.UUID `json:"branch_id,omitempty"`
	Roles      []string   `json:"roles"`
	MFAEnabled bool       `json:"mfa_enabled"`
}

type ProductResponse struct {
	ProductID   uuid.UUID `json:"product_id"`
	ProductCode string    `json:"product_code"`
	Roles       []string  `json:"roles,omitempty"`
	Permissions []string  `json:"permissions,omitempty"`
}

type TenantResponse struct {
	TenantID uuid.UUID         `json:"tenant_id"`
	Products []ProductResponse `json:"products,omitempty"`
}

type LoginUserResponse struct {
	ID       uuid.UUID        `json:"id"`
	Email    string           `json:"email"`
	FullName string           `json:"full_name"`
	Tenants  []TenantResponse `json:"tenants,omitempty"`
}

type RefreshTokenResponse struct {
	AccessToken  string            `json:"access_token"`
	RefreshToken string            `json:"refresh_token"`
	ExpiresIn    int               `json:"expires_in"`
	TokenType    string            `json:"token_type"`
	User         LoginUserResponse `json:"user"`
}

type OTPConfig struct {
	Length                int `json:"length"`
	ExpiresInMinutes      int `json:"expires_in_minutes"`
	MaxAttempts           int `json:"max_attempts"`
	ResendCooldownSeconds int `json:"resend_cooldown_seconds"`
	MaxResends            int `json:"max_resends"`
}

type InitiateRegistrationResponse struct {
	RegistrationID string    `json:"registration_id"`
	Email          string    `json:"email"`
	Status         string    `json:"status"`
	Message        string    `json:"message"`
	ExpiresAt      time.Time `json:"expires_at"`
	OTPConfig      OTPConfig `json:"otp_config"`
}

type NextStep struct {
	Action         string   `json:"action"`
	Endpoint       string   `json:"endpoint"`
	RequiredFields []string `json:"required_fields"`
}

type VerifyRegistrationOTPResponse struct {
	RegistrationID    string    `json:"registration_id"`
	Status            string    `json:"status"`
	Message           string    `json:"message"`
	RegistrationToken string    `json:"registration_token"`
	TokenExpiresAt    time.Time `json:"token_expires_at"`
	NextStep          NextStep  `json:"next_step"`
}

type SetPasswordResponse struct {
	RegistrationID    string   `json:"registration_id"`
	Status            string   `json:"status"`
	Message           string   `json:"message"`
	RegistrationToken string   `json:"registration_token"`
	NextStep          NextStep `json:"next_step"`
}

type CompleteProfileRegistrationResponse struct {
	UserID       uuid.UUID               `json:"user_id"`
	Email        string                  `json:"email"`
	Status       string                  `json:"status"`
	Message      string                  `json:"message"`
	Profile      RegistrationUserProfile `json:"profile"`
	AccessToken  string                  `json:"access_token"`
	RefreshToken string                  `json:"refresh_token"`
	TokenType    string                  `json:"token_type"`
	ExpiresIn    int                     `json:"expires_in"`
}

type ResendRegistrationOTPResponse struct {
	RegistrationID        string    `json:"registration_id"`
	Message               string    `json:"message"`
	ExpiresAt             time.Time `json:"expires_at"`
	ResendsRemaining      int       `json:"resends_remaining"`
	NextResendAvailableAt time.Time `json:"next_resend_available_at"`
}

type RegistrationUserProfile struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type CompleteRegistrationResponse struct {
	UserID  uuid.UUID               `json:"user_id"`
	Email   string                  `json:"email"`
	Status  string                  `json:"status"`
	Message string                  `json:"message"`
	Profile RegistrationUserProfile `json:"profile"`

	AccessToken  *string `json:"access_token,omitempty"`
	RefreshToken *string `json:"refresh_token,omitempty"`
	TokenType    *string `json:"token_type,omitempty"`
	ExpiresIn    *int    `json:"expires_in,omitempty"`
}

type RegistrationStatusResponse struct {
	RegistrationID       string    `json:"registration_id"`
	Email                string    `json:"email"`
	Status               string    `json:"status"`
	ExpiresAt            time.Time `json:"expires_at"`
	OTPAttemptsRemaining int       `json:"otp_attempts_remaining"`
	ResendsRemaining     int       `json:"resends_remaining"`
}
