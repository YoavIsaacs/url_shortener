package config

import (
	"github.com/YoavIsaacs/url_shortener/internal/internal/db/sqlc"
)

type ApiConfig struct {
	database *sqlc.Queries
}
