package testutil

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// TestDB представляет тестовую базу данных
type TestDB struct {
	Container testcontainers.Container
	Pool      *pgxpool.Pool
	DSN       string
}

// SetupTestDB настраивает тестовую базу данных с PostgreSQL контейнером
func SetupTestDB(t *testing.T) *TestDB {
	ctx := context.Background()

	// Создаем PostgreSQL контейнер
	postgresContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:15-alpine"),
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second)),
	)
	require.NoError(t, err)

	// Получаем DSN
	dsn, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	// Создаем пул соединений
	pool, err := pgxpool.New(ctx, dsn)
	require.NoError(t, err)

	// Проверяем соединение
	err = pool.Ping(ctx)
	require.NoError(t, err)

	return &TestDB{
		Container: postgresContainer,
		Pool:      pool,
		DSN:       dsn,
	}
}

// Cleanup очищает ресурсы тестовой базы данных
func (tdb *TestDB) Cleanup(t *testing.T) {
	if tdb.Pool != nil {
		tdb.Pool.Close()
	}
	if tdb.Container != nil {
		err := tdb.Container.Terminate(context.Background())
		require.NoError(t, err)
	}
}

// RunMigrations выполняет миграции базы данных автоматически
func (tdb *TestDB) RunMigrations(t *testing.T) {
	ctx := context.Background()
	logger := GetTestLogger()

	// Определяем путь к папке миграций относительно текущего файла
	migrationsPath := "../../../migrations"

	// Создаем фасад миграций
	migrationFacade := NewMigrationFacade(migrationsPath, logger)

	// Применяем все миграции автоматически
	err := migrationFacade.ApplyAllMigrations(ctx, tdb.Pool)
	require.NoError(t, err, "Failed to apply migrations automatically")
}

// GetTestLogger возвращает тестовый логгер
func GetTestLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
}

// CleanupTable очищает таблицу после теста
func (tdb *TestDB) CleanupTable(t *testing.T, tableName string) {
	ctx := context.Background()
	_, err := tdb.Pool.Exec(ctx, fmt.Sprintf("TRUNCATE TABLE %s CASCADE", tableName))
	require.NoError(t, err)
}

// CleanupAllTables очищает все таблицы после теста автоматически
func (tdb *TestDB) CleanupAllTables(t *testing.T) {
	ctx := context.Background()

	// Получаем список всех таблиц из базы данных
	query := `
		SELECT table_name 
		FROM information_schema.tables 
		WHERE table_schema = 'public' 
		AND table_type = 'BASE TABLE'
		ORDER BY table_name
	`

	rows, err := tdb.Pool.Query(ctx, query)
	require.NoError(t, err, "Failed to get table list")
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var tableName string
		err := rows.Scan(&tableName)
		require.NoError(t, err, "Failed to scan table name")
		tables = append(tables, tableName)
	}

	// Очищаем каждую таблицу
	for _, table := range tables {
		_, err := tdb.Pool.Exec(ctx, fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table))
		if err != nil {
			// Логируем ошибки, но не останавливаем выполнение
			t.Logf("Warning: Failed to truncate table %s: %v", table, err)
		}
	}
}
