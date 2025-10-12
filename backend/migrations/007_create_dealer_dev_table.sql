-- ============================================
-- ТАБЛИЦА: dealer_development
-- ============================================
CREATE TABLE IF NOT EXISTS dealer_development (
    id SERIAL PRIMARY KEY,
    dealer_id INTEGER REFERENCES dealers(dealer_id) ON DELETE CASCADE,
    period DATE NOT NULL,
    
    -- Dealer Development метрики
    check_list_score DECIMAL(5,2),  -- Check List Score % (0-100)
    dealership_class VARCHAR(1) CHECK (dealership_class IN ('A', 'B', 'C', 'D')),
    brands TEXT[],  -- Массив брендов ["Foton", "Kamaz", "Sitrak"]
    branding VARCHAR(3) CHECK (branding IN ('Y', 'N', 'Yes', 'No')),
    marketing_investments DECIMAL(15,2),  -- Marketing Investments Rub
    by_side_businesses TEXT,  -- By-side businesses описание
    
    -- Recommendation (из Excel)
    dd_recommendation TEXT,
    
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    
    UNIQUE(dealer_id, period)
);

CREATE INDEX idx_dealer_dev_dealer_period ON dealer_development(dealer_id, period);
CREATE INDEX idx_dealer_dev_class ON dealer_development(dealership_class);

