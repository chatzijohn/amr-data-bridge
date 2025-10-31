package repository

import "amr-data-bridge/internal/db"

type Repository struct {
	WaterMeter  *WaterMeterRepository
	WaterSupply *WaterSupplyRepository
}

func New(q *db.Queries) *Repository {
	return &Repository{
		WaterMeter:  NewWaterMeterRepository(q),
		WaterSupply: NewWaterSupplyRepository(q),
	}
}
