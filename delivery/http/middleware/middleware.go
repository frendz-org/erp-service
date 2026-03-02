package middleware

import (
	"runtime/debug"
	"time"

	"erp-service/config"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Middleware struct {
	config *config.Config
	logger *zap.Logger
}

func New(cfg *config.Config, logger *zap.Logger) *Middleware {
	return &Middleware{
		config: cfg,
		logger: logger,
	}
}

func (m *Middleware) Setup(app *fiber.App) {
	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(c *fiber.Ctx, e interface{}) {
			m.logger.Error("panic recovered",
				zap.Any("error", e),
				zap.String("stack", string(debug.Stack())),
				zap.String("path", c.Path()),
				zap.String("method", c.Method()),
				zap.String("request_id", c.GetRespHeader("X-Request-ID")),
			)
		},
	}))

	app.Use(func(c *fiber.Ctx) error {
		c.Set("X-Content-Type-Options", "nosniff")
		c.Set("X-Frame-Options", "DENY")
		c.Set("X-XSS-Protection", "0")
		c.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		c.Set("Cache-Control", "no-store, no-cache, must-revalidate")
		c.Set("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
		c.Set("Content-Security-Policy", "default-src 'none'; frame-ancestors 'none'")
		return c.Next()
	})

	app.Use(requestid.New(requestid.Config{
		Generator: func() string {
			return uuid.New().String()
		},
	}))

	app.Use(RequestContext())

	app.Use(RequestLogger(m.logger))

	corsOrigins := m.config.Server.CORSOrigins
	if corsOrigins == "" {
		corsOrigins = "*"
	}
	allowCredentials := corsOrigins != "*"
	app.Use(cors.New(cors.Config{
		AllowOrigins:     corsOrigins,
		AllowMethods:     "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization,X-Request-ID,X-Tenant-ID",
		AllowCredentials: allowCredentials,
		MaxAge:           300,
	}))

	if m.config.IsProduction() {
		app.Use(limiter.New(limiter.Config{
			Max:               10,
			Expiration:        1 * time.Minute,
			LimiterMiddleware: limiter.SlidingWindow{},
			KeyGenerator: func(c *fiber.Ctx) string {
				return c.IP()
			},
			LimitReached: func(c *fiber.Ctx) error {
				return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
					"success": false,
					"error":   "too many requests",
				})
			},
		}))
	}
}
