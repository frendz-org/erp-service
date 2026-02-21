package iamclient

import (
	"strings"

	"erp-service/pkg/jwt"

	"github.com/gofiber/fiber/v2"
)

func (c *Client) JWTAuthMiddleware() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		authHeader := ctx.Get("Authorization")
		if authHeader == "" {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing authorization header",
			})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid authorization header format",
			})
		}

		tokenString := parts[1]

		claims, err := c.ValidateToken(tokenString)
		if err != nil {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired token",
			})
		}

		ctx.Locals("user_claims", claims)
		ctx.Locals("user_id", claims.UserID)
		if claims.TenantID != nil {
			ctx.Locals("tenant_id", *claims.TenantID)
		}
		if claims.ProductID != nil {
			ctx.Locals("product_id", *claims.ProductID)
		}

		return ctx.Next()
	}
}

func (c *Client) RequirePermission(permissionCode string) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		claims := getClaims(ctx)
		if claims == nil {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Authentication required",
			})
		}

		if !claims.HasPermission(permissionCode) {
			return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Insufficient permissions",
			})
		}

		return ctx.Next()
	}
}

func (c *Client) RequireRole(roleCode string) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		claims := getClaims(ctx)
		if claims == nil {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Authentication required",
			})
		}

		if !claims.HasRole(roleCode) {
			return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Insufficient role",
			})
		}

		return ctx.Next()
	}
}

func (c *Client) RequireAnyRole(roleCodes []string) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		claims := getClaims(ctx)
		if claims == nil {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Authentication required",
			})
		}

		hasRole := false
		for _, role := range roleCodes {
			if claims.HasRole(role) {
				hasRole = true
				break
			}
		}

		if !hasRole {
			return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Insufficient role",
			})
		}

		return ctx.Next()
	}
}

func getClaims(ctx *fiber.Ctx) *jwt.JWTClaims {
	claims := ctx.Locals("user_claims")
	if claims == nil {
		return nil
	}
	if jwtClaims, ok := claims.(*jwt.JWTClaims); ok {
		return jwtClaims
	}
	return nil
}

func GetClaimsFromContext(ctx *fiber.Ctx) *jwt.JWTClaims {
	return getClaims(ctx)
}
