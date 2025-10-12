package model

import "time"

// PerformanceSales отвечает за блок Performance Sales (финансовая производительность продаж).
// Содержит информацию о выручке, прибыли и марже по продажам автомобилей.
type PerformanceSales struct {
	ID       int       `json:"id" db:"id"`
	DealerID int       `json:"dealer_id" db:"dealer_id"`
	Period   time.Time `json:"period" db:"period"`

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
	ID       int       `json:"id" db:"id"`
	DealerID int       `json:"dealer_id" db:"dealer_id"`
	Period   time.Time `json:"period" db:"period"`

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
	DealerID     int     `json:"dealer_id"`
	DealerName   string  `json:"dealer_name"`
	City         string  `json:"city"`
	Region       string  `json:"region"`
	Manager      string  `json:"manager"`
	FotonRank    int     `json:"foton_rank"`
	SalesRevenue float64 `json:"sales_revenue"`
	SalesProfit  float64 `json:"sales_profit"`
	SalesMargin  float64 `json:"sales_margin"`
	AsRevenue    float64 `json:"as_revenue"`
	AsProfit     float64 `json:"as_profit"`
	AsMargin     float64 `json:"as_margin"`
	Marketing    float64 `json:"marketing"`
	Decision     string  `json:"decision"`
}
