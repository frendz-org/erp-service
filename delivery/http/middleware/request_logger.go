package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func RequestLogger(zapLogger *zap.Logger) fiber.Handler {
	log := zapLogger.WithOptions(zap.AddStacktrace(zap.DPanicLevel))
	return func(c *fiber.Ctx) error {
		start := time.Now()
		requestID := c.GetRespHeader("X-Request-ID")

		log.Info("request_started",
			zap.String("request_id", requestID),
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.String("ip", c.IP()),
			zap.String("user_agent", c.Get("User-Agent")),
			zap.String("query", string(c.Request().URI().QueryString())),
		)

		err := c.Next()

		if err != nil {
			if handlerErr := c.App().ErrorHandler(c, err); handlerErr != nil {
				return handlerErr
			}
		}

		status := c.Response().StatusCode()
		fields := []zap.Field{
			zap.String("request_id", requestID),
			zap.Int("status", status),
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.String("ip", c.IP()),
			zap.Duration("duration", time.Since(start)),
		}

		if status >= 500 {
			log.Error("request_completed", fields...)
		} else if status >= 400 {
			log.Warn("request_completed", fields...)
		} else {
			log.Info("request_completed", fields...)
		}

		return nil
	}
}
