package middleware

import (
	"erp-service/pkg/errors"

	"github.com/gofiber/fiber/v2"
)

func RequirePlatformAdmin() fiber.Handler {
	return func(c *fiber.Ctx) error {
		claims, err := GetUserClaims(c)
		if err != nil {

			appErr := errors.ErrUnauthorized("authentication required")
			return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
				"success": false,
				"error":   appErr.Message,
				"code":    appErr.Code,
			})
		}

		if !claims.IsPlatformAdmin() {
			appErr := errors.ErrPlatformAdminRequired()
			return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
				"success": false,
				"error":   appErr.Message,
				"code":    appErr.Code,
			})
		}

		return c.Next()
	}
}
func RequireRole(roleCode string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		claims, err := GetUserClaims(c)
		if err != nil {
			appErr := errors.ErrUnauthorized("authentication required")
			return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
				"success": false,
				"error":   appErr.Message,
				"code":    appErr.Code,
			})
		}

		if !claims.HasRole(roleCode) {
			appErr := errors.ErrAccessForbidden("insufficient permissions")
			return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
				"success": false,
				"error":   appErr.Message,
				"code":    appErr.Code,
			})
		}

		return c.Next()
	}
}

func RequireTenantRole(roleCodes ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		multiClaims, err := GetMultiTenantClaims(c)
		if err != nil {
			appErr := errors.ErrUnauthorized("authentication required")
			return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
				"success": false,
				"error":   appErr.Message,
				"code":    appErr.Code,
			})
		}

		tenantID, err := GetTenantIDFromHeader(c)
		if err != nil {
			appErr := errors.ErrBadRequest("X-Tenant-ID header is required")
			return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
				"success": false,
				"error":   appErr.Message,
				"code":    appErr.Code,
			})
		}

		if !multiClaims.HasTenant(tenantID) {
			appErr := errors.ErrForbidden("access denied to this tenant")
			return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
				"success": false,
				"error":   appErr.Message,
				"code":    appErr.Code,
			})
		}

		tenantClaim := multiClaims.GetTenantClaim(tenantID)
		if tenantClaim == nil {
			appErr := errors.ErrForbidden("tenant not found in claims")
			return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
				"success": false,
				"error":   appErr.Message,
				"code":    appErr.Code,
			})
		}

		hasRole := false
		for _, product := range tenantClaim.Products {
			for _, role := range product.Roles {
				for _, requiredRole := range roleCodes {
					if role == requiredRole {
						hasRole = true
						break
					}
				}
				if hasRole {
					break
				}
			}
			if hasRole {
				break
			}
		}

		if !hasRole {
			if multiClaims.IsPlatformAdmin() {
				return c.Next()
			}
			appErr := errors.ErrAccessForbidden("insufficient permissions for this tenant")
			return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
				"success": false,
				"error":   appErr.Message,
				"code":    appErr.Code,
			})
		}

		return c.Next()
	}
}

func RequireAnyTenantRole(roleCodes ...string) fiber.Handler {
	return RequireTenantRole(roleCodes...)
}

func RequireProductRole(roleCodes ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		multiClaims, err := GetMultiTenantClaims(c)
		if err != nil {
			appErr := errors.ErrUnauthorized("authentication required")
			return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
				"success": false,
				"error":   appErr.Message,
				"code":    appErr.Code,
			})
		}

		tenantID, err := GetTenantIDFromContext(c)
		if err != nil {
			appErr := errors.ErrBadRequest("tenant context is required")
			return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
				"success": false,
				"error":   appErr.Message,
				"code":    appErr.Code,
			})
		}

		productID, err := GetProductIDFromContext(c)
		if err != nil {
			appErr := errors.ErrBadRequest("product context is required")
			return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
				"success": false,
				"error":   appErr.Message,
				"code":    appErr.Code,
			})
		}

		for _, roleCode := range roleCodes {
			if multiClaims.HasRoleInProduct(tenantID, productID, roleCode) {
				return c.Next()
			}
		}

		if multiClaims.IsPlatformAdmin() {
			return c.Next()
		}

		appErr := errors.ErrAccessForbidden("insufficient permissions for this product")
		return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
			"success": false,
			"error":   appErr.Message,
			"code":    appErr.Code,
		})
	}
}

func RequireProductPermission(permissionCodes ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		multiClaims, err := GetMultiTenantClaims(c)
		if err != nil {
			appErr := errors.ErrUnauthorized("authentication required")
			return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
				"success": false,
				"error":   appErr.Message,
				"code":    appErr.Code,
			})
		}

		tenantID, err := GetTenantIDFromContext(c)
		if err != nil {
			appErr := errors.ErrBadRequest("tenant context is required")
			return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
				"success": false,
				"error":   appErr.Message,
				"code":    appErr.Code,
			})
		}

		productID, err := GetProductIDFromContext(c)
		if err != nil {
			appErr := errors.ErrBadRequest("product context is required")
			return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
				"success": false,
				"error":   appErr.Message,
				"code":    appErr.Code,
			})
		}

		for _, permissionCode := range permissionCodes {
			if multiClaims.HasPermissionInProduct(tenantID, productID, permissionCode) {
				return c.Next()
			}
		}

		if multiClaims.IsPlatformAdmin() {
			return c.Next()
		}

		appErr := errors.ErrAccessForbidden("insufficient permissions for this product")
		return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
			"success": false,
			"error":   appErr.Message,
			"code":    appErr.Code,
		})
	}
}

func RequireTenantPermission(permissionCode string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		multiClaims, err := GetMultiTenantClaims(c)
		if err != nil {
			appErr := errors.ErrUnauthorized("authentication required")
			return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
				"success": false,
				"error":   appErr.Message,
				"code":    appErr.Code,
			})
		}

		tenantID, err := GetTenantIDFromHeader(c)
		if err != nil {
			appErr := errors.ErrBadRequest("X-Tenant-ID header is required")
			return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
				"success": false,
				"error":   appErr.Message,
				"code":    appErr.Code,
			})
		}

		if !multiClaims.HasTenant(tenantID) {
			appErr := errors.ErrForbidden("access denied to this tenant")
			return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
				"success": false,
				"error":   appErr.Message,
				"code":    appErr.Code,
			})
		}

		tenantClaim := multiClaims.GetTenantClaim(tenantID)
		if tenantClaim == nil {
			appErr := errors.ErrForbidden("tenant not found in claims")
			return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
				"success": false,
				"error":   appErr.Message,
				"code":    appErr.Code,
			})
		}

		hasPermission := false
		for _, product := range tenantClaim.Products {
			for _, perm := range product.Permissions {
				if perm == permissionCode {
					hasPermission = true
					break
				}
			}
			if hasPermission {
				break
			}
		}

		if !hasPermission {
			if multiClaims.IsPlatformAdmin() {
				return c.Next()
			}
			appErr := errors.ErrAccessForbidden("insufficient permissions for this operation")
			return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
				"success": false,
				"error":   appErr.Message,
				"code":    appErr.Code,
			})
		}

		return c.Next()
	}
}
