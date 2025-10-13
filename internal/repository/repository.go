package repository

import "amr-data-bridge/internal/db"

type Repository struct {
	WaterMeter *WaterMeterRepository
}

func New(q *db.Queries) *Repository {
	return &Repository{
		WaterMeter: NewWaterMeterRepository(q),
	}
}
