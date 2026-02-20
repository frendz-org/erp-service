package router

import (
	"iam-service/delivery/http/controller"

	"github.com/gofiber/fiber/v2"
)

func SetupHealthRoutes(api fiber.Router, healthController *controller.HealthController) {
	health := api.Group("/health")
	health.Get("/", healthController.Check)
	health.Get("/ready", healthController.Ready)
	health.Get("/live", healthController.Live)
}
