package auth

import (
	"context"
	"time"

	"erp-service/pkg/errors"
)

func (uc *usecase) Logout(ctx context.Context, req *LogoutRequest) error {
	if req.RefreshToken == "" {
		return errors.ErrBadRequest("refresh token is required")
	}

	tokenHash := hashToken(req.RefreshToken)

	refreshToken, err := uc.RefreshTokenRepo.GetByTokenHash(ctx, tokenHash)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return errors.ErrInternal("failed to verify token").WithError(err)
	}

	if refreshToken.UserID != req.UserID {
		return nil
	}

	if refreshToken.RevokedAt != nil {
		return nil
	}

	if refreshToken.IsExpired() {
		return nil
	}

	session, _ := uc.UserSessionRepo.GetByRefreshTokenID(ctx, refreshToken.ID)

	err = uc.TxManager.WithTransaction(ctx, func(txCtx context.Context) error {
		if err := uc.RefreshTokenRepo.Revoke(txCtx, refreshToken.ID, "User logout"); err != nil {
			return err
		}
		if session != nil && session.IsActive() {
			if err := uc.UserSessionRepo.Revoke(txCtx, session.ID); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return errors.ErrInternal("failed to revoke session").WithError(err)
	}

	if req.AccessTokenJTI != "" {
		ttl := time.Until(req.AccessTokenExp)
		if ttl > 0 {
			_ = uc.InMemoryStore.BlacklistToken(context.WithoutCancel(ctx), req.AccessTokenJTI, ttl)
		}
	}

	return nil
}
