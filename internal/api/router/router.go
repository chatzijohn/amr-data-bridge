package router

import (
	"amr-data-bridge/internal/api/handler"
	"net/http"
)

func SetupRouter() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", handler.HealthCheck)

	return mux
}
