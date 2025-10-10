// @title AMR Data Bridge API
// @version 1.0
// @description Minimal Go API with clean architecture
// @BasePath /api

package main

import (
	"amr-data-bridge/config"
	"amr-data-bridge/internal/db"
	"amr-data-bridge/internal/observability/metrics"
	httpServer "amr-data-bridge/internal/transport/http"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.Load()

	// Create root context with cancellation on interrupt
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	pool, err := db.NewPGPool(ctx, &cfg.DB)
	if err != nil {
		log.Fatalf("Unable to connect to DB: %v", err)
	}
	defer pool.Close()

	var metricsHandler http.Handler
	if cfg.TELEMETRY {
		metricsHandler = metrics.Init()
	}

	// sqlc Queries instance
	queries := db.New(pool)

	// Start HTTP server
	if err := httpServer.Start(ctx, &cfg.SERVER, queries, metricsHandler); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
