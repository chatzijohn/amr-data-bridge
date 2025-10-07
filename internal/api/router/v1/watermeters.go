package v1

import (
	"amr-data-bridge/internal/api/handler"
	"amr-data-bridge/internal/api/middleware"
	"net/http"
)

func RegisterWatermeterRoutes(mux *http.ServeMux, h *handler.Handlers) {
	mux.HandleFunc("/api/v1/watermeters/active", middleware.HandleErrors(
		middleware.OnlyAllow(http.MethodGet, h.WaterMeter.GetActiveWaterMeters),
	))

}
