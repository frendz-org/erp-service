package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"time"

	"erp-service/entity"
	"erp-service/pkg/errors"
	jwtpkg "erp-service/pkg/jwt"
	"erp-service/pkg/logger"

	"github.com/google/uuid"
)

type transferTokenData struct {
	UserID          uuid.UUID `json:"user_id"`
	SourceSessionID uuid.UUID `json:"source_session_id"`
	ProductCode     string    `json:"product_code"`
	CreatedAt       time.Time `json:"created_at"`
}

func (uc *usecase) CreateTransferToken(ctx context.Context, req *CreateTransferTokenRequest) (*CreateTransferTokenResponse, error) {
	count, err := uc.InMemoryStore.IncrementTransferTokenRateLimit(
		ctx, req.UserID,
		time.Duration(TransferTokenRateLimitWindow)*time.Minute,
	)
	if err != nil {
		return nil, errors.ErrInternal("failed to check rate limit").WithError(err)
	}
	if count > TransferTokenRateLimitPerMin {
		return nil, errors.New("RATE_LIMIT_EXCEEDED", "Too many transfer token requests. Please try again later.", http.StatusTooManyRequests)
	}

	hasAccess, err := uc.validateProductAccess(ctx, req.UserID, req.ProductCode)
	if err != nil {
		return nil, errors.ErrInternal("failed to validate product access").WithError(err)
	}
	if !hasAccess {
		return nil, errors.New("NO_PRODUCT_ACCESS", "User does not have active registration for the target product", http.StatusForbidden)
	}

	session, err := uc.UserSessionRepo.GetByID(ctx, req.SessionID)
	if err != nil || !session.IsActive() {
		return nil, errors.New("SESSION_INVALID", "Source session is no longer active", http.StatusUnauthorized)
	}

	codeBytes := make([]byte, TransferTokenCodeBytes)
	if _, err := rand.Read(codeBytes); err != nil {
		return nil, errors.ErrInternal("failed to generate transfer code").WithError(err)
	}
	code := hex.EncodeToString(codeBytes)

	data := transferTokenData{
		UserID:          req.UserID,
		SourceSessionID: req.SessionID,
		ProductCode:     req.ProductCode,
		CreatedAt:       time.Now(),
	}
	dataJSON, err := json.Marshal(data)
	if err != nil {
		return nil, errors.ErrInternal("failed to marshal transfer data").WithError(err)
	}

	ttl := time.Duration(TransferTokenTTLSeconds) * time.Second
	if err := uc.InMemoryStore.StoreTransferToken(ctx, code, dataJSON, ttl); err != nil {
		return nil, errors.ErrInternal("failed to store transfer token").WithError(err)
	}

	uc.AuditLogger.Log(ctx, logger.AuditEvent{
		Domain:  "auth",
		Action:  "transfer_token_created",
		ActorID: req.UserID.String(),
		Success: true,
		Reason:  "transfer token created for product " + req.ProductCode,
	})

	return &CreateTransferTokenResponse{
		Code:      code,
		ExpiresIn: TransferTokenTTLSeconds,
	}, nil
}

func (uc *usecase) ExchangeTransferToken(ctx context.Context, req *ExchangeTransferTokenRequest) (*ExchangeTransferTokenResponse, error) {
	dataJSON, err := uc.InMemoryStore.GetAndDeleteTransferToken(ctx, req.Code)
	if err != nil {
		return nil, errors.ErrInternal("failed to retrieve transfer token").WithError(err)
	}
	if dataJSON == nil {
		return nil, errors.New("INVALID_CODE", "Invalid or expired transfer code", http.StatusUnauthorized)
	}

	var data transferTokenData
	if err := json.Unmarshal(dataJSON, &data); err != nil {
		return nil, errors.ErrInternal("failed to parse transfer data").WithError(err)
	}

	parentSession, err := uc.UserSessionRepo.GetByID(ctx, data.SourceSessionID)
	if err != nil || !parentSession.IsActive() {
		return nil, errors.New("SESSION_REVOKED", "Source session has been revoked", http.StatusUnauthorized)
	}

	user, err := uc.UserRepo.GetByID(ctx, data.UserID)
	if err != nil {
		return nil, errors.New("USER_NOT_FOUND", "User not found", http.StatusUnauthorized)
	}
	if user.Status != entity.UserStatusActive {
		return nil, errors.New("USER_INACTIVE", "User account is not active", http.StatusUnauthorized)
	}

	blacklistTS, _ := uc.InMemoryStore.GetUserBlacklistTimestamp(ctx, data.UserID)
	if blacklistTS != nil && data.CreatedAt.Before(*blacklistTS) {
		return nil, errors.New("USER_BLACKLISTED", "Invalid or expired transfer code", http.StatusUnauthorized)
	}

	tenantClaims, userTenants, platformRoles, err := uc.buildMultiTenantClaims(ctx, data.UserID)
	if err != nil {
		return nil, errors.ErrInternal("failed to build tenant claims").WithError(err)
	}

	sessionID := uuid.New()
	tokenFamily := uuid.New()

	tokenConfig, err := uc.buildTokenConfig()
	if err != nil {
		return nil, err
	}

	accessToken, err := jwtpkg.GenerateMultiTenantAccessToken(
		data.UserID,
		user.Email,
		platformRoles,
		tenantClaims,
		sessionID,
		tokenConfig,
	)
	if err != nil {
		return nil, errors.ErrInternal("failed to generate access token").WithError(err)
	}

	refreshToken, err := jwtpkg.GenerateRefreshToken(data.UserID, sessionID, tokenConfig)
	if err != nil {
		return nil, errors.ErrInternal("failed to generate refresh token").WithError(err)
	}

	refreshTokenHash := hashToken(refreshToken)
	now := time.Now()
	refreshTokenEntity := &entity.RefreshToken{
		UserID:      data.UserID,
		TokenHash:   refreshTokenHash,
		TokenFamily: tokenFamily,
		ExpiresAt:   now.Add(uc.Config.JWT.RefreshExpiry),
		IPAddress:   req.IPAddress,
		UserAgent:   req.UserAgent,
		CreatedAt:   now,
	}

	userSession := &entity.UserSession{
		UserID:          data.UserID,
		ParentSessionID: &data.SourceSessionID,
		IPAddress:       req.IPAddress,
		LoginMethod:     entity.UserSessionLoginMethodTransferToken,
		Status:          entity.UserSessionStatusActive,
		LastActiveAt:    now,
		ExpiresAt:       now.Add(uc.Config.JWT.RefreshExpiry),
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	if req.UserAgent != "" {
		userSession.UserAgent = &req.UserAgent
	}
	if req.DeviceFingerprint != nil {
		userSession.DeviceFingerprint = req.DeviceFingerprint
	}

	if err := uc.TxManager.WithTransaction(ctx, func(txCtx context.Context) error {
		if err := uc.RefreshTokenRepo.Create(txCtx, refreshTokenEntity); err != nil {
			return err
		}
		userSession.RefreshTokenID = &refreshTokenEntity.ID
		if err := uc.UserSessionRepo.Create(txCtx, userSession); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, errors.ErrInternal("failed to create transfer session").WithError(err)
	}

	profile, _ := uc.UserProfileRepo.GetByUserID(ctx, data.UserID)
	fullName := ""
	if profile != nil {
		fullName = profile.FirstName
		if profile.LastName != "" {
			fullName += " " + profile.LastName
		}
	}

	uc.AuditLogger.Log(ctx, logger.AuditEvent{
		Domain:  "auth",
		Action:  "transfer_token_exchanged",
		ActorID: data.UserID.String(),
		Success: true,
		Reason:  "session created via transfer token for product " + data.ProductCode,
	})

	return &ExchangeTransferTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int(uc.Config.JWT.AccessExpiry.Seconds()),
		TokenType:    "Bearer",
		User: LoginUserResponse{
			ID:       data.UserID,
			Email:    user.Email,
			FullName: fullName,
			Tenants:  userTenants,
		},
	}, nil
}

func (uc *usecase) validateProductAccess(ctx context.Context, userID uuid.UUID, productCode string) (bool, error) {
	registrations, err := uc.UserTenantRegRepo.ListActiveByUserID(ctx, userID)
	if err != nil {
		return false, err
	}

	tenantIDs := make(map[uuid.UUID]bool)
	for _, reg := range registrations {
		tenantIDs[reg.TenantID] = true
	}

	for tenantID := range tenantIDs {
		products, err := uc.ProductsByTenantRepo.ListActiveByTenantID(ctx, tenantID)
		if err != nil {
			return false, err
		}
		for _, p := range products {
			if p.Code == productCode {
				for _, reg := range registrations {
					if reg.ProductID != nil && *reg.ProductID == p.ID {
						return true, nil
					}
				}
			}
		}
	}

	return false, nil
}
