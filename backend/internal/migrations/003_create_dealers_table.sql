-- Создание таблицы дилеров
CREATE TABLE IF NOT EXISTS dealers (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    city VARCHAR(100) NOT NULL,
    region VARCHAR(50) NOT NULL,
    manager VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

CREATE INDEX idx_dealers_region ON dealers(region);
CREATE INDEX idx_dealers_manager ON dealers(manager);

