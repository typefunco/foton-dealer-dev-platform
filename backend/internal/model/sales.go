package model

import "time"

// SalesDecision представляет решение по продажам дилера.
type SalesDecision string

const (
	SalesDecisionPlannedResult    SalesDecision = "Planned Result"
	SalesDecisionNeedsDevelopment SalesDecision = "Needs development"
	SalesDecisionFindNewCandidate SalesDecision = "Find New Candidate"
	SalesDecisionCloseDown        SalesDecision = "Close Down"
)

// Sales отвечает за блок Dealer Sales.
// Содержит информацию о команде продаж, целях и обучении.
type Sales struct {
	ID                    int64         `json:"id" db:"id"`
	DealerID              int64         `json:"dealer_id" db:"dealer_id"`
	Quarter               string        `json:"quarter" db:"quarter"`                                 // q1, q2, q3, q4
	Year                  int           `json:"year" db:"year"`                                       // 2024, 2025 и т.д.
	SalesTarget           string        `json:"sales_target" db:"sales_target"`                       // Формат: "40/100" (выполнено/план)
	StockHDT              int16         `json:"stock_hdt" db:"stock_hdt"`                             // Heavy Duty Truck
	StockMDT              int16         `json:"stock_mdt" db:"stock_mdt"`                             // Medium Duty Truck
	StockLDT              int16         `json:"stock_ldt" db:"stock_ldt"`                             // Light Duty Truck
	BuyoutHDT             int16         `json:"buyout_hdt" db:"buyout_hdt"`                           // Heavy Duty Truck
	BuyoutMDT             int16         `json:"buyout_mdt" db:"buyout_mdt"`                           // Medium Duty Truck
	BuyoutLDT             int16         `json:"buyout_ldt" db:"buyout_ldt"`                           // Light Duty Truck
	FotonSalesmen         int16         `json:"foton_salesmen" db:"foton_salesmen"`                   // Количество продавцов
	SalesTrainings        bool          `json:"sales_trainings" db:"sales_trainings"`                 // Прошли ли обучение
	ServiceContractsSales int16         `json:"service_contracts_sales" db:"service_contracts_sales"` // Количество сервисных контрактов
	SalesDecision         SalesDecision `json:"sales_decision" db:"sales_decision"`
	CreatedAt             time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt             time.Time     `json:"updated_at" db:"updated_at"`
}

// SalesWithDetails содержит полную информацию о продажах дилера.
// Используется для API ответов с объединёнными данными.
type SalesWithDetails struct {
	Sales
	DealerName      string `json:"dealer_name"`
	City            string `json:"city"`
	Region          string `json:"region"`
	Manager         string `json:"manager"`
	StockHdtMdtLdt  string `json:"stock_hdt_mdt_ldt"`  // Формат: "5/2/3"
	BuyoutHdtMdtLdt string `json:"buyout_hdt_mdt_ldt"` // Формат: "5/2/3"
}
