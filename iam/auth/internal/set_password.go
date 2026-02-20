package internal

import (
	"context"
	"net/http"

	"iam-service/iam/auth/authdto"
	"iam-service/pkg/errors"

	"golang.org/x/crypto/bcrypt"
)

func (uc *usecase) SetPassword(
	ctx context.Context,
	req *authdto.SetPasswordRequest,
) (*authdto.SetPasswordResponse, error) {
	_, err := uc.validateRegistrationCompleteToken(req.RegistrationToken, req.RegistrationID)
	if err != nil {
		return nil, err
	}

	session, err := uc.InMemoryStore.GetRegistrationSession(ctx, req.RegistrationID)
	if err != nil {
		return nil, err
	}

	if session.IsExpired() {
		return nil, errors.New("REGISTRATION_EXPIRED", "Registration session has expired", http.StatusGone)
	}

	if !session.CanSetPassword() {
		return nil, errors.ErrForbidden("Email has not been verified")
	}

	if req.Password != req.ConfirmationPassword {
		return nil, errors.ErrValidation("Passwords do not match")
	}

	if err := uc.validatePassword(req.Password); err != nil {
		return nil, errors.ErrValidation(err.Error())
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.ErrInternal("failed to hash password").WithError(err)
	}
	passwordHashStr := string(passwordHash)

	newToken, tokenHash, err := uc.generateRegistrationCompleteToken(req.RegistrationID, session.Email)
	if err != nil {
		return nil, errors.ErrInternal("failed to generate new registration token").WithError(err)
	}

	if err := uc.InMemoryStore.MarkRegistrationPasswordSet(ctx, req.RegistrationID, passwordHashStr, tokenHash); err != nil {
		return nil, errors.ErrInternal("failed to update registration session").WithError(err)
	}

	response := &authdto.SetPasswordResponse{
		RegistrationID:    req.RegistrationID.String(),
		Status:            "PASSWORD_SET",
		Message:           "Password set successfully. Please proceed to complete your profile.",
		RegistrationToken: newToken,
		NextStep: authdto.NextStep{
			Action:         "set-profile",
			Endpoint:       "/api/v1/auth/registration/complete-profile",
			RequiredFields: []string{"full_name", "gender", "date_of_birth"},
		},
	}

	return response, nil
}
