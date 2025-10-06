package db

import (
	"amr-data-bridge/config"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPGPool(ctx context.Context, cfg *config.DBConfig) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		cfg.USER,
		cfg.PASSWORD,
		cfg.HOST,
		cfg.PORT,
		cfg.DB)

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to create pgx pool: %w", err)
	}

	// Test connection with Ping
	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}
	return pool, nil
}
