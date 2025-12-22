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
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// Setup structured logging (JSON is better for production, Text for local dev)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cfg := config.Load()

	// Create root context with cancellation on interrupt
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	var pool *pgxpool.Pool
	pool, err := db.NewPGPool(ctx, &cfg.DB)
	if err != nil {
		logger.Error("Unable to connect to DB", "error", err)
		os.Exit(1)
	}
	defer pool.Close()
	logger.Info("Connected to Database", "pool_config", cfg.DB)

	var metricsHandler http.Handler
	if cfg.TELEMETRY {
		metricsHandler = metrics.Init()
		logger.Info("Telemetry enabled")
	}

	// Load User Defined Preferences
	// Eg. export fields
	prefs, err := config.LoadPreferences(cfg.SERVER.PREFERENCES)
	if err != nil {
		logger.Error("Failed to load preferences", "path", cfg.SERVER.PREFERENCES, "error", err)
		os.Exit(1)
	}

	// Initialize the server (but don't block yet)
	srv := httpServer.New(ctx, &cfg.SERVER, pool, prefs, metricsHandler)

	// Run server in a goroutine so main can listen for signals
	go func() {
		logger.Info("Starting HTTP Server", "port", cfg.SERVER.PORT)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("HTTP server failed", "error", err)
			stop() // Stop the context if server crashes
		}
	}()

	// Block until we receive a signal (ctx.Done)
	<-ctx.Done()
	logger.Info("Shutdown signal received")

	// Create a separate context for the shutdown timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("Server forced to shutdown", "error", err)
	} else {
		logger.Info("Server exited gracefully")
	}
}
