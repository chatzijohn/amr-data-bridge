package service

import (
	"amr-data-bridge/config"
	"amr-data-bridge/internal/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Service interface {
	WaterMeter() WaterMeterService
	WaterSupply() WaterSupplyService
}

type service struct {
	waterMeter  WaterMeterService
	waterSupply WaterSupplyService
}

func (s *service) WaterMeter() WaterMeterService {
	return s.waterMeter
}

func (s *service) WaterSupply() WaterSupplyService {
	return s.waterSupply
}

func New(pool *pgxpool.Pool, prefs *config.Preferences) Service {
	store := repository.New(pool)
	return &service{
		waterMeter:  NewWaterMeterService(store.WaterMeter(), prefs),
		waterSupply: NewWaterSupplyService(store.WaterSupply()),
	}
}