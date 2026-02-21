package controller

import (
	"net"
	"time"

	"erp-service/pkg/errors"
	jwtpkg "erp-service/pkg/jwt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const (
	UserClaimsKey      = "user_claims"
	ClientIPKey        = "client_ip"
	UserAgentKey       = "user_agent"
	TenantIDFromHdrKey = "tenant_id_from_header"
)

func getUserClaims(c *fiber.Ctx) (*jwtpkg.JWTClaims, error) {
	claims := c.Locals(UserClaimsKey)
	if claims == nil {
		return nil, errors.ErrUnauthorized("user claims not found in context")
	}

	jwtClaims, ok := claims.(*jwtpkg.JWTClaims)
	if !ok {
		return nil, errors.ErrInternal("invalid claims type in context")
	}

	return jwtClaims, nil
}

func getUserID(c *fiber.Ctx) (uuid.UUID, error) {
	claims, err := getUserClaims(c)
	if err != nil {
		return uuid.Nil, err
	}
	return claims.UserID, nil
}

func getTenantID(c *fiber.Ctx) (*uuid.UUID, error) {
	claims, err := getUserClaims(c)
	if err != nil {
		return nil, err
	}
	return claims.TenantID, nil
}

func getBranchID(c *fiber.Ctx) (*uuid.UUID, error) {
	claims, err := getUserClaims(c)
	if err != nil {
		return nil, err
	}
	return claims.BranchID, nil
}

func getSessionID(c *fiber.Ctx) (uuid.UUID, error) {
	claims, err := getUserClaims(c)
	if err != nil {
		return uuid.Nil, err
	}
	return claims.SessionID, nil
}

func getUserRoles(c *fiber.Ctx) ([]string, error) {
	claims, err := getUserClaims(c)
	if err != nil {
		return nil, err
	}
	return claims.Roles, nil
}

func isPlatformAdmin(c *fiber.Ctx) (bool, error) {
	claims, err := getUserClaims(c)
	if err != nil {
		return false, err
	}
	return claims.IsPlatformAdmin(), nil
}

func getClientIP(c *fiber.Ctx) net.IP {
	if ip := c.Locals(ClientIPKey); ip != nil {
		if netIP, ok := ip.(net.IP); ok {
			return netIP
		}
	}

	return net.ParseIP(c.IP())
}

func getUserAgent(c *fiber.Ctx) string {
	if ua := c.Locals(UserAgentKey); ua != nil {
		if userAgent, ok := ua.(string); ok {
			return userAgent
		}
	}

	return c.Get("User-Agent")
}

func extractClaimsForLogout(c *fiber.Ctx) (userID uuid.UUID, jti string, tokenExp time.Time, err error) {
	multiClaims := c.Locals("multi_tenant_claims")
	if multiClaims != nil {
		mc, ok := multiClaims.(*jwtpkg.MultiTenantClaims)
		if ok {
			userID = mc.UserID
			jti = mc.RegisteredClaims.ID
			if mc.ExpiresAt != nil {
				tokenExp = mc.ExpiresAt.Time
			}
			return userID, jti, tokenExp, nil
		}
	}

	claims, claimsErr := getUserClaims(c)
	if claimsErr != nil {
		return uuid.Nil, "", time.Time{}, claimsErr
	}
	userID = claims.UserID
	jti = claims.RegisteredClaims.ID
	if claims.ExpiresAt != nil {
		tokenExp = claims.ExpiresAt.Time
	}
	return userID, jti, tokenExp, nil
}

func getTenantIDFromHeader(c *fiber.Ctx) (uuid.UUID, error) {
	if tid := c.Locals(TenantIDFromHdrKey); tid != nil {
		if tenantID, ok := tid.(uuid.UUID); ok {
			return tenantID, nil
		}
	}

	tenantIDStr := c.Get("X-Tenant-ID")
	if tenantIDStr == "" {
		return uuid.Nil, errors.ErrBadRequest("X-Tenant-ID header is required")
	}

	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		return uuid.Nil, errors.ErrBadRequest("Invalid X-Tenant-ID format")
	}

	return tenantID, nil
}
