package auth

import (
	"context"
	"time"

	"erp-service/entity"
	"erp-service/pkg/errors"
	jwtpkg "erp-service/pkg/jwt"

	"github.com/google/uuid"
)

func (uc *usecase) generateAuthTokensForRegistration(ctx context.Context, userID uuid.UUID, email, ipAddress, userAgent string) (string, string, int, error) {
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
		IPAddress:   ipAddress,
		UserAgent:   userAgent,
		ExpiresAt:   time.Now().Add(uc.Config.JWT.RefreshExpiry),
		CreatedAt:   time.Now(),
	}

	if err := uc.RefreshTokenRepo.Create(ctx, refreshTokenEntity); err != nil {
		return "", "", 0, errors.ErrInternal("failed to create refresh token").WithError(err)
	}

	expiresIn := int(uc.Config.JWT.AccessExpiry.Seconds())

	return accessToken, refreshToken, expiresIn, nil
}
