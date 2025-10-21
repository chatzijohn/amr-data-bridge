-- name: GetWaterSupplyByNumber :one
SELECT * FROM "waterSupplies" WHERE "supplyNumber" = $1 LIMIT 1;

-- name: InsertWaterSupply :one
INSERT INTO "waterSupplies" ("supplyNumber", geometry, "waterMeterDevEUI", "createdAt", "updatedAt")
VALUES ($1, ST_SetSRID(ST_MakePoint($2, $3), 4326), $4, NOW(), NOW())
RETURNING *;

-- name: UpdateWaterSupply :one
UPDATE "waterSupplies"
SET geometry = ST_SetSRID(ST_MakePoint($2, $3), 4326),
    "waterMeterDevEUI" = $4,
    "updatedAt" = NOW()
WHERE "supplyNumber" = $1
RETURNING *;
