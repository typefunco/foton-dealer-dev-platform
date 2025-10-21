package delivery

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/typefunco/dealer_dev_platform/internal/model"
)

// UploadExcelFile обрабатывает загрузку Excel файла.
// @Summary Upload Excel file
// @Description Загружает Excel файл и конвертирует его в таблицы PostgreSQL
// @Tags excel
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Excel file (.xlsx)"
// @Success 200 {object} model.ExcelUploadResponse
// @Failure 400 {object} ErrorResponse
// @Failure 413 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/excel/upload [post]
func (s *Server) UploadExcelFile(c echo.Context) error {
	// Получаем файл из формы
	file, err := c.FormFile("file")
	if err != nil {
		s.logger.Error("Failed to get file from form", slog.String("error", err.Error()))
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "No file provided",
		})
	}

	// Проверяем размер файла
	if file.Size > s.maxFileSize {
		s.logger.Error("File too large",
			slog.String("file_name", file.Filename),
			slog.Int64("file_size", file.Size),
			slog.Int64("max_size", s.maxFileSize),
		)
		return c.JSON(http.StatusRequestEntityTooLarge, ErrorResponse{
			Error: fmt.Sprintf("File size exceeds maximum allowed size of %d MB", s.maxFileSize/(1024*1024)),
		})
	}

	// Проверяем расширение файла
	if !isExcelFile(file.Filename) {
		s.logger.Error("Invalid file type", slog.String("file_name", file.Filename))
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Only .xlsx files are supported",
		})
	}

	// Открываем файл
	src, err := file.Open()
	if err != nil {
		s.logger.Error("Failed to open uploaded file",
			slog.String("file_name", file.Filename),
			slog.String("error", err.Error()),
		)
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to open uploaded file",
		})
	}
	defer src.Close()

	s.logger.Info("Processing Excel file",
		slog.String("file_name", file.Filename),
		slog.Int64("file_size", file.Size),
	)

	// Обрабатываем файл
	result, err := s.excelService.ProcessExcelFile(c.Request().Context(), src, file.Filename)
	if err != nil {
		s.logger.Error("Failed to process Excel file",
			slog.String("file_name", file.Filename),
			slog.String("error", err.Error()),
		)
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to process Excel file",
		})
	}

	// Формируем ответ
	response := model.ExcelUploadResponse{
		Status:         "success",
		Message:        "Файл успешно обработан и данные загружены в БД",
		TablesCreated:  make([]string, len(result.TablesCreated)),
		RowsInserted:   result.TotalRows,
		ProcessingTime: result.ProcessingTime,
	}

	for i, table := range result.TablesCreated {
		response.TablesCreated[i] = table.TableName
	}

	s.logger.Info("Excel file processed successfully",
		slog.String("file_name", file.Filename),
		slog.Int("tables_created", len(result.TablesCreated)),
		slog.Int("rows_inserted", result.TotalRows),
		slog.Duration("processing_time", result.ProcessingTime),
	)

	return c.JSON(http.StatusOK, response)
}

// GetExcelTables возвращает список созданных таблиц из Excel файлов.
// @Summary Get Excel tables
// @Description Возвращает список всех динамически созданных таблиц из Excel файлов
// @Tags excel
// @Produce json
// @Success 200 {array} model.ExcelTableMetadata
// @Failure 500 {object} ErrorResponse
// @Router /api/excel/tables [get]
func (s *Server) GetExcelTables(c echo.Context) error {
	tables, err := s.excelService.GetDynamicTables(c.Request().Context())
	if err != nil {
		s.logger.Error("Failed to get dynamic tables", slog.String("error", err.Error()))
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to get dynamic tables",
		})
	}

	return c.JSON(http.StatusOK, tables)
}

// GetExcelTableMetadata возвращает метаданные конкретной таблицы.
// @Summary Get Excel table metadata
// @Description Возвращает метаданные указанной таблицы
// @Tags excel
// @Produce json
// @Param tableName path string true "Table Name"
// @Success 200 {object} model.ExcelTableMetadata
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/excel/tables/{tableName} [get]
func (s *Server) GetExcelTableMetadata(c echo.Context) error {
	tableName := c.Param("tableName")
	if tableName == "" {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Table name is required",
		})
	}

	metadata, err := s.excelService.GetTableMetadata(c.Request().Context(), tableName)
	if err != nil {
		s.logger.Error("Failed to get table metadata",
			slog.String("table_name", tableName),
			slog.String("error", err.Error()),
		)
		return c.JSON(http.StatusNotFound, ErrorResponse{
			Error: "Table not found",
		})
	}

	return c.JSON(http.StatusOK, metadata)
}

// GetExcelTableData возвращает данные из указанной таблицы.
// @Summary Get Excel table data
// @Description Возвращает данные из указанной таблицы с пагинацией
// @Tags excel
// @Produce json
// @Param tableName path string true "Table Name"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(50)
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/excel/tables/{tableName}/data [get]
func (s *Server) GetExcelTableData(c echo.Context) error {
	tableName := c.Param("tableName")
	if tableName == "" {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Table name is required",
		})
	}

	// Получаем параметры пагинации
	pageStr := c.QueryParam("page")
	page := 1
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	limitStr := c.QueryParam("limit")
	limit := 50
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 1000 {
			limit = l
		}
	}

	offset := (page - 1) * limit

	// Получаем данные из таблицы
	data, err := s.getTableData(c.Request().Context(), tableName, limit, offset)
	if err != nil {
		s.logger.Error("Failed to get table data",
			slog.String("table_name", tableName),
			slog.String("error", err.Error()),
		)
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to get table data",
		})
	}

	// Получаем общее количество записей
	totalCount, err := s.getTableCount(c.Request().Context(), tableName)
	if err != nil {
		s.logger.Error("Failed to get table count",
			slog.String("table_name", tableName),
			slog.String("error", err.Error()),
		)
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to get table count",
		})
	}

	response := map[string]interface{}{
		"data": data,
		"pagination": map[string]interface{}{
			"page":        page,
			"limit":       limit,
			"total_count": totalCount,
			"total_pages": (totalCount + limit - 1) / limit,
		},
	}

	return c.JSON(http.StatusOK, response)
}

// DeleteExcelTable удаляет указанную таблицу.
// @Summary Delete Excel table
// @Description Удаляет указанную динамически созданную таблицу
// @Tags excel
// @Produce json
// @Param tableName path string true "Table Name"
// @Success 200 {object} map[string]string
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/excel/tables/{tableName} [delete]
func (s *Server) DeleteExcelTable(c echo.Context) error {
	tableName := c.Param("tableName")
	if tableName == "" {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Table name is required",
		})
	}

	// Проверяем, что таблица существует
	exists, err := s.dynamicRepo.TableExists(c.Request().Context(), tableName)
	if err != nil {
		s.logger.Error("Failed to check table existence",
			slog.String("table_name", tableName),
			slog.String("error", err.Error()),
		)
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to check table existence",
		})
	}

	if !exists {
		return c.JSON(http.StatusNotFound, ErrorResponse{
			Error: "Table not found",
		})
	}

	// Удаляем таблицу
	query := fmt.Sprintf("DROP TABLE IF EXISTS %s", tableName)
	_, err = s.pool.Exec(c.Request().Context(), query)
	if err != nil {
		s.logger.Error("Failed to delete table",
			slog.String("table_name", tableName),
			slog.String("error", err.Error()),
		)
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to delete table",
		})
	}

	s.logger.Info("Table deleted successfully", slog.String("table_name", tableName))

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Table deleted successfully",
	})
}

// isExcelFile проверяет, является ли файл Excel файлом.
func isExcelFile(filename string) bool {
	return len(filename) > 5 && filename[len(filename)-5:] == ".xlsx"
}

// getTableData получает данные из таблицы с пагинацией.
func (s *Server) getTableData(ctx context.Context, tableName string, limit, offset int) ([]map[string]interface{}, error) {
	query := fmt.Sprintf("SELECT * FROM %s ORDER BY id LIMIT $1 OFFSET $2", tableName)

	rows, err := s.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []map[string]interface{}
	columns := rows.FieldDescriptions()

	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}

		row := make(map[string]interface{})
		for i, col := range columns {
			row[col.Name] = values[i]
		}
		result = append(result, row)
	}

	return result, rows.Err()
}

// getTableCount получает общее количество записей в таблице.
func (s *Server) getTableCount(ctx context.Context, tableName string) (int, error) {
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", tableName)

	var count int
	err := s.pool.QueryRow(ctx, query).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
