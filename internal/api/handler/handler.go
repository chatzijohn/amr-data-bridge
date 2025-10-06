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

// HealthCheck is a simple liveness probe
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
