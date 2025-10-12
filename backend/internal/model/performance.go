package model

import "time"

// PerformanceSales отвечает за блок Performance Sales (финансовая производительность продаж).
// Содержит информацию о выручке, прибыли и марже по продажам автомобилей.
type PerformanceSales struct {
	ID       int    `json:"id" db:"id"`
	DealerID int    `json:"dealer_id" db:"dealer_id"`
	Quarter  string `json:"quarter" db:"quarter"`
	Year     int    `json:"year" db:"year"`

	// Sales financial metrics
	QuantitySold      *int     `json:"quantity_sold" db:"quantity_sold"`               // Количество проданных автомобилей
	SalesRevenue      *float64 `json:"sales_revenue" db:"sales_revenue"`               // Выручка (с НДС)
	SalesRevenueNoVat *float64 `json:"sales_revenue_no_vat" db:"sales_revenue_no_vat"` // Выручка без НДС
	SalesCost         *float64 `json:"sales_cost" db:"sales_cost"`                     // Стоимость
	SalesMargin       *float64 `json:"sales_margin" db:"sales_margin"`                 // Валовая прибыль (в рублях)
	SalesMarginPct    *float64 `json:"sales_margin_pct" db:"sales_margin_pct"`         // Маржа % = (margin / revenue) * 100
	SalesProfitPct    *float64 `json:"sales_profit_pct" db:"sales_profit_pct"`         // Рентабельность %

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// PerformanceAfterSales отвечает за блок Performance AfterSales (финансовая производительность запчастей).
// Содержит информацию о выручке, прибыли и марже по продажам запчастей.
type PerformanceAfterSales struct {
	ID       int    `json:"id" db:"id"`
	DealerID int    `json:"dealer_id" db:"dealer_id"`
	Quarter  string `json:"quarter" db:"quarter"`
	Year     int    `json:"year" db:"year"`

	// AfterSales financial metrics
	ASRevenue      *float64 `json:"as_revenue" db:"as_revenue"`               // Выручка (с НДС)
	ASRevenueNoVat *float64 `json:"as_revenue_no_vat" db:"as_revenue_no_vat"` // Выручка без НДС
	ASCost         *float64 `json:"as_cost" db:"as_cost"`                     // Стоимость
	ASMargin       *float64 `json:"as_margin" db:"as_margin"`                 // Валовая прибыль (в рублях)
	ASMarginPct    *float64 `json:"as_margin_pct" db:"as_margin_pct"`         // Маржа % = (margin / revenue) * 100
	ASProfitPct    *float64 `json:"as_profit_pct" db:"as_profit_pct"`         // Рентабельность %

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// PerformanceSalesWithDetails содержит полную информацию о производительности продаж дилера.
// Используется для API ответов с объединёнными данными.
type PerformanceSalesWithDetails struct {
	PerformanceSales
	DealerNameRu string `json:"dealer_name_ru"`
	DealerNameEn string `json:"dealer_name_en"`
	City         string `json:"city"`
	Region       string `json:"region"`
	Manager      string `json:"manager"`
	Ruft         string `json:"ruft"`
}

// PerformanceAfterSalesWithDetails содержит полную информацию о производительности запчастей дилера.
// Используется для API ответов с объединёнными данными.
type PerformanceAfterSalesWithDetails struct {
	PerformanceAfterSales
	DealerNameRu string `json:"dealer_name_ru"`
	DealerNameEn string `json:"dealer_name_en"`
	City         string `json:"city"`
	Region       string `json:"region"`
	Manager      string `json:"manager"`
	Ruft         string `json:"ruft"`
}

// PerformanceWithDetails - алиас для совместимости со старым кодом.
type PerformanceWithDetails struct {
	DealerID                int     `json:"dealer_id"`
	DealerNameRu            string  `json:"dealer_name_ru"`
	City                    string  `json:"city"`
	Region                  string  `json:"region"`
	Manager                 string  `json:"manager"`
	FotonRank               int     `json:"foton_rank"`
	SalesRevenueRub         float64 `json:"sales_revenue_rub"`
	SalesProfitRub          float64 `json:"sales_profit_rub"`
	SalesMarginPercent      float64 `json:"sales_margin_percent"`
	AfterSalesRevenueRub    float64 `json:"after_sales_revenue_rub"`
	AfterSalesProfitRub     float64 `json:"after_sales_profit_rub"`
	AfterSalesMarginPercent float64 `json:"after_sales_margin_percent"`
	MarketingInvestment     float64 `json:"marketing_investment"`
	PerformanceDecision     string  `json:"performance_decision"`
}
