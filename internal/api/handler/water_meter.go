package handler

import (
	"amr-data-bridge/internal/api/middleware"
	"amr-data-bridge/internal/service"
	"encoding/json"
	"net/http"
)

type WaterMeterHandler struct {
	svc *service.WaterMeterService
}

func NewWaterMeterHandler(svc *service.WaterMeterService) *WaterMeterHandler {
	return &WaterMeterHandler{svc: svc}
}

// GetActiveWaterMeters godoc
// @Summary      Get active water meters
// @Description  Returns a list of all active water meters from the database.
// @Tags         water-meters
// @Produce      json
// @Success      200  {array}  db.WaterMeter
// @Failure      500  {object}  middleware.HttpError
// @Router       /water-meters/active [get]
func (h *WaterMeterHandler) GetActiveWaterMeters(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	meters, err := h.svc.GetActiveWaterMeters(ctx)
	if err != nil {
		return middleware.NewHttpError(http.StatusInternalServerError, "failed to fetch water meters")
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(meters); err != nil {
		return middleware.NewHttpError(http.StatusInternalServerError, "failed to encode response")
	}

	return nil
}
