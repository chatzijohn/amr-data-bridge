package db

import (
	"context"
)

type Querier interface {
	GetWaterMeterBySerial(ctx context.Context, serialnumber string) (WaterMeter, error)
	GetWaterMeters(ctx context.Context, arg GetWaterMetersParams) ([]GetWaterMetersRow, error)
	UpdateWaterMeterActiveStatus(ctx context.Context, arg UpdateWaterMeterActiveStatusParams) error

	GetWaterSupplyByNumber(ctx context.Context, supplyNumber string) (GetWaterSupplyByNumberRow, error)
	InsertWaterSupply(ctx context.Context, arg InsertWaterSupplyParams) (InsertWaterSupplyRow, error)
	UpdateWaterSupply(ctx context.Context, arg UpdateWaterSupplyParams) error
}

var _ Querier = (*Queries)(nil)
