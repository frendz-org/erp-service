package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"erp-service/config"
	iamhttp "erp-service/delivery/http"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	server := iamhttp.NewServer(cfg)

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

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("failed to shutdown server: %v", err)
	}

	log.Println("server stopped")
}
