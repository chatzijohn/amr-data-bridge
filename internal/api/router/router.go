package router

import (
	"amr-data-bridge/internal/api/handler"
	v1 "amr-data-bridge/internal/api/router/v1"
	"amr-data-bridge/internal/db"
	"net/http"
)

func SetupRouter(queries *db.Queries) http.Handler {
	mux := http.NewServeMux()

	h := handler.New(queries)
	mux.HandleFunc("/health", handler.HealthCheck)

	v1.RegisterWatermeterRoutes(mux, h)
	return mux
}
