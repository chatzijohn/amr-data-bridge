-- name: GetWaterMeters :many
-- Optional filters:
--  - limit: int (nil = unlimited)
--  - active: boolean (nil = all)

SELECT *
FROM public."waterMeters"
WHERE (
  sqlc.narg(active) IS NULL
  OR "isActive" = sqlc.arg(active)
)
ORDER BY "lastSeen" DESC NULLS LAST
LIMIT $1;

