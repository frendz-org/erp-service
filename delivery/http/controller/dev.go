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
				"role_assignments_queue",
				"mfa_devices",
				"user_sessions",
				"refresh_tokens",
				"password_reset_tokens",
				"email_verifications",
				"pin_verification_logs",
				"auth_logs",
				"permission_checks",
				"admin_audit_logs",
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

			coreTables := []string{
				"user_tenant_registrations",
				"registrations",
				"user_auth_methods",
				"user_security_states",
				"user_profiles",
				"users",
			}
			for _, table := range coreTables {
				if err := tx.Exec(fmt.Sprintf("DELETE FROM %s WHERE id = ? OR user_id = ?", table), userID, userID).Error; err != nil {

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

	return c.JSON(response.SuccessResponse("user reset complete", fiber.Map{
		"email":      email,
		"user_found": userFound,
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
