package model

import "time"

// UserRole представляет роль пользователя в системе.
type UserRole string

const (
	UserRoleAdmin   UserRole = "admin"   // Администратор системы
	UserRoleManager UserRole = "manager" // Менеджер региона
	UserRoleSales   UserRole = "sales"   // Сотрудник отдела продаж
	UserRoleViewer  UserRole = "viewer"  // Пользователь только для просмотра
)

// User структура пользователя системы.
// Содержит информацию для аутентификации и авторизации.
type User struct {
	ID        int64     `json:"id" db:"id"`
	Login     string    `json:"login" db:"login"`
	Password  string    `json:"password,omitempty" db:"password"` // omitempty для исключения из JSON ответов
	IsAdmin   bool      `json:"is_admin" db:"is_admin"`
	Role      UserRole  `json:"role" db:"role"`     // manager, sales, admin, viewer
	Region    string    `json:"region" db:"region"` // Регион, за который отвечает пользователь
	FirstName string    `json:"first_name" db:"first_name"`
	LastName  string    `json:"last_name" db:"last_name"`
	Email     string    `json:"email" db:"email"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// UserResponse представляет данные пользователя для API ответов.
// Не содержит пароль и другие чувствительные данные.
type UserResponse struct {
	ID        int64     `json:"id"`
	Login     string    `json:"login"`
	IsAdmin   bool      `json:"is_admin"`
	Role      UserRole  `json:"role"`
	Region    string    `json:"region"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}
