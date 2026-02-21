package internal

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"

	"erp-service/entity"
	"erp-service/iam/auth/authdto"
	"erp-service/pkg/errors"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (uc *usecase) VerifyRegistrationOTP(
	ctx context.Context,
	req *authdto.VerifyRegistrationOTPRequest,
) (*authdto.VerifyRegistrationOTPResponse, error) {
	session, err := uc.InMemoryStore.GetRegistrationSession(ctx, req.RegistrationID)
	if err != nil {
		return nil, err
	}

	if !strings.EqualFold(session.Email, req.Email) {
		return nil, errors.ErrValidation("Email does not match registration")
	}

	if session.IsExpired() {
		return nil, errors.New("REGISTRATION_EXPIRED", "Registration session has expired", http.StatusGone)
	}

	if session.Status == entity.RegistrationSessionStatusVerified {
		return nil, errors.ErrConflict("Registration is already verified")
	}

	if session.Status == entity.RegistrationSessionStatusFailed {
		return nil, errors.ErrTooManyRequests("Too many failed attempts. Please start a new registration.")
	}

	if session.Status != entity.RegistrationSessionStatusPendingVerification {
		return nil, errors.ErrBadRequest("Registration is not in a verifiable state")
	}

	if session.IsOTPExpired() {
		return nil, errors.New("OTP_EXPIRED", "Verification code has expired. Please request a new one.", http.StatusGone)
	}

	if !session.CanAttemptOTP() {
		return nil, errors.ErrTooManyRequests("Too many failed attempts. Please start a new registration.")
	}

	err = bcrypt.CompareHashAndPassword([]byte(session.OTPHash), []byte(req.OTPCode))
	if err != nil {

		attempts, incErr := uc.InMemoryStore.IncrementRegistrationAttempts(ctx, req.RegistrationID)
		if incErr != nil {
			return nil, incErr
		}

		remaining := session.MaxAttempts - attempts
		if remaining <= 0 {
			return nil, errors.ErrTooManyRequests("Too many failed attempts. Registration has been invalidated.")
		}

		return nil, errors.ErrUnauthorized("The verification code is incorrect").
			WithDetails(map[string]interface{}{
				"attempts_remaining": remaining,
			})
	}

	token, tokenHash, err := uc.generateRegistrationCompleteToken(req.RegistrationID, req.Email)
	if err != nil {
		return nil, errors.ErrInternal("failed to generate registration token").WithError(err)
	}

	if err := uc.InMemoryStore.MarkRegistrationVerified(ctx, req.RegistrationID, tokenHash); err != nil {
		return nil, err
	}

	tokenExpiry := time.Now().Add(time.Duration(RegistrationCompleteTokenExpiryMinutes) * time.Minute)

	return &authdto.VerifyRegistrationOTPResponse{
		RegistrationID:    req.RegistrationID.String(),
		Status:            string(entity.RegistrationSessionStatusVerified),
		Message:           "Email verified successfully",
		RegistrationToken: token,
		TokenExpiresAt:    tokenExpiry,
		NextStep: authdto.NextStep{
			Action:         "set_password",
			Endpoint:       fmt.Sprintf("/api/iam/v1/registrations/%s/set-password", req.RegistrationID.String()),
			RequiredFields: []string{"password", "confirmation_password"},
		},
	}, nil
}

func (uc *usecase) generateRegistrationCompleteToken(registrationID uuid.UUID, email string) (string, string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"registration_id": registrationID.String(),
		"email":           email,
		"purpose":         RegistrationCompleteTokenPurpose,
		"exp":             now.Add(time.Duration(RegistrationCompleteTokenExpiryMinutes) * time.Minute).Unix(),
		"iat":             now.Unix(),
		"jti":             uuid.New().String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(uc.registrationSigningSecret()))
	if err != nil {
		return "", "", err
	}

	hash := sha256.Sum256([]byte(tokenString))
	tokenHash := hex.EncodeToString(hash[:])

	return tokenString, tokenHash, nil
}

func (uc *usecase) validateRegistrationCompleteToken(tokenString string, expectedRegistrationID uuid.UUID) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.ErrTokenInvalid()
		}
		return []byte(uc.registrationSigningSecret()), nil
	})

	if err != nil {
		return nil, errors.ErrUnauthorized("Registration token is invalid or expired")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.ErrUnauthorized("Registration token is invalid")
	}

	if purpose, ok := claims["purpose"].(string); !ok || purpose != RegistrationCompleteTokenPurpose {
		return nil, errors.ErrUnauthorized("Token is not a registration completion token")
	}

	if regID, ok := claims["registration_id"].(string); !ok || regID != expectedRegistrationID.String() {
		return nil, errors.ErrUnauthorized("Token does not match this registration")
	}

	return claims, nil
}
