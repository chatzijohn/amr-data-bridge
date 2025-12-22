package repository

import (
	"amr-data-bridge/internal/db"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Storer interface {
	WaterMeter() WaterMeterRepository
	WaterSupply() WaterSupplyRepository
}

type store struct {
	waterMeter  *waterMeterRepository
	waterSupply *waterSupplyRepository
}

func (s *store) WaterMeter() WaterMeterRepository {
	return s.waterMeter
}

func (s *store) WaterSupply() WaterSupplyRepository {
	return s.waterSupply
}

func New(pool *pgxpool.Pool) Storer {
	q := db.New(pool)
	return &store{
		waterMeter:  &waterMeterRepository{q: q, pool: pool},
		waterSupply: &waterSupplyRepository{q: q, pool: pool},
	}
}

// execTx executes a function within a database transaction.
func execTx(ctx context.Context, pool *pgxpool.Pool, fn func(q *db.Queries) error) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	q := db.New(tx)
	err = fn(q)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}