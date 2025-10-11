-- Создание таблицы производительности (Performance)
CREATE TABLE IF NOT EXISTS performance (
    id BIGSERIAL PRIMARY KEY,
    dealer_id BIGINT NOT NULL REFERENCES dealers(id) ON DELETE CASCADE,
    quarter VARCHAR(10) NOT NULL,
    year INT NOT NULL,
    sales_revenue_rub BIGINT NOT NULL,
    sales_profit_rub BIGINT NOT NULL,
    sales_profit_percent DOUBLE PRECISION NOT NULL,
    sales_margin_percent DOUBLE PRECISION NOT NULL,
    after_sales_revenue_rub BIGINT NOT NULL,
    after_sales_profit_rub BIGINT NOT NULL,
    after_sales_margin_percent DOUBLE PRECISION NOT NULL,
    marketing_investment DOUBLE PRECISION NOT NULL,
    foton_rank SMALLINT NOT NULL,
    performance_decision VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    UNIQUE(dealer_id, quarter, year)
);

CREATE INDEX idx_performance_dealer_id ON performance(dealer_id);
CREATE INDEX idx_performance_quarter_year ON performance(quarter, year);
CREATE INDEX idx_performance_rank ON performance(foton_rank);

