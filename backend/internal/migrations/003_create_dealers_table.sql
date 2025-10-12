-- ============================================
-- ТАБЛИЦА: dealers (Главная справочная таблица дилеров)
-- ============================================
CREATE TABLE IF NOT EXISTS dealers (
    dealer_id SERIAL PRIMARY KEY,
    ruft VARCHAR(50) UNIQUE NOT NULL,  -- Уникальный идентификатор дилера (0.1, 0.2 и т.д.)
    dealer_name_ru VARCHAR(255) NOT NULL,  -- Название на русском
    dealer_name_en VARCHAR(255) NOT NULL,  -- Название на английском
    region VARCHAR(100),
    city VARCHAR(100),
    manager VARCHAR(100),
    
    -- Joint Decision (заполняется вручную через UI)
    joint_decision TEXT,
    
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_dealers_ruft ON dealers(ruft);
CREATE INDEX idx_dealers_name_ru ON dealers(dealer_name_ru);
CREATE INDEX idx_dealers_name_en ON dealers(dealer_name_en);
CREATE INDEX idx_dealers_region ON dealers(region);
CREATE INDEX idx_dealers_manager ON dealers(manager);

