-- name: AddQRCode :one
UPDATE urls
  SET qr_code = $1
  WHERE short_domain = $2
  RETURNING *;
