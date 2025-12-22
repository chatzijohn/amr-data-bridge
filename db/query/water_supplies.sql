-- name: GetWaterSupplyByNumber :one
SELECT id, "supplyNumber", geometry, "waterMeterDevEUI", "createdAt", "updatedAt" FROM public."waterSupplies"
WHERE "supplyNumber" = @supply_number
LIMIT 1;

-- name: InsertWaterSupply :one
INSERT INTO public."waterSupplies" (
    "supplyNumber",
    geometry,
    "waterMeterDevEUI"
)
VALUES (
    @supply_number,
    ST_SetSRID(ST_MakePoint(@longitude, @latitude), 4326),
    @water_meter_dev_eui
)
RETURNING id, "supplyNumber", geometry, "waterMeterDevEUI", "createdAt", "updatedAt";

-- name: UpdateWaterSupply :exec
UPDATE public."waterSupplies"
SET
    geometry = ST_SetSRID(ST_MakePoint(@longitude, @latitude), 4326),
    "waterMeterDevEUI" = @water_meter_dev_eui
WHERE "supplyNumber" = @supply_number;
