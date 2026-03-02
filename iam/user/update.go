package user

import (
	"context"

	"erp-service/entity"
	"erp-service/pkg/errors"

	"github.com/google/uuid"
)

func (uc *usecase) Update(ctx context.Context, callerTenantID *uuid.UUID, id uuid.UUID, req *UpdateRequest) (*UserDetailResponse, error) {
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

	var gender *GenderResponse
	if profile != nil && profile.Gender != nil {
		item, err := uc.MasterdataUsecase.GetItemByCode(ctx, "GENDER", nil, string(*profile.Gender))
		if err == nil {
			gender = &GenderResponse{Code: item.Code, Name: item.Name}
		} else {
			gender = &GenderResponse{Code: string(*profile.Gender)}
		}
	}

	return mapUserToDetailResponse(user, profile, authMethod, securityState, gender), nil
}
