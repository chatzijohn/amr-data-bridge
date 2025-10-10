package service

import (
	"amr-data-bridge/internal/db"
	"context"
)

type WaterMeterStore interface {
	GetWaterMeters(ctx context.Context, arg db.GetWaterMetersParams) ([]db.WaterMeter, error)
}

type WaterMeterService struct {
	store WaterMeterStore
}

func NewWaterMeterService(store WaterMeterStore) *WaterMeterService {
	return &WaterMeterService{store: store}
}

// GetWaterMeters returns all active water meters.
// any business logic, validation, filtering, or caching belongs here.
func (s *WaterMeterService) GetWaterMeters(ctx context.Context, arg db.GetWaterMetersParams) ([]db.WaterMeter, error) {
	return s.store.GetWaterMeters(ctx, arg)
}
