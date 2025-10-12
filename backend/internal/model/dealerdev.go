package model

import "time"

// DealershipClass представляет класс дилера (A, B, C, D).
type DealershipClass string

const (
	DealershipClassA DealershipClass = "A"
	DealershipClassB DealershipClass = "B"
	DealershipClassC DealershipClass = "C"
	DealershipClassD DealershipClass = "D"
)

// BrandingStatus представляет статус брендинга.
type BrandingStatus string

const (
	BrandingYes BrandingStatus = "Yes"
	BrandingNo  BrandingStatus = "No"
	BrandingY   BrandingStatus = "Y"
	BrandingN   BrandingStatus = "N"
)

// DealerDevelopment отвечает за блок Dealer Development.
// Содержит информацию о развитии дилера, брендах, классе и чек-листе.
type DealerDevelopment struct {
	ID                   int       `json:"id" db:"id"`
	DealerID             int       `json:"dealer_id" db:"dealer_id"`
	Quarter              string    `json:"quarter" db:"quarter"`
	Year                 int       `json:"year" db:"year"`
	CheckListScore       int       `json:"check_list_score" db:"check_list_score"`                   // Check List Score % (0-100)
	DealershipClass      string    `json:"dealership_class" db:"dealer_ship_class"`                  // A, B, C, D
	Branding             bool      `json:"branding" db:"branding"`                                   // Y, N, Yes, No
	MarketingInvestments int64     `json:"marketing_investments" db:"marketing_investments"`         // Marketing Investments Rub
	DDRecommendation     string    `json:"dealer_dev_recommendation" db:"dealer_dev_recommendation"` // Recommendation (из Excel)
	CreatedAt            time.Time `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time `json:"updated_at" db:"updated_at"`
}

// DealerDevelopmentWithDetails содержит полную информацию о дилере для отображения.
// Используется для API ответов с объединёнными данными.
type DealerDevelopmentWithDetails struct {
	DealerDevelopment
	DealerNameRu string `json:"dealer_name_ru"`
	DealerNameEn string `json:"dealer_name_en"`
	City         string `json:"city"`
	Region       string `json:"region"`
	Manager      string `json:"manager"`
	Ruft         string `json:"ruft"`
}

// DealerDevWithDetails - алиас для совместимости со старым кодом.
type DealerDevWithDetails = DealerDevelopmentWithDetails
