package router

import (
	"erp-service/delivery/http/controller"

	"github.com/gofiber/fiber/v2"
)

func SetupPensionRoutes(api fiber.Router, ctrl *controller.PensionController, jwtMiddleware fiber.Handler) {
	pension := api.Group("/pension")
	pension.Use(jwtMiddleware)

	pension.Get("/amount-summary", ctrl.GetAmountSummary)
}
