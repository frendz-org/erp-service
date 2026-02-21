package logger

import (
	"erp-service/config"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	With(args ...interface{}) Logger
	WithFields(fields map[string]interface{}) Logger
	WithError(err error) Logger
	WithRequestID(requestID string) Logger
	WithTenant(tenantID, tenantSlug string) Logger
	WithUser(userID string) Logger
	Sync() error
	Named(name string) Logger

	Debug(msg string, keysAndValues ...interface{})
	Info(msg string, keysAndValues ...interface{})
	Warn(msg string, keysAndValues ...interface{})
	Error(msg string, keysAndValues ...interface{})
	Fatal(msg string, keysAndValues ...interface{})
	Panic(msg string, keysAndValues ...interface{})
}

type logger struct {
	*zap.SugaredLogger
	base *zap.Logger
}

func NewLogger(environment string) (Logger, error) {
	var base *zap.Logger
	var err error

	if environment == "production" {
		base, err = zap.NewProduction()
	} else {
		base, err = zap.NewDevelopment()
	}

	if err != nil {
		return nil, err
	}

	return &logger{
		SugaredLogger: base.Sugar(),
		base:          base,
	}, nil
}

func NewZapLogger(environment string) (*zap.Logger, error) {
	if environment == "production" {
		return zap.NewProduction()
	}
	return zap.NewDevelopment()
}

func NewZapLoggerWithConfig(logCfg config.LogConfig, environment string) (*zap.Logger, error) {
	if environment == "development" {
		return zap.NewDevelopment()
	}

	dir := filepath.Dir(logCfg.FilePath)
	prefix := strings.TrimSuffix(filepath.Base(logCfg.FilePath), filepath.Ext(logCfg.FilePath))

	dfw, err := newDailyFileWriter(dir, prefix, logCfg.Compress, logCfg.MaxAgeDays, logCfg.RetainAll)
	if err != nil {
		return nil, fmt.Errorf("failed to create daily file writer: %w", err)
	}

	fileWriter := zapcore.AddSync(dfw)
	stdoutWriter := zapcore.AddSync(os.Stdout)

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	core := zapcore.NewTee(
		zapcore.NewCore(zapcore.NewJSONEncoder(encoderCfg), fileWriter, zap.InfoLevel),
		zapcore.NewCore(zapcore.NewJSONEncoder(encoderCfg), stdoutWriter, zap.InfoLevel),
	)

	return zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel)), nil
}

func (l logger) With(args ...interface{}) Logger {
	return &logger{
		SugaredLogger: l.SugaredLogger.With(args...),
		base:          l.base,
	}
}

func (l logger) WithFields(fields map[string]interface{}) Logger {
	args := make([]interface{}, 0, len(fields)*2)
	for k, v := range fields {
		args = append(args, k, v)
	}
	return l.With(args...)
}

func (l logger) WithError(err error) Logger {
	return l.With("error", err)
}

func (l logger) WithRequestID(requestID string) Logger {
	return l.With("request_id", requestID)
}

func (l logger) WithTenant(tenantID, tenantSlug string) Logger {
	return l.With("tenant_id", tenantID, "tenant_slug", tenantSlug)
}

func (l logger) WithUser(userID string) Logger {
	return l.With("user_id", userID)
}

func (l logger) Sync() error {
	return l.base.Sync()
}

func (l logger) Named(name string) Logger {
	return &logger{
		SugaredLogger: l.base.Named(name).Sugar(),
		base:          l.base.Named(name),
	}
}

func (l logger) Debug(msg string, keysAndValues ...interface{}) {
	l.SugaredLogger.Debugw(msg, keysAndValues...)
}

func (l logger) Info(msg string, keysAndValues ...interface{}) {
	l.SugaredLogger.Infow(msg, keysAndValues...)
}

func (l logger) Warn(msg string, keysAndValues ...interface{}) {
	l.SugaredLogger.Warnw(msg, keysAndValues...)
}

func (l logger) Error(msg string, keysAndValues ...interface{}) {
	l.SugaredLogger.Errorw(msg, keysAndValues...)
}

func (l logger) Fatal(msg string, keysAndValues ...interface{}) {
	l.SugaredLogger.Fatalw(msg, keysAndValues...)
}

func (l logger) Panic(msg string, keysAndValues ...interface{}) {
	l.SugaredLogger.Panicw(msg, keysAndValues...)
}
