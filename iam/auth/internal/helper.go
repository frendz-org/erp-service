package internal

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"math/big"
	"strings"
	"time"

	"erp-service/pkg/errors"
	"erp-service/pkg/logger"

	"golang.org/x/crypto/bcrypt"
)

var emailSendConcurrency = make(chan struct{}, 50)

func (uc *usecase) sendEmailAsync(ctx context.Context, fn func(ctx context.Context) error) {
	bgCtx := context.WithoutCancel(ctx)
	select {
	case emailSendConcurrency <- struct{}{}:
		go func() {
			defer func() { <-emailSendConcurrency }()
			sendCtx, cancel := context.WithTimeout(bgCtx, 30*time.Second)
			defer cancel()
			if err := fn(sendCtx); err != nil {
				uc.AuditLogger.Log(sendCtx, logger.AuditEvent{
					Domain:  "auth",
					Action:  "email_send_failed",
					Success: false,
					Reason:  err.Error(),
				})
			}
		}()
	default:
		uc.AuditLogger.Log(bgCtx, logger.AuditEvent{
			Domain:  "auth",
			Action:  "email_send_shed",
			Success: false,
			Reason:  "email concurrency limit reached",
		})
	}
}

func (uc *usecase) validatePassword(password string) error {
	if len(password) < PasswordMinLength {
		return errors.ErrValidation("Password must be at least 8 characters long")
	}

	hasUpper := false
	hasLower := false
	hasNumber := false
	hasSpecial := false

	for _, char := range password {
		switch {
		case char >= 'A' && char <= 'Z':
			hasUpper = true
		case char >= 'a' && char <= 'z':
			hasLower = true
		case char >= '0' && char <= '9':
			hasNumber = true
		case char >= '!' && char <= '/' || char >= ':' && char <= '@' || char >= '[' && char <= '`' || char >= '{' && char <= '~':
			hasSpecial = true
		}
	}

	if !hasUpper {
		return errors.ErrValidation("Password must contain at least one uppercase letter")
	}
	if !hasLower {
		return errors.ErrValidation("Password must contain at least one lowercase letter")
	}
	if !hasNumber {
		return errors.ErrValidation("Password must contain at least one number")
	}
	if !hasSpecial {
		return errors.ErrValidation("Password must contain at least one special character")
	}

	return nil
}

func (uc *usecase) generateOTP() (otp string, otpHash string, err error) {
	digits := make([]byte, OTPLength)
	for i := 0; i < OTPLength; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", "", err
		}
		digits[i] = byte('0' + n.Int64())
	}
	otp = string(digits)

	hash, err := bcrypt.GenerateFromPassword([]byte(otp), bcrypt.DefaultCost)
	if err != nil {
		return "", "", err
	}

	return otp, string(hash), nil
}

func (uc *usecase) registrationSigningSecret() string {
	if uc.Config.JWT.RegistrationSecret != "" {
		return uc.Config.JWT.RegistrationSecret
	}
	return uc.Config.JWT.AccessSecret
}

func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

func maskEmail(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return "***"
	}

	local := parts[0]
	domain := parts[1]

	if len(local) == 0 {
		return "***@" + domain
	}

	if len(local) == 1 {
		return local + "***@" + domain
	}

	return string(local[0]) + "***@" + domain
}
