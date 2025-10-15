-- +goose Up
CREATE TABLE IF NOT EXISTS dealer_businesses (
    id BIGSERIAL PRIMARY KEY,
    dealer_id INTEGER REFERENCES dealers(id) ON DELETE CASCADE,
    business_type VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS dealer_businesses;