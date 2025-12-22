package service

import (
	"amr-data-bridge/internal/dto"
	"amr-data-bridge/internal/repository"
	"context"
)

type WaterSupplyService interface {
	ImportWaterSupplies(ctx context.Context, req []dto.WaterSupplyRequest) ([]dto.WaterSupplyResponse, error)
}

type waterSupplyService struct {
	repo repository.WaterSupplyRepository
}

func NewWaterSupplyService(repo repository.WaterSupplyRepository) WaterSupplyService {
	return &waterSupplyService{
		repo: repo,
	}
}

func (s *waterSupplyService) ImportWaterSupplies(ctx context.Context, req []dto.WaterSupplyRequest) ([]dto.WaterSupplyResponse, error) {
	return s.repo.ImportWaterSupplies(ctx, req)
}
