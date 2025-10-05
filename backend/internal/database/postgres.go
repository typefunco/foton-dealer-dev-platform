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
}

// DefaultPostgresConfig возвращает конфигурацию пула по умолчанию.
func DefaultPostgresConfig() *PostgresConfig {
	return &PostgresConfig{
		MaxConns:        25,
		MinConns:        5,
		MaxConnLifetime: time.Hour,
		MaxConnIdleTime: 30 * time.Minute,
	}
}

// NewPostgresPool создает новый пул соединений с PostgreSQL.
func NewPostgresPool(ctx context.Context, databaseURL string, cfg *PostgresConfig) (*pgxpool.Pool, error) {
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

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Проверка подключения
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return pool, nil
}
