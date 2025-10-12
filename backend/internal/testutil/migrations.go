package testutil

import (
	"context"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

// MigrationFacade предоставляет интерфейс для автоматического применения миграций
type MigrationFacade struct {
	migrationsPath string
	logger         *slog.Logger
}

// NewMigrationFacade создает новый экземпляр фасада миграций
func NewMigrationFacade(migrationsPath string, logger *slog.Logger) *MigrationFacade {
	return &MigrationFacade{
		migrationsPath: migrationsPath,
		logger:         logger,
	}
}

// MigrationFile представляет информацию о файле миграции
type MigrationFile struct {
	Number int
	Name   string
	Path   string
}

// ApplyAllMigrations применяет все миграции из папки в порядке их номеров
func (mf *MigrationFacade) ApplyAllMigrations(ctx context.Context, pool *pgxpool.Pool) error {
	mf.logger.Info("Starting automatic migration application", "path", mf.migrationsPath)

	// Получаем список всех файлов миграций
	migrationFiles, err := mf.getMigrationFiles()
	if err != nil {
		return fmt.Errorf("failed to get migration files: %w", err)
	}

	if len(migrationFiles) == 0 {
		mf.logger.Warn("No migration files found", "path", mf.migrationsPath)
		return nil
	}

	mf.logger.Info("Found migration files", "count", len(migrationFiles))

	// Применяем каждую миграцию
	for _, migration := range migrationFiles {
		mf.logger.Info("Applying migration", "file", migration.Name, "number", migration.Number)

		if err := mf.applyMigration(ctx, pool, migration); err != nil {
			return fmt.Errorf("failed to apply migration %s: %w", migration.Name, err)
		}

		mf.logger.Info("Migration applied successfully", "file", migration.Name)
	}

	mf.logger.Info("All migrations applied successfully", "total", len(migrationFiles))
	return nil
}

// getMigrationFiles сканирует папку миграций и возвращает отсортированный список файлов
func (mf *MigrationFacade) getMigrationFiles() ([]MigrationFile, error) {
	var migrationFiles []MigrationFile

	err := filepath.WalkDir(mf.migrationsPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Пропускаем директории
		if d.IsDir() {
			return nil
		}

		// Проверяем, что это SQL файл
		if !strings.HasSuffix(strings.ToLower(path), ".sql") {
			return nil
		}

		// Извлекаем номер миграции из имени файла
		fileName := d.Name()
		number, err := mf.extractMigrationNumber(fileName)
		if err != nil {
			mf.logger.Warn("Skipping file - invalid migration number", "file", fileName, "error", err)
			return nil
		}

		migrationFiles = append(migrationFiles, MigrationFile{
			Number: number,
			Name:   fileName,
			Path:   path,
		})

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk migrations directory: %w", err)
	}

	// Сортируем по номеру миграции
	sort.Slice(migrationFiles, func(i, j int) bool {
		return migrationFiles[i].Number < migrationFiles[j].Number
	})

	return migrationFiles, nil
}

// extractMigrationNumber извлекает номер миграции из имени файла
// Ожидается формат: 001_description.sql
func (mf *MigrationFacade) extractMigrationNumber(fileName string) (int, error) {
	// Убираем расширение
	nameWithoutExt := strings.TrimSuffix(fileName, ".sql")

	// Ищем первый символ подчеркивания
	underscoreIndex := strings.Index(nameWithoutExt, "_")
	if underscoreIndex == -1 {
		return 0, fmt.Errorf("no underscore found in filename")
	}

	// Извлекаем номер
	numberStr := nameWithoutExt[:underscoreIndex]
	number, err := strconv.Atoi(numberStr)
	if err != nil {
		return 0, fmt.Errorf("invalid migration number: %s", numberStr)
	}

	return number, nil
}

// applyMigration применяет одну миграцию
func (mf *MigrationFacade) applyMigration(ctx context.Context, pool *pgxpool.Pool, migration MigrationFile) error {
	// Читаем содержимое файла миграции
	content, err := os.ReadFile(migration.Path)
	if err != nil {
		return fmt.Errorf("failed to read migration file: %w", err)
	}

	// Выполняем SQL
	_, err = pool.Exec(ctx, string(content))
	if err != nil {
		return fmt.Errorf("failed to execute migration SQL: %w", err)
	}

	return nil
}

// GetMigrationStatus возвращает информацию о примененных миграциях
func (mf *MigrationFacade) GetMigrationStatus(ctx context.Context, pool *pgxpool.Pool) (map[string]bool, error) {
	// Создаем таблицу для отслеживания миграций, если её нет
	createTableSQL := `
		CREATE TABLE IF NOT EXISTS migration_status (
			migration_name VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`

	_, err := pool.Exec(ctx, createTableSQL)
	if err != nil {
		return nil, fmt.Errorf("failed to create migration status table: %w", err)
	}

	// Получаем список примененных миграций
	query := "SELECT migration_name FROM migration_status"
	rows, err := pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query migration status: %w", err)
	}
	defer rows.Close()

	appliedMigrations := make(map[string]bool)
	for rows.Next() {
		var migrationName string
		if err := rows.Scan(&migrationName); err != nil {
			return nil, fmt.Errorf("failed to scan migration name: %w", err)
		}
		appliedMigrations[migrationName] = true
	}

	return appliedMigrations, nil
}
