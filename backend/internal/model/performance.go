package model

import "time"

// PerformanceDecision представляет решение по производительности дилера.
type PerformanceDecision string

const (
	PerformanceDecisionPlannedResult    PerformanceDecision = "Planned Result"
	PerformanceDecisionNeedsDevelopment PerformanceDecision = "Needs development"
	PerformanceDecisionFindNewCandidate PerformanceDecision = "Find New Candidate"
	PerformanceDecisionCloseDown        PerformanceDecision = "Close Down"
)

// Performance отвечает за блок Performance (финансовая производительность).
// Содержит информацию о выручке, прибыли и марже по продажам и послепродажному обслуживанию.
type Performance struct {
	ID       int64  `json:"id" db:"id"`
	DealerID int64  `json:"dealer_id" db:"dealer_id"`
	Quarter  string `json:"quarter" db:"quarter"` // q1, q2, q3, q4
	Year     int    `json:"year" db:"year"`       // 2024, 2025 и т.д.

	// Продажи автомобилей
	SalesRevenueRub    int64   `json:"sales_revenue_rub" db:"sales_revenue_rub"`       // Выручка в рублях
	SalesProfitRub     int64   `json:"sales_profit_rub" db:"sales_profit_rub"`         // Прибыль в рублях
	SalesProfitPercent float64 `json:"sales_profit_percent" db:"sales_profit_percent"` // Прибыль в процентах
	SalesMarginPercent float64 `json:"sales_margin_percent" db:"sales_margin_percent"` // Маржа в процентах

	// Послепродажное обслуживание
	AfterSalesRevenueRub    int64   `json:"after_sales_revenue_rub" db:"after_sales_revenue_rub"`       // Выручка в рублях
	AfterSalesProfitRub     int64   `json:"after_sales_profit_rub" db:"after_sales_profit_rub"`         // Прибыль в рублях
	AfterSalesMarginPercent float64 `json:"after_sales_margin_percent" db:"after_sales_margin_percent"` // Маржа в процентах

	// Маркетинг и рейтинг
	MarketingInvestment float64 `json:"marketing_investment" db:"marketing_investment"` // Инвестиции в маркетинг (в млн рублей)
	FotonRank           int16   `json:"foton_rank" db:"foton_rank"`                     // Рейтинг от 1 до 10

	PerformanceDecision PerformanceDecision `json:"performance_decision" db:"performance_decision"`
	CreatedAt           time.Time           `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time           `json:"updated_at" db:"updated_at"`
}

// PerformanceWithDetails содержит полную информацию о производительности дилера.
// Используется для API ответов с объединёнными данными.
type PerformanceWithDetails struct {
	Performance
	DealerName                 string `json:"dealer_name"`
	City                       string `json:"city"`
	Region                     string `json:"region"`
	Manager                    string `json:"manager"`
	SalesRevenueFormatted      string `json:"sales_revenue_formatted"`       // Форматированная строка "5 555 555"
	SalesProfitFormatted       string `json:"sales_profit_formatted"`        // Форматированная строка "5 000 000"
	AfterSalesRevenueFormatted string `json:"after_sales_revenue_formatted"` // Форматированная строка
	AfterSalesProfitFormatted  string `json:"after_sales_profit_formatted"`  // Форматированная строка
}
