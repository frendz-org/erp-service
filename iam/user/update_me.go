package user

import (
	"context"

	"erp-service/pkg/errors"

	"github.com/google/uuid"
)

func (uc *usecase) UpdateMe(ctx context.Context, userID uuid.UUID, req *UpdateMeRequest) (*UserDetailResponse, error) {
	user, err := uc.UserRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, errors.ErrUserNotFound()
		}
		return nil, err
	}

	profile, err := uc.UserProfileRepo.GetByUserID(ctx, userID)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, errors.ErrInternal("user profile not found")
		}
		return nil, err
	}

	if req.FirstName != nil {
		profile.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		profile.LastName = *req.LastName
	}
	if req.PhoneNumber != nil {
		profile.PhoneNumber = req.PhoneNumber
	}
	if req.Address != nil {
		profile.Address = req.Address
	}

	if err := uc.UserProfileRepo.Update(ctx, profile); err != nil {
		return nil, errors.ErrInternal("failed to update user profile").WithError(err)
	}

	authMethod, err := uc.UserAuthMethodRepo.GetByUserID(ctx, userID)
	if err != nil && !errors.IsNotFound(err) {
		return nil, err
	}

	securityState, err := uc.UserSecurityStateRepo.GetByUserID(ctx, userID)
	if err != nil && !errors.IsNotFound(err) {
		return nil, err
	}

	return mapUserToDetailResponse(user, profile, authMethod, securityState), nil
}
