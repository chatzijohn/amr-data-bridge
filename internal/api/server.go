package api

import (
	"amr-data-bridge/config"
	"amr-data-bridge/internal/api/router"
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

func StartServer(ctx context.Context, cfg *config.ServerConfig) error {
	addr := fmt.Sprintf("%s:%s", cfg.HOST, cfg.PORT)

	srv := &http.Server{
		Addr:    addr,
		Handler: router.SetupRouter(),
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
