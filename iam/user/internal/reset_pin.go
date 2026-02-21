package internal

import (
	"context"

	"erp-service/iam/user/userdto"
	"erp-service/pkg/errors"

	"github.com/google/uuid"
)

func (uc *usecase) ResetPIN(ctx context.Context, id uuid.UUID) (*userdto.ResetPINResponse, error) {
	_, err := uc.UserRepo.GetByID(ctx, id)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, errors.ErrUserNotFound()
		}
		return nil, err
	}

	securityState, err := uc.UserSecurityStateRepo.GetByUserID(ctx, id)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, errors.ErrInternal("user security state not found")
		}
		return nil, err
	}

	if !securityState.PINVerified {
		return nil, errors.ErrBadRequest("user does not have a PIN set")
	}

	securityState.PINVerified = false
	securityState.FailedPINAttempts = 0

	if err := uc.UserSecurityStateRepo.Update(ctx, securityState); err != nil {
		return nil, errors.ErrInternal("failed to reset PIN").WithError(err)
	}

	return &userdto.ResetPINResponse{
		UserID:  id,
		Message: "PIN reset successfully. User will need to set a new PIN.",
	}, nil
}
