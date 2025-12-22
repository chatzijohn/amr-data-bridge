package repository

import (
	"amr-data-bridge/internal/db"
	"amr-data-bridge/internal/dto"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type WaterSupplyRepository interface {
	ImportWaterSupplies(ctx context.Context, req []dto.WaterSupplyRequest) ([]dto.WaterSupplyResponse, error)
}

type waterSupplyRepository struct {
	q    db.Querier
	pool *pgxpool.Pool
}

func (r *waterSupplyRepository) ImportWaterSupplies(ctx context.Context, req []dto.WaterSupplyRequest) ([]dto.WaterSupplyResponse, error) {
	var responses []dto.WaterSupplyResponse

	txErr := execTx(ctx, r.pool, func(q *db.Queries) error {
		for _, row := range req {
			// Check if water supply already exists
			existing, err := q.GetWaterSupplyByNumber(ctx, row.SupplyNumber)
			if err != nil {
				if err == pgx.ErrNoRows {
					// Insert new
					insertArg := db.InsertWaterSupplyParams{
						SupplyNumber:           row.SupplyNumber,
						Longitude:              row.Longitude,
						Latitude:               row.Latitude,
						WaterMeterSerialNumber: pgtype.Text{String: row.SerialNumber, Valid: row.SerialNumber != ""},
					}
					ws, err := q.InsertWaterSupply(ctx, insertArg)
					if err != nil {
						return fmt.Errorf("insert supply %s: %w", row.SupplyNumber, err)
					}
					responses = append(responses, dto.WaterSupplyResponse{
						ID:           int64(ws.ID),
						SupplyNumber: ws.SupplyNumber,
						Latitude:     row.Latitude,  // From request, as geometry isn't returned directly
						Longitude:    row.Longitude, // From request
						SerialNumber: ws.WaterMeterSerialNumber.String,
						CreatedAt:    ws.CreatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
						UpdatedAt:    ws.UpdatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
					})
					continue
				}
				return fmt.Errorf("get supply %s: %w", row.SupplyNumber, err)
			}

			// Update existing record
			updateArg := db.UpdateWaterSupplyParams{
				Longitude:              row.Longitude,
				Latitude:               row.Latitude,
				WaterMeterSerialNumber: pgtype.Text{String: row.SerialNumber, Valid: row.SerialNumber != ""},
				SupplyNumber:           existing.SupplyNumber,
			}
			err = q.UpdateWaterSupply(ctx, updateArg)
			if err != nil {
				return fmt.Errorf("update supply %s: %w", row.SupplyNumber, err)
			}

			// We need to fetch the updated record to get the new UpdatedAt timestamp
			updated, err := q.GetWaterSupplyByNumber(ctx, row.SupplyNumber)
			if err != nil {
				return fmt.Errorf("get updated supply %s: %w", row.SupplyNumber, err)
			}

			responses = append(responses, dto.WaterSupplyResponse{
				ID:           int64(updated.ID),
				SupplyNumber: updated.SupplyNumber,
				Latitude:     row.Latitude,
				Longitude:    row.Longitude,
				SerialNumber: updated.WaterMeterSerialNumber.String,
				CreatedAt:    updated.CreatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
				UpdatedAt:    updated.UpdatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
			})
		}
		return nil
	})

	if txErr != nil {
		return nil, txErr
	}

	return responses, nil
}
