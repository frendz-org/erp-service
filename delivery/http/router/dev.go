package router

import (
	"erp-service/config"
	"erp-service/delivery/http/controller"

	"github.com/gofiber/fiber/v2"
)

func SetupDevRoutes(api fiber.Router, cfg *config.Config, devController *controller.DevController) {
	if cfg.IsProduction() {
		return
	}

	dev := api.Group("/dev")
	dev.Delete("/users/:email", devController.ResetUserByEmail)
	dev.Delete("/users/:email/sessions", devController.ResetUserSessionsByEmail)
}
