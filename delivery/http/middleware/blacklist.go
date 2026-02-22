package middleware

import (
	"erp-service/iam/auth"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func checkBlacklist(c *fiber.Ctx, store auth.TokenBlacklistStore, jti string, userID uuid.UUID, claims jwt.RegisteredClaims) error {
	if jti == "" {
		return nil
	}

	blacklisted, err := store.IsTokenBlacklisted(c.UserContext(), jti)
	if err == nil && blacklisted {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   "token has been revoked",
			"code":    "TOKEN_REVOKED",
		})
	}

	blacklistTS, err := store.GetUserBlacklistTimestamp(c.UserContext(), userID)
	if err == nil && blacklistTS != nil && claims.IssuedAt != nil {
		if claims.IssuedAt.Time.Before(*blacklistTS) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"error":   "token has been revoked",
				"code":    "TOKEN_REVOKED",
			})
		}
	}

	return nil
}
