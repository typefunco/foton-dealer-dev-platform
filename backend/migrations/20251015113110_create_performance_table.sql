-- +goose Up
CREATE TABLE IF NOT EXISTS performance (
    id BIGSERIAL PRIMARY KEY,
    dealer_id INTEGER REFERENCES dealers(id) ON DELETE CASCADE,
    quarter VARCHAR(2) NOT NULL CHECK (quarter IN ('Q1', 'Q2', 'Q3', 'Q4')),
    year INTEGER NOT NULL,
    sales_revenue_rub DECIMAL(15,2),
    sales_profit_rub DECIMAL(15,2),
    sales_profit_percent DECIMAL(5,2),
    sales_margin_percent DECIMAL(5,2),
    after_sales_revenue_rub DECIMAL(15,2),
    after_sales_profit_rub DECIMAL(15,2),
    after_sales_margin_percent DECIMAL(5,2),
    marketing_investment DECIMAL(15,2),
    foton_rank INTEGER,
    performance_decision TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(dealer_id, quarter, year)
);

-- +goose Down
DROP TABLE IF EXISTS performance;