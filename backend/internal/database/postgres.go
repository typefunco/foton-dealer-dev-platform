package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresConfig содержит настройки пула соединений PostgreSQL.
type PostgresConfig struct {
	MaxConns        int32
	MinConns        int32
	MaxConnLifetime time.Duration
	MaxConnIdleTime time.Duration
	ConnectTimeout  time.Duration // Таймаут подключения
	QueryTimeout    time.Duration // Таймаут выполнения запросов
}

// DefaultPostgresConfig возвращает конфигурацию пула по умолчанию.
func DefaultPostgresConfig() *PostgresConfig {
	return &PostgresConfig{
		MaxConns:        25,
		MinConns:        5,
		MaxConnLifetime: time.Hour,
		MaxConnIdleTime: 30 * time.Minute,
		ConnectTimeout:  30 * time.Second, // 30 секунд на подключение
		QueryTimeout:    60 * time.Second, // 60 секунд на выполнение запроса
	}
}

// NewPostgresPool создает новый пул соединений с PostgreSQL.
func NewPostgresPool(ctx context.Context, databaseURL string, cfg *PostgresConfig) (*pgxpool.Pool, error) {
	// Создаем контекст с таймаутом для подключения
	connectCtx, cancel := context.WithTimeout(ctx, cfg.ConnectTimeout)
	defer cancel()

	poolConfig, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database URL: %w", err)
	}

	// Применяем настройки пула
	if cfg != nil {
		poolConfig.MaxConns = cfg.MaxConns
		poolConfig.MinConns = cfg.MinConns
		poolConfig.MaxConnLifetime = cfg.MaxConnLifetime
		poolConfig.MaxConnIdleTime = cfg.MaxConnIdleTime
	}

	pool, err := pgxpool.NewWithConfig(connectCtx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Проверка подключения с таймаутом
	if err := pool.Ping(connectCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return pool, nil
}

// WithQueryTimeout создает контекст с таймаутом для выполнения запросов
func (cfg *PostgresConfig) WithQueryTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, cfg.QueryTimeout)
}
