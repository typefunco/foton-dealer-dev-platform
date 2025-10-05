package model

// Region представляет регион России.
type Region struct {
	ID   int64  `json:"id" db:"id"`
	Code string `json:"code" db:"code"` // all-russia, center, north-west, volga и т.д.
	Name string `json:"name" db:"name"` // All Russia, Center, North West, Volga и т.д.
}
