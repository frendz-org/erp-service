package publickey

import (
	"erp-service/config"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	Config *config.Config
}

func NewHandler(cfg *config.Config) *Handler {
	return &Handler{
		Config: cfg,
	}
}

func (h *Handler) GetPublicKeyPEM(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
		"error": "Public key endpoint not available when using HS256 signing method",
	})
}

type JWKSResponse struct {
	Keys []JWK `json:"keys"`
}

type JWK struct {
	Kty string `json:"kty"`
	Use string `json:"use"`
	Kid string `json:"kid"`
	N   string `json:"n"`
	E   string `json:"e"`
}

func (h *Handler) GetJWKS(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
		"error": "JWKS endpoint not available when using HS256 signing method",
	})
}

func (h *Handler) RegisterRoutes(router fiber.Router) {
	wellKnown := router.Group("/.well-known")
	wellKnown.Get("/public-key.pem", h.GetPublicKeyPEM)
	wellKnown.Get("/jwks.json", h.GetJWKS)
}
