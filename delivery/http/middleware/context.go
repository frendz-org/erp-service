package middleware

import (
	"net"
	"strings"

	"erp-service/iam/product"
	"erp-service/pkg/errors"
	jwtpkg "erp-service/pkg/jwt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const (
	UserClaimsKey        = "user_claims"
	MultiTenantClaimsKey = "multi_tenant_claims"
)

func GetUserClaims(c *fiber.Ctx) (*jwtpkg.JWTClaims, error) {
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

func GetMultiTenantClaims(c *fiber.Ctx) (*jwtpkg.MultiTenantClaims, error) {
	claims := c.Locals(MultiTenantClaimsKey)
	if claims == nil {
		return nil, errors.ErrUnauthorized("multi-tenant claims not found in context")
	}

	multiClaims, ok := claims.(*jwtpkg.MultiTenantClaims)
	if !ok {
		return nil, errors.ErrInternal("invalid multi-tenant claims type in context")
	}

	return multiClaims, nil
}

func GetUserID(c *fiber.Ctx) (uuid.UUID, error) {
	claims, err := GetUserClaims(c)
	if err != nil {
		return uuid.Nil, err
	}
	return claims.UserID, nil
}

func GetTenantID(c *fiber.Ctx) (*uuid.UUID, error) {
	claims, err := GetUserClaims(c)
	if err != nil {
		return nil, err
	}
	return claims.TenantID, nil
}

func GetBranchID(c *fiber.Ctx) (*uuid.UUID, error) {
	claims, err := GetUserClaims(c)
	if err != nil {
		return nil, err
	}
	return claims.BranchID, nil
}

func GetSessionID(c *fiber.Ctx) (uuid.UUID, error) {
	claims, err := GetUserClaims(c)
	if err != nil {
		return uuid.Nil, err
	}
	return claims.SessionID, nil
}

func GetUserRoles(c *fiber.Ctx) ([]string, error) {
	claims, err := GetUserClaims(c)
	if err != nil {
		return nil, err
	}
	return claims.Roles, nil
}

func GetClientIP(c *fiber.Ctx) net.IP {
	forwarded := c.Get("X-Forwarded-For")
	if forwarded != "" {
		if idx := strings.IndexByte(forwarded, ','); idx != -1 {
			forwarded = strings.TrimSpace(forwarded[:idx])
		}
		if ip := net.ParseIP(forwarded); ip != nil {
			return ip
		}
	}
	return net.ParseIP(c.IP())
}

func GetUserAgent(c *fiber.Ctx) string {
	return c.Get("User-Agent")
}

func GetRequestID(c *fiber.Ctx) string {
	if id := c.GetRespHeader("X-Request-ID"); id != "" {
		return id
	}
	return c.Get("X-Request-ID")
}

func IsPlatformAdmin(c *fiber.Ctx) (bool, error) {
	claims, err := GetUserClaims(c)
	if err != nil {
		return false, err
	}
	if claims.IsPlatformAdmin() {
		return true, nil
	}

	multiClaims, err := GetMultiTenantClaims(c)
	if err == nil && multiClaims.IsPlatformAdmin() {
		return true, nil
	}

	return false, nil
}

func GetTenantIDFromHeader(c *fiber.Ctx) (uuid.UUID, error) {
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

func ExtractTenantContext() fiber.Handler {
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
			if isPlatformAdmin, err := IsPlatformAdmin(c); err != nil || !isPlatformAdmin {
				appErr := errors.ErrForbidden("access denied to this tenant")
				return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
					"success": false,
					"error":   appErr.Message,
					"code":    appErr.Code,
				})
			}
		}

		c.Locals("tenant_id", tenantID)

		return c.Next()
	}
}

func GetTenantIDFromContext(c *fiber.Ctx) (uuid.UUID, error) {
	tenantID := c.Locals("tenant_id")
	if tenantID == nil {
		return uuid.Nil, errors.ErrForbidden("tenant context not found")
	}

	tid, ok := tenantID.(uuid.UUID)
	if !ok {
		return uuid.Nil, errors.ErrInternal("invalid tenant ID type in context")
	}

	return tid, nil
}

func ExtractProductContext() fiber.Handler {
	return func(c *fiber.Ctx) error {
		productIDStr := c.Params("productId")
		if productIDStr == "" {
			appErr := errors.ErrBadRequest("product ID is required")
			return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
				"success": false,
				"error":   appErr.Message,
				"code":    appErr.Code,
			})
		}

		productID, err := uuid.Parse(productIDStr)
		if err != nil {
			appErr := errors.ErrBadRequest("invalid product ID format")
			return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
				"success": false,
				"error":   appErr.Message,
				"code":    appErr.Code,
			})
		}

		multiClaims, _ := GetMultiTenantClaims(c)
		if multiClaims != nil && !multiClaims.IsPlatformAdmin() {
			tenantID, err := GetTenantIDFromContext(c)
			if err == nil {
				tc := multiClaims.GetTenantClaim(tenantID)
				if tc == nil {
					appErr := errors.ErrForbidden("access denied to this tenant")
					return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
						"success": false,
						"error":   appErr.Message,
						"code":    appErr.Code,
					})
				}
				hasProduct := false
				for _, p := range tc.Products {
					if p.ProductID == productID {
						hasProduct = true
						break
					}
				}
				if !hasProduct {
					appErr := errors.ErrForbidden("access denied to this product")
					return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
						"success": false,
						"error":   appErr.Message,
						"code":    appErr.Code,
					})
				}
			}
		}

		c.Locals("product_id", productID)

		return c.Next()
	}
}

func ExtractFrendzSavingProduct(productUC product.Usecase) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tenantID, err := GetTenantIDFromContext(c)
		if err != nil {
			appErr := errors.ErrForbidden("tenant context not found")
			return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
				"success": false,
				"error":   appErr.Message,
				"code":    appErr.Code,
			})
		}

		product, err := productUC.GetFrendzSaving(c.UserContext(), tenantID)
		if err != nil {
			appErr := errors.ErrNotFound("product not available for this tenant")
			return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
				"success": false,
				"error":   appErr.Message,
				"code":    appErr.Code,
			})
		}

		multiClaims, _ := GetMultiTenantClaims(c)
		if multiClaims != nil && !multiClaims.IsPlatformAdmin() {
			tc := multiClaims.GetTenantClaim(tenantID)
			if tc == nil {
				appErr := errors.ErrForbidden("access denied to this tenant")
				return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
					"success": false,
					"error":   appErr.Message,
					"code":    appErr.Code,
				})
			}
			hasProduct := false
			for _, p := range tc.Products {
				if p.ProductID == product.ID {
					hasProduct = true
					break
				}
			}
			if !hasProduct {
				appErr := errors.ErrForbidden("access denied to this product")
				return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
					"success": false,
					"error":   appErr.Message,
					"code":    appErr.Code,
				})
			}
		}

		c.Locals("product_id", product.ID)
		return c.Next()
	}
}

func GetProductIDFromContext(c *fiber.Ctx) (uuid.UUID, error) {
	productID := c.Locals("product_id")
	if productID == nil {
		return uuid.Nil, errors.ErrBadRequest("product context not found")
	}

	pid, ok := productID.(uuid.UUID)
	if !ok {
		return uuid.Nil, errors.ErrInternal("invalid product ID type in context")
	}

	return pid, nil
}
