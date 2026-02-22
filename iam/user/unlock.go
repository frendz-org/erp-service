package user

import (
	"context"

	"erp-service/pkg/errors"

	"github.com/google/uuid"
)

func (uc *usecase) Unlock(ctx context.Context, id uuid.UUID) (*UnlockResponse, error) {
	_, err := uc.UserRepo.GetByID(ctx, id)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, errors.ErrUserNotFound()
		}
		return nil, err
	}

	security, err := uc.UserSecurityStateRepo.GetByUserID(ctx, id)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, errors.ErrInternal("user security state not found")
		}
		return nil, err
	}

	if security.LockedUntil == nil {
		return nil, errors.ErrBadRequest("user is not locked")
	}

	security.LockedUntil = nil
	security.FailedLoginAttempts = 0

	if err := uc.UserSecurityStateRepo.Update(ctx, security); err != nil {
		return nil, errors.ErrInternal("failed to unlock user").WithError(err)
	}

	return &UnlockResponse{
		UserID:  id,
		Message: "User unlocked successfully",
	}, nil
}
