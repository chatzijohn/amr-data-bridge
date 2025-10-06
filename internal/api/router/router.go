package router

import (
	"amr-data-bridge/internal/api/handler"
	"amr-data-bridge/internal/db"
	"net/http"
)

func SetupRouter(queries *db.Queries) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", handler.HealthCheck)

	return mux
}
