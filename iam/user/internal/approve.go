package internal

import (
	"context"
	"time"

	"iam-service/entity"
	"iam-service/iam/user/userdto"
	"iam-service/pkg/errors"

	"github.com/google/uuid"
)

func (uc *usecase) Approve(ctx context.Context, id uuid.UUID, approverID uuid.UUID) (*userdto.ApproveResponse, error) {
	user, err := uc.UserRepo.GetByID(ctx, id)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, errors.ErrUserNotFound()
		}
		return nil, err
	}

	if user.Status != entity.UserStatusPendingVerification {
		return nil, errors.ErrBadRequest("user is not pending verification")
	}

	now := time.Now()
	user.Status = entity.UserStatusActive
	user.StatusChangedAt = &now
	user.StatusChangedBy = &approverID

	if err := uc.UserRepo.Update(ctx, user); err != nil {
		return nil, errors.ErrInternal("failed to approve user").WithError(err)
	}

	return &userdto.ApproveResponse{
		UserID:  id,
		Message: "User approved successfully",
	}, nil
}
