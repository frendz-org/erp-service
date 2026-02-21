package middleware

import (
	"erp-service/config"
	"erp-service/iam/auth"
	"erp-service/pkg/errors"
	jwtpkg "erp-service/pkg/jwt"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func JWTAuth(cfg *config.Config, blacklistStore ...auth.TokenBlacklistStore) fiber.Handler {
	tokenConfig := &jwtpkg.TokenConfig{
		AccessSecret:  cfg.JWT.AccessSecret,
		RefreshSecret: cfg.JWT.RefreshSecret,
		AccessExpiry:  cfg.JWT.AccessExpiry,
		RefreshExpiry: cfg.JWT.RefreshExpiry,
		Issuer:        cfg.JWT.Issuer,
	}

	var store auth.TokenBlacklistStore
	if len(blacklistStore) > 0 {
		store = blacklistStore[0]
	}

	return func(c *fiber.Ctx) error {

		authHeader := c.Get("Authorization")
		if authHeader == "" {
			appErr := errors.ErrUnauthorized("missing authorization header")
			return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
				"success": false,
				"error":   appErr.Message,
				"code":    appErr.Code,
			})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			appErr := errors.ErrUnauthorized("invalid authorization header format")
			return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
				"success": false,
				"error":   appErr.Message,
				"code":    appErr.Code,
			})
		}

		tokenString := parts[1]

		multiClaims, multiErr := jwtpkg.ParseMultiTenantAccessToken(tokenString, tokenConfig)
		if multiErr == nil && len(multiClaims.Tenants) > 0 {
			legacyClaims := &jwtpkg.JWTClaims{
				UserID:           multiClaims.UserID,
				Email:            multiClaims.Email,
				Roles:            multiClaims.Roles,
				SessionID:        multiClaims.SessionID,
				RegisteredClaims: multiClaims.RegisteredClaims,
			}
			c.Locals(UserClaimsKey, legacyClaims)
			c.Locals(MultiTenantClaimsKey, multiClaims)

			c.Locals("userID", multiClaims.UserID.String())
			c.Locals("jti", multiClaims.RegisteredClaims.ID)

			if store != nil {
				if rejected := checkBlacklist(c, store, multiClaims.RegisteredClaims.ID, multiClaims.UserID, multiClaims.RegisteredClaims); rejected != nil {
					return rejected
				}
			}

			return c.Next()
		}

		claims, err := jwtpkg.ParseAccessToken(tokenString, tokenConfig)
		if err != nil {
			var appErr *errors.AppError
			switch err {
			case jwtpkg.ErrTokenExpired:
				appErr = errors.ErrTokenExpired()
			case jwtpkg.ErrTokenInvalid, jwtpkg.ErrTokenMalformed, jwtpkg.ErrTokenSignature:
				appErr = errors.ErrTokenInvalid()
			default:
				appErr = errors.ErrTokenInvalid()
			}

			return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
				"success": false,
				"error":   appErr.Message,
				"code":    appErr.Code,
			})
		}

		c.Locals(UserClaimsKey, claims)

		c.Locals("userID", claims.UserID.String())
		c.Locals("jti", claims.RegisteredClaims.ID)

		if store != nil {
			if rejected := checkBlacklist(c, store, claims.RegisteredClaims.ID, claims.UserID, claims.RegisteredClaims); rejected != nil {
				return rejected
			}
		}

		return c.Next()
	}
}
