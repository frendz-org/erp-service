package internal

import (
	"context"
	"time"

	"iam-service/entity"
	"iam-service/iam/auth/authdto"
	"iam-service/pkg/errors"

	"github.com/google/uuid"
)

func (uc *usecase) InitiateRegistration(
	ctx context.Context,
	req *authdto.InitiateRegistrationRequest,
) (*authdto.InitiateRegistrationResponse, error) {
	rateLimitTTL := time.Duration(RegistrationRateLimitWindow) * time.Minute
	count, err := uc.InMemoryStore.IncrementRegistrationRateLimit(ctx, req.Email, rateLimitTTL)
	if err != nil {
		return nil, err
	}
	if count > int64(RegistrationRateLimitPerHour) {
		return nil, errors.ErrTooManyRequests("Too many registration attempts. Please try again later.")
	}

	emailExists, err := uc.UserRepo.EmailExists(ctx, req.Email)
	if err != nil {
		return nil, errors.ErrInternal("failed to check email").WithError(err)
	}
	if emailExists {
		return &authdto.InitiateRegistrationResponse{
			RegistrationID: uuid.New().String(),
			Email:          req.Email,
			Status:         string(entity.RegistrationSessionStatusPendingVerification),
			Message:        "Verification code sent to your email",
			ExpiresAt:      time.Now().Add(time.Duration(RegistrationSessionExpiryMinutes) * time.Minute),
			OTPConfig: authdto.OTPConfig{
				Length:                RegistrationOTPLength,
				ExpiresInMinutes:      RegistrationOTPExpiryMinutes,
				MaxAttempts:           RegistrationOTPMaxAttempts,
				ResendCooldownSeconds: RegistrationOTPResendCooldown,
				MaxResends:            RegistrationOTPMaxResends,
			},
		}, nil
	}

	emailLocked, err := uc.InMemoryStore.IsRegistrationEmailLocked(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if emailLocked {
		return nil, errors.ErrConflict("An active registration already exists for this email")
	}

	otp, otpHash, err := uc.generateOTP()
	if err != nil {
		return nil, errors.ErrInternal("failed to generate OTP").WithError(err)
	}

	now := time.Now()
	sessionID := uuid.New()
	sessionTTL := time.Duration(RegistrationSessionExpiryMinutes) * time.Minute
	otpExpiry := now.Add(time.Duration(RegistrationOTPExpiryMinutes) * time.Minute)

	session := &entity.RegistrationSession{
		ID:                    sessionID,
		Email:                 req.Email,
		Status:                entity.RegistrationSessionStatusPendingVerification,
		OTPHash:               otpHash,
		OTPCreatedAt:          now,
		OTPExpiresAt:          otpExpiry,
		Attempts:              0,
		MaxAttempts:           RegistrationOTPMaxAttempts,
		ResendCount:           0,
		MaxResends:            RegistrationOTPMaxResends,
		ResendCooldownSeconds: RegistrationOTPResendCooldown,
		IPAddress:             req.IPAddress,
		UserAgent:             req.UserAgent,
		CreatedAt:             now,
		ExpiresAt:             now.Add(sessionTTL),
	}

	locked, err := uc.InMemoryStore.LockRegistrationEmail(ctx, req.Email, sessionTTL)
	if err != nil {
		return nil, errors.ErrInternal("failed to lock email").WithError(err)
	}
	if !locked {
		return nil, errors.ErrConflict("An active registration already exists for this email")
	}

	if err := uc.InMemoryStore.CreateRegistrationSession(ctx, session, sessionTTL); err != nil {

		_ = uc.InMemoryStore.UnlockRegistrationEmail(ctx, req.Email)
		return nil, err
	}

	uc.sendEmailAsync(ctx, func(ctx context.Context) error {
		return uc.EmailService.SendRegistrationOTP(ctx, req.Email, otp, RegistrationOTPExpiryMinutes)
	})

	return &authdto.InitiateRegistrationResponse{
		RegistrationID: sessionID.String(),
		Email:          req.Email,
		Status:         string(entity.RegistrationSessionStatusPendingVerification),
		Message:        "Verification code sent to your email",
		ExpiresAt:      session.ExpiresAt,
		OTPConfig: authdto.OTPConfig{
			Length:                RegistrationOTPLength,
			ExpiresInMinutes:      RegistrationOTPExpiryMinutes,
			MaxAttempts:           RegistrationOTPMaxAttempts,
			ResendCooldownSeconds: RegistrationOTPResendCooldown,
			MaxResends:            RegistrationOTPMaxResends,
		},
	}, nil
}
