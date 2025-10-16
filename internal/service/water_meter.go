package service

import (
	"amr-data-bridge/internal/db"
	"amr-data-bridge/internal/dto"
	"amr-data-bridge/internal/export"
	"amr-data-bridge/internal/mapper"
	"context"
	"log"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
)

type WaterMeterStore interface {
	GetWaterMeters(ctx context.Context, arg db.GetWaterMetersParams) ([]db.WaterMeter, error)
}

type WaterMeterService struct {
	store WaterMeterStore
}

func NewWaterMeterService(store WaterMeterStore) *WaterMeterService {
	return &WaterMeterService{store: store}
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

	response := mapper.WaterMetersToDTO(meters)

	// Normalize export type — not passed to repo
	exportType := strings.ToLower(strings.TrimSpace(req.Type))

	switch exportType {
	default:
		return response, nil
	case "csv":
		return export.ToCSV(response), nil

		// case "xlsx":
		// 	return exportToXLSX(meters)
	}

}
