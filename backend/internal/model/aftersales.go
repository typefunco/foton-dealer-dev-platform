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
	ID       int    `json:"id" db:"id"`
	DealerID int    `json:"dealer_id" db:"dealer_id"`
	Quarter  string `json:"quarter" db:"quarter"`
	Year     int    `json:"year" db:"year"`

	// Stock metrics
	RecommendedStock int `json:"recommended_stock" db:"recommended_stock"`
	WarrantyStock    int `json:"warranty_stock" db:"warranty_stock"`

	// Labor hours
	FotonLaborHours    int `json:"foton_labor_hours" db:"foton_labor_hours"`
	FotonWarrantyHours int `json:"foton_warranty_hours" db:"foton_warranty_hours"`
	ServiceContracts   int `json:"service_contracts" db:"service_contracts"`

	// Training
	ASTrainings bool `json:"as_trainings" db:"as_trainings"`

	// CSI and Decision
	CSI        *string `json:"csi" db:"csi"`
	ASDecision string  `json:"as_decision" db:"as_decision"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// AfterSalesWithDetails содержит полную информацию о послепродажном обслуживании дилера.
// Используется для API ответов с объединёнными данными.
type AfterSalesWithDetails struct {
	AfterSales
	DealerNameRu string `json:"dealer_name_ru"`
	City         string `json:"city"`
	Region       string `json:"region"`
	Manager      string `json:"manager"`
	// Дополнительные поля для Excel данных
	SparePartsSalesQ3     string `json:"spare_parts_sales_q3"`
	SparePartsSalesYtd    string `json:"spare_parts_sales_ytd_percent"`
	FotonLabourHoursShare string `json:"foton_labour_hours_share"`
}
