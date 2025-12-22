package api

import (
	"amr-data-bridge/config"
	"amr-data-bridge/internal/transport/http/router"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func New(ctx context.Context, cfg *config.ServerConfig, pool *pgxpool.Pool, prefs *config.Preferences, metricsHandler http.Handler) *http.Server {
	addr := fmt.Sprintf("%s:%s", cfg.HOST, cfg.PORT)

	return &http.Server{
		Addr:         addr,
		Handler:      router.New(pool, prefs, metricsHandler),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
}
