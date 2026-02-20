package router

import (
	"iam-service/config"
	"iam-service/delivery/http/controller"
	"iam-service/delivery/http/middleware"
	"iam-service/iam/auth/contract"

	"github.com/gofiber/fiber/v2"
)

func SetupRoleRoutes(api fiber.Router, cfg *config.Config, roleController *controller.RoleController, blacklistStore ...contract.TokenBlacklistStore) {
	roles := api.Group("/roles")

	roles.Use(middleware.JWTAuth(cfg, blacklistStore...))
	roles.Use(middleware.RequirePlatformAdmin())

	roles.Post("/", roleController.Create)
}
