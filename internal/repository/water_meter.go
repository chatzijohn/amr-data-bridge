package repository

import (
	"amr-data-bridge/internal/db"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type WaterMeterRepository interface {
	GetWaterMeters(ctx context.Context, arg db.GetWaterMetersParams) ([]db.GetWaterMetersRow, error)
	GetWaterMeterBySerial(ctx context.Context, serialNumber string) (db.WaterMeter, error)
	UpdateWaterMeterActiveStatus(ctx context.Context, arg db.UpdateWaterMeterActiveStatusParams) error
}

type waterMeterRepository struct {
	q    db.Querier
	pool *pgxpool.Pool
}

func (r *waterMeterRepository) GetWaterMeters(ctx context.Context, arg db.GetWaterMetersParams) ([]db.GetWaterMetersRow, error) {
	meters, err := r.q.GetWaterMeters(ctx, arg)
	if err != nil {
		return nil, fmt.Errorf("repository: failed to get water meters: %w", err)
	}
	return meters, nil
}

func (r *waterMeterRepository) GetWaterMeterBySerial(ctx context.Context, serialNumber string) (db.WaterMeter, error) {
	meter, err := r.q.GetWaterMeterBySerial(ctx, serialNumber)
	if err != nil {
		return db.WaterMeter{}, fmt.Errorf("repository: failed to get water meter by serial number: %w", err)
	}
	return meter, nil
}

func (r *waterMeterRepository) UpdateWaterMeterActiveStatus(ctx context.Context, arg db.UpdateWaterMeterActiveStatusParams) error {
	err := r.q.UpdateWaterMeterActiveStatus(ctx, arg)
	if err != nil {
		return fmt.Errorf("repository: failed to update water meter active status: %w", err)
	}
	return nil
}
