package v1

import (
	"amr-data-bridge/internal/transport/http/handler"
	"amr-data-bridge/internal/transport/http/middleware"
	"net/http"
)

func RegisterWatermeterRoutes(mux *http.ServeMux, h *handler.Handlers) {
	mux.HandleFunc("/api/v1/watermeters", middleware.HandleErrors(
		middleware.OnlyAllow(http.MethodGet, h.WaterMeter.GetWaterMeters),
	))

}
