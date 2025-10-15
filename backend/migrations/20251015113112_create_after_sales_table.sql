-- +goose Up
CREATE TABLE IF NOT EXISTS after_sales (
    id BIGSERIAL PRIMARY KEY,
    dealer_id INTEGER REFERENCES dealers(id) ON DELETE CASCADE,
    quarter VARCHAR(2) NOT NULL CHECK (quarter IN ('Q1', 'Q2', 'Q3', 'Q4')),
    year INTEGER NOT NULL,
    recommended_stock INTEGER,
    warranty_stock INTEGER,
    foton_labor_hours INTEGER,
    service_contracts INTEGER,
    as_trainings BOOLEAN DEFAULT FALSE,
    csi VARCHAR(50),
    foton_warranty_hours INTEGER,
    as_decision TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(dealer_id, quarter, year)
);

-- +goose Down
DROP TABLE IF EXISTS after_sales;