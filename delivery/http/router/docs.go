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
    <style>
        .tenant-input-container {
            max-width: 1460px;
            margin: 0 auto;
            padding: 10px 20px;
            display: flex;
            align-items: center;
            gap: 10px;
            font-family: sans-serif;
            font-size: 14px;
        }
        .tenant-input-container label {
            font-weight: 600;
            white-space: nowrap;
        }
        .tenant-input-container input {
            padding: 6px 10px;
            border: 1px solid #d9d9d9;
            border-radius: 4px;
            width: 320px;
            font-size: 14px;
        }
    </style>
</head>
<body>
    <div class="tenant-input-container">
        <label for="tenant-id-input">X-Tenant-ID:</label>
        <input type="text" id="tenant-id-input" placeholder="Enter tenant UUID (auto-injected into requests)">
    </div>
    <div id="swagger-ui"></div>
    <script>%s</script>
    <script>
        SwaggerUIBundle({
            url: "/api/v1/docs/openapi.yaml",
            dom_id: "#swagger-ui",
            persistAuthorization: true,
            requestInterceptor: function(req) {
                var input = document.getElementById('tenant-id-input');
                var tenantId = input ? input.value : '';
                if (tenantId) {
                    req.headers['X-Tenant-ID'] = tenantId;
                }
                return req;
            }
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
