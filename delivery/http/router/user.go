package router

import (
	"iam-service/config"
	"iam-service/delivery/http/controller"
	"iam-service/delivery/http/middleware"
	"iam-service/iam/auth/contract"

	"github.com/gofiber/fiber/v2"
)

func SetupUserRoutes(api fiber.Router, cfg *config.Config, userController *controller.UserController, blacklistStore ...contract.TokenBlacklistStore) {
	users := api.Group("/users")
	users.Use(middleware.JWTAuth(cfg, blacklistStore...))

	users.Get("/me", userController.GetMe)
	users.Put("/me", userController.UpdateMe)

	adminUsers := users.Group("")
	adminUsers.Use(middleware.RequirePlatformAdmin())

	adminUsers.Post("/", userController.Create)
	adminUsers.Get("/", userController.List)
	adminUsers.Get("/:id", userController.GetByID)
	adminUsers.Put("/:id", userController.Update)
	adminUsers.Delete("/:id", userController.Delete)
	adminUsers.Post("/:id/approve", userController.Approve)
	adminUsers.Post("/:id/reject", userController.Reject)
	adminUsers.Post("/:id/unlock", userController.Unlock)
	adminUsers.Post("/:id/reset-pin", userController.ResetPIN)
}
