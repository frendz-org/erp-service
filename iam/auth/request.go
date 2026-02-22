package auth

import (
	"time"

	"github.com/google/uuid"
)

type LogoutRequest struct {
	RefreshToken   string    `json:"refresh_token" validate:"required"`
	AccessTokenJTI string    `json:"-"`
	AccessTokenExp time.Time `json:"-"`
	UserID         uuid.UUID `json:"-"`
	IPAddress      string    `json:"-"`
	UserAgent      string    `json:"-"`
}

type LogoutAllRequest struct {
	UserID    uuid.UUID `json:"-"`
	IPAddress string    `json:"-"`
	UserAgent string    `json:"-"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
	IPAddress    string `json:"-"`
	UserAgent    string `json:"-"`
}

type InitiateRegistrationRequest struct {
	Email     string `json:"email" validate:"required,email,max=255"`
	IPAddress string `json:"-"`
	UserAgent string `json:"-"`
}

type VerifyRegistrationOTPRequest struct {
	RegistrationID uuid.UUID `json:"-"`
	Email          string    `json:"email" validate:"required,email"`
	OTPCode        string    `json:"otp_code" validate:"required,len=6,numeric"`
}

type ResendRegistrationOTPRequest struct {
	RegistrationID uuid.UUID `json:"-"`
	Email          string    `json:"email" validate:"required,email"`
}

type SetPasswordRequest struct {
	RegistrationID       uuid.UUID `json:"-"`
	RegistrationToken    string    `json:"-"`
	Password             string    `json:"password" validate:"required,min=8,max=128"`
	ConfirmationPassword string    `json:"confirmation_password" validate:"required,eqfield=Password"`
}

type CompleteProfileRegistrationRequest struct {
	RegistrationID    uuid.UUID `json:"-"`
	RegistrationToken string    `json:"-"`
	IPAddress         string    `json:"-"`
	UserAgent         string    `json:"-"`
	FullName          string    `json:"full_name"     validate:"required,min=1,max=200"`
	Gender            string    `json:"gender"        validate:"required,max=20"`
	DateOfBirth       string    `json:"date_of_birth" validate:"required,datetime=2006-01-02"`
}

