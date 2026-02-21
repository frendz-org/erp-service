package auth

import (
	"context"
	"net/http"
	"time"

	"erp-service/entity"
	"erp-service/pkg/errors"
	jwtpkg "erp-service/pkg/jwt"

	"github.com/google/uuid"
)

func (uc *usecase) RefreshToken(
	ctx context.Context,
	req *RefreshTokenRequest,
) (*RefreshTokenResponse, error) {
	tokenConfig, err := uc.buildTokenConfig()
	if err != nil {
		return nil, err
	}

	claims, err := jwtpkg.ParseRefreshToken(req.RefreshToken, tokenConfig)
	if err != nil {
		return nil, errors.New("INVALID_TOKEN", "Invalid or expired refresh token", http.StatusUnauthorized)
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return nil, errors.New("INVALID_TOKEN", "Invalid or expired refresh token", http.StatusUnauthorized)
	}

	tokenHash := hashToken(req.RefreshToken)
	oldToken, err := uc.RefreshTokenRepo.GetByTokenHash(ctx, tokenHash)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, errors.New("INVALID_TOKEN", "Invalid or expired refresh token", http.StatusUnauthorized)
		}
		return nil, errors.ErrInternal("failed to verify token").WithError(err)
	}

	if oldToken.UserID != userID {
		return nil, errors.New("INVALID_TOKEN", "Invalid or expired refresh token", http.StatusUnauthorized)
	}

	if oldToken.IsRevoked() {
		_ = uc.RefreshTokenRepo.RevokeByFamily(ctx, oldToken.TokenFamily, "Token reuse detected")
		return nil, errors.New("TOKEN_REUSED", "Invalid or expired refresh token", http.StatusUnauthorized)
	}

	if oldToken.IsExpired() {
		return nil, errors.New("INVALID_TOKEN", "Invalid or expired refresh token", http.StatusUnauthorized)
	}

	blacklisted, err := uc.InMemoryStore.GetUserBlacklistTimestamp(ctx, userID)
	if err == nil && blacklisted != nil && blacklisted.After(oldToken.CreatedAt) {
		return nil, errors.New("USER_BLACKLISTED", "Invalid or expired refresh token", http.StatusUnauthorized)
	}

	user, err := uc.UserRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.New("INVALID_TOKEN", "Invalid or expired refresh token", http.StatusUnauthorized)
	}
	if user.Status != entity.UserStatusActive {
		return nil, errors.New("USER_INACTIVE", "Invalid or expired refresh token", http.StatusForbidden)
	}

	tenantClaims, userTenants, platformRoles, err := uc.buildMultiTenantClaims(ctx, userID)
	if err != nil {
		return nil, errors.ErrInternal("failed to build tenant claims").WithError(err)
	}

	sessionID, err := uuid.Parse(claims.ID)
	if err != nil {
		sessionID = uuid.New()
	}

	accessToken, err := jwtpkg.GenerateMultiTenantAccessToken(
		userID,
		user.Email,
		platformRoles,
		tenantClaims,
		sessionID,
		tokenConfig,
	)
	if err != nil {
		return nil, errors.ErrInternal("failed to generate access token").WithError(err)
	}

	newRefreshTokenStr, err := jwtpkg.GenerateRefreshToken(userID, sessionID, tokenConfig)
	if err != nil {
		return nil, errors.ErrInternal("failed to generate refresh token").WithError(err)
	}

	newTokenHash := hashToken(newRefreshTokenStr)
	newRefreshToken := &entity.RefreshToken{
		UserID:      userID,
		TokenHash:   newTokenHash,
		TokenFamily: oldToken.TokenFamily,
		ExpiresAt:   time.Now().Add(uc.Config.JWT.RefreshExpiry),
		IPAddress:   req.IPAddress,
		UserAgent:   req.UserAgent,
		CreatedAt:   time.Now(),
	}

	session, _ := uc.UserSessionRepo.GetByRefreshTokenID(ctx, oldToken.ID)

	if err := uc.TxManager.WithTransaction(ctx, func(txCtx context.Context) error {
		if err := uc.RefreshTokenRepo.Create(txCtx, newRefreshToken); err != nil {
			return err
		}
		if err := uc.RefreshTokenRepo.Revoke(txCtx, oldToken.ID, "Token rotation"); err != nil {
			return err
		}
		if err := uc.RefreshTokenRepo.SetReplacedBy(txCtx, oldToken.ID, newRefreshToken.ID); err != nil {
			return err
		}
		if session != nil && session.IsActive() {
			if err := uc.UserSessionRepo.UpdateRefreshTokenID(txCtx, session.ID, newRefreshToken.ID); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return nil, errors.ErrInternal("failed to rotate token").WithError(err)
	}

	profile, _ := uc.UserProfileRepo.GetByUserID(ctx, userID)
	fullName := ""
	if profile != nil {
		fullName = profile.FirstName
		if profile.LastName != "" {
			fullName += " " + profile.LastName
		}
	}

	return &RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshTokenStr,
		ExpiresIn:    int(uc.Config.JWT.AccessExpiry.Seconds()),
		TokenType:    "Bearer",
		User: LoginUserResponse{
			ID:       userID,
			Email:    user.Email,
			FullName: fullName,
			Tenants:  userTenants,
		},
	}, nil
}
