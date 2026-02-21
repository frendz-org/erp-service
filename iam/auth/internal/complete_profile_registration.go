package internal

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"regexp"
	"strings"
	"time"

	"erp-service/entity"
	"erp-service/iam/auth/authdto"
	"erp-service/masterdata/masterdatadto"
	"erp-service/pkg/errors"
)

var genderCodePattern = regexp.MustCompile(`^GENDER_\d{3}$`)

func (uc *usecase) CompleteProfileRegistration(
	ctx context.Context,
	req *authdto.CompleteProfileRegistrationRequest,
) (*authdto.CompleteProfileRegistrationResponse, error) {
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

	if !session.IsPasswordSet() {
		return nil, errors.ErrForbidden("Registration session is not ready for profile completion. Password must be set first.")
	}

	tokenHash := sha256.Sum256([]byte(req.RegistrationToken))
	tokenHashStr := hex.EncodeToString(tokenHash[:])
	if session.RegistrationTokenHash == nil || *session.RegistrationTokenHash != tokenHashStr {
		return nil, errors.ErrUnauthorized("Registration token has already been used or is invalid")
	}

	if err := uc.validateProfileFields(ctx, req); err != nil {
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

	emailExists, err := uc.UserRepo.EmailExists(ctx, session.Email)
	if err != nil {
		return nil, errors.ErrInternal("failed to check email").WithError(err)
	}
	if emailExists {
		return nil, errors.ErrConflict("This email has already been registered")
	}

	passwordHashStr, err := uc.InMemoryStore.GetRegistrationPasswordHash(ctx, req.RegistrationID)
	if err != nil {
		return nil, errors.ErrForbidden("Password has not been set")
	}

	now := time.Now()
	var user *entity.User

	err = uc.TxManager.WithTransaction(ctx, func(txCtx context.Context) error {
		user = &entity.User{
			Email:              session.Email,
			Status:             entity.UserStatusActive,
			StatusChangedAt:    &now,
			RegistrationSource: "SELF",
		}
		if err := uc.UserRepo.Create(txCtx, user); err != nil {
			return err
		}

		authMethod := entity.NewPasswordAuthMethod(user.ID, passwordHashStr)
		if err := uc.UserAuthMethodRepo.Create(txCtx, authMethod); err != nil {
			return err
		}

		profile := &entity.UserProfile{
			UserID:      user.ID,
			FirstName:   firstName,
			LastName:    lastName,
			DateOfBirth: &dob,
			Gender:      &gender,
			UpdatedAt:   now,
		}
		if err := uc.UserProfileRepo.Create(txCtx, profile); err != nil {
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
		return nil, errors.ErrInternal("failed to create user").WithError(err)
	}

	_ = uc.InMemoryStore.DeleteRegistrationSession(ctx, req.RegistrationID)
	_ = uc.InMemoryStore.UnlockRegistrationEmail(ctx, session.Email)

	accessToken, refreshToken, expiresIn, err := uc.generateAuthTokensForRegistration(ctx, user.ID, session.Email)
	if err != nil {
		return nil, errors.ErrInternal("failed to generate auth tokens").WithError(err)
	}

	uc.sendEmailAsync(ctx, func(ctx context.Context) error {
		return uc.EmailService.SendWelcome(ctx, session.Email, firstName)
	})

	response := &authdto.CompleteProfileRegistrationResponse{
		UserID:  user.ID,
		Email:   session.Email,
		Status:  string(entity.UserStatusActive),
		Message: "Registration completed successfully. You are now logged in.",
		Profile: authdto.RegistrationUserProfile{
			FirstName: firstName,
			LastName:  lastName,
		},
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    expiresIn,
	}

	return response, nil
}

func (uc *usecase) validateProfileFields(ctx context.Context, req *authdto.CompleteProfileRegistrationRequest) error {
	if strings.TrimSpace(req.FullName) == "" {
		return errors.ErrValidation("full_name is required")
	}

	if req.DateOfBirth == "" {
		return errors.ErrValidation("date_of_birth is required")
	}

	if !genderCodePattern.MatchString(req.Gender) {
		return errors.ErrValidation("gender must be in format GENDER_NNN (e.g. GENDER_001)")
	}

	resp, err := uc.MasterdataUsecase.ValidateItemCode(ctx, &masterdatadto.ValidateCodeRequest{
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

func splitFullName(fullName string) (firstName, lastName string) {
	trimmed := strings.TrimSpace(fullName)
	lastSpaceIndex := strings.LastIndex(trimmed, " ")

	if lastSpaceIndex == -1 {
		return trimmed, ""
	}

	return strings.TrimSpace(trimmed[:lastSpaceIndex]), strings.TrimSpace(trimmed[lastSpaceIndex+1:])
}

func calculateAge(birthDate time.Time) int {
	return calculateAgeAt(birthDate, time.Now())
}

func calculateAgeAt(birthDate, now time.Time) int {
	age := now.Year() - birthDate.Year()
	birthdayThisYear := time.Date(now.Year(), birthDate.Month(), birthDate.Day(), 0, 0, 0, 0, now.Location())
	if now.Before(birthdayThisYear) {
		age--
	}
	return age
}
