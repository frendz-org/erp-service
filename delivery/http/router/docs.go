package router

import "github.com/gofiber/fiber/v2"

func SetupDocsRoutes(api fiber.Router) {
	docs := api.Group("/docs")
	docs.Get("/openapi.yaml", func(c *fiber.Ctx) error {
		return c.SendFile("doc/openapi/openapi-iam-v1.yaml")
	})
}
