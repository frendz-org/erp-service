package infrastructure

import (
	"fmt"
	"iam-service/config"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgres(cfg config.PostgresConfig, logger *zap.Logger) (*gorm.DB, error) {
	dsn := cfg.Platform.GetDSN()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	sqlDB.SetMaxOpenConns(cfg.Platform.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.Platform.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.Platform.ConnMaxLifetime)

	logger.Info("Successfully connected to Postgres",
		zap.Int("max_open_conns", cfg.Platform.MaxOpenConns),
		zap.Int("max_idle_conns", cfg.Platform.MaxIdleConns),
		zap.Duration("conn_max_lifetime", cfg.Platform.ConnMaxLifetime),
	)

	return db, nil
}
