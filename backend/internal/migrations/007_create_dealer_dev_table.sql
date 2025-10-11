-- Создание таблицы развития дилеров (Dealer Development)
CREATE TABLE IF NOT EXISTS dealer_dev (
    id BIGSERIAL PRIMARY KEY,
    dealer_id BIGINT NOT NULL REFERENCES dealers(id) ON DELETE CASCADE,
    quarter VARCHAR(10) NOT NULL,
    year INT NOT NULL,
    check_list_score SMALLINT NOT NULL,
    dealer_ship_class VARCHAR(1) NOT NULL,
    branding BOOLEAN NOT NULL,
    marketing_investments BIGINT NOT NULL,
    dealer_dev_recommendation VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    UNIQUE(dealer_id, quarter, year)
);

CREATE INDEX idx_dealer_dev_dealer_id ON dealer_dev(dealer_id);
CREATE INDEX idx_dealer_dev_quarter_year ON dealer_dev(quarter, year);
CREATE INDEX idx_dealer_dev_class ON dealer_dev(dealer_ship_class);

