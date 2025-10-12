package testutil

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestMigrationFacade_ApplyAllMigrations(t *testing.T) {
	ctx := context.Background()
	logger := GetTestLogger()

	// Создаем тестовую базу данных
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
	defer func() {
		err := postgresContainer.Terminate(ctx)
		require.NoError(t, err)
	}()

	// Получаем DSN и создаем пул соединений
	dsn, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	pool, err := pgxpool.New(ctx, dsn)
	require.NoError(t, err)
	defer pool.Close()

	// Создаем фасад миграций
	migrationsPath := "../../migrations"
	migrationFacade := NewMigrationFacade(migrationsPath, logger)

	// Применяем миграции
	err = migrationFacade.ApplyAllMigrations(ctx, pool)
	require.NoError(t, err)

	// Проверяем, что таблицы созданы
	tables := []string{
		"regions",
		"brands",
		"dealers",
		"dealer_brands",
		"dealer_businesses",
		"users",
		"dealer_development",
		"sales",
		"performance_sales",
		"aftersales",
		"dealer_performance",
		"performance_aftersales",
	}

	for _, table := range tables {
		query := `
			SELECT EXISTS (
				SELECT FROM information_schema.tables 
				WHERE table_schema = 'public' 
				AND table_name = $1
			)
		`
		var exists bool
		err := pool.QueryRow(ctx, query, table).Scan(&exists)
		require.NoError(t, err)
		assert.True(t, exists, "Table %s should exist", table)
	}
}

func TestMigrationFacade_GetMigrationFiles(t *testing.T) {
	logger := GetTestLogger()
	migrationsPath := "../../migrations"
	migrationFacade := NewMigrationFacade(migrationsPath, logger)

	migrationFiles, err := migrationFacade.getMigrationFiles()
	require.NoError(t, err)

	// Проверяем, что найдены файлы миграций
	assert.Greater(t, len(migrationFiles), 0, "Should find migration files")

	// Проверяем, что файлы отсортированы по номеру
	for i := 1; i < len(migrationFiles); i++ {
		assert.LessOrEqual(t, migrationFiles[i-1].Number, migrationFiles[i].Number,
			"Migration files should be sorted by number")
	}

	// Проверяем, что первый файл имеет номер 1
	assert.Equal(t, 1, migrationFiles[0].Number, "First migration should have number 1")
	assert.Contains(t, migrationFiles[0].Name, "001_", "First migration should start with 001_")

	// Проверяем, что последний файл имеет максимальный номер
	lastIndex := len(migrationFiles) - 1
	assert.Equal(t, 13, migrationFiles[lastIndex].Number, "Last migration should have number 13")
	assert.Contains(t, migrationFiles[lastIndex].Name, "013_", "Last migration should start with 013_")
}

func TestMigrationFacade_ExtractMigrationNumber(t *testing.T) {
	logger := GetTestLogger()
	migrationFacade := NewMigrationFacade("", logger)

	testCases := []struct {
		fileName string
		expected int
		hasError bool
	}{
		{"001_create_regions_table.sql", 1, false},
		{"002_create_brands_table.sql", 2, false},
		{"013_create_all_data_view.sql", 13, false},
		{"100_some_migration.sql", 100, false},
		{"invalid_file.sql", 0, true},
		{"no_number.sql", 0, true},
		{"", 0, true},
	}

	for _, tc := range testCases {
		t.Run(tc.fileName, func(t *testing.T) {
			number, err := migrationFacade.extractMigrationNumber(tc.fileName)

			if tc.hasError {
				assert.Error(t, err)
				assert.Equal(t, 0, number)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, number)
			}
		})
	}
}
