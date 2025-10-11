-- Создание таблицы производительности дилеров (отдельная таблица для аналитики)
CREATE TABLE IF NOT EXISTS dealer_performance
(
    id                        serial,
    dealer_name               varchar(100),
    region                    varchar(50),
    city                      varchar(30),
    manager                   varchar(30),
    sales_revenues            int,
    sales_profits             int,
    sales_margin              int,
    as_revenues               int,
    as_profits                int,
    as_margin                 int,
    foton_rank                smallint,
    dealer_dev_recommendation varchar(20)
);

