package repository

import (
	"amr-data-bridge/internal/db"
	"context"
	"fmt"
)

// WaterMeterRepository provides access to water meter data.
type WaterMeterRepository struct {
	queries *db.Queries
}

func NewWaterMeterRepository(q *db.Queries) *WaterMeterRepository {
	return &WaterMeterRepository{queries: q}
}

// GetWaterMeters returns all active water meters.
// any business logic, validation, filtering, or caching belongs here.
func (r *WaterMeterRepository) GetWaterMeters(ctx context.Context, arg db.GetWaterMetersParams) ([]db.WaterMeter, error) {
	meters, err := r.queries.GetWaterMeters(ctx, arg)
	if err != nil {
		return nil, fmt.Errorf("get water meters: %w", err)
	}

	return meters, nil
}
