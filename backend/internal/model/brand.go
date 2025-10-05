package model

// Brand представляет бренд грузовиков.
type Brand struct {
	ID       int64  `json:"id" db:"id"`
	Name     string `json:"name" db:"name"`           // FOTON, GAZ, KAMAZ и т.д.
	LogoPath string `json:"logo_path" db:"logo_path"` // Путь к логотипу бренда
}
