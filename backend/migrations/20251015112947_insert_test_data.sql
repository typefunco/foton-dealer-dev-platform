-- +goose Up
-- РЕГИОНЫ
INSERT INTO regions (code, name) VALUES
('CENTRAL', 'Central'),
('NORTH_WEST', 'North West'),
('VOLGA', 'Volga'),
('SOUTH', 'South'),
('URAL', 'Ural'),
('SIBERIA', 'Siberia'),
('FAR_EAST', 'Far East'),
('CAUCASUS', 'Caucasus');

-- БРЕНДЫ
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

-- ДИЛЕРЫ
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

-- ПРОДАЖИ (SALES)
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

-- +goose Down
DELETE FROM sales;
DELETE FROM dealers;
DELETE FROM brands;
DELETE FROM regions;