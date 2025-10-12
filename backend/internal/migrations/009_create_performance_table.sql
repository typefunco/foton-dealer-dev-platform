-- ============================================
-- ТАБЛИЦА: performance_sales
-- ============================================
CREATE TABLE IF NOT EXISTS performance_sales (
    id SERIAL PRIMARY KEY,
    dealer_id INTEGER REFERENCES dealers(dealer_id) ON DELETE CASCADE,
    period DATE NOT NULL,
    
    -- Sales financial metrics
    quantity_sold INTEGER,  -- Количество проданных автомобилей
    sales_revenue DECIMAL(15,2),  -- Выручка (с НДС)
    sales_revenue_no_vat DECIMAL(15,2),  -- Выручка без НДС
    sales_cost DECIMAL(15,2),  -- Стоимость
    sales_margin DECIMAL(15,2),  -- Валовая прибыль (в рублях)
    sales_margin_pct DECIMAL(5,2),  -- Маржа % = (margin / revenue) * 100
    sales_profit_pct DECIMAL(5,2),  -- Рентабельность %
    
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    
    UNIQUE(dealer_id, period)
);

CREATE INDEX idx_perf_sales_dealer_period ON performance_sales(dealer_id, period);

