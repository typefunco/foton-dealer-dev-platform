package repository

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/typefunco/dealer_dev_platform/internal/model"
)

// DynamicTableRepository интерфейс репозитория для работы с динамическими таблицами.
type DynamicTableRepository interface {
	// CreateTable создает таблицу с заданной схемой
	CreateTable(ctx context.Context, schema *model.DynamicTableSchema) error

	// InsertData вставляет данные в таблицу
	InsertData(ctx context.Context, tableName string, columns []string, rows [][]interface{}) error

	// TableExists проверяет существование таблицы
	TableExists(ctx context.Context, tableName string) (bool, error)

	// GetTableMetadata возвращает метаданные таблицы
	GetTableMetadata(ctx context.Context, tableName string) (*model.ExcelTableMetadata, error)

	// ListDynamicTables возвращает список всех динамически созданных таблиц
	ListDynamicTables(ctx context.Context) ([]*model.ExcelTableMetadata, error)

	// BeginTransaction начинает транзакцию
	BeginTransaction(ctx context.Context) (pgx.Tx, error)
}

// dynamicTableRepository реализация репозитория для динамических таблиц.
type dynamicTableRepository struct {
	pool   *pgxpool.Pool
	logger *slog.Logger
}

// NewDynamicTableRepository создает новый экземпляр репозитория.
func NewDynamicTableRepository(pool *pgxpool.Pool, logger *slog.Logger) DynamicTableRepository {
	return &dynamicTableRepository{
		pool:   pool,
		logger: logger,
	}
}

// CreateTable создает таблицу с заданной схемой.
func (r *dynamicTableRepository) CreateTable(ctx context.Context, schema *model.DynamicTableSchema) error {
	// Санитизация названия таблицы
	tableName := r.sanitizeTableName(schema.TableName)

	// Проверяем, что таблица не существует
	exists, err := r.TableExists(ctx, tableName)
	if err != nil {
		return fmt.Errorf("failed to check table existence: %w", err)
	}

	if exists {
		r.logger.Info("Table already exists, skipping creation", slog.String("table", tableName))
		return nil
	}

	// Строим SQL для создания таблицы
	query := r.buildCreateTableQuery(tableName, schema.Columns)

	r.logger.Info("Creating table",
		slog.String("table", tableName),
		slog.String("query", query),
	)

	_, err = r.pool.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to create table %s: %w", tableName, err)
	}

	r.logger.Info("Table created successfully", slog.String("table", tableName))
	return nil
}

// InsertData вставляет данные в таблицу.
func (r *dynamicTableRepository) InsertData(ctx context.Context, tableName string, columns []string, rows [][]interface{}) error {
	if len(rows) == 0 {
		r.logger.Info("No data to insert", slog.String("table", tableName))
		return nil
	}

	// Санитизация названия таблицы
	tableName = r.sanitizeTableName(tableName)

	// Санитизация названий колонок
	sanitizedColumns := make([]string, len(columns))
	for i, col := range columns {
		sanitizedColumns[i] = r.sanitizeColumnName(col)
	}

	// Строим SQL для вставки данных
	query := r.buildInsertQuery(tableName, sanitizedColumns, len(rows))

	r.logger.Info("Inserting data",
		slog.String("table", tableName),
		slog.Int("rows", len(rows)),
		slog.Int("columns", len(columns)),
	)

	// Подготавливаем данные для вставки
	args := make([]interface{}, 0, len(rows)*len(columns))
	for _, row := range rows {
		for i := range sanitizedColumns {
			if i < len(row) {
				args = append(args, row[i])
			} else {
				args = append(args, nil)
			}
		}
	}

	_, err := r.pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to insert data into table %s: %w", tableName, err)
	}

	r.logger.Info("Data inserted successfully",
		slog.String("table", tableName),
		slog.Int("rows", len(rows)),
	)

	return nil
}

// TableExists проверяет существование таблицы.
func (r *dynamicTableRepository) TableExists(ctx context.Context, tableName string) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT FROM information_schema.tables 
			WHERE table_schema = 'public' 
			AND table_name = $1
		)
	`

	var exists bool
	err := r.pool.QueryRow(ctx, query, tableName).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check table existence: %w", err)
	}

	return exists, nil
}

// GetTableMetadata возвращает метаданные таблицы.
func (r *dynamicTableRepository) GetTableMetadata(ctx context.Context, tableName string) (*model.ExcelTableMetadata, error) {
	// Проверяем существование таблицы
	exists, err := r.TableExists(ctx, tableName)
	if err != nil {
		return nil, fmt.Errorf("failed to check table existence: %w", err)
	}

	if !exists {
		return nil, fmt.Errorf("table %s does not exist", tableName)
	}

	// Получаем количество строк
	var rowsCount int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s", r.sanitizeTableName(tableName))
	err = r.pool.QueryRow(ctx, countQuery).Scan(&rowsCount)
	if err != nil {
		return nil, fmt.Errorf("failed to get row count: %w", err)
	}

	// Получаем информацию о колонках
	columnsQuery := `
		SELECT column_name 
		FROM information_schema.columns 
		WHERE table_schema = 'public' 
		AND table_name = $1 
		AND column_name != 'id' 
		AND column_name != 'created_at'
		ORDER BY ordinal_position
	`

	rows, err := r.pool.Query(ctx, columnsQuery, tableName)
	if err != nil {
		return nil, fmt.Errorf("failed to get column information: %w", err)
	}
	defer rows.Close()

	var columns []string
	for rows.Next() {
		var columnName string
		if err := rows.Scan(&columnName); err != nil {
			return nil, fmt.Errorf("failed to scan column name: %w", err)
		}
		columns = append(columns, columnName)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating columns: %w", err)
	}

	return &model.ExcelTableMetadata{
		TableName: tableName,
		RowsCount: rowsCount,
		CreatedAt: time.Now(),
		Columns:   columns,
	}, nil
}

// ListDynamicTables возвращает список всех динамически созданных таблиц.
func (r *dynamicTableRepository) ListDynamicTables(ctx context.Context) ([]*model.ExcelTableMetadata, error) {
	query := `
		SELECT table_name 
		FROM information_schema.tables 
		WHERE table_schema = 'public' 
		AND table_name LIKE '%_q%_%'
		ORDER BY table_name
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list dynamic tables: %w", err)
	}
	defer rows.Close()

	var tables []*model.ExcelTableMetadata
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return nil, fmt.Errorf("failed to scan table name: %w", err)
		}

		metadata, err := r.GetTableMetadata(ctx, tableName)
		if err != nil {
			r.logger.Warn("Failed to get metadata for table",
				slog.String("table", tableName),
				slog.String("error", err.Error()),
			)
			continue
		}

		tables = append(tables, metadata)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating tables: %w", err)
	}

	return tables, nil
}

// BeginTransaction начинает транзакцию.
func (r *dynamicTableRepository) BeginTransaction(ctx context.Context) (pgx.Tx, error) {
	return r.pool.Begin(ctx)
}

// sanitizeTableName санитизирует название таблицы для безопасности.
func (r *dynamicTableRepository) sanitizeTableName(name string) string {
	// Приводим к нижнему регистру
	name = strings.ToLower(name)

	// Заменяем спецсимволы на underscore
	name = strings.ReplaceAll(name, " ", "_")
	name = strings.ReplaceAll(name, "-", "_")
	name = strings.ReplaceAll(name, ".", "_")
	name = strings.ReplaceAll(name, "/", "_")
	name = strings.ReplaceAll(name, "\\", "_")

	// Убираем множественные underscore
	for strings.Contains(name, "__") {
		name = strings.ReplaceAll(name, "__", "_")
	}

	// Убираем начальные и конечные underscore
	name = strings.Trim(name, "_")

	// Ограничиваем длину
	if len(name) > 63 {
		name = name[:63]
	}

	return name
}

// sanitizeColumnName санитизирует название колонки для безопасности.
func (r *dynamicTableRepository) sanitizeColumnName(name string) string {
	// Приводим к нижнему регистру
	name = strings.ToLower(name)

	// Заменяем спецсимволы на underscore
	name = strings.ReplaceAll(name, " ", "_")
	name = strings.ReplaceAll(name, "-", "_")
	name = strings.ReplaceAll(name, ".", "_")
	name = strings.ReplaceAll(name, "/", "_")
	name = strings.ReplaceAll(name, "\\", "_")
	name = strings.ReplaceAll(name, "(", "_")
	name = strings.ReplaceAll(name, ")", "_")
	name = strings.ReplaceAll(name, "%", "_percent")

	// Убираем множественные underscore
	for strings.Contains(name, "__") {
		name = strings.ReplaceAll(name, "__", "_")
	}

	// Убираем начальные и конечные underscore
	name = strings.Trim(name, "_")

	// Если название пустое или начинается с цифры, добавляем префикс
	if name == "" || (len(name) > 0 && name[0] >= '0' && name[0] <= '9') {
		name = "column_" + name
	}

	// Ограничиваем длину
	if len(name) > 63 {
		name = name[:63]
	}

	return name
}

// buildCreateTableQuery строит SQL запрос для создания таблицы.
func (r *dynamicTableRepository) buildCreateTableQuery(tableName string, columns []model.DynamicTableColumn) string {
	var columnDefs []string

	// Добавляем стандартные колонки
	columnDefs = append(columnDefs, "id BIGSERIAL PRIMARY KEY")

	// Добавляем колонки из схемы
	for _, col := range columns {
		columnDefs = append(columnDefs, fmt.Sprintf("%s TEXT", col.Name))
	}

	// Добавляем колонку created_at
	columnDefs = append(columnDefs, "created_at TIMESTAMP DEFAULT NOW()")

	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (\n    %s\n)",
		tableName,
		strings.Join(columnDefs, ",\n    "))

	return query
}

// buildInsertQuery строит SQL запрос для вставки данных.
func (r *dynamicTableRepository) buildInsertQuery(tableName string, columns []string, rowCount int) string {
	// Создаем плейсхолдеры для значений
	var placeholders []string
	for i := 0; i < rowCount; i++ {
		var rowPlaceholders []string
		for j := 0; j < len(columns); j++ {
			rowPlaceholders = append(rowPlaceholders, fmt.Sprintf("$%d", i*len(columns)+j+1))
		}
		placeholders = append(placeholders, fmt.Sprintf("(%s)", strings.Join(rowPlaceholders, ", ")))
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES %s",
		tableName,
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "))

	return query
}
