-- +goose Up

CREATE TABLE webpages (
  id UUID PRIMARY KEY,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  name VARCHAR(255) NOT NULL,
  url TEXT UNIQUE NOT NULL,
  type VARCHAR(255) NOT NULL
);

-- +goose Down

DROP TABLE webpages;