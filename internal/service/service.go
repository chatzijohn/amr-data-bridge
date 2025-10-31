package service

import (
	"amr-data-bridge/config"
	"amr-data-bridge/internal/db"
)

// Services groups all service layer dependencies together.
type Services struct {
	WaterMeter  *WaterMeterService
	WaterSupply *WaterSupplyService
}

// New initializes the main Services struct with all dependencies.
// It now accepts preferences, which are shared across sub-services.
func New(q *db.Queries, prefs *config.Preferences) *Services {
	return &Services{
		WaterMeter:  NewWaterMeterService(q, prefs),
		WaterSupply: NewWaterSupplyService(NewWaterSupplyRepository(q)),
	}
}
