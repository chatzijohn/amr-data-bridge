package service

import (
	"amr-data-bridge/internal/db"
	"context"
)

type WaterMeterService struct {
	q *db.Queries
}

func NewWaterMeterService(q *db.Queries) *WaterMeterService {
	return &WaterMeterService{q: q}
}

// GetActiveWaterMeters returns all active water meters.
// any business logic, validation, filtering, or caching belongs here.
func (s *WaterMeterService) GetActiveWaterMeters(ctx context.Context) ([]db.WaterMeter, error) {
	return s.q.GetActiveWaterMeters(ctx)
}
