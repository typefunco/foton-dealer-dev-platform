-- +goose Up
CREATE TABLE IF NOT EXISTS brands (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    logo_path VARCHAR(255)
);

-- +goose Down
DROP TABLE IF EXISTS brands;