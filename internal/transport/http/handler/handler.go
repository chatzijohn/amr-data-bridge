package handler

import (
	"amr-data-bridge/config"
	"amr-data-bridge/internal/service"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type Handlers struct {
	WaterMeter  *WaterMeterHandler
	WaterSupply *WaterSupplyHandler
}

func New(svc *service.Services, prefs *config.Preferences) *Handlers {
	v := validator.New()

	return &Handlers{
		WaterMeter:  NewWaterMeterHandler(svc.WaterMeter, prefs, v),
		WaterSupply: NewWaterSupplyHandler(svc.WaterSupply, prefs, v),
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
