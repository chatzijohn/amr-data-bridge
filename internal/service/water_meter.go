package service

import (
	"amr-data-bridge/config"
	"amr-data-bridge/internal/db"
	"amr-data-bridge/internal/dto"
	"amr-data-bridge/internal/mapper"
	"amr-data-bridge/internal/repository"
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

type WaterMeterService interface {
	GetWaterMeters(ctx context.Context, req dto.WaterMetersRequest) ([]dto.WaterMeterResponse, error)
}

type waterMeterService struct {
	repo  repository.WaterMeterRepository
	prefs *config.Preferences
}

func NewWaterMeterService(repo repository.WaterMeterRepository, prefs *config.Preferences) WaterMeterService {
	return &waterMeterService{
		repo:  repo,
		prefs: prefs,
	}
}

func (s *waterMeterService) GetWaterMeters(ctx context.Context, req dto.WaterMetersRequest) ([]dto.WaterMeterResponse, error) {
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
		active = pgtype.Bool{
			Valid: false,
		}
	}

	params := db.GetWaterMetersParams{
		Limit:  int32(limit),
		Active: active,
	}

	meters, err := s.repo.GetWaterMeters(ctx, params)
	if err != nil {
		return nil, err
	}

	response := mapper.WaterMetersToDTO(meters, s.prefs)

	return response, nil
}
