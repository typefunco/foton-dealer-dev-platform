-- ============================================
-- ТАБЛИЦА: sales
-- ============================================
CREATE TABLE IF NOT EXISTS sales (
    id SERIAL PRIMARY KEY,
    dealer_id INTEGER REFERENCES dealers(dealer_id) ON DELETE CASCADE,
    period DATE NOT NULL,
    
    -- Stock (остатки на складе)
    stock_hdt INTEGER DEFAULT 0,
    stock_mdt INTEGER DEFAULT 0,
    stock_ldt INTEGER DEFAULT 0,
    
    -- Buyout (выкуп)
    buyout_hdt INTEGER DEFAULT 0,
    buyout_mdt INTEGER DEFAULT 0,
    buyout_ldt INTEGER DEFAULT 0,
    
    -- Sales targets and personnel
    foton_sales_personnel INTEGER,  -- Количество продавцов Foton
    sales_target_plan INTEGER,  -- План продаж
    sales_target_fact INTEGER,  -- Факт продаж
    
    -- Service contracts
    service_contracts_sales DECIMAL(15,2),
    
    -- Training
    sales_trainings VARCHAR(3) CHECK (sales_trainings IN ('Y', 'N', 'Yes', 'No')),
    
    -- Recommendation (из Excel)
    sales_recommendation TEXT,
    
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    
    UNIQUE(dealer_id, period)
);

CREATE INDEX idx_sales_dealer_period ON sales(dealer_id, period);

