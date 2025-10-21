-- name: GetWaterMeters :many
-- Optional filters:
--  - limit: int (nil = unlimited)
--  - active: boolean (nil = all)

SELECT
    wm."devEUI",
    wm."serialNumber",
    wm."brandName",
    wm."ltPerPulse",
    wm."isActive",
    wm."alarmStatus",
    wm."noFlow",
    wm."currentReading",
    wm."lastSeen",
    ws."supplyNumber"
FROM public."waterMeters" AS wm
LEFT JOIN public."waterSupplies" AS ws
    ON wm."devEUI" = ws."waterMeterDevEUI"
WHERE (
    sqlc.narg(active)::boolean IS NULL
    OR wm."isActive" = sqlc.arg(active)::boolean
)
ORDER BY wm."lastSeen" DESC NULLS LAST
LIMIT $1;

-- name: GetWaterMeterBySerial :one
SELECT * FROM public."waterMeters"
WHERE "serialNumber" = $1
LIMIT 1;

-- name: UpdateWaterMeterActiveStatus :exec
UPDATE public."waterMeters"
SET "isActive" = $2,
    "updatedAt" = NOW()
WHERE "serialNumber" = $1;