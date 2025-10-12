package model

import "time"

// Dealer представляет основную информацию о дилере.
// Это центральная модель, которая связывает все остальные данные.
type Dealer struct {
	DealerID      int       `json:"dealer_id" db:"dealer_id"`
	Ruft          string    `json:"ruft" db:"ruft"`                     // Уникальный идентификатор дилера (0.1, 0.2 и т.д.)
	DealerNameRu  string    `json:"dealer_name_ru" db:"dealer_name_ru"` // Название на русском
	DealerNameEn  string    `json:"dealer_name_en" db:"dealer_name_en"` // Название на английском
	Region        string    `json:"region" db:"region"`
	City          string    `json:"city" db:"city"`
	Manager       string    `json:"manager" db:"manager"`               // Менеджер, отвечающий за дилера
	JointDecision *string   `json:"joint_decision" db:"joint_decision"` // Joint Decision (заполняется вручную через UI)
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

// DealerBrand представляет связь между дилером и брендами в его портфеле.
type DealerBrand struct {
	ID        int64     `json:"id" db:"id"`
	DealerID  int64     `json:"dealer_id" db:"dealer_id"`
	BrandName string    `json:"brand_name" db:"brand_name"` // FOTON, GAZ, KAMAZ и т.д.
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// DealerBusiness представляет побочный бизнес дилера (by-side business).
type DealerBusiness struct {
	ID           int64     `json:"id" db:"id"`
	DealerID     int64     `json:"dealer_id" db:"dealer_id"`
	BusinessType string    `json:"business_type" db:"business_type"` // Logistics, Warehousing, Retail и т.д.
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}
