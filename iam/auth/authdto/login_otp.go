package authdto

import (
	"time"

	"github.com/google/uuid"
)

type InitiateLoginRequest struct {
	Email     string `json:"email" validate:"required,email"`
	IPAddress string `json:"-"`
	UserAgent string `json:"-"`
}

type VerifyLoginOTPRequest struct {
	LoginSessionID uuid.UUID `json:"-"`
	Email          string    `json:"email" validate:"required,email"`
	OTPCode        string    `json:"otp_code" validate:"required,len=6,numeric"`
	IPAddress      string    `json:"-"`
	UserAgent      string    `json:"-"`
}

type ResendLoginOTPRequest struct {
	LoginSessionID uuid.UUID `json:"-"`
	Email          string    `json:"email" validate:"required,email"`
}

type GetLoginStatusRequest struct {
	LoginSessionID uuid.UUID `json:"-"`
	Email          string    `json:"email" validate:"required,email"`
}

type LoginOTPRequiredResponse struct {
	Status          string    `json:"status"`
	LoginSessionID  uuid.UUID `json:"login_session_id"`
	Email           string    `json:"email"`
	ExpiresAt       time.Time `json:"expires_at"`
	OTPExpiresAt    time.Time `json:"otp_expires_at"`
	AttemptsAllowed int       `json:"attempts_allowed"`
	ResendsAllowed  int       `json:"resends_allowed"`
}

type VerifyLoginOTPResponse struct {
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

type LoginResultType string

const (
	LoginResultSuccess     LoginResultType = "SUCCESS"
	LoginResultOTPRequired LoginResultType = "OTP_REQUIRED"
)

type UnifiedLoginResponse struct {
	Status LoginResultType `json:"status"`

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

func NewOTPRequiredResponse(sessionID uuid.UUID, email string, sessionExpires, otpExpires time.Time, maxAttempts, maxResends int) *UnifiedLoginResponse {
	return &UnifiedLoginResponse{
		Status:          LoginResultOTPRequired,
		LoginSessionID:  &sessionID,
		Email:           email,
		OTPExpiresAt:    &otpExpires,
		SessionExpires:  &sessionExpires,
		AttemptsAllowed: &maxAttempts,
		ResendsAllowed:  &maxResends,
	}
}

func NewLoginSuccessResponse(accessToken, refreshToken string, expiresIn int, user LoginUserResponse) *UnifiedLoginResponse {
	return &UnifiedLoginResponse{
		Status:       LoginResultSuccess,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
		TokenType:    "Bearer",
		User:         &user,
	}
}
