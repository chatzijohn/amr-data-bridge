package handler

import (
	"amr-data-bridge/internal/db"
	"encoding/json"
	"net/http"
)

type WaterMeterHandler struct {
	q *db.Queries
}

func NewWaterMeterHandler(q *db.Queries) *WaterMeterHandler {
	return &WaterMeterHandler{q: q}
}

func (h *WaterMeterHandler) GetActiveWaterMeters(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	meters, err := h.q.GetActiveWaterMeters(ctx)
	if err != nil {
		http.Error(w, "failed to fetch water meters", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(meters)
}
