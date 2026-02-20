package internal

import (
	"context"

	"iam-service/entity"
	"iam-service/iam/user/userdto"
	"iam-service/pkg/errors"

	"github.com/google/uuid"
)

func (uc *usecase) Update(ctx context.Context, callerTenantID *uuid.UUID, id uuid.UUID, req *userdto.UpdateRequest) (*userdto.UserDetailResponse, error) {
	user, err := uc.UserRepo.GetByID(ctx, id)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, errors.ErrUserNotFound()
		}
		return nil, err
	}

	profile, err := uc.UserProfileRepo.GetByUserID(ctx, id)
	if err != nil && !errors.IsNotFound(err) {
		return nil, err
	}

	userUpdated := false
	profileUpdated := false

	if req.Status != nil {
		user.Status = entity.UserStatus(*req.Status)
		userUpdated = true
	}

	if userUpdated {
		if err := uc.UserRepo.Update(ctx, user); err != nil {
			return nil, errors.ErrInternal("failed to update user").WithError(err)
		}
	}

	if profile != nil {
		if req.FirstName != nil {
			profile.FirstName = *req.FirstName
			profileUpdated = true
		}
		if req.LastName != nil {
			profile.LastName = *req.LastName
			profileUpdated = true
		}
		if req.Phone != nil {
			profile.PhoneNumber = req.Phone
			profileUpdated = true
		}
		if req.Address != nil {
			profile.Address = req.Address
			profileUpdated = true
		}

		if profileUpdated {
			if err := uc.UserProfileRepo.Update(ctx, profile); err != nil {
				return nil, errors.ErrInternal("failed to update user profile").WithError(err)
			}
		}
	}

	authMethod, err := uc.UserAuthMethodRepo.GetByUserID(ctx, id)
	if err != nil && !errors.IsNotFound(err) {
		return nil, err
	}

	securityState, err := uc.UserSecurityStateRepo.GetByUserID(ctx, id)
	if err != nil && !errors.IsNotFound(err) {
		return nil, err
	}

	return mapUserToDetailResponse(user, profile, authMethod, securityState), nil
}
