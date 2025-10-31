package service

import (
	"amr-data-bridge/internal/db"
	"amr-data-bridge/internal/dto"
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5"
)

type WaterSupplyStore interface {
	// GetWaterSupplys(ctx context.Context, arg db.GetWaterSupplysParams) ([]db.GetWaterSupplysRow, error)
	GetWaterSupplyByNumber(ctx context.Context, q *db.Queries, supplyNumber string) (db.WaterSupply, error)
	InsertWaterSupply(ctx context.Context, q *db.Queries, arg db.InsertWaterSupplyParams) (db.WaterSupply, error)
	UpdateWaterSupply(ctx context.Context, q *db.Queries, arg db.UpdateWaterSupplyParams) error
	BeginTx(ctx context.Context) (pgx.Tx, error)
	WithTx(tx pgx.Tx) *db.Queries
}

// WaterSupplyService provides business logic for water supplies.
type WaterSupplyService struct {
	store WaterSupplyStore
	// prefs *config.Preferences
}

// NewWaterSupplyService creates a new WaterSupplyService.
func NewWaterSupplyService(store WaterSupplyStore) *WaterSupplyService {
	return &WaterSupplyService{
		store: store,
		// prefs: prefs,
	}
}

// ImportWaterSupplies inserts or updates multiple water supplies.
// Each row is validated in the handler layer before reaching here.
func (s *WaterSupplyService) ImportWaterSupplies(ctx context.Context, req []dto.WaterSupplyRequest) ([]dto.WaterSupplyResponse, error) {
	tx, err := s.store.BeginTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx) // safe rollback; commit will override if successful
	}()

	q := s.store.WithTx(tx)

	var responses []dto.WaterSupplyResponse

	for _, row := range req {
		// Check if water supply already exists
		existing, err := s.store.GetWaterSupplyByNumber(ctx, q, row.SupplyNumber)
		if err != nil {
			log.Printf("supply %s not found, inserting new one", row.SupplyNumber)

			// Insert new
			insertArg := db.InsertWaterSupplyParams{
				SupplyNumber:           row.SupplyNumber,
				Longitude:              row.Longitude,
				Latitude:               row.Latitude,
				WaterMeterSerialNumber: pgtype.Text{String: row.SerialNumber, Valid: true},
			}
			ws, err := s.store.InsertWaterSupply(ctx, q, insertArg)
			if err != nil {
				return nil, fmt.Errorf("insert supply %s: %w", row.SupplyNumber, err)
			}
			responses = append(responses, dto.WaterSupplyResponse{
				ID:           int64(ws.ID),
				SupplyNumber: ws.SupplyNumber,
				Latitude:     row.Latitude,
				Longitude:    row.Longitude,
				SerialNumber: ws.WaterMeterSerialNumber.String,
				CreatedAt:    ws.CreatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
				UpdatedAt:    ws.UpdatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
			})
			continue
		}

		// Update existing record
		updateArg := db.UpdateWaterSupplyParams{
			Longitude:              row.Longitude,
			Latitude:               row.Latitude,
			WaterMeterSerialNumber: pgtype.Text{String: row.SerialNumber, Valid: true},
			SupplyNumber:           existing.SupplyNumber,
		}
		err = s.store.UpdateWaterSupply(ctx, q, updateArg)
		if err != nil {
			return nil, fmt.Errorf("update supply %s: %w", row.SupplyNumber, err)
		}

		responses = append(responses, dto.WaterSupplyResponse{
			ID:           int64(existing.ID),
			SupplyNumber: existing.SupplyNumber,
			Latitude:     row.Latitude,
			Longitude:    row.Longitude,
			SerialNumber: existing.WaterMeterSerialNumber.String,
			CreatedAt:    existing.CreatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:    existing.UpdatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return responses, nil
}
