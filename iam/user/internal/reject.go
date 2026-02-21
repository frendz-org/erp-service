package internal

import (
	"context"
	"time"

	"erp-service/entity"
	"erp-service/iam/user/userdto"
	"erp-service/pkg/errors"

	"github.com/google/uuid"
)

func (uc *usecase) Reject(ctx context.Context, id uuid.UUID, approverID uuid.UUID, req *userdto.RejectRequest) (*userdto.RejectResponse, error) {
	user, err := uc.UserRepo.GetByID(ctx, id)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, errors.ErrUserNotFound()
		}
		return nil, err
	}

	if user.IsActive() {
		return nil, errors.ErrBadRequest("cannot reject an active user")
	}

	now := time.Now()
	user.Status = entity.UserStatusInactive
	user.StatusChangedAt = &now
	user.StatusChangedBy = &approverID

	if err := uc.UserRepo.Update(ctx, user); err != nil {
		return nil, errors.ErrInternal("failed to reject user").WithError(err)
	}

	return &userdto.RejectResponse{
		UserID:  id,
		Message: "User rejected successfully",
	}, nil
}
