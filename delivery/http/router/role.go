package router

import (
	"erp-service/config"
	"erp-service/delivery/http/controller"
	"erp-service/delivery/http/middleware"
	"erp-service/iam/auth"

	"github.com/gofiber/fiber/v2"
)

func SetupRoleRoutes(api fiber.Router, cfg *config.Config, roleController *controller.RoleController, blacklistStore ...auth.TokenBlacklistStore) {
	roles := api.Group("/roles")

	roles.Use(middleware.JWTAuth(cfg, blacklistStore...))
	roles.Use(middleware.RequirePlatformAdmin())

	roles.Post("/", roleController.Create)
}
