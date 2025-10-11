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

// UserFilter представляет фильтры для поиска пользователей.
// Все поля являются опциональными и могут быть nil.
// Если поле nil, то оно не участвует в фильтрации.
type UserFilter struct {
	ID        *int64    `json:"id,omitempty"`
	Login     *string   `json:"login,omitempty"`
	Email     *string   `json:"email,omitempty"`
	IsAdmin   *bool     `json:"is_admin,omitempty"`
	Role      *UserRole `json:"role,omitempty"`
	Region    *string   `json:"region,omitempty"`
	FirstName *string   `json:"first_name,omitempty"`
	LastName  *string   `json:"last_name,omitempty"`
}

// UserUpdate представляет поля для обновления пользователя.
// Только не-nil поля будут обновлены в базе данных.
type UserUpdate struct {
	Login     *string   `json:"login,omitempty"`
	Password  *string   `json:"password,omitempty"`
	IsAdmin   *bool     `json:"is_admin,omitempty"`
	Role      *UserRole `json:"role,omitempty"`
	Region    *string   `json:"region,omitempty"`
	FirstName *string   `json:"first_name,omitempty"`
	LastName  *string   `json:"last_name,omitempty"`
	Email     *string   `json:"email,omitempty"`
}

// UserCreateRequest представляет запрос на создание нового пользователя.
type UserCreateRequest struct {
	Login     string   `json:"login" binding:"required"`
	Password  string   `json:"password" binding:"required"`
	IsAdmin   bool     `json:"is_admin"`
	Role      UserRole `json:"role" binding:"required"`
	Region    string   `json:"region,omitempty"`
	FirstName string   `json:"first_name,omitempty"`
	LastName  string   `json:"last_name,omitempty"`
	Email     string   `json:"email,omitempty"`
}
