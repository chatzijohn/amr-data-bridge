package handler

import (
	"amr-data-bridge/internal/db"
	"net/http"
)

type Handlers struct {
	WaterMeter *WaterMeterHandler
}

// New initializes the main handler struct
func New(q *db.Queries) *Handlers {
	return &Handlers{
		WaterMeter: NewWaterMeterHandler(q),
	}
}

// @Summary Health check
// @Description Returns OK if the service is running
// @Tags health
// @Produce json
// @Success 200 {object} map[string]string
// @Router /healthz [get]
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}
