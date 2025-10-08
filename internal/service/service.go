package service

import (
	"amr-data-bridge/internal/db"
)

type Services struct {
	WaterMeter *WaterMeterService
}

// New initializes the main handler struct
func New(q *db.Queries) *Services {
	return &Services{
		WaterMeter: NewWaterMeterService(q),
	}
}
