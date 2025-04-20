-- name: AddOneToHits :one
UPDATE urls
  SET hits = hits + 1
  WHERE short_domain = $2;
  RETURNING *;
