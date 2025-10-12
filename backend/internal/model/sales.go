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
	ID       int       `json:"id" db:"id"`
	DealerID int       `json:"dealer_id" db:"dealer_id"`
	Period   time.Time `json:"period" db:"period"`

	// Stock (остатки на складе)
	StockHDT *int `json:"stock_hdt" db:"stock_hdt"` // Heavy Duty Truck
	StockMDT *int `json:"stock_mdt" db:"stock_mdt"` // Medium Duty Truck
	StockLDT *int `json:"stock_ldt" db:"stock_ldt"` // Light Duty Truck

	// Buyout (выкуп)
	BuyoutHDT *int `json:"buyout_hdt" db:"buyout_hdt"` // Heavy Duty Truck
	BuyoutMDT *int `json:"buyout_mdt" db:"buyout_mdt"` // Medium Duty Truck
	BuyoutLDT *int `json:"buyout_ldt" db:"buyout_ldt"` // Light Duty Truck

	// Sales targets and personnel
	FotonSalesPersonnel *int `json:"foton_sales_personnel" db:"foton_sales_personnel"` // Количество продавцов Foton
	SalesTargetPlan     *int `json:"sales_target_plan" db:"sales_target_plan"`         // План продаж
	SalesTargetFact     *int `json:"sales_target_fact" db:"sales_target_fact"`         // Факт продаж

	// Service contracts
	ServiceContractsSales *float64 `json:"service_contracts_sales" db:"service_contracts_sales"` // Service contracts sales

	// Training
	SalesTrainings *SalesTrainingsStatus `json:"sales_trainings" db:"sales_trainings"` // Y, N, Yes, No

	// Recommendation (из Excel)
	SalesRecommendation *string `json:"sales_recommendation" db:"sales_recommendation"` // Recommendation (из Excel)

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
