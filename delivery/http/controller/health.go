package controller

import (
	"iam-service/config"
	"iam-service/delivery/http/dto/response"
	"iam-service/health"

	"github.com/gofiber/fiber/v2"
)

type HealthController struct {
	config        *config.Config
	healthUsecase health.Usecase
}

func NewHealthController(cfg *config.Config, healthUsecase health.Usecase) *HealthController {
	return &HealthController{
		config:        cfg,
		healthUsecase: healthUsecase,
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
	if err := h.healthUsecase.CheckHealth(); err != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(response.ErrorResponseWithDetails(
			"SERVICE_UNAVAILABLE",
			"Service not ready",
			fiber.Map{"status": "not ready"},
		))
	}

	return c.JSON(response.SuccessResponse("Service ready", fiber.Map{
		"status": "ready",
	}))
}

func (h *HealthController) Live(c *fiber.Ctx) error {
	return c.JSON(response.SuccessResponse("Service alive", fiber.Map{
		"status": "alive",
	}))
}
