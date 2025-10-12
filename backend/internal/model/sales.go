package model

import "time"

// SalesTrainingsStatus представляет статус тренингов по продажам.
type SalesTrainingsStatus string

const (
	SalesTrainingsYes SalesTrainingsStatus = "Yes"
	SalesTrainingsNo  SalesTrainingsStatus = "No"
	SalesTrainingsY   SalesTrainingsStatus = "Y"
	SalesTrainingsN   SalesTrainingsStatus = "N"
)

// Sales отвечает за блок Dealer Sales.
// Содержит информацию о команде продаж, целях и обучении.
type Sales struct {
	ID       int    `json:"id" db:"id"`
	DealerID int    `json:"dealer_id" db:"dealer_id"`
	Quarter  string `json:"quarter" db:"quarter"`
	Year     int    `json:"year" db:"year"`

	// Sales target
	SalesTarget string `json:"sales_target" db:"sales_target"`

	// Stock (остатки на складе)
	StockHDT int `json:"stock_hdt" db:"stock_hdt"` // Heavy Duty Truck
	StockMDT int `json:"stock_mdt" db:"stock_mdt"` // Medium Duty Truck
	StockLDT int `json:"stock_ldt" db:"stock_ldt"` // Light Duty Truck

	// Buyout (выкуп)
	BuyoutHDT int `json:"buyout_hdt" db:"buyout_hdt"` // Heavy Duty Truck
	BuyoutMDT int `json:"buyout_mdt" db:"buyout_mdt"` // Medium Duty Truck
	BuyoutLDT int `json:"buyout_ldt" db:"buyout_ldt"` // Light Duty Truck

	// Sales personnel
	FotonSalesmen int `json:"foton_salesmen" db:"foton_salesmen"` // Количество продавцов Foton

	// Service contracts
	ServiceContractsSales int `json:"service_contracts_sales" db:"service_contracts_sales"` // Service contracts sales

	// Training
	SalesTrainings bool `json:"sales_trainings" db:"sales_trainings"` // Training status

	// Recommendation (из Excel)
	SalesDecision string `json:"sales_decision" db:"sales_decision"` // Sales decision

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// SalesWithDetails содержит полную информацию о продажах дилера.
// Используется для API ответов с объединёнными данными.
type SalesWithDetails struct {
	Sales
	DealerNameRu string `json:"dealer_name_ru"`
	DealerNameEn string `json:"dealer_name_en"`
	City         string `json:"city"`
	Region       string `json:"region"`
	Manager      string `json:"manager"`
	Ruft         string `json:"ruft"`
}
