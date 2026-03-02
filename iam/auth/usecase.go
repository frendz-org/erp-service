package auth

import (
	"context"

	"github.com/google/uuid"
)

type SessionManager interface {
	Logout(ctx context.Context, req *LogoutRequest) error
	LogoutAll(ctx context.Context, req *LogoutAllRequest) error
	RefreshToken(ctx context.Context, req *RefreshTokenRequest) (*RefreshTokenResponse, error)
}

type RegistrationFlow interface {
	InitiateRegistration(ctx context.Context, req *InitiateRegistrationRequest) (*InitiateRegistrationResponse, error)
	VerifyRegistrationOTP(ctx context.Context, req *VerifyRegistrationOTPRequest) (*VerifyRegistrationOTPResponse, error)
	ResendRegistrationOTP(ctx context.Context, req *ResendRegistrationOTPRequest) (*ResendRegistrationOTPResponse, error)
	SetPassword(ctx context.Context, req *SetPasswordRequest) (*SetPasswordResponse, error)
	CompleteProfileRegistration(ctx context.Context, req *CompleteProfileRegistrationRequest) (*CompleteProfileRegistrationResponse, error)
	GetRegistrationStatus(ctx context.Context, registrationID uuid.UUID, email string) (*RegistrationStatusResponse, error)
}

type LoginFlow interface {
	InitiateLogin(ctx context.Context, req *InitiateLoginRequest) (*UnifiedLoginResponse, error)
	VerifyLoginOTP(ctx context.Context, req *VerifyLoginOTPRequest) (*VerifyLoginOTPResponse, error)
	ResendLoginOTP(ctx context.Context, req *ResendLoginOTPRequest) (*ResendLoginOTPResponse, error)
	GetLoginStatus(ctx context.Context, req *GetLoginStatusRequest) (*LoginStatusResponse, error)
}

type GoogleOAuthFlow interface {
	GetGoogleAuthURL(ctx context.Context) (*GoogleAuthURLResponse, error)
	HandleGoogleCallback(ctx context.Context, req *GoogleCallbackRequest) (*GoogleCallbackResponse, error)
}

type TransferTokenFlow interface {
	CreateTransferToken(ctx context.Context, req *CreateTransferTokenRequest) (*CreateTransferTokenResponse, error)
	ExchangeTransferToken(ctx context.Context, req *ExchangeTransferTokenRequest) (*ExchangeTransferTokenResponse, error)
	LogoutTree(ctx context.Context, req *LogoutTreeRequest) error
}

type Usecase interface {
	SessionManager
	RegistrationFlow
	LoginFlow
	GoogleOAuthFlow
	TransferTokenFlow
}
