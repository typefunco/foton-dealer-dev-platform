package model

// User структура пользователя.
type User struct {
	ID        int64  `json:"id" db:"id"`
	Login     string `json:"login" db:"login"`
	Password  string `json:"password" db:"password"`
	IsAdmin   bool   `json:"is_admin" db:"is_admin"`
	Role      string `json:"role" db:"role"` // manager, sales and etc.
	Region    string `json:"region" db:"region"`
	CreatedAt string `json:"created_at" db:"created_at"`
	UpdatedAt string `json:"updated_at" db:"updated_at"`
}
