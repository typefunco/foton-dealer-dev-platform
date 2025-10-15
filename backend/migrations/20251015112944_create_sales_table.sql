-- +goose Up
CREATE TABLE IF NOT EXISTS sales (
    id BIGSERIAL PRIMARY KEY,
    dealer_id INTEGER REFERENCES dealers(id) ON DELETE CASCADE,
    quarter VARCHAR(2) NOT NULL CHECK (quarter IN ('Q1', 'Q2', 'Q3', 'Q4')),
    year INTEGER NOT NULL,
    
    -- Sales target
    sales_target VARCHAR(50),
    
    -- Stock (остатки на складе)
    stock_hdt INTEGER DEFAULT 0,
    stock_mdt INTEGER DEFAULT 0,
    stock_ldt INTEGER DEFAULT 0,
    
    -- Buyout (выкуп)
    buyout_hdt INTEGER DEFAULT 0,
    buyout_mdt INTEGER DEFAULT 0,
    buyout_ldt INTEGER DEFAULT 0,
    
    -- Sales personnel
    foton_salesmen INTEGER DEFAULT 0,
    
    -- Service contracts
    service_contracts_sales INTEGER DEFAULT 0,
    
    -- Training
    sales_trainings BOOLEAN DEFAULT FALSE,
    
    -- Sales decision
    sales_decision TEXT,
    
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    
    UNIQUE(dealer_id, quarter, year)
);

CREATE INDEX idx_sales_dealer_period ON sales(dealer_id, quarter, year);

-- +goose Down
DROP TABLE IF EXISTS sales;