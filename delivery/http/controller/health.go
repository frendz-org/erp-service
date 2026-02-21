package controller

import (
	"erp-service/config"
	"erp-service/delivery/http/dto/response"

	"github.com/gofiber/fiber/v2"
)

type HealthController struct {
	config *config.Config
}

func NewHealthController(cfg *config.Config) *HealthController {
	return &HealthController{
		config: cfg,
	}
}

func (h *HealthController) Check(c *fiber.Ctx) error {
	return c.JSON(response.SuccessResponse("Health check passed", fiber.Map{
		"status":      "healthy",
		"app":         h.config.App.Name,
		"version":     h.config.App.Version,
		"environment": h.config.App.Environment,
	}))
}

func (h *HealthController) Ready(c *fiber.Ctx) error {
	return c.JSON(response.SuccessResponse("Service ready", fiber.Map{
		"status": "ready",
	}))
}

func (h *HealthController) Live(c *fiber.Ctx) error {
	return c.JSON(response.SuccessResponse("Service alive", fiber.Map{
		"status": "alive",
	}))
}
