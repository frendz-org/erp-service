package user

import (
	"context"

	"erp-service/pkg/errors"

	"github.com/google/uuid"
)

func (uc *usecase) GetMe(ctx context.Context, userID uuid.UUID) (*UserDetailResponse, error) {
	user, err := uc.UserRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, errors.ErrUserNotFound()
		}
		return nil, err
	}

	profile, err := uc.UserProfileRepo.GetByUserID(ctx, userID)
	if err != nil && !errors.IsNotFound(err) {
		return nil, err
	}

	authMethod, err := uc.UserAuthMethodRepo.GetByUserID(ctx, userID)
	if err != nil && !errors.IsNotFound(err) {
		return nil, err
	}

	securityState, err := uc.UserSecurityStateRepo.GetByUserID(ctx, userID)
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
