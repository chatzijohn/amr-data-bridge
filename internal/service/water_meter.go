package service

import (
	"amr-data-bridge/internal"
	"amr-data-bridge/internal/db"
	"amr-data-bridge/internal/dto"
	"amr-data-bridge/internal/mapper"
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgtype"
)

type WaterMeterStore interface {
	GetWaterMeters(ctx context.Context, arg db.GetWaterMetersParams) ([]db.GetWaterMetersRow, error)
}

// WaterMeterService provides business logic for water meters.
type WaterMeterService struct {
	store WaterMeterStore
	prefs *internal.Preferences
}

// NewWaterMeterService creates a new WaterMeterService.
func NewWaterMeterService(store WaterMeterStore, prefs *internal.Preferences) *WaterMeterService {
	return &WaterMeterService{
		store: store,
		prefs: prefs,
	}
}

func (s *WaterMeterService) GetWaterMeters(ctx context.Context, req dto.WaterMetersRequest) ([]dto.WaterMeterResponse, error) {

	const defaultLimit = 10000

	limit := req.Limit
	if limit <= 0 {
		limit = defaultLimit
	}

	var active pgtype.Bool
	if req.Active != nil {
		active = pgtype.Bool{
			Bool:  *req.Active,
			Valid: true,
		}
	} else {
		active = pgtype.Bool{Valid: false}
	}

	params := db.GetWaterMetersParams{
		Limit:  int32(limit),
		Active: active,
	}

	// Call repository
	meters, err := s.store.GetWaterMeters(ctx, params)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	// Map DB results to DTOs with preferences applied
	response := mapper.WaterMetersToDTO(meters, s.prefs)

	return response, nil

}
