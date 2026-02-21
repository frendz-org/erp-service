package internal

import (
	"context"
	"net/http"
	"strings"
	"time"

	"erp-service/iam/auth/authdto"
	"erp-service/pkg/errors"
)

func (uc *usecase) ResendLoginOTP(
	ctx context.Context,
	req *authdto.ResendLoginOTPRequest,
) (*authdto.ResendLoginOTPResponse, error) {
	session, err := uc.InMemoryStore.GetLoginSession(ctx, req.LoginSessionID)
	if err != nil {
		return nil, errors.New("SESSION_NOT_FOUND", "Login session not found or expired", http.StatusNotFound)
	}

	if !strings.EqualFold(session.Email, req.Email) {
		return nil, errors.New("SESSION_MISMATCH", "Email does not match login session", http.StatusBadRequest)
	}

	if session.IsExpired() {
		return nil, errors.New("SESSION_EXPIRED", "Login session has expired. Please start a new login.", http.StatusGone)
	}

	if !session.CanResend() {
		return nil, errors.New("RESEND_LIMIT", "Maximum resend attempts reached", http.StatusTooManyRequests)
	}

	if session.IsInCooldown() {
		remaining := session.CooldownRemainingSeconds()
		return nil, errors.New("RESEND_COOLDOWN", "Please wait before requesting a new OTP", http.StatusTooManyRequests).
			WithDetails(map[string]interface{}{"cooldown_remaining": remaining})
	}

	otp, otpHash, err := uc.generateOTP()
	if err != nil {
		return nil, errors.ErrInternal("failed to generate OTP").WithError(err)
	}

	otpExpiry := time.Duration(LoginOTPExpiryMinutes) * time.Minute
	newOTPExpiresAt := time.Now().Add(otpExpiry)

	if err := uc.InMemoryStore.UpdateLoginOTP(ctx, req.LoginSessionID, otpHash, newOTPExpiresAt); err != nil {
		return nil, errors.ErrInternal("failed to update OTP").WithError(err)
	}

	uc.sendEmailAsync(ctx, func(ctx context.Context) error {
		return uc.EmailService.SendLoginOTP(ctx, session.Email, otp, LoginOTPExpiryMinutes)
	})

	updatedSession, err := uc.InMemoryStore.GetLoginSession(ctx, req.LoginSessionID)
	if err != nil {
		return nil, errors.ErrInternal("failed to get updated session").WithError(err)
	}

	return &authdto.ResendLoginOTPResponse{
		Status:           "OTP_RESENT",
		LoginSessionID:   req.LoginSessionID,
		Email:            maskEmail(session.Email),
		OTPExpiresAt:     newOTPExpiresAt,
		ResendsRemaining: updatedSession.RemainingResends(),
		CooldownSeconds:  LoginOTPResendCooldown,
	}, nil
}
