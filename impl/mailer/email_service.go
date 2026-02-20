package mailer

import (
	"context"
	"fmt"
	"iam-service/config"
	"log"
	"strings"

	"gopkg.in/gomail.v2"
)

const (
	ProviderConsole = "console"
	ProviderSMTP    = "smtp"
)

type EmailService struct {
	config *config.EmailConfig
	dialer *gomail.Dialer
}

func NewEmailService(cfg *config.EmailConfig) *EmailService {
	svc := &EmailService{config: cfg}

	if cfg.Provider == ProviderSMTP {
		svc.dialer = gomail.NewDialer(
			cfg.SMTPHost,
			cfg.SMTPPort,
			cfg.SMTPUser,
			cfg.SMTPPass,
		)
	}

	return svc
}

func (s *EmailService) SendRegistrationOTP(ctx context.Context, email, otp string, expiryMinutes int) error {
	subject := "Kode Verifikasi Registrasi - Frendz"

	htmlBody, err := renderRegistrationOTPEmail(otp, expiryMinutes)
	if err != nil {
		return fmt.Errorf("failed to render registration OTP email: %w", err)
	}

	return s.send(ctx, email, subject, htmlBody)
}

func (s *EmailService) SendLoginOTP(ctx context.Context, email, otp string, expiryMinutes int) error {
	subject := "Kode Verifikasi Login - Frendz"

	htmlBody, err := renderLoginOTPEmail(otp, expiryMinutes)
	if err != nil {
		return fmt.Errorf("failed to render login OTP email: %w", err)
	}

	return s.send(ctx, email, subject, htmlBody)
}

func (s *EmailService) SendWelcome(ctx context.Context, email, firstName string) error {
	subject := "Selamat Datang di Frendz"

	htmlBody, err := renderWelcomeEmail(firstName)
	if err != nil {
		return fmt.Errorf("failed to render welcome email: %w", err)
	}

	return s.send(ctx, email, subject, htmlBody)
}

func (s *EmailService) SendPasswordReset(ctx context.Context, email, token string, expiryMinutes int) error {
	subject := "Reset Password - Frendz"

	htmlBody, err := renderPasswordResetEmail(token, expiryMinutes)
	if err != nil {
		return fmt.Errorf("failed to render password reset email: %w", err)
	}

	return s.send(ctx, email, subject, htmlBody)
}

func (s *EmailService) SendPINReset(ctx context.Context, email, otp string, expiryMinutes int) error {
	subject := "Reset PIN - Frendz"

	htmlBody, err := renderPINResetEmail(otp, expiryMinutes)
	if err != nil {
		return fmt.Errorf("failed to render PIN reset email: %w", err)
	}

	return s.send(ctx, email, subject, htmlBody)
}

func (s *EmailService) SendAdminInvitation(ctx context.Context, email, token string, expiryMinutes int) error {
	subject := "Undangan Administrator - Frendz"

	htmlBody, err := renderAdminInvitationEmail(token, expiryMinutes)
	if err != nil {
		return fmt.Errorf("failed to render admin invitation email: %w", err)
	}

	return s.send(ctx, email, subject, htmlBody)
}

func (s *EmailService) send(ctx context.Context, to, subject, htmlBody string) error {
	if s.config.Provider == ProviderConsole {
		return s.sendConsole(to, subject, htmlBody)
	}
	return s.sendSMTP(ctx, to, subject, htmlBody)
}

func (s *EmailService) sendConsole(to, subject, htmlBody string) error {
	maskedTo := maskEmail(to)
	log.Printf(`
========================================
EMAIL (Console Mode)
========================================
To: %s
From: %s <%s>
Subject: %s
Content-Type: text/html

[HTML Content - %d bytes]
========================================
`, maskedTo, s.config.FromName, s.config.FromAddress, subject, len(htmlBody))
	return nil
}

func (s *EmailService) sendSMTP(ctx context.Context, to, subject, htmlBody string) error {
	if s.dialer == nil {
		return fmt.Errorf("SMTP dialer not configured")
	}

	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress(s.config.FromAddress, s.config.FromName))
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", htmlBody)

	if err := s.dialer.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email via SMTP: %w", err)
	}

	maskedTo := maskEmail(to)
	log.Printf("[Email] Sent to %s: %s", maskedTo, subject)

	return nil
}

func maskEmail(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return "***"
	}

	local := parts[0]
	domain := parts[1]

	if len(local) <= 2 {
		return local[:1] + "***@" + domain
	}

	return local[:2] + "***@" + domain
}
