-- Создание таблицы послепродажного обслуживания (After Sales)
CREATE TABLE IF NOT EXISTS after_sales (
    id BIGSERIAL PRIMARY KEY,
    dealer_id BIGINT NOT NULL REFERENCES dealers(id) ON DELETE CASCADE,
    quarter VARCHAR(10) NOT NULL,
    year INT NOT NULL,
    recommended_stock SMALLINT NOT NULL,
    warranty_stock SMALLINT NOT NULL,
    foton_labor_hours SMALLINT NOT NULL,
    service_contracts SMALLINT NOT NULL,
    as_trainings BOOLEAN NOT NULL,
    csi VARCHAR(50),
    foton_warranty_hours INT NOT NULL,
    as_decision VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    UNIQUE(dealer_id, quarter, year)
);

CREATE INDEX idx_after_sales_dealer_id ON after_sales(dealer_id);
CREATE INDEX idx_after_sales_quarter_year ON after_sales(quarter, year);

