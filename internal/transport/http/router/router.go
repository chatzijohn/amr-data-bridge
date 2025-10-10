package router

import (
	"amr-data-bridge/internal/db"
	"amr-data-bridge/internal/service"
	"amr-data-bridge/internal/transport/http/handler"
	v1 "amr-data-bridge/internal/transport/http/router/v1"
	"net/http"
)

func New(queries *db.Queries, metricsHandler http.Handler) http.Handler {
	mux := http.NewServeMux()

	svcs := service.New(queries)
	h := handler.New(svcs)

	mux.HandleFunc("/health", handler.HealthCheck)

	if metricsHandler != nil {
		// Add Prometheus /metrics
		mux.Handle("/metrics", metricsHandler)
	}

	v1.RegisterWatermeterRoutes(mux, h)
	return mux
}
