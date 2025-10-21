package api

import (
	"amr-data-bridge/config"
	"amr-data-bridge/internal"
	"amr-data-bridge/internal/db"
	"amr-data-bridge/internal/transport/http/router"
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

func Start(ctx context.Context, cfg *config.ServerConfig, queries *db.Queries, metricsHandler http.Handler) error {
	addr := fmt.Sprintf("%s:%s", cfg.HOST, cfg.PORT)

	// Load User Defined Prefernces
	// Eg. export fields
	prefs, err := internal.LoadPreferences("preferences.yaml")
	if err != nil {
		return fmt.Errorf("failed to load preferences: %w", err)
	}

	srv := &http.Server{
		Addr:         addr,
		Handler:      router.New(queries, prefs, metricsHandler),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		log.Printf("Starting server on %s\n", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Wait for shutdown signal
	<-ctx.Done()
	log.Println("Shutting down gracefully...")

	ctxTimeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return srv.Shutdown(ctxTimeout)
}
