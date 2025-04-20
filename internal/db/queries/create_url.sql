-- name: CreateNewURL :one
INSERT INTO urls (id, created_at, updated_at, hits, original_domain, short_domain)
  VALUES ($1, $2, $3, 0, $4, $5)
  RETURNING *;
