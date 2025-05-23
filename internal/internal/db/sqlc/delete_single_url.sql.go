// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: delete_single_url.sql

package sqlc

import (
	"context"
	"database/sql"
)

const deleteSingleURL = `-- name: DeleteSingleURL :execresult
DELETE FROM urls
  WHERE short_domain = $1
`

func (q *Queries) DeleteSingleURL(ctx context.Context, shortDomain string) (sql.Result, error) {
	return q.db.ExecContext(ctx, deleteSingleURL, shortDomain)
}
