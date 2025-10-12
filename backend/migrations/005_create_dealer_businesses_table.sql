-- Создание таблицы побочного бизнеса дилеров
CREATE TABLE IF NOT EXISTS dealer_businesses (
    id BIGSERIAL PRIMARY KEY,
    dealer_id BIGINT NOT NULL REFERENCES dealers(dealer_id) ON DELETE CASCADE,
    business_type VARCHAR(100) NOT NULL,
    created_at TIMESTAMP NOT NULL
);

CREATE INDEX idx_dealer_businesses_dealer_id ON dealer_businesses(dealer_id);

