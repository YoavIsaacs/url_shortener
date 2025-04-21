-- name: DeleteSingleURL :execresult
DELETE FROM urls
  WHERE short_domain = $1;
