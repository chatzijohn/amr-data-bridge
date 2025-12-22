package router

import (
	"amr-data-bridge/config"
	"amr-data-bridge/internal/service"
	"amr-data-bridge/internal/transport/http/handler"
	v1 "amr-data-bridge/internal/transport/http/router/v1"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

// New creates and configures the main HTTP router.
// It now accepts preferences so that handlers and services can respect user settings.
func New(pool *pgxpool.Pool, prefs *config.Preferences, metricsHandler http.Handler) http.Handler {
	mux := http.NewServeMux()

	// Initialize service layer with preferences
	svcs := service.New(pool, prefs)

	// Initialize handlers with preferences
	h := handler.New(svcs, prefs)

	// Health check route
	mux.HandleFunc("/health", handler.HealthCheck)

	// Optional Prometheus metrics
	if metricsHandler != nil {
		mux.Handle("/metrics", metricsHandler)
	}

	// Register versioned routes (v1)
	v1.RegisterWatermeterRoutes(mux, h)
	v1.RegisterWaterSupplyRoutes(mux, h)

	return mux
}