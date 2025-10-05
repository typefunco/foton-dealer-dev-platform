package model

import "time"

// AfterSalesDecision представляет решение по послепродажному обслуживанию.
type AfterSalesDecision string

const (
	AfterSalesDecisionPlannedResult    AfterSalesDecision = "Planned Result"
	AfterSalesDecisionNeedsDevelopment AfterSalesDecision = "Needs development"
	AfterSalesDecisionFindNewCandidate AfterSalesDecision = "Find New Candidate"
	AfterSalesDecisionCloseDown        AfterSalesDecision = "Close Down"
)

// AfterSales отвечает за блок AfterSales (послепродажное обслуживание).
// Содержит информацию о запчастях, обучении и сервисе.
type AfterSales struct {
	ID                 int64              `json:"id" db:"id"`
	DealerID           int64              `json:"dealer_id" db:"dealer_id"`
	Quarter            string             `json:"quarter" db:"quarter"`                           // q1, q2, q3, q4
	Year               int                `json:"year" db:"year"`                                 // 2024, 2025 и т.д.
	RecommendedStock   int16              `json:"recommended_stock" db:"recommended_stock"`       // В процентах
	WarrantyStock      int16              `json:"warranty_stock" db:"warranty_stock"`             // В процентах
	FotonLaborHours    int16              `json:"foton_labor_hours" db:"foton_labor_hours"`       // В процентах
	ServiceContracts   int16              `json:"service_contracts" db:"service_contracts"`       // Абсолютное количество
	ASTrainings        bool               `json:"as_trainings" db:"as_trainings"`                 // Прошли ли обучение
	CSI                string             `json:"csi" db:"csi"`                                   // Customer Satisfaction Index
	FotonWarrantyHours int32              `json:"foton_warranty_hours" db:"foton_warranty_hours"` // Гарантийные часы Foton
	AfterSalesDecision AfterSalesDecision `json:"as_decision" db:"as_decision"`
	CreatedAt          time.Time          `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time          `json:"updated_at" db:"updated_at"`
}

// AfterSalesWithDetails содержит полную информацию о послепродажном обслуживании дилера.
// Используется для API ответов с объединёнными данными.
type AfterSalesWithDetails struct {
	AfterSales
	DealerName string `json:"dealer_name"`
	City       string `json:"city"`
	Region     string `json:"region"`
	Manager    string `json:"manager"`
}
