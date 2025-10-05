package config

import (
	"fmt"
	"os"
)

// Config содержит конфигурацию приложения.
type Config struct {
	DatabaseURL string
	JWTSecret   string
	ServerPort  string
}

// Load загружает конфигурацию из переменных окружения.
// Возвращает ошибку, если обязательные переменные не установлены.
func Load() (*Config, error) {
	cfg := &Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		JWTSecret:   os.Getenv("JWT_SECRET"),
		ServerPort:  os.Getenv("SERVER_PORT"),
	}

	// Валидация обязательных полей
	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable is required")
	}

	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET environment variable is required")
	}

	if cfg.ServerPort == "" {
		return nil, fmt.Errorf("SERVER_PORT environment variable is required")
	}

	return cfg, nil
}

// MaskDSN маскирует пароль в DSN для безопасного логирования.
func MaskDSN(dsn string) string {
	if len(dsn) > 20 {
		return dsn[:20] + "***"
	}
	return "***"
}
