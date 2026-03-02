package auth

import (
	"context"
	"net/http"
	"strings"
	"time"

	"erp-service/entity"
	"erp-service/pkg/errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (uc *usecase) InitiateLogin(
	ctx context.Context,
	req *InitiateLoginRequest,
) (*UnifiedLoginResponse, error) {
	start := time.Now()
	defer func() {
		elapsed := time.Since(start)
		if elapsed < 500*time.Millisecond {
			time.Sleep(500*time.Millisecond - elapsed)
		}
	}()

	email := strings.ToLower(strings.TrimSpace(req.Email))

	count, err := uc.InMemoryStore.IncrementLoginRateLimit(ctx, email, time.Duration(LoginRateLimitWindow)*time.Minute)
	if err != nil {
		return nil, errors.ErrInternal("failed to check rate limit").WithError(err)
	}
	if count > int64(LoginRateLimitPerHour) {
		return nil, errors.New("RATE_LIMITED", "Too many login attempts. Please try again later.", http.StatusTooManyRequests)
	}

	user, err := uc.UserRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, errors.New("INVALID_CREDENTIALS", "Invalid email or password.", http.StatusUnauthorized)
		}
		return nil, errors.ErrInternal("failed to look up user").WithError(err)
	}

	if !user.IsActive() {
		return nil, errors.New("INVALID_CREDENTIALS", "Invalid email or password.", http.StatusUnauthorized)
	}

	authMethod, err := uc.UserAuthMethodRepo.GetByUserID(ctx, user.ID)
	if err != nil || authMethod == nil {
		return dummyOTPResponse(email), nil
	}
	passwordHash := authMethod.GetPasswordHash()
	if passwordHash == "" {
		return dummyOTPResponse(email), nil
	}
	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password)); err != nil {
		return nil, errors.New("INVALID_CREDENTIALS", "Invalid email or password.", http.StatusUnauthorized)
	}

	uc.upgradePasswordHashIfNeeded(ctx, authMethod, req.Password)

	otp, otpHash, err := uc.generateOTP()
	if err != nil {
		return nil, errors.ErrInternal("failed to generate OTP").WithError(err)
	}

	now := time.Now()
	sessionExpiry := time.Duration(LoginSessionExpiryMinutes) * time.Minute
	otpExpiry := time.Duration(LoginOTPExpiryMinutes) * time.Minute

	session := &entity.LoginSession{
		ID:                    uuid.New(),
		UserID:                user.ID,
		Email:                 email,
		Status:                entity.LoginSessionStatusPendingVerification,
		OTPHash:               otpHash,
		OTPCreatedAt:          now,
		OTPExpiresAt:          now.Add(otpExpiry),
		Attempts:              0,
		MaxAttempts:           LoginOTPMaxAttempts,
		ResendCount:           0,
		MaxResends:            LoginOTPMaxResends,
		ResendCooldownSeconds: LoginOTPResendCooldown,
		IPAddress:             req.IPAddress,
		UserAgent:             req.UserAgent,
		CreatedAt:             now,
		ExpiresAt:             now.Add(sessionExpiry),
	}

	if err := uc.InMemoryStore.CreateLoginSession(ctx, session, sessionExpiry); err != nil {
		return nil, errors.ErrInternal("failed to create login session").WithError(err)
	}

	uc.sendEmailAsync(ctx, func(ctx context.Context) error {
		return uc.EmailService.SendLoginOTP(ctx, email, otp, LoginOTPExpiryMinutes)
	})

	return NewOTPRequiredResponse(
		session.ID,
		MaskEmail(email),
		session.ExpiresAt,
		session.OTPExpiresAt,
		LoginOTPMaxAttempts,
		LoginOTPMaxResends,
	), nil
}

func dummyOTPResponse(email string) *UnifiedLoginResponse {
	now := time.Now()
	sessionExpires := now.Add(time.Duration(LoginSessionExpiryMinutes) * time.Minute)
	otpExpires := now.Add(time.Duration(LoginOTPExpiryMinutes) * time.Minute)
	attempts := LoginOTPMaxAttempts
	resends := LoginOTPMaxResends
	fakeID := uuid.New()
	return NewOTPRequiredResponse(
		fakeID,
		MaskEmail(email),
		sessionExpires,
		otpExpires,
		attempts,
		resends,
	)
}
