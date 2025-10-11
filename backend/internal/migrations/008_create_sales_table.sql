-- Создание таблицы продаж (Sales)
CREATE TABLE IF NOT EXISTS sales (
    id BIGSERIAL PRIMARY KEY,
    dealer_id BIGINT NOT NULL REFERENCES dealers(id) ON DELETE CASCADE,
    quarter VARCHAR(10) NOT NULL,
    year INT NOT NULL,
    sales_target VARCHAR(50) NOT NULL,
    stock_hdt SMALLINT NOT NULL,
    stock_mdt SMALLINT NOT NULL,
    stock_ldt SMALLINT NOT NULL,
    buyout_hdt SMALLINT NOT NULL,
    buyout_mdt SMALLINT NOT NULL,
    buyout_ldt SMALLINT NOT NULL,
    foton_salesmen SMALLINT NOT NULL,
    sales_trainings BOOLEAN NOT NULL,
    service_contracts_sales SMALLINT NOT NULL,
    sales_decision VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    UNIQUE(dealer_id, quarter, year)
);

CREATE INDEX idx_sales_dealer_id ON sales(dealer_id);
CREATE INDEX idx_sales_quarter_year ON sales(quarter, year);

