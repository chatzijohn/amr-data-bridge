package handler

import (
	"amr-data-bridge/internal/service"
	"net/http"
)

type Handlers struct {
	WaterMeter *WaterMeterHandler
}

// New initializes the main handler struct
func New(svc *service.Services) *Handlers {
	return &Handlers{
		WaterMeter: NewWaterMeterHandler(svc.WaterMeter),
	}
}

// @Summary Health check
// @Description Returns OK if the service is running
// @Tags health
// @Produce json
// @Success 200 {object} map[string]string
// @Router /healthz [get]
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	_ = ctx // placeholder, useful if you add checks later

	w.Write([]byte(`{"status":"ok"}`))
}
