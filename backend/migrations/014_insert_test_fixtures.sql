-- ============================================
-- ТЕСТОВЫЕ ФИКСТУРЫ ДЛЯ ВСЕХ ТАБЛИЦ
-- ============================================

-- Очистка существующих данных (в порядке зависимостей)
DELETE FROM after_sales;
DELETE FROM performance;
DELETE FROM sales;
DELETE FROM dealer_dev;
DELETE FROM dealer_businesses;
DELETE FROM dealer_brands;
DELETE FROM users;
DELETE FROM dealers;
DELETE FROM brands;
DELETE FROM regions;

-- Сброс последовательностей
ALTER SEQUENCE dealers_id_seq RESTART WITH 1;
ALTER SEQUENCE users_id_seq RESTART WITH 1;
ALTER SEQUENCE brands_id_seq RESTART WITH 1;
ALTER SEQUENCE regions_id_seq RESTART WITH 1;
ALTER SEQUENCE dealer_brands_id_seq RESTART WITH 1;
ALTER SEQUENCE dealer_businesses_id_seq RESTART WITH 1;
ALTER SEQUENCE dealer_dev_id_seq RESTART WITH 1;
ALTER SEQUENCE sales_id_seq RESTART WITH 1;
ALTER SEQUENCE performance_id_seq RESTART WITH 1;
ALTER SEQUENCE after_sales_id_seq RESTART WITH 1;

-- ============================================
-- РЕГИОНЫ
-- ============================================
INSERT INTO regions (code, name) VALUES
('CENTRAL', 'Central'),
('NORTH_WEST', 'North West'),
('VOLGA', 'Volga'),
('SOUTH', 'South'),
('URAL', 'Ural'),
('SIBERIA', 'Siberia'),
('FAR_EAST', 'Far East'),
('CAUCASUS', 'Caucasus');

-- ============================================
-- БРЕНДЫ
-- ============================================
INSERT INTO brands (name, logo_path) VALUES
('FOTON', '/brands/Foton.png'),
('DONGFENG', '/brands/Dongfeng.png'),
('GAZ', '/brands/GAZ.png'),
('KAMAZ', '/brands/KAMAZ.png'),
('SHACMAN', '/brands/Shacman.png'),
('FAW', '/brands/Faw.png'),
('JAC', '/brands/Jac.png'),
('MAZ', '/brands/MAZ.png'),
('SANY', '/brands/SANY.png'),
('SITRAK', '/brands/Sitrak.png'),
('SOLLERS', '/brands/SOLLERS.png'),
('VALDAI', '/brands/VALDAI.png'),
('ISUZU', '/brands/ISUZU.png'),
('CHENLONG', '/brands/CHENLONG.png'),
('AMBERTRUCK', '/brands/AMBERTRUCK.png');

-- ============================================
-- ДИЛЕРЫ
-- ============================================
INSERT INTO dealers (name, city, region, manager, created_at, updated_at) VALUES
('Автофургон', 'Moscow', 'Central', 'Иван Петров', NOW(), NOW()),
('Автокуб', 'Moscow', 'Central', 'Мария Сидорова', NOW(), NOW()),
('Авто-М', 'Moscow', 'Central', 'Алексей Козлов', NOW(), NOW()),
('БТС Белгород', 'Moscow', 'Central', 'Елена Волкова', NOW(), NOW()),
('БТС Смоленск', 'Noginsk', 'Central', 'Дмитрий Соколов', NOW(), NOW()),
('Центр Трак Групп', 'Solnechnogorsk', 'Central', 'Анна Морозова', NOW(), NOW()),
('Экомтех', 'Ecomtekh', 'Central', 'Сергей Лебедев', NOW(), NOW()),
('ГАЗ 36', 'Yaroslavl', 'Central', 'Ольга Новикова', NOW(), NOW()),
('Глобал Трак Сейлс', 'Ryazan', 'Central', 'Павел Кузнецов', NOW(), NOW()),
('Гус Техцентр', 'Tambov', 'Central', 'Татьяна Попова', NOW(), NOW()),
('КомДорАвто', 'Tula', 'Central', 'Михаил Васильев', NOW(), NOW()),
('Мейджор Трак Центр', 'Lipeck', 'Central', 'Наталья Федорова', NOW(), NOW());

-- ============================================
-- ПОЛЬЗОВАТЕЛИ
-- ============================================
INSERT INTO users (login, password, is_admin, role, region, first_name, last_name, email, created_at, updated_at) VALUES
('admin', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', true, 'admin', 'Central', 'Администратор', 'Системы', 'admin@dealer-platform.com', NOW(), NOW()),
('manager1', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', false, 'manager', 'Central', 'Иван', 'Петров', 'ivan.petrov@dealer-platform.com', NOW(), NOW()),
('manager2', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', false, 'manager', 'North West', 'Мария', 'Сидорова', 'maria.sidorova@dealer-platform.com', NOW(), NOW()),
('manager3', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', false, 'manager', 'Volga', 'Алексей', 'Козлов', 'alexey.kozlov@dealer-platform.com', NOW(), NOW()),
('analyst1', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', false, 'analyst', 'South', 'Елена', 'Волкова', 'elena.volkova@dealer-platform.com', NOW(), NOW()),
('analyst2', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', false, 'analyst', 'Ural', 'Дмитрий', 'Соколов', 'dmitry.sokolov@dealer-platform.com', NOW(), NOW()),
('sales1', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', false, 'sales', 'Siberia', 'Анна', 'Морозова', 'anna.morozova@dealer-platform.com', NOW(), NOW()),
('sales2', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', false, 'sales', 'Far East', 'Сергей', 'Лебедев', 'sergey.lebedev@dealer-platform.com', NOW(), NOW());

-- ============================================
-- СВЯЗИ ДИЛЕРОВ И БРЕНДОВ
-- ============================================
INSERT INTO dealer_brands (dealer_id, brand_name, created_at) VALUES
-- Автофургон (1)
(1, 'FOTON', NOW()),
(1, 'DONGFENG', NOW()),
(1, 'GAZ', NOW()),
(1, 'KAMAZ', NOW()),
(1, 'SHACMAN', NOW()),

-- Автокуб (2)
(2, 'FOTON', NOW()),
(2, 'FAW', NOW()),

-- Авто-М (3)
(3, 'FOTON', NOW()),
(3, 'JAC', NOW()),
(3, 'MAZ', NOW()),

-- БТС Белгород (4)
(4, 'FOTON', NOW()),
(4, 'SANY', NOW()),
(4, 'SITRAK', NOW()),
(4, 'SOLLERS', NOW()),
(4, 'VALDAI', NOW()),
(4, 'ISUZU', NOW()),
(4, 'CHENLONG', NOW()),
(4, 'AMBERTRUCK', NOW()),

-- БТС Смоленск (5)
(5, 'FOTON', NOW()),
(5, 'FAW', NOW()),

-- Центр Трак Групп (6)
(6, 'FOTON', NOW()),
(6, 'DONGFENG', NOW()),
(6, 'GAZ', NOW()),

-- Экомтех (7)
(7, 'FOTON', NOW()),
(7, 'DONGFENG', NOW()),
(7, 'GAZ', NOW()),

-- ГАЗ 36 (8)
(8, 'FOTON', NOW()),
(8, 'DONGFENG', NOW()),
(8, 'GAZ', NOW()),

-- Глобал Трак Сейлс (9)
(9, 'FOTON', NOW()),
(9, 'DONGFENG', NOW()),
(9, 'GAZ', NOW()),

-- Гус Техцентр (10)
(10, 'FOTON', NOW()),
(10, 'DONGFENG', NOW()),
(10, 'GAZ', NOW()),

-- КомДорАвто (11)
(11, 'FOTON', NOW()),
(11, 'DONGFENG', NOW()),
(11, 'GAZ', NOW()),

-- Мейджор Трак Центр (12)
(12, 'FOTON', NOW()),
(12, 'DONGFENG', NOW()),
(12, 'GAZ', NOW());

-- ============================================
-- ПОБОЧНЫЙ БИЗНЕС ДИЛЕРОВ
-- ============================================
INSERT INTO dealer_businesses (dealer_id, business_type, created_at) VALUES
-- Автофургон (1)
(1, 'Logistics', NOW()),
(1, 'Warehousing', NOW()),

-- Автокуб (2)
(2, 'Transport', NOW()),

-- Авто-М (3)
(3, 'Logistics', NOW()),
(3, 'Retail', NOW()),

-- БТС Белгород (4)
(4, 'Logistics', NOW()),
(4, 'Warehousing', NOW()),
(4, 'Retail', NOW()),

-- БТС Смоленск (5)
(5, 'Logistics', NOW()),
(5, 'Retail', NOW()),

-- Центр Трак Групп (6)
(6, 'Logistics', NOW()),
(6, 'Warehousing', NOW()),
(6, 'Retail', NOW()),
(6, 'Service', NOW()),

-- Экомтех (7)
(7, 'Logistics', NOW()),
(7, 'Retail', NOW()),

-- ГАЗ 36 (8)
(8, 'Transport', NOW()),

-- Глобал Трак Сейлс (9)
(9, 'Transport', NOW()),

-- Гус Техцентр (10)
(10, 'Transport', NOW());

-- ============================================
-- РАЗВИТИЕ ДИЛЕРОВ (DEALER DEV)
-- ============================================
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

-- ============================================
-- ПРОДАЖИ (SALES)
-- ============================================
INSERT INTO sales (dealer_id, quarter, year, sales_target, stock_hdt, stock_mdt, stock_ldt, buyout_hdt, buyout_mdt, buyout_ldt, foton_salesmen, sales_trainings, service_contracts_sales, sales_decision, created_at, updated_at) VALUES
(1, 'Q1', 2024, '45', 15, 25, 10, 8, 12, 5, 3, true, 2500, 'Needs Development', NOW(), NOW()),
(2, 'Q1', 2024, '38', 12, 20, 8, 6, 10, 4, 2, false, 1800, 'Needs Development', NOW(), NOW()),
(3, 'Q1', 2024, '55', 18, 30, 12, 10, 15, 6, 4, true, 2200, 'Needs Development', NOW(), NOW()),
(4, 'Q1', 2024, '78', 25, 40, 18, 15, 25, 10, 6, true, 4100, 'Planned Result', NOW(), NOW()),
(5, 'Q1', 2024, '68', 20, 35, 15, 12, 20, 8, 5, true, 2900, 'Planned Result', NOW(), NOW()),
(6, 'Q1', 2024, '72', 22, 38, 16, 14, 22, 9, 5, true, 3700, 'Planned Result', NOW(), NOW()),
(7, 'Q1', 2024, '52', 16, 28, 12, 9, 18, 7, 4, true, 2300, 'Planned Result', NOW(), NOW()),
(8, 'Q1', 2024, '32', 10, 18, 8, 5, 12, 4, 2, false, 1500, 'Find New Candidate', NOW(), NOW()),
(9, 'Q1', 2024, '42', 14, 22, 10, 7, 15, 5, 3, true, 2100, 'Find New Candidate', NOW(), NOW()),
(10, 'Q1', 2024, '38', 12, 20, 8, 6, 13, 4, 2, false, 1900, 'Find New Candidate', NOW(), NOW()),
(11, 'Q1', 2024, '28', 8, 15, 6, 4, 10, 3, 1, true, 1200, 'Close Down', NOW(), NOW()),
(12, 'Q1', 2024, '23', 6, 12, 5, 3, 8, 2, 1, false, 2800, 'Close Down', NOW(), NOW());

-- ============================================
-- ПРОИЗВОДИТЕЛЬНОСТЬ (PERFORMANCE)
-- ============================================
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

-- ============================================
-- ПОСЛЕПРОДАЖНОЕ ОБСЛУЖИВАНИЕ (AFTER SALES)
-- ============================================
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

-- ============================================
-- ДОПОЛНИТЕЛЬНЫЕ ДАННЫЕ ДЛЯ ТЕСТИРОВАНИЯ
-- ============================================

-- Добавляем данные за предыдущий квартал для сравнения
INSERT INTO dealer_dev (dealer_id, quarter, year, check_list_score, dealer_ship_class, branding, marketing_investments, dealer_dev_recommendation, created_at, updated_at) VALUES
(1, 'Q4', 2023, 75, 'B', false, 2000000, 'Needs Development', NOW(), NOW()),
(2, 'Q4', 2023, 80, 'B', false, 1500000, 'Needs Development', NOW(), NOW()),
(3, 'Q4', 2023, 78, 'B', false, 1800000, 'Needs Development', NOW(), NOW()),
(4, 'Q4', 2023, 88, 'A', true, 3500000, 'Planned Result', NOW(), NOW()),
(5, 'Q4', 2023, 90, 'A', true, 2500000, 'Planned Result', NOW(), NOW());

INSERT INTO sales (dealer_id, quarter, year, sales_target, stock_hdt, stock_mdt, stock_ldt, buyout_hdt, buyout_mdt, buyout_ldt, foton_salesmen, sales_trainings, service_contracts_sales, sales_decision, created_at, updated_at) VALUES
(1, 'Q4', 2023, '42', 12, 20, 8, 6, 10, 4, 2, false, 2000, 'Needs Development', NOW(), NOW()),
(2, 'Q4', 2023, '33', 10, 18, 6, 5, 8, 3, 2, false, 1500, 'Needs Development', NOW(), NOW()),
(3, 'Q4', 2023, '48', 15, 25, 10, 8, 12, 5, 3, false, 1800, 'Needs Development', NOW(), NOW()),
(4, 'Q4', 2023, '68', 20, 35, 15, 12, 20, 8, 5, true, 3500, 'Planned Result', NOW(), NOW()),
(5, 'Q4', 2023, '58', 18, 30, 12, 10, 18, 7, 4, true, 2500, 'Planned Result', NOW(), NOW());

INSERT INTO performance (dealer_id, quarter, year, sales_revenue_rub, sales_profit_rub, sales_profit_percent, sales_margin_percent, after_sales_revenue_rub, after_sales_profit_rub, after_sales_margin_percent, marketing_investment, foton_rank, performance_decision, created_at, updated_at) VALUES
(1, 'Q4', 2023, 5000000, 225000, 4.5, 20.0, 5000000, 225000, 4.5, 2.0, 4, 'Needs Development', NOW(), NOW()),
(2, 'Q4', 2023, 4000000, 180000, 4.5, 20.0, 4000000, 180000, 4.5, 1.5, 4, 'Needs Development', NOW(), NOW()),
(3, 'Q4', 2023, 6000000, 270000, 4.5, 20.0, 6000000, 270000, 4.5, 1.8, 4, 'Needs Development', NOW(), NOW()),
(4, 'Q4', 2023, 8000000, 360000, 4.5, 20.0, 8000000, 360000, 4.5, 3.5, 4, 'Needs Development', NOW(), NOW()),
(5, 'Q4', 2023, 7000000, 315000, 4.5, 20.0, 7000000, 315000, 4.5, 2.5, 4, 'Needs Development', NOW(), NOW());

INSERT INTO after_sales (dealer_id, quarter, year, recommended_stock, warranty_stock, foton_labor_hours, service_contracts, as_trainings, csi, foton_warranty_hours, as_decision, created_at, updated_at) VALUES
(1, 'Q4', 2023, 4, 4, 4, 90, false, 'Poor', 110, 'Needs Development', NOW(), NOW()),
(2, 'Q4', 2023, 4, 4, 4, 90, false, 'Poor', 90, 'Needs Development', NOW(), NOW()),
(3, 'Q4', 2023, 4, 4, 4, 90, false, 'Poor', 100, 'Needs Development', NOW(), NOW()),
(4, 'Q4', 2023, 4, 4, 4, 90, false, 'Poor', 130, 'Needs Development', NOW(), NOW()),
(5, 'Q4', 2023, 4, 4, 4, 90, false, 'Poor', 120, 'Needs Development', NOW(), NOW());

-- ============================================
-- ЗАВЕРШЕНИЕ
-- ============================================
-- Обновляем статистику таблиц
ANALYZE regions;
ANALYZE brands;
ANALYZE dealers;
ANALYZE users;
ANALYZE dealer_brands;
ANALYZE dealer_businesses;
ANALYZE dealer_dev;
ANALYZE sales;
ANALYZE performance;
ANALYZE after_sales;

-- Выводим статистику
SELECT 'regions' as table_name, COUNT(*) as records FROM regions
UNION ALL
SELECT 'brands', COUNT(*) FROM brands
UNION ALL
SELECT 'dealers', COUNT(*) FROM dealers
UNION ALL
SELECT 'users', COUNT(*) FROM users
UNION ALL
SELECT 'dealer_brands', COUNT(*) FROM dealer_brands
UNION ALL
SELECT 'dealer_businesses', COUNT(*) FROM dealer_businesses
UNION ALL
SELECT 'dealer_dev', COUNT(*) FROM dealer_dev
UNION ALL
SELECT 'sales', COUNT(*) FROM sales
UNION ALL
SELECT 'performance', COUNT(*) FROM performance
UNION ALL
SELECT 'after_sales', COUNT(*) FROM after_sales;