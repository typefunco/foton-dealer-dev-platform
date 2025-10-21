package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config содержит конфигурацию приложения.
type Config struct {
	DatabaseURL string
	JWTSecret   string
	ServerPort  string
	MaxFileSize int64  // Максимальный размер файла в байтах (по умолчанию 100MB)
	LogLevel    string // Уровень логирования (по умолчанию INFO)
	DBMaxConns  int32  // Максимальное количество соединений с БД (по умолчанию 25)
}

// Load загружает конфигурацию из переменных окружения.
// Возвращает ошибку, если обязательные переменные не установлены.
func Load() (*Config, error) {
	cfg := &Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		JWTSecret:   os.Getenv("JWT_SECRET"),
		ServerPort:  os.Getenv("SERVER_PORT"),
		MaxFileSize: 100 * 1024 * 1024, // 100MB по умолчанию
		LogLevel:    "INFO",
		DBMaxConns:  25,
	}

	// Парсим MaxFileSize из переменной окружения
	if maxFileSizeStr := os.Getenv("MAX_FILE_SIZE"); maxFileSizeStr != "" {
		if maxFileSize, err := strconv.ParseInt(maxFileSizeStr, 10, 64); err == nil {
			cfg.MaxFileSize = maxFileSize * 1024 * 1024 // Конвертируем МБ в байты
		}
	}

	// Парсим LogLevel из переменной окружения
	if logLevel := os.Getenv("LOG_LEVEL"); logLevel != "" {
		cfg.LogLevel = logLevel
	}

	// Парсим DBMaxConns из переменной окружения
	if dbMaxConnsStr := os.Getenv("DB_MAX_CONNECTIONS"); dbMaxConnsStr != "" {
		if dbMaxConns, err := strconv.ParseInt(dbMaxConnsStr, 10, 32); err == nil {
			cfg.DBMaxConns = int32(dbMaxConns)
		}
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
