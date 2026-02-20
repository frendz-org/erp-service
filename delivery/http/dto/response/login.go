package response

import (
	"time"

	"github.com/google/uuid"
)

type LoginProductResponse struct {
	ProductID   uuid.UUID `json:"product_id"`
	ProductCode string    `json:"product_code"`
	Roles       []string  `json:"roles,omitempty"`
	Permissions []string  `json:"permissions,omitempty"`
}

type LoginTenantResponse struct {
	TenantID uuid.UUID              `json:"tenant_id"`
	Products []LoginProductResponse `json:"products,omitempty"`
}

type LoginUserResponse struct {
	ID       uuid.UUID             `json:"id"`
	Email    string                `json:"email"`
	FullName string                `json:"full_name"`
	Tenants  []LoginTenantResponse `json:"tenants,omitempty"`
}

type VerifyLoginOTPResponse struct {
	AccessToken  string            `json:"access_token"`
	RefreshToken string            `json:"refresh_token"`
	ExpiresIn    int               `json:"expires_in"`
	TokenType    string            `json:"token_type"`
	User         LoginUserResponse `json:"user"`
}

type UnifiedLoginResponse struct {
	Status string `json:"status"`

	LoginSessionID  *uuid.UUID `json:"login_session_id,omitempty"`
	Email           string     `json:"email,omitempty"`
	OTPExpiresAt    *time.Time `json:"otp_expires_at,omitempty"`
	AttemptsAllowed *int       `json:"attempts_allowed,omitempty"`
	ResendsAllowed  *int       `json:"resends_allowed,omitempty"`
	SessionExpires  *time.Time `json:"session_expires_at,omitempty"`

	AccessToken  string             `json:"access_token,omitempty"`
	RefreshToken string             `json:"refresh_token,omitempty"`
	ExpiresIn    int                `json:"expires_in,omitempty"`
	TokenType    string             `json:"token_type,omitempty"`
	User         *LoginUserResponse `json:"user,omitempty"`
}

type RefreshTokenResponse struct {
	AccessToken  string            `json:"access_token"`
	RefreshToken string            `json:"refresh_token"`
	ExpiresIn    int               `json:"expires_in"`
	TokenType    string            `json:"token_type"`
	User         LoginUserResponse `json:"user"`
}

type ResendLoginOTPResponse struct {
	Status           string    `json:"status"`
	LoginSessionID   uuid.UUID `json:"login_session_id"`
	Email            string    `json:"email"`
	OTPExpiresAt     time.Time `json:"otp_expires_at"`
	ResendsRemaining int       `json:"resends_remaining"`
	CooldownSeconds  int       `json:"cooldown_seconds"`
}

type LoginStatusResponse struct {
	Status            string    `json:"status"`
	LoginSessionID    uuid.UUID `json:"login_session_id"`
	Email             string    `json:"email"`
	AttemptsRemaining int       `json:"attempts_remaining"`
	ResendsRemaining  int       `json:"resends_remaining"`
	ExpiresAt         time.Time `json:"expires_at"`
	CooldownRemaining int       `json:"cooldown_remaining,omitempty"`
}
