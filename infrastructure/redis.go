package infrastructure

import (
	"context"
	"fmt"
	"log"
	"time"

	"iam-service/config"

	"github.com/redis/go-redis/v9"
)

func NewRedis(cfg config.RedisConfig) (*redis.Client, error) {
	poolSize := cfg.PoolSize
	if poolSize <= 0 {
		poolSize = 20
	}

	minIdleConns := cfg.MinIdleConns
	if minIdleConns <= 0 {
		minIdleConns = 5
	}

	connMaxIdleTime := cfg.ConnMaxIdleTime
	if connMaxIdleTime <= 0 {
		connMaxIdleTime = 5 * time.Minute
	}

	connMaxLifetime := cfg.ConnMaxLifetime
	if connMaxLifetime <= 0 {
		connMaxLifetime = 30 * time.Minute
	}

	readTimeout := cfg.ReadTimeout
	if readTimeout <= 0 {
		readTimeout = 3 * time.Second
	}

	writeTimeout := cfg.WriteTimeout
	if writeTimeout <= 0 {
		writeTimeout = 3 * time.Second
	}

	client := redis.NewClient(&redis.Options{
		Addr:            cfg.GetAddress(),
		Password:        cfg.Password,
		DB:              cfg.DB,
		PoolSize:        poolSize,
		MinIdleConns:    minIdleConns,
		ConnMaxIdleTime: connMaxIdleTime,
		ConnMaxLifetime: connMaxLifetime,
		ReadTimeout:     readTimeout,
		WriteTimeout:    writeTimeout,
		PoolTimeout:     readTimeout + time.Second,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	log.Printf("Connected to Redis at %s (pool_size=%d, min_idle=%d)\n", cfg.GetAddress(), poolSize, minIdleConns)

	return client, nil
}
