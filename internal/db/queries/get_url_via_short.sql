-- name: GetURLViaShortURL :one
SELECT original_domain FROM urls
  WHERE (short_domain = $1);
