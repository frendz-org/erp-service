package auth

import (
	"context"
	"net/http"
	"strings"
	"time"

	"erp-service/entity"
	"erp-service/pkg/errors"
	jwtpkg "erp-service/pkg/jwt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (uc *usecase) VerifyLoginOTP(
	ctx context.Context,
	req *VerifyLoginOTPRequest,
) (*VerifyLoginOTPResponse, error) {
	session, err := uc.InMemoryStore.GetLoginSession(ctx, req.LoginSessionID)
	if err != nil {
		return nil, errors.New("SESSION_NOT_FOUND", "Login session not found or expired", http.StatusNotFound)
	}

	if !strings.EqualFold(session.Email, req.Email) {
		return nil, errors.New("SESSION_MISMATCH", "Email does not match login session", http.StatusBadRequest)
	}

	if !session.CanAttemptOTP() {
		if session.IsExpired() {
			return nil, errors.New("SESSION_EXPIRED", "Login session has expired. Please start a new login.", http.StatusGone)
		}
		if session.IsOTPExpired() {
			return nil, errors.New("OTP_EXPIRED", "OTP has expired. Please request a new one.", http.StatusGone)
		}
		if session.IsLocked() {
			return nil, errors.New("SESSION_LOCKED", "Too many failed attempts. Please start a new login.", http.StatusForbidden)
		}
		return nil, errors.New("OTP_INVALID", "Unable to verify OTP", http.StatusBadRequest)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(session.OTPHash), []byte(req.OTPCode)); err != nil {
		_, _ = uc.InMemoryStore.IncrementLoginAttempts(ctx, req.LoginSessionID)
		remaining := session.RemainingAttempts() - 1
		if remaining <= 0 {
			return nil, errors.New("SESSION_LOCKED", "Too many failed attempts. Please start a new login.", http.StatusForbidden)
		}
		return nil, errors.New("OTP_INVALID", "Invalid OTP code", http.StatusBadRequest)
	}

	if err := uc.InMemoryStore.MarkLoginVerified(ctx, req.LoginSessionID); err != nil {
		return nil, errors.ErrInternal("failed to mark session verified").WithError(err)
	}

	tenantClaims, userTenants, platformRoles, err := uc.buildMultiTenantClaims(ctx, session.UserID)
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
		session.UserID,
		session.Email,
		platformRoles,
		tenantClaims,
		sessionID,
		tokenConfig,
	)
	if err != nil {
		return nil, errors.ErrInternal("failed to generate access token").WithError(err)
	}

	refreshToken, err := jwtpkg.GenerateRefreshToken(session.UserID, sessionID, tokenConfig)
	if err != nil {
		return nil, errors.ErrInternal("failed to generate refresh token").WithError(err)
	}

	refreshTokenHash := hashToken(refreshToken)
	refreshTokenEntity := &entity.RefreshToken{
		UserID:      session.UserID,
		TokenHash:   refreshTokenHash,
		TokenFamily: tokenFamily,
		ExpiresAt:   time.Now().Add(uc.Config.JWT.RefreshExpiry),
		IPAddress:   req.IPAddress,
		UserAgent:   req.UserAgent,
		CreatedAt:   time.Now(),
	}
	now := time.Now()
	userSession := &entity.UserSession{
		UserID:       session.UserID,
		IPAddress:    req.IPAddress,
		LoginMethod:  entity.UserSessionLoginMethodEmailOTP,
		Status:       entity.UserSessionStatusActive,
		LastActiveAt: now,
		ExpiresAt:    now.Add(uc.Config.JWT.RefreshExpiry),
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if req.UserAgent != "" {
		userSession.UserAgent = &req.UserAgent
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
		return nil, errors.ErrInternal("failed to complete login").WithError(err)
	}

	_ = uc.InMemoryStore.DeleteLoginSession(ctx, req.LoginSessionID)

	profile, _ := uc.UserProfileRepo.GetByUserID(ctx, session.UserID)
	fullName := ""
	if profile != nil {
		fullName = profile.FirstName
		if profile.LastName != "" {
			fullName += " " + profile.LastName
		}
	}

	return &VerifyLoginOTPResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int(uc.Config.JWT.AccessExpiry.Seconds()),
		TokenType:    "Bearer",
		User: LoginUserResponse{
			ID:       session.UserID,
			Email:    session.Email,
			FullName: fullName,
			Tenants:  userTenants,
		},
	}, nil
}

func (uc *usecase) buildMultiTenantClaims(ctx context.Context, userID uuid.UUID) ([]jwtpkg.TenantClaim, []TenantResponse, []string, error) {
	registrations, err := uc.UserTenantRegRepo.ListByUserIDForClaims(ctx, userID)
	if err != nil {
		return nil, nil, nil, err
	}

	var jwtClaims []jwtpkg.TenantClaim
	var dtoTenants []TenantResponse

	for _, reg := range registrations {
		products, err := uc.ProductsByTenantRepo.ListActiveByTenantID(ctx, reg.TenantID)
		if err != nil {
			return nil, nil, nil, err
		}

		var jwtProducts []jwtpkg.ProductClaim
		var dtoProducts []ProductResponse

		for _, product := range products {
			userRoles, err := uc.UserRoleRepo.ListActiveByUserID(ctx, userID, &product.ID)
			if err != nil {
				return nil, nil, nil, err
			}

			var roleIDs []uuid.UUID
			var roleNames []string
			for _, ur := range userRoles {
				roleIDs = append(roleIDs, ur.RoleID)
			}

			if len(roleIDs) > 0 {
				roles, err := uc.RoleRepo.GetByIDs(ctx, roleIDs)
				if err != nil {
					return nil, nil, nil, err
				}
				for _, r := range roles {
					roleNames = append(roleNames, r.Code)
				}
			}

			var permissions []string
			if len(roleIDs) > 0 {
				permissions, err = uc.PermissionRepo.GetCodesByRoleIDs(ctx, roleIDs)
				if err != nil {
					return nil, nil, nil, err
				}
			}

			jwtProducts = append(jwtProducts, jwtpkg.ProductClaim{
				ProductID:   product.ID,
				ProductCode: product.Code,
				Roles:       roleNames,
				Permissions: permissions,
			})

			dtoProducts = append(dtoProducts, ProductResponse{
				ProductID:   product.ID,
				ProductCode: product.Code,
				Roles:       roleNames,
				Permissions: permissions,
			})
		}

		jwtClaims = append(jwtClaims, jwtpkg.TenantClaim{
			TenantID: reg.TenantID,
			Products: jwtProducts,
		})

		dtoTenants = append(dtoTenants, TenantResponse{
			TenantID: reg.TenantID,
			Products: dtoProducts,
		})
	}

	platformRoles, err := uc.UserRoleRepo.ListActiveByUserID(ctx, userID, nil)
	if err != nil {
		return nil, nil, nil, err
	}

	var platformRoleIDs []uuid.UUID
	for _, ur := range platformRoles {
		if ur.ProductID == nil {
			platformRoleIDs = append(platformRoleIDs, ur.RoleID)
		}
	}

	var platformRoleNames []string
	if len(platformRoleIDs) > 0 {
		roles, err := uc.RoleRepo.GetByIDs(ctx, platformRoleIDs)
		if err != nil {
			return nil, nil, nil, err
		}
		for _, r := range roles {
			platformRoleNames = append(platformRoleNames, r.Code)
		}
	}

	return jwtClaims, dtoTenants, platformRoleNames, nil
}

func (uc *usecase) buildTokenConfig() (*jwtpkg.TokenConfig, error) {
	tokenConfig := &jwtpkg.TokenConfig{
		SigningMethod: uc.Config.JWT.SigningMethod,
		AccessSecret:  uc.Config.JWT.AccessSecret,
		RefreshSecret: uc.Config.JWT.RefreshSecret,
		AccessExpiry:  uc.Config.JWT.AccessExpiry,
		RefreshExpiry: uc.Config.JWT.RefreshExpiry,
		Issuer:        uc.Config.JWT.Issuer,
		Audience:      uc.Config.JWT.Audience,
	}

	return tokenConfig, nil
}
