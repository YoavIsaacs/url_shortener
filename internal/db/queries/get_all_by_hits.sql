-- name: GetAllHitsDesc :many
SELECT original_domain, hits FROM urls
  ORDER BY hits DESC;
