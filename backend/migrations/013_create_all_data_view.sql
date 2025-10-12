-- ============================================
-- VIEW: all_data
-- Объединяет все данные для отображения на странице AllData
-- ============================================
CREATE VIEW all_data AS
SELECT 
    -- General Info (из dealers)
    d.dealer_id,
    d.ruft,
    d.dealer_name_ru,
    d.dealer_name_en,
    d.region,
    d.city,
    d.manager,
    d.joint_decision,  -- заполняется вручную
    
    -- Dealer Development
    dd.check_list_score,
    dd.dealership_class,
    dd.brands,
    dd.branding,
    dd.marketing_investments,
    dd.by_side_businesses,
    dd.dd_recommendation,
    
    -- Sales
    s.stock_hdt,
    s.stock_mdt,
    s.stock_ldt,
    s.buyout_hdt,
    s.buyout_mdt,
    s.buyout_ldt,
    s.foton_sales_personnel,
    s.sales_target_plan,
    s.sales_target_fact,
    s.service_contracts_sales,
    s.sales_trainings,
    s.sales_recommendation,
    
    -- AfterSales
    a.recommended_stock_pct,
    a.warranty_stock_pct,
    a.foton_labor_hours_pct,
    a.warranty_hours,
    a.service_contracts_hours,
    a.as_trainings,
    a.spare_parts_sales_q,
    a.spare_parts_sales_ytd_pct,
    a.as_recommendation,
    
    -- Performance Sales
    ps.quantity_sold,
    ps.sales_revenue,
    ps.sales_margin,
    ps.sales_margin_pct,
    ps.sales_profit_pct,
    
    -- Performance AfterSales
    pa.as_revenue,
    pa.as_margin,
    pa.as_margin_pct,
    pa.as_profit_pct,
    
    -- Period (берем из любой таблицы, они должны совпадать)
    COALESCE(dd.period, s.period, a.period, ps.period, pa.period) as period,
    
    -- Timestamps
    d.updated_at
    
FROM dealers d
LEFT JOIN dealer_development dd ON d.dealer_id = dd.dealer_id
LEFT JOIN sales s ON d.dealer_id = s.dealer_id AND dd.period = s.period
LEFT JOIN aftersales a ON d.dealer_id = a.dealer_id AND dd.period = a.period
LEFT JOIN performance_sales ps ON d.dealer_id = ps.dealer_id AND dd.period = ps.period
LEFT JOIN performance_aftersales pa ON d.dealer_id = pa.dealer_id AND dd.period = pa.period;
