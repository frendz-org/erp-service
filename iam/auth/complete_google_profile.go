package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"time"

	"erp-service/entity"
	"erp-service/masterdata"
	"erp-service/pkg/errors"
)

func (uc *usecase) CompleteGoogleProfile(
	ctx context.Context,
	req *CompleteGoogleProfileRequest,
) (*CompleteGoogleProfileResponse, error) {
	_, err := uc.validateGoogleRegistrationToken(req.RegistrationToken, req.RegistrationID)
	if err != nil {
		return nil, err
	}

	session, err := uc.InMemoryStore.GetAndDeleteGoogleRegistrationSession(ctx, req.RegistrationID)
	if err != nil {
		return nil, err
	}

	if session.IsExpired() {
		return nil, errors.New("GOOGLE_REGISTRATION_EXPIRED", "Google registration session has expired", http.StatusGone)
	}

	if !session.IsPendingProfile() {
		return nil, errors.ErrConflict("Google registration session is not in the correct state")
	}

	tokenHash := sha256.Sum256([]byte(req.RegistrationToken))
	tokenHashStr := hex.EncodeToString(tokenHash[:])
	if session.RegistrationTokenHash != tokenHashStr {
		return nil, errors.ErrUnauthorized("Google registration token has already been used or is invalid")
	}

	if err := uc.validateGoogleProfileFields(ctx, req); err != nil {
		return nil, err
	}

	firstName, lastName := splitFullName(req.FullName)

	dob, err := time.Parse("2006-01-02", req.DateOfBirth)
	if err != nil {
		return nil, errors.ErrValidation("Invalid date_of_birth format. Use YYYY-MM-DD")
	}

	age := calculateAge(dob)
	if age < 18 {
		return nil, errors.ErrValidation("You must be at least 18 years old to register")
	}

	gender := entity.Gender(req.Gender)

	now := time.Now()
	var user *entity.User

	err = uc.TxManager.WithTransaction(ctx, func(txCtx context.Context) error {

		emailExists, err := uc.UserRepo.EmailExists(txCtx, session.Email)
		if err != nil {
			return errors.ErrInternal("failed to check email").WithError(err)
		}
		if emailExists {
			return errors.ErrConflict("This email has already been registered")
		}

		user = &entity.User{
			Email:              session.Email,
			Status:             entity.UserStatusActive,
			StatusChangedAt:    &now,
			RegistrationSource: "GOOGLE",
		}
		if err := uc.UserRepo.Create(txCtx, user); err != nil {
			return err
		}

		profile := &entity.UserProfile{
			UserID:            user.ID,
			FirstName:         firstName,
			LastName:          lastName,
			DateOfBirth:       &dob,
			Gender:            &gender,
			ProfilePictureURL: nilIfEmpty(session.Picture),
			Metadata:          json.RawMessage("{}"),
			UpdatedAt:         now,
		}
		if err := uc.UserProfileRepo.Create(txCtx, profile); err != nil {
			return err
		}

		googleAuth := entity.NewGoogleAuthMethod(user.ID, entity.GoogleCredentialData{
			GoogleID:      session.GoogleID,
			Email:         session.Email,
			EmailVerified: session.GoogleEmailVerified,
			Name:          session.Name,
			Picture:       session.Picture,
		})
		if err := uc.UserAuthMethodRepo.Create(txCtx, googleAuth); err != nil {
			return err
		}

		securityState := &entity.UserSecurityState{
			UserID:          user.ID,
			EmailVerified:   true,
			EmailVerifiedAt: &now,
			UpdatedAt:       now,
		}
		if err := uc.UserSecurityStateRepo.Create(txCtx, securityState); err != nil {
			return err
		}

		return nil
	})

	if err != nil {

		if errors.IsAppError(err) {
			return nil, err
		}
		return nil, errors.ErrInternal("failed to create Google user").WithError(err)
	}

	accessToken, refreshToken, expiresIn, err := uc.generateAuthTokensForRegistration(
		ctx, user.ID, session.Email, req.IPAddress, req.UserAgent,
	)
	if err != nil {
		return nil, errors.ErrInternal("failed to generate auth tokens").WithError(err)
	}

	uc.sendEmailAsync(ctx, func(ctx context.Context) error {
		return uc.EmailService.SendWelcome(ctx, session.Email, firstName)
	})

	return &CompleteGoogleProfileResponse{
		UserID:  user.ID,
		Email:   session.Email,
		Status:  string(entity.UserStatusActive),
		Message: "Profile completed successfully. You are now logged in.",
		Profile: RegistrationUserProfile{
			FirstName: firstName,
			LastName:  lastName,
		},
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    expiresIn,
	}, nil
}
func (uc *usecase) validateGoogleProfileFields(ctx context.Context, req *CompleteGoogleProfileRequest) error {
	if !genderCodePattern.MatchString(req.Gender) {
		return errors.ErrValidation("gender must be in format GENDER_NNN (e.g. GENDER_001)")
	}

	resp, err := uc.MasterdataUsecase.ValidateItemCode(ctx, &masterdata.ValidateCodeRequest{
		CategoryCode: "GENDER",
		ItemCode:     req.Gender,
	})
	if err != nil {
		return errors.ErrInternal("failed to validate gender").WithError(err)
	}
	if !resp.Valid {
		return errors.ErrValidation("gender is not a valid value")
	}

	return nil
}
