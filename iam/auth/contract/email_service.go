package contract

import "context"

type EmailService interface {
	SendRegistrationOTP(ctx context.Context, email, otp string, expiryMinutes int) error
	SendLoginOTP(ctx context.Context, email, otp string, expiryMinutes int) error
	SendWelcome(ctx context.Context, email, firstName string) error
	SendPasswordReset(ctx context.Context, email, token string, expiryMinutes int) error
	SendPINReset(ctx context.Context, email, otp string, expiryMinutes int) error
	SendAdminInvitation(ctx context.Context, email, token string, expiryMinutes int) error
}
