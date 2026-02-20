package logger

import (
	"context"
	"time"

	"go.uber.org/zap"
)

type ContextKey string

const (
	CtxRequestID ContextKey = "request_id"
	CtxIPAddress ContextKey = "ip_address"
	CtxUserAgent ContextKey = "user_agent"
	CtxActorID   ContextKey = "actor_id"
	CtxTenantID  ContextKey = "tenant_id"
	CtxSessionID ContextKey = "session_id"
)

type AuditEvent struct {
	Domain     string
	Action     string
	ActorID    string
	ActorType  string
	TargetID   string
	TargetType string
	TenantID   string
	Success    bool
	Reason     string
	Metadata   map[string]any
}

type AuditLogger interface {
	Log(ctx context.Context, event AuditEvent)

	Sync() error
}

type AuditConfig struct {
	Enabled bool
}

type zapAuditLogger struct {
	logger *zap.Logger
}

func NewAuditLogger(base *zap.Logger, cfg AuditConfig) AuditLogger {
	if !cfg.Enabled {
		return &NoopAuditLogger{}
	}
	return &zapAuditLogger{
		logger: base.Named("audit"),
	}
}

func (l *zapAuditLogger) Log(ctx context.Context, event AuditEvent) {
	fields := make([]zap.Field, 0, 16)

	fields = append(fields,
		zap.String("domain", event.Domain),
		zap.String("action", event.Action),
		zap.Bool("success", event.Success),
		zap.Time("timestamp", time.Now().UTC()),
	)

	if event.ActorID != "" {
		fields = append(fields, zap.String("actor_id", event.ActorID))
	}
	if event.ActorType != "" {
		fields = append(fields, zap.String("actor_type", event.ActorType))
	}

	if event.TargetID != "" {
		fields = append(fields, zap.String("target_id", event.TargetID))
	}
	if event.TargetType != "" {
		fields = append(fields, zap.String("target_type", event.TargetType))
	}

	if event.TenantID != "" {
		fields = append(fields, zap.String("tenant_id", event.TenantID))
	}
	if event.Reason != "" {
		fields = append(fields, zap.String("reason", event.Reason))
	}

	if reqID := extractString(ctx, CtxRequestID); reqID != "" {
		fields = append(fields, zap.String("request_id", reqID))
	}
	if ip := extractString(ctx, CtxIPAddress); ip != "" {
		fields = append(fields, zap.String("ip_address", ip))
	}
	if ua := extractString(ctx, CtxUserAgent); ua != "" {
		fields = append(fields, zap.String("user_agent", ua))
	}
	if sessionID := extractString(ctx, CtxSessionID); sessionID != "" {
		fields = append(fields, zap.String("session_id", sessionID))
	}

	if len(event.Metadata) > 0 {
		fields = append(fields, zap.Any("metadata", event.Metadata))
	}

	l.logger.Info("audit", fields...)
}

func (l *zapAuditLogger) Sync() error {
	return l.logger.Sync()
}

func extractString(ctx context.Context, key ContextKey) string {
	if val := ctx.Value(key); val != nil {
		if s, ok := val.(string); ok {
			return s
		}
	}
	return ""
}

type NoopAuditLogger struct{}

func NewNoopAuditLogger() AuditLogger {
	return &NoopAuditLogger{}
}

func (l *NoopAuditLogger) Log(ctx context.Context, event AuditEvent) {}
func (l *NoopAuditLogger) Sync() error                               { return nil }
