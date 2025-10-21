package excel

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/typefunco/dealer_dev_platform/internal/model"
	"github.com/typefunco/dealer_dev_platform/internal/repository"
	"github.com/xuri/excelize/v2"
)

// Service сервис для работы с Excel файлами.
type Service struct {
	dynamicRepo repository.DynamicTableRepository
	logger      *slog.Logger
}

// NewService создает новый экземпляр сервиса Excel.
func NewService(dynamicRepo repository.DynamicTableRepository, logger *slog.Logger) *Service {
	return &Service{
		dynamicRepo: dynamicRepo,
		logger:      logger,
	}
}

// ProcessExcelFile обрабатывает Excel файл и создает единую таблицу в БД.
func (s *Service) ProcessExcelFile(ctx context.Context, file io.Reader, fileName string) (*model.ExcelProcessingResult, error) {
	startTime := time.Now()

	s.logger.Info("Starting Excel file processing",
		slog.String("file_name", fileName),
		slog.Time("start_time", startTime),
	)

	// Читаем Excel файл
	f, err := excelize.OpenReader(file)
	if err != nil {
		return nil, fmt.Errorf("failed to open Excel file: %w", err)
	}
	defer f.Close()

	// Получаем список листов
	sheetList := f.GetSheetList()
	if len(sheetList) == 0 {
		return nil, fmt.Errorf("Excel file contains no sheets")
	}

	s.logger.Info("Found sheets in Excel file",
		slog.String("file_name", fileName),
		slog.Int("sheets_count", len(sheetList)),
		slog.String("sheet_names", strings.Join(sheetList, ", ")),
	)

	// Начинаем транзакцию
	tx, err := s.dynamicRepo.BeginTransaction(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Парсим метаданные из названия файла
	fileInfo, err := s.parseFileName(fileName)
	if err != nil {
		s.logger.Error("Failed to parse file name",
			slog.String("file_name", fileName),
			slog.String("error", err.Error()),
		)
		return nil, fmt.Errorf("failed to parse file name: %w", err)
	}

	s.logger.Info("File info parsed",
		slog.String("file_name", fileName),
		slog.String("table_name", fileInfo.TableName),
		slog.String("quarter", fileInfo.Quarter),
		slog.Int("year", fileInfo.Year),
	)

	// Обрабатываем все листы и собираем данные
	var allData []model.ExcelRowWithRegion
	var commonColumns []string
	var errors []model.ExcelError

	for _, sheetName := range sheetList {
		s.logger.Info("Processing sheet",
			slog.String("sheet_name", sheetName),
		)

		// Извлекаем регион из названия листа
		region, err := s.extractRegionFromSheetName(sheetName)
		if err != nil {
			s.logger.Error("Failed to extract region from sheet name",
				slog.String("sheet_name", sheetName),
				slog.String("error", err.Error()),
			)
			errors = append(errors, model.ExcelError{
				SheetName: sheetName,
				Message:   "Failed to extract region",
				Error:     err.Error(),
			})
			continue
		}

		// Обрабатываем лист
		sheetData, columns, err := s.processSheetForUnifiedTable(ctx, tx, f, sheetName, region)
		if err != nil {
			s.logger.Error("Failed to process sheet",
				slog.String("sheet_name", sheetName),
				slog.String("error", err.Error()),
			)
			errors = append(errors, model.ExcelError{
				SheetName: sheetName,
				Message:   "Failed to process sheet",
				Error:     err.Error(),
			})
			continue
		}

		s.logger.Info("Sheet processed successfully",
			slog.String("sheet_name", sheetName),
			slog.String("region", region),
			slog.Int("rows", len(sheetData)),
			slog.Int("columns", len(columns)),
		)

		// Проверяем структуру колонок
		if len(commonColumns) == 0 {
			commonColumns = columns
		} else {
			if !s.compareColumnStructures(commonColumns, columns) {
				return nil, fmt.Errorf("sheet %s has different column structure", sheetName)
			}
		}

		allData = append(allData, sheetData...)
	}

	// Если есть ошибки, откатываем транзакцию
	if len(errors) > 0 {
		s.logger.Error("Errors occurred during processing, rolling back transaction",
			slog.Int("errors_count", len(errors)),
		)
		return &model.ExcelProcessingResult{
			Success:        false,
			TablesCreated:  []model.ExcelTableMetadata{},
			Errors:         errors,
			TotalRows:      0,
			ProcessingTime: time.Since(startTime),
		}, nil
	}

	// Создаем единую таблицу
	tableName := fileInfo.TableName
	err = s.createUnifiedTable(ctx, tx, tableName, commonColumns)
	if err != nil {
		return nil, fmt.Errorf("failed to create unified table: %w", err)
	}

	// Вставляем все данные
	err = s.insertUnifiedData(ctx, tx, tableName, commonColumns, allData)
	if err != nil {
		return nil, fmt.Errorf("failed to insert unified data: %w", err)
	}

	// Коммитим транзакцию
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	processingTime := time.Since(startTime)

	s.logger.Info("Excel file processing completed successfully",
		slog.String("file_name", fileName),
		slog.String("table_name", tableName),
		slog.Int("total_rows", len(allData)),
		slog.Duration("processing_time", processingTime),
	)

	return &model.ExcelProcessingResult{
		Success: true,
		TablesCreated: []model.ExcelTableMetadata{
			{
				TableName: tableName,
				RowsCount: len(allData),
				CreatedAt: time.Now(),
				Columns:   commonColumns,
			},
		},
		Errors:         []model.ExcelError{},
		TotalRows:      len(allData),
		ProcessingTime: processingTime,
	}, nil
}

// processSheet обрабатывает отдельный лист Excel.
func (s *Service) processSheet(ctx context.Context, tx pgx.Tx, f *excelize.File, sheetName string) (*model.ExcelTableMetadata, error) {
	// Парсим метаданные из названия листа
	sheetInfo, err := s.parseSheetName(sheetName)
	if err != nil {
		return nil, fmt.Errorf("failed to parse sheet name: %w", err)
	}

	// Читаем данные листа
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("failed to get rows from sheet: %w", err)
	}

	if len(rows) == 0 {
		s.logger.Warn("Sheet is empty", slog.String("sheet_name", sheetName))
		return nil, nil
	}

	// Первая строка - заголовки
	headers := rows[0]
	if len(headers) == 0 {
		return nil, fmt.Errorf("sheet has no headers")
	}

	// Санитизируем заголовки
	sanitizedHeaders := s.sanitizeHeaders(headers)

	// Создаем схему таблицы
	schema := &model.DynamicTableSchema{
		TableName: sheetInfo.TableName,
		Columns:   make([]model.DynamicTableColumn, len(sanitizedHeaders)),
	}

	for i, header := range sanitizedHeaders {
		schema.Columns[i] = model.DynamicTableColumn{
			Name:     header,
			DataType: "TEXT",
			Nullable: true,
		}
	}

	// Создаем таблицу
	createQuery := s.buildCreateTableQuery(schema.TableName, schema.Columns)
	_, err = tx.Exec(ctx, createQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	// Подготавливаем данные для вставки
	dataRows := rows[1:] // Пропускаем заголовки
	if len(dataRows) == 0 {
		s.logger.Info("No data rows to insert", slog.String("sheet_name", sheetName))
		return &model.ExcelTableMetadata{
			TableName: schema.TableName,
			RowsCount: 0,
			CreatedAt: time.Now(),
			Columns:   sanitizedHeaders,
		}, nil
	}

	// Вставляем данные
	err = s.insertSheetData(ctx, tx, schema.TableName, sanitizedHeaders, dataRows)
	if err != nil {
		return nil, fmt.Errorf("failed to insert data: %w", err)
	}

	return &model.ExcelTableMetadata{
		TableName: schema.TableName,
		RowsCount: len(dataRows),
		CreatedAt: time.Now(),
		Columns:   sanitizedHeaders,
	}, nil
}

// parseSheetName парсит название листа и извлекает метаданные.
func (s *Service) parseSheetName(sheetName string) (*model.ExcelSheetInfo, error) {
	// Приводим к нижнему регистру для обработки
	name := strings.ToLower(sheetName)

	// Извлекаем квартал и год
	quarter, year := s.extractQuarterAndYear(name)

	// Определяем категорию
	category := s.extractCategory(name)

	// Генерируем название таблицы
	tableName := s.generateTableName(category, quarter, year)

	return &model.ExcelSheetInfo{
		Name:      sheetName,
		Category:  category,
		Quarter:   quarter,
		Year:      year,
		TableName: tableName,
	}, nil
}

// extractQuarterAndYear извлекает квартал и год из названия листа.
func (s *Service) extractQuarterAndYear(sheetName string) (string, int) {
	// Регулярное выражение для поиска квартала и года (регистронезависимое)
	quarterYearRegex := regexp.MustCompile(`(?i)q([1-4])[_\-\s]*(\d{4})`)
	matches := quarterYearRegex.FindStringSubmatch(sheetName)

	if len(matches) >= 3 {
		quarter := "Q" + matches[1]
		year, _ := strconv.Atoi(matches[2])
		return quarter, year
	}

	// Если год не найден, используем текущий год
	currentYear := time.Now().Year()

	// Ищем только квартал (регистронезависимо)
	quarterRegex := regexp.MustCompile(`(?i)q([1-4])`)
	quarterMatches := quarterRegex.FindStringSubmatch(sheetName)

	if len(quarterMatches) >= 2 {
		quarter := "Q" + quarterMatches[1]
		return quarter, currentYear
	}

	// По умолчанию возвращаем текущий квартал и год
	return s.getCurrentQuarter(), currentYear
}

// extractCategory извлекает категорию из названия листа.
func (s *Service) extractCategory(sheetName string) string {
	// Убираем квартал и год из названия
	name := strings.ToLower(sheetName)
	name = regexp.MustCompile(`q[1-4][_\-\s]*\d{4}`).ReplaceAllString(name, "")
	name = regexp.MustCompile(`q[1-4]`).ReplaceAllString(name, "")

	// Убираем лишние символы и заменяем пробелы на underscore
	name = strings.ReplaceAll(name, "_", " ")
	name = strings.ReplaceAll(name, "-", " ")
	name = strings.TrimSpace(name)
	name = strings.ReplaceAll(name, " ", "_")

	// Если название пустое, используем "data"
	if name == "" {
		return "data"
	}

	return name
}

// generateTableName генерирует название таблицы для БД.
func (s *Service) generateTableName(category, quarter string, year int) string {
	// Санитизируем категорию
	category = strings.ToLower(category)
	category = strings.ReplaceAll(category, " ", "_")
	category = strings.ReplaceAll(category, "-", "_")
	category = strings.Trim(category, "_")

	// Формируем название таблицы
	tableName := fmt.Sprintf("%s_%s_%d", category, strings.ToLower(quarter), year)

	// Ограничиваем длину
	if len(tableName) > 63 {
		tableName = tableName[:63]
	}

	return tableName
}

// sanitizeHeaders санитизирует заголовки колонок.
func (s *Service) sanitizeHeaders(headers []string) []string {
	sanitized := make([]string, len(headers))
	usedNames := make(map[string]int)

	for i, header := range headers {
		// Санитизируем название
		name := strings.ToLower(header)
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

		// Если название пустое или начинается с цифры
		if name == "" || (len(name) > 0 && name[0] >= '0' && name[0] <= '9') {
			name = fmt.Sprintf("column_%d", i+1)
		}

		// Обрабатываем дубликаты
		if count, exists := usedNames[name]; exists {
			usedNames[name] = count + 1
			name = fmt.Sprintf("%s_%d", name, count+1)
		} else {
			usedNames[name] = 1
		}

		// Ограничиваем длину
		if len(name) > 63 {
			name = name[:63]
		}

		sanitized[i] = name
	}

	return sanitized
}

// insertSheetData вставляет данные листа в таблицу.
func (s *Service) insertSheetData(ctx context.Context, tx pgx.Tx, tableName string, columns []string, rows [][]string) error {
	if len(rows) == 0 {
		return nil
	}

	// Строим SQL запрос
	query := s.buildInsertQuery(tableName, columns, len(rows))

	// Подготавливаем аргументы
	args := make([]interface{}, 0, len(rows)*len(columns))
	for _, row := range rows {
		for i := range columns {
			var value interface{}
			if i < len(row) && row[i] != "" {
				value = row[i]
			} else {
				value = nil
			}
			args = append(args, value)
		}
	}

	// Выполняем запрос
	_, err := tx.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to insert data: %w", err)
	}

	return nil
}

// buildCreateTableQuery строит SQL запрос для создания таблицы.
func (s *Service) buildCreateTableQuery(tableName string, columns []model.DynamicTableColumn) string {
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
func (s *Service) buildInsertQuery(tableName string, columns []string, rowCount int) string {
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

// getCurrentQuarter возвращает текущий квартал.
func (s *Service) getCurrentQuarter() string {
	month := int(time.Now().Month())
	switch {
	case month >= 1 && month <= 3:
		return "Q1"
	case month >= 4 && month <= 6:
		return "Q2"
	case month >= 7 && month <= 9:
		return "Q3"
	default:
		return "Q4"
	}
}

// GetDynamicTables возвращает список динамически созданных таблиц.
func (s *Service) GetDynamicTables(ctx context.Context) ([]*model.ExcelTableMetadata, error) {
	return s.dynamicRepo.ListDynamicTables(ctx)
}

// GetTableMetadata возвращает метаданные таблицы.
func (s *Service) GetTableMetadata(ctx context.Context, tableName string) (*model.ExcelTableMetadata, error) {
	return s.dynamicRepo.GetTableMetadata(ctx, tableName)
}

// parseFileName парсит название файла и извлекает метаданные.
func (s *Service) parseFileName(fileName string) (*model.ExcelFileInfo, error) {
	// Убираем расширение
	name := strings.TrimSuffix(fileName, ".xlsx")
	name = strings.TrimSuffix(name, ".XLSX")

	// Извлекаем квартал и год
	quarter, year := s.extractQuarterAndYear(name)

	// Генерируем название таблицы
	tableName := fmt.Sprintf("dealer_net_%d_%s", year, strings.ToLower(quarter))

	return &model.ExcelFileInfo{
		FileName:  fileName,
		Quarter:   quarter,
		Year:      year,
		TableName: tableName,
	}, nil
}

// extractRegionFromSheetName извлекает регион из названия листа.
func (s *Service) extractRegionFromSheetName(sheetName string) (string, error) {
	// Формат: Q3-Central, Q3-NW, Q3-Volga, etc.
	parts := strings.Split(sheetName, "-")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid sheet name format: %s", sheetName)
	}

	region := strings.TrimSpace(parts[1])

	// Маппинг сокращений
	regionMap := map[string]string{
		"NW":      "North West",
		"Central": "Central",
		"Volga":   "Volga",
		"South":   "South",
		"Ural":    "Ural",
		"Siberia": "Siberia",
		"FE":      "Far East",
	}

	if mappedRegion, exists := regionMap[region]; exists {
		return mappedRegion, nil
	}

	return region, nil
}

// processSheetForUnifiedTable обрабатывает лист для объединенной таблицы.
func (s *Service) processSheetForUnifiedTable(ctx context.Context, tx pgx.Tx, f *excelize.File, sheetName string, region string) ([]model.ExcelRowWithRegion, []string, error) {
	s.logger.Info("Processing sheet for unified table",
		slog.String("sheet_name", sheetName),
		slog.String("region", region),
	)

	// Читаем данные листа
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get rows from sheet: %w", err)
	}

	s.logger.Info("Sheet rows read",
		slog.String("sheet_name", sheetName),
		slog.Int("total_rows", len(rows)),
	)

	if len(rows) < 3 {
		s.logger.Warn("Sheet has less than 3 rows", slog.String("sheet_name", sheetName))
		return []model.ExcelRowWithRegion{}, []string{}, nil
	}

	// Строка 2 (индекс 1) - заголовки колонок
	headers := rows[1]
	s.logger.Info("Headers row",
		slog.String("sheet_name", sheetName),
		slog.Int("headers_count", len(headers)),
		slog.String("headers", strings.Join(headers, "|")),
	)

	if len(headers) == 0 {
		return nil, nil, fmt.Errorf("sheet has no headers in row 2")
	}

	// Санитизируем заголовки и фильтруем пустые
	sanitizedHeaders := s.sanitizeAndFilterHeaders(headers)

	// Данные начинаются с строки 3 (индекс 2)
	dataRows := rows[2:]
	s.logger.Info("Data rows",
		slog.String("sheet_name", sheetName),
		slog.Int("data_rows_count", len(dataRows)),
	)

	var result []model.ExcelRowWithRegion
	for rowIndex, row := range dataRows {
		if len(row) == 0 {
			continue // Пропускаем пустые строки
		}

		values := make(map[string]interface{})
		for i, header := range sanitizedHeaders {
			var value interface{}
			if i < len(row) && row[i] != "" {
				value = row[i]
			} else {
				value = nil
			}
			values[header] = value
		}

		// Перезаписываем значение region из названия листа
		values["region"] = region

		// Проверяем, что дилер не пустой - если пустой, пропускаем запись
		dealerValue, exists := values["dealer"]
		if !exists || dealerValue == nil || dealerValue == "" {
			s.logger.Warn("Skipping row with empty dealer",
				slog.String("sheet_name", sheetName),
				slog.Int("row_index", rowIndex),
			)
			continue
		}

		result = append(result, model.ExcelRowWithRegion{
			Region: region,
			Values: values,
		})

		// Логируем первые несколько строк для отладки
		if rowIndex < 3 {
			s.logger.Info("Sample data row",
				slog.String("sheet_name", sheetName),
				slog.Int("row_index", rowIndex),
				slog.String("row_data", strings.Join(row, "|")),
			)
		}
	}

	s.logger.Info("Sheet processing completed",
		slog.String("sheet_name", sheetName),
		slog.Int("processed_rows", len(result)),
		slog.Int("columns_count", len(sanitizedHeaders)),
	)

	return result, sanitizedHeaders, nil
}

// sanitizeAndFilterHeaders санитизирует заголовки и фильтрует пустые.
func (s *Service) sanitizeAndFilterHeaders(headers []string) []string {
	var sanitized []string
	usedNames := make(map[string]int)

	for i, header := range headers {
		// Пропускаем пустые заголовки
		if strings.TrimSpace(header) == "" {
			continue
		}

		// Санитизируем название
		name := strings.ToLower(header)
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

		// Если название пустое или начинается с цифры
		if name == "" || (len(name) > 0 && name[0] >= '0' && name[0] <= '9') {
			name = fmt.Sprintf("column_%d", i+1)
		}

		// Обрабатываем дубликаты
		if count, exists := usedNames[name]; exists {
			usedNames[name] = count + 1
			name = fmt.Sprintf("%s_%d", name, count+1)
		} else {
			usedNames[name] = 1
		}

		// Ограничиваем длину
		if len(name) > 63 {
			name = name[:63]
		}

		sanitized = append(sanitized, name)
	}

	s.logger.Info("Headers sanitized",
		slog.Int("original_count", len(headers)),
		slog.Int("sanitized_count", len(sanitized)),
		slog.String("sanitized_headers", strings.Join(sanitized, ", ")),
	)

	return sanitized
}

// compareColumnStructures сравнивает структуры колонок.
func (s *Service) compareColumnStructures(cols1, cols2 []string) bool {
	if len(cols1) != len(cols2) {
		return false
	}

	for i, col1 := range cols1 {
		if col1 != cols2[i] {
			return false
		}
	}

	return true
}

// createUnifiedTable создает единую таблицу.
func (s *Service) createUnifiedTable(ctx context.Context, tx pgx.Tx, tableName string, columns []string) error {
	var columnDefs []string

	// Добавляем стандартные колонки
	columnDefs = append(columnDefs, "id BIGSERIAL PRIMARY KEY")

	// Добавляем колонки из схемы
	for _, col := range columns {
		if col == "dealer" {
			// Колонка dealer должна быть NOT NULL
			columnDefs = append(columnDefs, fmt.Sprintf("%s TEXT NOT NULL", col))
		} else {
			columnDefs = append(columnDefs, fmt.Sprintf("%s TEXT", col))
		}
	}

	// Добавляем колонку created_at
	columnDefs = append(columnDefs, "created_at TIMESTAMP DEFAULT NOW()")

	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (\n    %s\n)",
		tableName,
		strings.Join(columnDefs, ",\n    "))

	s.logger.Info("Creating unified table",
		slog.String("table", tableName),
		slog.String("query", query),
	)

	_, err := tx.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to create unified table %s: %w", tableName, err)
	}

	s.logger.Info("Unified table created successfully", slog.String("table", tableName))
	return nil
}

// insertUnifiedData вставляет данные в единую таблицу.
func (s *Service) insertUnifiedData(ctx context.Context, tx pgx.Tx, tableName string, columns []string, data []model.ExcelRowWithRegion) error {
	if len(data) == 0 {
		s.logger.Info("No data to insert", slog.String("table", tableName))
		return nil
	}

	// Используем колонки как есть
	allColumns := columns

	// Строим SQL запрос
	query := s.buildUnifiedInsertQuery(tableName, allColumns, len(data))

	s.logger.Info("Inserting unified data",
		slog.String("table", tableName),
		slog.Int("rows", len(data)),
		slog.Int("columns", len(allColumns)),
	)

	// Подготавливаем аргументы
	args := make([]interface{}, 0, len(data)*len(allColumns))
	for _, row := range data {
		// Добавляем значения колонок
		for _, col := range columns {
			value, exists := row.Values[col]
			if !exists || value == "" {
				args = append(args, nil)
			} else {
				args = append(args, value)
			}
		}
	}

	// Выполняем запрос
	_, err := tx.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to insert unified data into table %s: %w", tableName, err)
	}

	s.logger.Info("Unified data inserted successfully",
		slog.String("table", tableName),
		slog.Int("rows", len(data)),
	)

	return nil
}

// buildUnifiedInsertQuery строит SQL запрос для вставки данных в единую таблицу.
func (s *Service) buildUnifiedInsertQuery(tableName string, columns []string, rowCount int) string {
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
