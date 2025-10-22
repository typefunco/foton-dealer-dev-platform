package database

import (
	"context"
	"time"
)

// DBConfig глобальная конфигурация базы данных
var DBConfig *PostgresConfig

// SetGlobalConfig устанавливает глобальную конфигурацию БД
func SetGlobalConfig(cfg *PostgresConfig) {
	DBConfig = cfg
}

// WithTimeout создает контекст с таймаутом для запросов к БД
func WithTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	if DBConfig != nil {
		return context.WithTimeout(ctx, DBConfig.QueryTimeout)
	}
	// Fallback на 30 секунд если конфигурация не установлена
	return context.WithTimeout(ctx, 30*time.Second)
}
