package router

import (
	_ "embed"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

//go:embed swagger/swagger-ui.css
var swaggerCSS string

//go:embed swagger/swagger-ui-bundle.js
var swaggerJS string

func SetupDocsRoutes(api fiber.Router) {
	docs := api.Group("/docs")
	docs.Get("/openapi.yaml", func(c *fiber.Ctx) error {
		return c.SendFile("doc/openapi/openapi-iam-v1.yaml")
	})

	page := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>ERP Service API</title>
    <style>%s</style>
</head>
<body>
    <div id="swagger-ui"></div>
    <script>%s</script>
    <script>
        SwaggerUIBundle({
            url: "/api/v1/docs/openapi.yaml",
            dom_id: "#swagger-ui",
        });
    </script>
</body>
</html>`, swaggerCSS, swaggerJS)

	docs.Get("/swagger", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html")
		c.Set("Content-Security-Policy", "default-src 'self'; style-src 'unsafe-inline'; script-src 'unsafe-inline'; connect-src 'self'; img-src 'self' data:; font-src data:")
		return c.SendString(page)
	})
}
