package service

import (
	"amr-data-bridge/internal/db"
	"context"
)

type WaterMeterStore interface {
	GetActiveWaterMeters(ctx context.Context) ([]db.WaterMeter, error)
}

type WaterMeterService struct {
	store WaterMeterStore
}

func NewWaterMeterService(store WaterMeterStore) *WaterMeterService {
	return &WaterMeterService{store: store}
}

// GetActiveWaterMeters returns all active water meters.
// any business logic, validation, filtering, or caching belongs here.
func (s *WaterMeterService) GetActiveWaterMeters(ctx context.Context) ([]db.WaterMeter, error) {
	return s.store.GetActiveWaterMeters(ctx)
}
