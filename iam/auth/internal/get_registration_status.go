package internal

import (
	"context"
	"strings"

	"erp-service/iam/auth/authdto"
	"erp-service/pkg/errors"

	"github.com/google/uuid"
)

func (uc *usecase) GetRegistrationStatus(
	ctx context.Context,
	registrationID uuid.UUID,
	email string,
) (*authdto.RegistrationStatusResponse, error) {
	session, err := uc.InMemoryStore.GetRegistrationSession(ctx, registrationID)
	if err != nil {
		return nil, err
	}

	if !strings.EqualFold(session.Email, email) {
		return nil, errors.ErrNotFound("Registration session not found")
	}

	maskedEmail := maskEmail(session.Email)

	return &authdto.RegistrationStatusResponse{
		RegistrationID:       registrationID.String(),
		Email:                maskedEmail,
		Status:               string(session.Status),
		ExpiresAt:            session.ExpiresAt,
		OTPAttemptsRemaining: session.RemainingAttempts(),
		ResendsRemaining:     session.RemainingResends(),
	}, nil
}
