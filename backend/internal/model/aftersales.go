package model

import "time"

// ASTrainingsStatus представляет статус тренингов по послепродажному обслуживанию.
type ASTrainingsStatus string

const (
	ASTrainingsYes ASTrainingsStatus = "Yes"
	ASTrainingsNo  ASTrainingsStatus = "No"
	ASTrainingsY   ASTrainingsStatus = "Y"
	ASTrainingsN   ASTrainingsStatus = "N"
)

// AfterSales отвечает за блок AfterSales (послепродажное обслуживание).
// Содержит информацию о запчастях, обучении и сервисе.
type AfterSales struct {
	ID       int       `json:"id" db:"id"`
	DealerID int       `json:"dealer_id" db:"dealer_id"`
	Period   time.Time `json:"period" db:"period"`

	// Stock metrics
	RecommendedStockPct *float64 `json:"recommended_stock_pct" db:"recommended_stock_pct"` // В процентах
	WarrantyStockPct    *float64 `json:"warranty_stock_pct" db:"warranty_stock_pct"`       // В процентах

	// Labor hours
	FotonLaborHoursPct    *float64 `json:"foton_labor_hours_pct" db:"foton_labor_hours_pct"`     // В процентах
	WarrantyHours         *float64 `json:"warranty_hours" db:"warranty_hours"`                   // Гарантийные часы
	ServiceContractsHours *float64 `json:"service_contracts_hours" db:"service_contracts_hours"` // Часы сервисных контрактов

	// Training
	ASTrainings *ASTrainingsStatus `json:"as_trainings" db:"as_trainings"` // Y, N, Yes, No

	// Spare parts sales (revenue)
	SparePartsSalesQ      *float64 `json:"spare_parts_sales_q" db:"spare_parts_sales_q"`             // За квартал
	SparePartsSalesYtdPct *float64 `json:"spare_parts_sales_ytd_pct" db:"spare_parts_sales_ytd_pct"` // YTD динамика %

	// Recommendation (из Excel)
	ASRecommendation *string `json:"as_recommendation" db:"as_recommendation"` // Recommendation (из Excel)

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// AfterSalesWithDetails содержит полную информацию о послепродажном обслуживании дилера.
// Используется для API ответов с объединёнными данными.
type AfterSalesWithDetails struct {
	AfterSales
	DealerNameRu string `json:"dealer_name_ru"`
	DealerNameEn string `json:"dealer_name_en"`
	City         string `json:"city"`
	Region       string `json:"region"`
	Manager      string `json:"manager"`
	Ruft         string `json:"ruft"`
}
