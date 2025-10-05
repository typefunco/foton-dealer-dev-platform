package model

import "time"

// DealerDevClass представляет класс дилера (A, B, C, D).
type DealerDevClass string

const (
	DealerDevClassA DealerDevClass = "A"
	DealerDevClassB DealerDevClass = "B"
	DealerDevClassC DealerDevClass = "C"
	DealerDevClassD DealerDevClass = "D"
)

// DealerDevRecommendation представляет рекомендацию по развитию дилера.
type DealerDevRecommendation string

const (
	RecommendationPlannedResult    DealerDevRecommendation = "Planned Result"
	RecommendationNeedsDevelopment DealerDevRecommendation = "Needs Development"
	RecommendationFindNewCandidate DealerDevRecommendation = "Find New Candidate"
	RecommendationCloseDown        DealerDevRecommendation = "Close Down"
)

// DealerDev отвечает за блок Dealer Development.
// Содержит информацию о развитии дилера, брендах, классе и чек-листе.
type DealerDev struct {
	ID                   int64                   `json:"id" db:"id"`
	DealerID             int64                   `json:"dealer_id" db:"dealer_id"`
	Quarter              string                  `json:"quarter" db:"quarter"`                             // q1, q2, q3, q4
	Year                 int                     `json:"year" db:"year"`                                   // 2024, 2025 и т.д.
	CheckListScore       int16                   `json:"check_list_score" db:"check_list_score"`           // 0-100
	DealerShipClass      DealerDevClass          `json:"dealer_ship_class" db:"dealer_ship_class"`         // A, B, C, D
	Branding             bool                    `json:"branding" db:"branding"`                           // Наличие брендинга
	MarketingInvestments int64                   `json:"marketing_investments" db:"marketing_investments"` // В рублях
	Recommendation       DealerDevRecommendation `json:"dealer_dev_recommendation" db:"dealer_dev_recommendation"`
	CreatedAt            time.Time               `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time               `json:"updated_at" db:"updated_at"`
}

// DealerDevWithDetails содержит полную информацию о дилере для отображения.
// Используется для API ответов с объединёнными данными.
type DealerDevWithDetails struct {
	DealerDev
	DealerName       string   `json:"dealer_name"`
	City             string   `json:"city"`
	Region           string   `json:"region"`
	Manager          string   `json:"manager"`
	Brands           []string `json:"brands"`             // Массив названий брендов
	BrandsCount      int      `json:"brands_count"`       // Количество брендов
	BySideBusinesses []string `json:"by_side_businesses"` // Массив типов побочного бизнеса
}
