-- +goose Up
CREATE TABLE urls (
  id UUID PRIMARY KEY,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  hits INTEGER NOT NULL,
  original_domain TEXT NOT NULL,
  short_domain TEXT NOT NULL
);

-- +goose Down
DROP TABLE urls;
