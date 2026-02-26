package controller

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"erp-service/delivery/http/dto/response"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type DevRedisStore interface {
	Delete(ctx context.Context, keys ...string) error
}

type DevController struct {
	db    *gorm.DB
	redis DevRedisStore
}

func NewDevController(db *gorm.DB, redis DevRedisStore) *DevController {
	return &DevController{db: db, redis: redis}
}

func (d *DevController) ResetUserByEmail(c *fiber.Ctx) error {
	rawEmail, _ := url.QueryUnescape(c.Params("email"))
	email := strings.ToLower(strings.TrimSpace(rawEmail))
	if email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse(
			"INVALID_EMAIL", "email parameter is required",
		))
	}

	ctx := c.UserContext()

	var userID string
	err := d.db.WithContext(ctx).
		Raw("SELECT id FROM users WHERE LOWER(email) = ?", email).
		Scan(&userID).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse(
			"DB_ERROR", "failed to look up user",
		))
	}

	userFound := userID != ""

	if userFound {
		if err := d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

			leafTables := []string{
				"user_role_assignments",
				"user_sessions",
				"refresh_tokens",
				"verification_challenges",
				"mfa_enrollments",
				"recovery_codes",
				"password_history",
				"user_branches",
			}
			for _, table := range leafTables {
				if err := tx.Exec(fmt.Sprintf("DELETE FROM %s WHERE user_id = ?", table), userID).Error; err != nil {
					return fmt.Errorf("delete from %s: %w", table, err)
				}
			}

			var participantIDs []string
			if err := tx.Raw("SELECT id FROM participants WHERE user_id = ?", userID).Scan(&participantIDs).Error; err != nil {
				return fmt.Errorf("select participant IDs: %w", err)
			}

			if len(participantIDs) > 0 {
				participantChildTables := []string{
					"participant_beneficiaries",
					"participant_family_members",
					"participant_identities",
					"participant_employments",
					"participant_addresses",
					"participant_status_history",
					"participant_bank_accounts",
					"participant_pensions",
				}
				for _, table := range participantChildTables {
					if err := tx.Exec(fmt.Sprintf("DELETE FROM %s WHERE participant_id IN ?", table), participantIDs).Error; err != nil {
						return fmt.Errorf("delete from %s: %w", table, err)
					}
				}
				if err := tx.Exec("DELETE FROM participants WHERE user_id = ?", userID).Error; err != nil {
					return fmt.Errorf("delete from participants: %w", err)
				}
			}

			if err := tx.Exec("DELETE FROM files WHERE uploaded_by = ?", userID).Error; err != nil {
				return fmt.Errorf("delete from files: %w", err)
			}

			coreTablesByUserID := []string{
				"user_tenant_registrations",
				"user_auth_methods",
				"user_security_states",
				"user_profiles",
			}
			for _, table := range coreTablesByUserID {
				if err := tx.Exec(fmt.Sprintf("DELETE FROM %s WHERE user_id = ?", table), userID).Error; err != nil {
					return fmt.Errorf("delete from %s: %w", table, err)
				}
			}

			if err := tx.Exec("DELETE FROM users WHERE id = ?", userID).Error; err != nil {
				return fmt.Errorf("delete from users: %w", err)
			}

			return nil
		}); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse(
				"RESET_FAILED", fmt.Sprintf("transaction failed: %s", err.Error()),
			))
		}
	}

	redisKeys := []string{
		"reg_email:" + email,
		"reg_rate:" + email,
		"login_rate:" + email,
	}
	if userFound {
		redisKeys = append(redisKeys, "blacklist:user:"+userID)
	}
	_ = d.redis.Delete(ctx, redisKeys...)

	return c.JSON(response.SuccessResponse("user reset complete", fiber.Map{
		"email":      email,
		"user_found": userFound,
	}))
}

type AssignProductAdminRequest struct {
	TenantCode  string `json:"tenant_code"`
	ProductCode string `json:"product_code"`
}

func (d *DevController) AssignProductAdmin(c *fiber.Ctx) error {
	rawEmail, _ := url.QueryUnescape(c.Params("email"))
	email := strings.ToLower(strings.TrimSpace(rawEmail))
	if email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse(
			"INVALID_EMAIL", "email parameter is required",
		))
	}

	var req AssignProductAdminRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse(
			"INVALID_BODY", "invalid request body",
		))
	}
	if req.TenantCode == "" || req.ProductCode == "" {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse(
			"MISSING_FIELDS", "tenant_code and product_code are required",
		))
	}

	ctx := c.UserContext()

	var userID string
	if err := d.db.WithContext(ctx).
		Raw("SELECT id FROM users WHERE LOWER(email) = ? AND deleted_at IS NULL", email).
		Scan(&userID).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse(
			"DB_ERROR", "failed to look up user",
		))
	}
	if userID == "" {
		return c.Status(fiber.StatusNotFound).JSON(response.ErrorResponse(
			"USER_NOT_FOUND", fmt.Sprintf("user with email %s not found", email),
		))
	}

	var tenantID string
	if err := d.db.WithContext(ctx).
		Raw("SELECT id FROM tenants WHERE code = ? AND deleted_at IS NULL", req.TenantCode).
		Scan(&tenantID).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse(
			"DB_ERROR", "failed to look up tenant",
		))
	}
	if tenantID == "" {
		return c.Status(fiber.StatusNotFound).JSON(response.ErrorResponse(
			"TENANT_NOT_FOUND", fmt.Sprintf("tenant with code %s not found", req.TenantCode),
		))
	}

	var productID string
	if err := d.db.WithContext(ctx).
		Raw("SELECT id FROM products WHERE tenant_id = ? AND code = ? AND deleted_at IS NULL", tenantID, req.ProductCode).
		Scan(&productID).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse(
			"DB_ERROR", "failed to look up product",
		))
	}
	if productID == "" {
		return c.Status(fiber.StatusNotFound).JSON(response.ErrorResponse(
			"PRODUCT_NOT_FOUND", fmt.Sprintf("product with code %s not found for tenant %s", req.ProductCode, req.TenantCode),
		))
	}

	var roleID string
	if err := d.db.WithContext(ctx).
		Raw("SELECT id FROM roles WHERE product_id = ? AND code = ? AND deleted_at IS NULL", productID, "TENANT_PRODUCT_ADMIN").
		Scan(&roleID).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse(
			"DB_ERROR", "failed to look up role",
		))
	}
	if roleID == "" {
		return c.Status(fiber.StatusNotFound).JSON(response.ErrorResponse(
			"ROLE_NOT_FOUND", "TENANT_PRODUCT_ADMIN role not found for this product (not seeded?)",
		))
	}

	var existingID string
	if err := d.db.WithContext(ctx).
		Raw(`SELECT id FROM user_role_assignments
			 WHERE user_id = ? AND role_id = ? AND product_id = ?
			 AND status = 'ACTIVE' AND deleted_at IS NULL`, userID, roleID, productID).
		Scan(&existingID).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse(
			"DB_ERROR", "failed to check existing assignment",
		))
	}
	if existingID != "" {
		return c.JSON(response.SuccessResponse("TENANT_PRODUCT_ADMIN role already assigned", fiber.Map{
			"email":            email,
			"tenant_code":      req.TenantCode,
			"product_code":     req.ProductCode,
			"already_assigned": true,
		}))
	}

	if err := d.db.WithContext(ctx).
		Exec(`INSERT INTO user_role_assignments
			  (id, user_id, role_id, product_id, branch_id, status, assigned_at, assigned_by, expires_at, created_at, updated_at)
			  VALUES (gen_random_uuid(), ?, ?, ?, NULL, 'ACTIVE', NOW(), NULL, NULL, NOW(), NOW())`,
			userID, roleID, productID).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse(
			"ASSIGNMENT_FAILED", fmt.Sprintf("failed to create assignment: %s", err.Error()),
		))
	}

	return c.JSON(response.SuccessResponse("TENANT_PRODUCT_ADMIN role assigned", fiber.Map{
		"email":            email,
		"tenant_code":      req.TenantCode,
		"product_code":     req.ProductCode,
		"already_assigned": false,
	}))
}

func (d *DevController) ResetUserSessionsByEmail(c *fiber.Ctx) error {
	rawEmail, _ := url.QueryUnescape(c.Params("email"))
	email := strings.ToLower(strings.TrimSpace(rawEmail))
	if email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse(
			"INVALID_EMAIL", "email parameter is required",
		))
	}

	ctx := c.UserContext()

	var userID string
	err := d.db.WithContext(ctx).
		Raw("SELECT id FROM users WHERE LOWER(email) = ?", email).
		Scan(&userID).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse(
			"DB_ERROR", "failed to look up user",
		))
	}

	userFound := userID != ""

	if userFound {
		if err := d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			sessionTables := []string{
				"user_sessions",
				"refresh_tokens",
			}
			for _, table := range sessionTables {
				if err := tx.Exec(fmt.Sprintf("DELETE FROM %s WHERE user_id = ?", table), userID).Error; err != nil {
					return fmt.Errorf("delete from %s: %w", table, err)
				}
			}
			return nil
		}); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse(
				"RESET_FAILED", fmt.Sprintf("transaction failed: %s", err.Error()),
			))
		}
	}

	redisKeys := []string{
		"reg_email:" + email,
		"reg_rate:" + email,
		"login_rate:" + email,
	}
	if userFound {
		redisKeys = append(redisKeys, "blacklist:user:"+userID)
	}
	_ = d.redis.Delete(ctx, redisKeys...)

	return c.JSON(response.SuccessResponse("user sessions reset complete", fiber.Map{
		"email":      email,
		"user_found": userFound,
	}))
}
