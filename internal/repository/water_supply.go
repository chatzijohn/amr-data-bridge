package repository

import (
	"amr-data-bridge/internal/db"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

// WaterSupplyRepository provides access to water supply data.
type WaterSupplyRepository struct {
	q *db.Queries
}

// NewWaterSupplyRepository creates a new WaterSupplyRepository.
func NewWaterSupplyRepository(q *db.Queries) *WaterSupplyRepository {
	return &WaterSupplyRepository{q: q}
}

// BeginTx starts a new transaction and returns a *db.Queries bound to it.
func (r *WaterSupplyRepository) BeginTx(ctx context.Context) (*db.Queries, pgx.Tx, error) {
	tx, err := r.q.BeginTx(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("begin tx: %w", err)
	}
	return r.q.WithTx(tx), tx, nil
}

// GetWaterSupplyByNumber retrieves a supply by its number.
func (r *WaterSupplyRepository) GetWaterSupplyByNumber(ctx context.Context, q *db.Queries, supplyNumber string) (db.WaterSupply, error) {
	ws, err := q.GetWaterSupplyByNumber(ctx, supplyNumber)
	if err != nil {
		return db.WaterSupply{}, fmt.Errorf("get water supply by number: %w", err)
	}
	return ws, nil
}

// InsertWaterSupply inserts a new water supply record.
func (r *WaterSupplyRepository) InsertWaterSupply(ctx context.Context, q *db.Queries, arg db.InsertWaterSupplyParams) (db.WaterSupply, error) {
	ws, err := q.InsertWaterSupply(ctx, arg)
	if err != nil {
		return db.WaterSupply{}, fmt.Errorf("insert water supply: %w", err)
	}
	return ws, nil
}

// UpdateWaterSupply updates an existing water supply record.
func (r *WaterSupplyRepository) UpdateWaterSupply(ctx context.Context, q *db.Queries, arg db.UpdateWaterSupplyParams) error {
	if err := q.UpdateWaterSupply(ctx, arg); err != nil {
		return fmt.Errorf("update water supply: %w", err)
	}
	return nil
}
