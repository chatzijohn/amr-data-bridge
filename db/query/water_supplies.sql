-- name: GetWaterSupplyByNumber :one
SELECT id, "supplyNumber", geometry, "waterMeterSerialNumber", "currentImage", "previousImage", "createdAt", "updatedAt" FROM public."waterSupplies"
WHERE "supplyNumber" = $1
LIMIT 1;

-- name: InsertWaterSupply :one
INSERT INTO public."waterSupplies" (
    "supplyNumber",
    geometry,
    "waterMeterSerialNumber"
)
VALUES (
    @supply_number,
    ST_SetSRID(ST_MakePoint(@longitude, @latitude), 4326),
    @water_meter_serial_number
)
RETURNING id, "supplyNumber", geometry, "waterMeterSerialNumber", "currentImage", "previousImage", "createdAt", "updatedAt";

-- name: UpdateWaterSupply :exec
UPDATE public."waterSupplies"
SET
    geometry = ST_SetSRID(ST_MakePoint(@longitude, @latitude), 4326),
    "waterMeterSerialNumber" = @water_meter_serial_number
WHERE "supplyNumber" = @supply_number;
