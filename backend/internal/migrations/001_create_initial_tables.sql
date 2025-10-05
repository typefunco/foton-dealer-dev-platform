-- Создание таблицы регионов
CREATE TABLE IF NOT EXISTS regions (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(50) NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL
);

-- Создание таблицы брендов
CREATE TABLE IF NOT EXISTS brands (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    logo_path VARCHAR(255)
);

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

-- Создание таблицы связи дилеров и брендов
CREATE TABLE IF NOT EXISTS dealer_brands (
    id BIGSERIAL PRIMARY KEY,
    dealer_id BIGINT NOT NULL REFERENCES dealers(id) ON DELETE CASCADE,
    brand_name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP NOT NULL
);

CREATE INDEX idx_dealer_brands_dealer_id ON dealer_brands(dealer_id);

-- Создание таблицы побочного бизнеса дилеров
CREATE TABLE IF NOT EXISTS dealer_businesses (
    id BIGSERIAL PRIMARY KEY,
    dealer_id BIGINT NOT NULL REFERENCES dealers(id) ON DELETE CASCADE,
    business_type VARCHAR(100) NOT NULL,
    created_at TIMESTAMP NOT NULL
);

CREATE INDEX idx_dealer_businesses_dealer_id ON dealer_businesses(dealer_id);

-- Создание таблицы пользователей
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    login VARCHAR(100) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    is_admin BOOLEAN NOT NULL,
    role VARCHAR(50) NOT NULL,
    region VARCHAR(50),
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    email VARCHAR(255),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

CREATE INDEX idx_users_login ON users(login);
CREATE INDEX idx_users_region ON users(region);

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

