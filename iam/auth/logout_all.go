package auth

import (
	"context"
	"fmt"
	"time"

	"erp-service/pkg/errors"
)

func (uc *usecase) LogoutAll(ctx context.Context, req *LogoutAllRequest) error {
	err := uc.TxManager.WithTransaction(ctx, func(txCtx context.Context) error {
		if err := uc.RefreshTokenRepo.RevokeAllByUserID(txCtx, req.UserID, "User logout all"); err != nil {
			return fmt.Errorf("revoke all refresh tokens: %w", err)
		}
		if err := uc.UserSessionRepo.RevokeAllByUserID(txCtx, req.UserID); err != nil {
			return fmt.Errorf("revoke all sessions: %w", err)
		}
		return nil
	})
	if err != nil {
		return errors.ErrInternal("failed to revoke all sessions").WithError(err)
	}

	ttl := uc.Config.JWT.AccessExpiry
	if ttl <= 0 {
		ttl = 15 * time.Minute
	}

	_ = uc.InMemoryStore.BlacklistUser(context.WithoutCancel(ctx), req.UserID, time.Now(), ttl)

	return nil
}
