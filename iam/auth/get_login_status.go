package auth

import (
	"context"
	"net/http"
	"strings"

	"erp-service/pkg/errors"
)

func (uc *usecase) GetLoginStatus(
	ctx context.Context,
	req *GetLoginStatusRequest,
) (*LoginStatusResponse, error) {
	session, err := uc.InMemoryStore.GetLoginSession(ctx, req.LoginSessionID)
	if err != nil {
		return nil, errors.New("SESSION_NOT_FOUND", "Login session not found or expired", http.StatusNotFound)
	}

	if !strings.EqualFold(session.Email, req.Email) {
		return nil, errors.New("SESSION_MISMATCH", "Email does not match login session", http.StatusBadRequest)
	}

	return &LoginStatusResponse{
		Status:            string(session.Status),
		LoginSessionID:    req.LoginSessionID,
		Email:             MaskEmail(session.Email),
		AttemptsRemaining: session.RemainingAttempts(),
		ResendsRemaining:  session.RemainingResends(),
		ExpiresAt:         session.ExpiresAt,
		CooldownRemaining: session.CooldownRemainingSeconds(),
	}, nil
}
