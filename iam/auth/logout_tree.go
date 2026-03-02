package auth

import (
	"context"
	"fmt"
	"time"

	"erp-service/entity"
	"erp-service/pkg/errors"
	"erp-service/pkg/logger"

	"github.com/google/uuid"
)

func (uc *usecase) LogoutTree(ctx context.Context, req *LogoutTreeRequest) error {
	tokenHash := hashToken(req.RefreshToken)
	refreshToken, err := uc.RefreshTokenRepo.GetByTokenHash(ctx, tokenHash)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return errors.ErrInternal("failed to look up refresh token").WithError(err)
	}

	if refreshToken.UserID != req.UserID {
		return nil
	}

	if refreshToken.RevokedAt != nil || refreshToken.IsExpired() {
		return nil
	}

	session, err := uc.UserSessionRepo.GetByRefreshTokenID(ctx, refreshToken.ID)
	if err != nil {
		return nil
	}

	rootSessionID := uc.findRootSession(ctx, session)

	descendantIDs, err := uc.UserSessionRepo.GetDescendantSessionIDs(ctx, rootSessionID)
	if err != nil {
		return errors.ErrInternal("failed to get session tree").WithError(err)
	}

	var refreshTokenIDs []uuid.UUID
	for _, sid := range descendantIDs {
		s, sErr := uc.UserSessionRepo.GetByID(ctx, sid)
		if sErr != nil || s.RefreshTokenID == nil {
			continue
		}
		refreshTokenIDs = append(refreshTokenIDs, *s.RefreshTokenID)
	}

	if err := uc.TxManager.WithTransaction(ctx, func(txCtx context.Context) error {
		if err := uc.UserSessionRepo.RevokeByIDs(txCtx, descendantIDs); err != nil {
			return err
		}
		if err := uc.RefreshTokenRepo.RevokeByIDs(txCtx, refreshTokenIDs, "Session tree logout"); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return errors.ErrInternal("failed to revoke session tree").WithError(err)
	}

	if req.AccessTokenJTI != "" {
		ttl := time.Until(req.AccessTokenExp)
		if ttl > 0 {
			_ = uc.InMemoryStore.BlacklistToken(context.WithoutCancel(ctx), req.AccessTokenJTI, ttl)
		}
	}

	uc.AuditLogger.Log(ctx, logger.AuditEvent{
		Domain:  "auth",
		Action:  "session_tree_logout",
		ActorID: req.UserID.String(),
		Success: true,
		Reason:  fmt.Sprintf("revoked session tree with %d sessions", len(descendantIDs)),
	})

	return nil
}

func (uc *usecase) findRootSession(ctx context.Context, session *entity.UserSession) uuid.UUID {
	current := session
	for depth := 0; depth < TransferTokenMaxTreeDepth; depth++ {
		if current.ParentSessionID == nil {
			return current.ID
		}
		parent, err := uc.UserSessionRepo.GetByID(ctx, *current.ParentSessionID)
		if err != nil {
			return current.ID
		}
		current = parent
	}
	return current.ID
}
