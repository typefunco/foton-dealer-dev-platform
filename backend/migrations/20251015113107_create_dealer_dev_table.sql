-- +goose Up
CREATE TABLE IF NOT EXISTS dealer_dev (
    id BIGSERIAL PRIMARY KEY,
    dealer_id INTEGER REFERENCES dealers(id) ON DELETE CASCADE,
    quarter VARCHAR(2) NOT NULL CHECK (quarter IN ('Q1', 'Q2', 'Q3', 'Q4')),
    year INTEGER NOT NULL,
    check_list_score INTEGER,
    dealer_ship_class VARCHAR(1) CHECK (dealer_ship_class IN ('A', 'B', 'C', 'D')),
    branding BOOLEAN DEFAULT FALSE,
    marketing_investments DECIMAL(15,2),
    dealer_dev_recommendation TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(dealer_id, quarter, year)
);

-- +goose Down
DROP TABLE IF EXISTS dealer_dev;