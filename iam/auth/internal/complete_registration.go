package internal

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"time"

	"erp-service/entity"
	"erp-service/iam/auth/authdto"
	"erp-service/pkg/errors"
	jwtpkg "erp-service/pkg/jwt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (uc *usecase) CompleteRegistration(
	ctx context.Context,
	req *authdto.CompleteRegistrationRequest,
) (*authdto.CompleteRegistrationResponse, error) {
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

	if session.Status != entity.RegistrationSessionStatusVerified {
		return nil, errors.ErrForbidden("Email has not been verified")
	}

	tokenHash := sha256.Sum256([]byte(req.RegistrationToken))
	tokenHashStr := hex.EncodeToString(tokenHash[:])
	if session.RegistrationTokenHash == nil || *session.RegistrationTokenHash != tokenHashStr {
		return nil, errors.ErrUnauthorized("Registration token has already been used or is invalid")
	}

	if err := uc.validatePassword(req.Password); err != nil {
		return nil, errors.ErrValidation(err.Error())
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.ErrInternal("failed to hash password").WithError(err)
	}
	passwordHashStr := string(passwordHash)

	emailExists, err := uc.UserRepo.EmailExists(ctx, session.Email)
	if err != nil {
		return nil, errors.ErrInternal("failed to check email").WithError(err)
	}
	if emailExists {
		return nil, errors.ErrConflict("This email has already been registered")
	}

	requiresApproval := false
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
			UserID:    user.ID,
			FirstName: req.FirstName,
			LastName:  req.LastName,
			UpdatedAt: now,
		}
		if req.PhoneNumber != nil {
			profile.PhoneNumber = req.PhoneNumber
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

	response := &authdto.CompleteRegistrationResponse{
		UserID: user.ID,
		Email:  session.Email,
		Status: string(entity.UserStatusActive),
		Profile: authdto.RegistrationUserProfile{
			FirstName: req.FirstName,
			LastName:  req.LastName,
		},
	}

	if !requiresApproval {
		response.Message = "Registration completed successfully. You are now logged in."

		accessToken, refreshToken, expiresIn, err := uc.generateAuthTokensForRegistration(ctx, user.ID, session.Email)
		if err != nil {
			response.Message = "Registration completed successfully. Please login to continue."
		} else {
			tokenType := "Bearer"
			response.AccessToken = &accessToken
			response.RefreshToken = &refreshToken
			response.TokenType = &tokenType
			response.ExpiresIn = &expiresIn
		}

		uc.sendEmailAsync(ctx, func(ctx context.Context) error {
			return uc.EmailService.SendWelcome(ctx, session.Email, req.FirstName)
		})
	} else {
		response.Message = "Registration submitted. Your account is pending administrator approval. You will receive an email when your account is activated."
	}

	return response, nil
}

func (uc *usecase) generateAuthTokensForRegistration(ctx context.Context, userID uuid.UUID, email string) (string, string, int, error) {
	sessionID := uuid.New()
	tokenFamily := uuid.New()

	tokenConfig := &jwtpkg.TokenConfig{
		SigningMethod: uc.Config.JWT.SigningMethod,
		AccessSecret:  uc.Config.JWT.AccessSecret,
		RefreshSecret: uc.Config.JWT.RefreshSecret,
		AccessExpiry:  uc.Config.JWT.AccessExpiry,
		RefreshExpiry: uc.Config.JWT.RefreshExpiry,
		Issuer:        uc.Config.JWT.Issuer,
		Audience:      uc.Config.JWT.Audience,
	}

	accessToken, err := jwtpkg.GenerateAccessToken(
		userID,
		email,
		nil,
		nil,
		[]string{},
		[]string{},
		nil,
		sessionID,
		tokenConfig,
	)
	if err != nil {
		return "", "", 0, errors.ErrInternal("failed to generate access token").WithError(err)
	}

	refreshToken, err := jwtpkg.GenerateRefreshToken(userID, sessionID, tokenConfig)
	if err != nil {
		return "", "", 0, errors.ErrInternal("failed to generate refresh token").WithError(err)
	}

	refreshTokenHash := hashToken(refreshToken)
	refreshTokenEntity := &entity.RefreshToken{
		UserID:      userID,
		TokenHash:   refreshTokenHash,
		TokenFamily: tokenFamily,
		ExpiresAt:   time.Now().Add(uc.Config.JWT.RefreshExpiry),
		CreatedAt:   time.Now(),
	}

	if err := uc.RefreshTokenRepo.Create(ctx, refreshTokenEntity); err != nil {
		return "", "", 0, errors.ErrInternal("failed to create refresh token").WithError(err)
	}

	expiresIn := int(uc.Config.JWT.AccessExpiry.Seconds())

	return accessToken, refreshToken, expiresIn, nil
}
