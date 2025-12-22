package api

import (
	"amr-data-bridge/config"
	"amr-data-bridge/internal/transport/http/router"
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Start(ctx context.Context, cfg *config.ServerConfig, pool *pgxpool.Pool, prefs *config.Preferences, metricsHandler http.Handler) error {
	addr := fmt.Sprintf("%s:%s", cfg.HOST, cfg.PORT)

	srv := &http.Server{
		Addr:         addr,
		Handler:      router.New(pool, prefs, metricsHandler),
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