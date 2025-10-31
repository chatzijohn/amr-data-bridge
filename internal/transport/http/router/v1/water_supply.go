package v1

import (
	"amr-data-bridge/internal/transport/http/handler"
	"amr-data-bridge/internal/transport/http/middleware"
	"net/http"
)

func RegisterWaterSupplyRoutes(mux *http.ServeMux, h *handler.Handlers) {
	mux.HandleFunc("/api/v1/watersupplies", middleware.HandleErrors(
		middleware.OnlyAllow(http.MethodPost, h.WaterSupply.ImportWaterSupplies),
	))
}
