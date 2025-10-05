package model

import "time"

// Dealer представляет основную информацию о дилере.
// Это центральная модель, которая связывает все остальные данные.
type Dealer struct {
	ID        int64     `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	City      string    `json:"city" db:"city"`
	Region    string    `json:"region" db:"region"`
	Manager   string    `json:"manager" db:"manager"` // Менеджер, отвечающий за дилера
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
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
