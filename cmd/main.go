package main

import (
	"amr-data-bridge/config"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.Load()

	// Create root context with cancellation on interrupt
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Start HTTP server
	if err := http.StartServer(ctx, cfg); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
