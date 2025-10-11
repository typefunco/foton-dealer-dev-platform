-- Создание таблицы связи дилеров и брендов
CREATE TABLE IF NOT EXISTS dealer_brands (
    id BIGSERIAL PRIMARY KEY,
    dealer_id BIGINT NOT NULL REFERENCES dealers(id) ON DELETE CASCADE,
    brand_name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP NOT NULL
);

CREATE INDEX idx_dealer_brands_dealer_id ON dealer_brands(dealer_id);

