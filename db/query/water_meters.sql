-- name: GetActiveWaterMeters :many
SELECT * FROM public."waterMeters"
WHERE "isActive" = true
ORDER BY "id" DESC;
