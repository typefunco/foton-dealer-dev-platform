-- ============================================
-- ТАБЛИЦА: performance_aftersales
-- ============================================
CREATE TABLE IF NOT EXISTS performance_aftersales (
    id SERIAL PRIMARY KEY,
    dealer_id INTEGER REFERENCES dealers(dealer_id) ON DELETE CASCADE,
    period DATE NOT NULL,
    
    -- AfterSales financial metrics
    as_revenue DECIMAL(15,2),  -- Выручка (с НДС)
    as_revenue_no_vat DECIMAL(15,2),  -- Выручка без НДС
    as_cost DECIMAL(15,2),  -- Стоимость
    as_margin DECIMAL(15,2),  -- Валовая прибыль (в рублях)
    as_margin_pct DECIMAL(5,2),  -- Маржа % = (margin / revenue) * 100
    as_profit_pct DECIMAL(5,2),  -- Рентабельность %
    
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    
    UNIQUE(dealer_id, period)
);

CREATE INDEX idx_perf_as_dealer_period ON performance_aftersales(dealer_id, period);
