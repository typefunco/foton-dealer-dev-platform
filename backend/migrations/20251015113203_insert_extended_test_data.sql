-- +goose Up
-- ПОЛЬЗОВАТЕЛИ
INSERT INTO users (login, password, is_admin, role, region, first_name, last_name, email, created_at, updated_at) VALUES
('admin', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', true, 'admin', 'Central', 'Администратор', 'Системы', 'admin@dealer-platform.com', NOW(), NOW()),
('manager1', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', false, 'manager', 'Central', 'Иван', 'Петров', 'ivan.petrov@dealer-platform.com', NOW(), NOW()),
('manager2', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', false, 'manager', 'North West', 'Мария', 'Сидорова', 'maria.sidorova@dealer-platform.com', NOW(), NOW()),
('manager3', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', false, 'manager', 'Volga', 'Алексей', 'Козлов', 'alexey.kozlov@dealer-platform.com', NOW(), NOW()),
('analyst1', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', false, 'analyst', 'South', 'Елена', 'Волкова', 'elena.volkova@dealer-platform.com', NOW(), NOW()),
('analyst2', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', false, 'analyst', 'Ural', 'Дмитрий', 'Соколов', 'dmitry.sokolov@dealer-platform.com', NOW(), NOW()),
('sales1', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', false, 'sales', 'Siberia', 'Анна', 'Морозова', 'anna.morozova@dealer-platform.com', NOW(), NOW()),
('sales2', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', false, 'sales', 'Far East', 'Сергей', 'Лебедев', 'sergey.lebedev@dealer-platform.com', NOW(), NOW());

-- СВЯЗИ ДИЛЕРОВ И БРЕНДОВ
INSERT INTO dealer_brands (dealer_id, brand_name, created_at) VALUES
(1, 'FOTON', NOW()),
(1, 'DONGFENG', NOW()),
(1, 'GAZ', NOW()),
(1, 'KAMAZ', NOW()),
(1, 'SHACMAN', NOW()),
(2, 'FOTON', NOW()),
(2, 'FAW', NOW()),
(3, 'FOTON', NOW()),
(3, 'JAC', NOW()),
(3, 'MAZ', NOW()),
(4, 'FOTON', NOW()),
(4, 'SANY', NOW()),
(4, 'SITRAK', NOW()),
(4, 'SOLLERS', NOW()),
(4, 'VALDAI', NOW()),
(4, 'ISUZU', NOW()),
(4, 'CHENLONG', NOW()),
(4, 'AMBERTRUCK', NOW()),
(5, 'FOTON', NOW()),
(5, 'FAW', NOW()),
(6, 'FOTON', NOW()),
(6, 'DONGFENG', NOW()),
(6, 'GAZ', NOW()),
(7, 'FOTON', NOW()),
(7, 'DONGFENG', NOW()),
(7, 'GAZ', NOW()),
(8, 'FOTON', NOW()),
(8, 'DONGFENG', NOW()),
(8, 'GAZ', NOW()),
(9, 'FOTON', NOW()),
(9, 'DONGFENG', NOW()),
(9, 'GAZ', NOW()),
(10, 'FOTON', NOW()),
(10, 'DONGFENG', NOW()),
(10, 'GAZ', NOW()),
(11, 'FOTON', NOW()),
(11, 'DONGFENG', NOW()),
(11, 'GAZ', NOW()),
(12, 'FOTON', NOW()),
(12, 'DONGFENG', NOW()),
(12, 'GAZ', NOW());

-- ПОБОЧНЫЙ БИЗНЕС ДИЛЕРОВ
INSERT INTO dealer_businesses (dealer_id, business_type, created_at) VALUES
(1, 'Logistics', NOW()),
(1, 'Warehousing', NOW()),
(2, 'Transport', NOW()),
(3, 'Logistics', NOW()),
(3, 'Retail', NOW()),
(4, 'Logistics', NOW()),
(4, 'Warehousing', NOW()),
(4, 'Retail', NOW()),
(5, 'Logistics', NOW()),
(5, 'Retail', NOW()),
(6, 'Logistics', NOW()),
(6, 'Warehousing', NOW()),
(6, 'Retail', NOW()),
(6, 'Service', NOW()),
(7, 'Logistics', NOW()),
(7, 'Retail', NOW()),
(8, 'Transport', NOW()),
(9, 'Transport', NOW()),
(10, 'Transport', NOW());

-- РАЗВИТИЕ ДИЛЕРОВ (DEALER DEV)
INSERT INTO dealer_dev (dealer_id, quarter, year, check_list_score, dealer_ship_class, branding, marketing_investments, dealer_dev_recommendation, created_at, updated_at) VALUES
(1, 'Q1', 2024, 80, 'B', true, 2500000, 'Needs Development', NOW(), NOW()),
(2, 'Q1', 2024, 85, 'B', false, 1800000, 'Needs Development', NOW(), NOW()),
(3, 'Q1', 2024, 82, 'B', true, 2200000, 'Needs Development', NOW(), NOW()),
(4, 'Q1', 2024, 92, 'A', true, 4100000, 'Planned Result', NOW(), NOW()),
(5, 'Q1', 2024, 95, 'A', true, 2900000, 'Planned Result', NOW(), NOW()),
(6, 'Q1', 2024, 96, 'A', true, 3700000, 'Planned Result', NOW(), NOW()),
(7, 'Q1', 2024, 91, 'A', true, 2300000, 'Planned Result', NOW(), NOW()),
(8, 'Q1', 2024, 76, 'C', false, 1500000, 'Find New Candidate', NOW(), NOW()),
(9, 'Q1', 2024, 73, 'C', false, 2100000, 'Find New Candidate', NOW(), NOW()),
(10, 'Q1', 2024, 72, 'C', false, 1900000, 'Find New Candidate', NOW(), NOW()),
(11, 'Q1', 2024, 68, 'D', false, 1200000, 'Close Down', NOW(), NOW()),
(12, 'Q1', 2024, 66, 'D', false, 2800000, 'Close Down', NOW(), NOW());

-- ПРОИЗВОДИТЕЛЬНОСТЬ (PERFORMANCE)
INSERT INTO performance (dealer_id, quarter, year, sales_revenue_rub, sales_profit_rub, sales_profit_percent, sales_margin_percent, after_sales_revenue_rub, after_sales_profit_rub, after_sales_margin_percent, marketing_investment, foton_rank, performance_decision, created_at, updated_at) VALUES
(1, 'Q1', 2024, 5555555, 277777, 5.0, 20.0, 5555555, 277777, 5.0, 2.5, 5, 'Needs Development', NOW(), NOW()),
(2, 'Q1', 2024, 4444444, 222222, 5.0, 20.0, 4444444, 222222, 5.0, 3.2, 5, 'Needs Development', NOW(), NOW()),
(3, 'Q1', 2024, 6666666, 333333, 5.0, 20.0, 6666666, 333333, 5.0, 1.8, 5, 'Needs Development', NOW(), NOW()),
(4, 'Q1', 2024, 8888888, 444444, 5.0, 20.0, 8888888, 444444, 5.0, 4.1, 5, 'Needs Development', NOW(), NOW()),
(5, 'Q1', 2024, 7777777, 388888, 5.0, 20.0, 7777777, 388888, 5.0, 2.9, 5, 'Needs Development', NOW(), NOW()),
(6, 'Q1', 2024, 8888888, 444444, 5.0, 20.0, 8888888, 444444, 5.0, 3.7, 5, 'Needs Development', NOW(), NOW()),
(7, 'Q1', 2024, 6666666, 333333, 5.0, 20.0, 6666666, 333333, 5.0, 2.3, 5, 'Needs Development', NOW(), NOW()),
(8, 'Q1', 2024, 3333333, 166666, 5.0, 20.0, 3333333, 166666, 5.0, 1.5, 5, 'Needs Development', NOW(), NOW()),
(9, 'Q1', 2024, 4444444, 222222, 5.0, 20.0, 4444444, 222222, 5.0, 2.1, 5, 'Needs Development', NOW(), NOW()),
(10, 'Q1', 2024, 4444444, 222222, 5.0, 20.0, 4444444, 222222, 5.0, 1.9, 5, 'Needs Development', NOW(), NOW()),
(11, 'Q1', 2024, 2222222, 111111, 5.0, 20.0, 2222222, 111111, 5.0, 1.2, 5, 'Needs Development', NOW(), NOW()),
(12, 'Q1', 2024, 2222222, 111111, 5.0, 20.0, 2222222, 111111, 5.0, 2.8, 5, 'Needs Development', NOW(), NOW());

-- ПОСЛЕПРОДАЖНОЕ ОБСЛУЖИВАНИЕ (AFTER SALES)
INSERT INTO after_sales (dealer_id, quarter, year, recommended_stock, warranty_stock, foton_labor_hours, service_contracts, as_trainings, csi, foton_warranty_hours, as_decision, created_at, updated_at) VALUES
(1, 'Q1', 2024, 5, 5, 5, 95, true, 'Good', 120, 'Needs Development', NOW(), NOW()),
(2, 'Q1', 2024, 5, 5, 5, 95, false, 'Poor', 95, 'Needs Development', NOW(), NOW()),
(3, 'Q1', 2024, 5, 5, 5, 95, true, 'Good', 110, 'Needs Development', NOW(), NOW()),
(4, 'Q1', 2024, 5, 5, 5, 95, false, 'Poor', 140, 'Needs Development', NOW(), NOW()),
(5, 'Q1', 2024, 5, 5, 5, 95, true, 'Good', 125, 'Needs Development', NOW(), NOW()),
(6, 'Q1', 2024, 5, 5, 5, 95, false, 'Poor', 130, 'Needs Development', NOW(), NOW()),
(7, 'Q1', 2024, 5, 5, 5, 95, true, 'Good', 115, 'Needs Development', NOW(), NOW()),
(8, 'Q1', 2024, 5, 5, 5, 95, false, 'Poor', 100, 'Needs Development', NOW(), NOW()),
(9, 'Q1', 2024, 5, 5, 5, 95, true, 'Good', 105, 'Needs Development', NOW(), NOW()),
(10, 'Q1', 2024, 5, 5, 5, 95, false, 'Poor', 98, 'Needs Development', NOW(), NOW()),
(11, 'Q1', 2024, 5, 5, 5, 95, true, 'Good', 90, 'Needs Development', NOW(), NOW()),
(12, 'Q1', 2024, 5, 5, 5, 95, false, 'Poor', 85, 'Needs Development', NOW(), NOW());

-- +goose Down
DELETE FROM after_sales;
DELETE FROM performance;
DELETE FROM dealer_dev;
DELETE FROM dealer_businesses;
DELETE FROM dealer_brands;
DELETE FROM users;