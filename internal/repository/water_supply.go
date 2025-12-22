package repository

import (
	"amr-data-bridge/internal/db"
	"amr-data-bridge/internal/dto"
	"context"
	"errors"
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
			var waterMeterDevEUI pgtype.Text
			// If a serial number is provided, look up the water meter to get its devEUI.
			// This is the translation step.
			if row.SerialNumber != "" {
				meter, err := q.GetWaterMeterBySerial(ctx, row.SerialNumber)
				if err != nil {
					if errors.Is(err, pgx.ErrNoRows) {
						// Return a specific, user-friendly error if the meter isn't found.
						return fmt.Errorf("water meter with serial number '%s' not found for supply '%s'", row.SerialNumber, row.SupplyNumber)
					}
					// For any other error, wrap and return it.
					return fmt.Errorf("get water meter by serial '%s': %w", row.SerialNumber, err)
				}
				waterMeterDevEUI = pgtype.Text{String: meter.DevEUI, Valid: true}
			}

			// Check if water supply already exists to decide between INSERT and UPDATE.
			_, err := q.GetWaterSupplyByNumber(ctx, row.SupplyNumber)
			if err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					// --- INSERT new water supply ---
					insertArg := db.InsertWaterSupplyParams{
						SupplyNumber:     row.SupplyNumber,
						Longitude:        row.Longitude,
						Latitude:         row.Latitude,
						WaterMeterDevEui: waterMeterDevEUI,
					}
					ws, err := q.InsertWaterSupply(ctx, insertArg)
					if err != nil {
						return fmt.Errorf("insert supply %s: %w", row.SupplyNumber, err)
					}
					responses = append(responses, dto.WaterSupplyResponse{
						ID:           int64(ws.ID),
						SupplyNumber: ws.SupplyNumber,
						Latitude:     row.Latitude,
						Longitude:    row.Longitude,
						SerialNumber: row.SerialNumber, // Return the serial number from the request for consistency.
						CreatedAt:    ws.CreatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
						UpdatedAt:    ws.UpdatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
					})
				} else {
					// A different error occurred while checking for the supply.
					return fmt.Errorf("get supply %s: %w", row.SupplyNumber, err)
				}
			} else {
				// --- UPDATE existing water supply ---
				updateArg := db.UpdateWaterSupplyParams{
					Longitude:        row.Longitude,
					Latitude:         row.Latitude,
					WaterMeterDevEui: waterMeterDevEUI,
					SupplyNumber:     row.SupplyNumber,
				}
				if err := q.UpdateWaterSupply(ctx, updateArg); err != nil {
					return fmt.Errorf("update supply %s: %w", row.SupplyNumber, err)
				}

				// We need to fetch the updated record to get the new UpdatedAt timestamp.
				updated, err := q.GetWaterSupplyByNumber(ctx, row.SupplyNumber)
				if err != nil {
					return fmt.Errorf("get updated supply %s: %w", row.SupplyNumber, err)
				}
				responses = append(responses, dto.WaterSupplyResponse{
					ID:           int64(updated.ID),
					SupplyNumber: updated.SupplyNumber,
					Latitude:     row.Latitude,
					Longitude:    row.Longitude,
					SerialNumber: row.SerialNumber, // Return the serial number from the request.
					CreatedAt:    updated.CreatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
					UpdatedAt:    updated.UpdatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
				})
			}
		}
		return nil
	})

	if txErr != nil {
		return nil, txErr
	}

	return responses, nil
}
