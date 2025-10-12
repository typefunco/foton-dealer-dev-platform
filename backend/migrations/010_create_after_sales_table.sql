-- ============================================
-- ТАБЛИЦА: aftersales
-- ============================================
CREATE TABLE IF NOT EXISTS aftersales (
    id SERIAL PRIMARY KEY,
    dealer_id INTEGER REFERENCES dealers(dealer_id) ON DELETE CASCADE,
    period DATE NOT NULL,
    
    -- Stock metrics
    recommended_stock_pct DECIMAL(5,2),
    warranty_stock_pct DECIMAL(5,2),
    
    -- Labor hours
    foton_labor_hours_pct DECIMAL(5,2),
    warranty_hours DECIMAL(10,2),
    service_contracts_hours DECIMAL(10,2),
    
    -- Training
    as_trainings VARCHAR(3) CHECK (as_trainings IN ('Y', 'N', 'Yes', 'No')),
    
    -- Spare parts sales (revenue)
    spare_parts_sales_q DECIMAL(15,2),  -- За квартал
    spare_parts_sales_ytd_pct DECIMAL(5,2),  -- YTD динамика %
    
    -- Recommendation (из Excel)
    as_recommendation TEXT,
    
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    
    UNIQUE(dealer_id, period)
);

CREATE INDEX idx_aftersales_dealer_period ON aftersales(dealer_id, period);

