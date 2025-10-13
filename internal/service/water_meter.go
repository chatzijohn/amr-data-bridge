package service

import (
	"amr-data-bridge/internal/db"
	"amr-data-bridge/internal/dto"
	"context"
	"strings"
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

func (s *WaterMeterService) GetWaterMeters(ctx context.Context, req dto.GetWaterMetersRequest) ([]db.WaterMeter, error) {

	// Map DTO to DB params
	params := db.GetWaterMetersParams{
		Limit: req.Limit,
	}

	// Handle Active pointer (if nil, default false or whatever your DB expects)
	if req.Active != nil {
		params.Active = *req.Active
	}

	// Call repository
	meters, err := s.store.GetWaterMeters(ctx, params)
	if err != nil {
		return nil, err
	}

	// Normalize export type — not passed to repo
	exportType := strings.ToLower(strings.TrimSpace(req.Type))
	if exportType == "" {
		exportType = "json" // default
	}

	// 🔜 Future logic:
	// switch exportType {
	// case "csv":
	//     return exportToCSV(meters)
	// case "xlsx":
	//     return exportToXLSX(meters)
	// }

	return meters, nil
}
