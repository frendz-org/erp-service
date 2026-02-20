package middleware

import (
	"context"
	"net"
	"strings"

	"iam-service/pkg/logger"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const (
	ClientIPKey        = "client_ip"
	UserAgentKey       = "user_agent"
	TenantIDFromHdrKey = "tenant_id_from_header"
)

func RequestContext() fiber.Handler {
	return func(c *fiber.Ctx) error {
		clientIP := extractClientIP(c)
		c.Locals(ClientIPKey, clientIP)

		userAgent := c.Get("User-Agent")
		c.Locals(UserAgentKey, userAgent)

		if tenantIDStr := c.Get("X-Tenant-ID"); tenantIDStr != "" {
			if tenantID, err := uuid.Parse(tenantIDStr); err == nil {
				c.Locals(TenantIDFromHdrKey, tenantID)
			}
		}

		ctx := c.UserContext()
		ctx = context.WithValue(ctx, logger.CtxRequestID, c.GetRespHeader("X-Request-ID"))
		if clientIP != nil {
			ctx = context.WithValue(ctx, logger.CtxIPAddress, clientIP.String())
		}
		ctx = context.WithValue(ctx, logger.CtxUserAgent, userAgent)
		c.SetUserContext(ctx)

		return c.Next()
	}
}

func extractClientIP(c *fiber.Ctx) net.IP {
	if forwarded := c.Get("X-Forwarded-For"); forwarded != "" {
		parts := strings.SplitN(forwarded, ",", 2)
		ip := strings.TrimSpace(parts[0])
		if parsed := net.ParseIP(ip); parsed != nil {
			return parsed
		}
	}

	return net.ParseIP(c.IP())
}
