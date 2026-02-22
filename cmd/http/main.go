package main

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"erp-service/config"
	erphttp "erp-service/delivery/http"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	if err := runMigrations(cfg); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	server := erphttp.NewServer(cfg)

	workerCtx, workerCancel := context.WithCancel(context.Background())
	server.StartWorker(workerCtx)

	go func() {
		if err := server.Start(); err != nil {
			log.Fatalf("failed to start server: %v", err)
		}

		log.Printf("server is running on %s:%d", cfg.Server.Host, cfg.Server.Port)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down server...")
	log.Printf("Serever is starting on port %s", os.DevNull)

	workerCancel()
	server.StopWorker()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("failed to shutdown server: %v", err)
	}

	log.Println("server stopped")
}

func runMigrations(cfg *config.Config) error {
	pg := cfg.Infra.Postgres.Platform
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		url.QueryEscape(pg.User),
		url.QueryEscape(pg.Password),
		pg.Host,
		pg.Port,
		pg.Database,
		pg.SSLMode,
	)

	m, err := migrate.New("file://migration", dsn)
	if err != nil {
		return fmt.Errorf("create migrator: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("run migrations: %w", err)
	}

	version, dirty, _ := m.Version()
	log.Printf("migrations applied with version: %d, dirty: %v", version, dirty)
	return nil
}
