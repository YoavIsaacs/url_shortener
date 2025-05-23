// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: add_one_to_hit.sql

package sqlc

import (
	"context"
)

const addOneToHits = `-- name: AddOneToHits :one
UPDATE urls
  SET hits = hits + 1
  WHERE short_domain = $1
  RETURNING id, created_at, updated_at, hits, original_domain, short_domain, qr_code
`

func (q *Queries) AddOneToHits(ctx context.Context, shortDomain string) (Url, error) {
	row := q.db.QueryRowContext(ctx, addOneToHits, shortDomain)
	var i Url
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Hits,
		&i.OriginalDomain,
		&i.ShortDomain,
		&i.QrCode,
	)
	return i, err
}
