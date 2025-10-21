package model

import (
	"time"
)

// ExcelSheetInfo содержит метаданные о листе Excel.
type ExcelSheetInfo struct {
	Name      string `json:"name"`       // Название листа
	Category  string `json:"category"`   // Категория (извлеченная из названия)
	Quarter   string `json:"quarter"`    // Квартал (Q1, Q2, Q3, Q4)
	Year      int    `json:"year"`       // Год
	TableName string `json:"table_name"` // Название таблицы для БД
}

// ExcelColumnInfo содержит информацию о колонке Excel.
type ExcelColumnInfo struct {
	Name          string `json:"name"`           // Оригинальное название
	SanitizedName string `json:"sanitized_name"` // Санитизированное название для БД
	Index         int    `json:"index"`          // Индекс колонки
}

// ExcelRowWithRegion представляет строку данных из Excel с регионом.
type ExcelRowWithRegion struct {
	Region string                 `json:"region"`
	Values map[string]interface{} `json:"values"` // Значения по названиям колонок
}

// ExcelFileInfo содержит метаданные о файле Excel.
type ExcelFileInfo struct {
	FileName  string `json:"file_name"`
	Quarter   string `json:"quarter"`    // Квартал (Q1, Q2, Q3, Q4)
	Year      int    `json:"year"`       // Год
	TableName string `json:"table_name"` // Название таблицы для БД
}

// ExcelTableData содержит данные таблицы для создания в БД.
type ExcelTableData struct {
	SheetInfo ExcelSheetInfo       `json:"sheet_info"`
	Columns   []ExcelColumnInfo    `json:"columns"`
	Rows      []ExcelRowWithRegion `json:"rows"`
}

// ExcelUploadRequest представляет запрос на загрузку Excel файла.
type ExcelUploadRequest struct {
	FileName string `json:"file_name"`
	FileSize int64  `json:"file_size"`
}

// ExcelUploadResponse представляет ответ на загрузку Excel файла.
type ExcelUploadResponse struct {
	Status         string           `json:"status"`
	Message        string           `json:"message"`
	TablesCreated  []string         `json:"tables_created"`
	RowsInserted   int              `json:"rows_inserted"`
	ProcessingTime time.Duration    `json:"processing_time"`
	Details        []ExcelTableData `json:"details,omitempty"`
}

// ExcelTableMetadata содержит метаданные о созданной таблице.
type ExcelTableMetadata struct {
	TableName string    `json:"table_name"`
	RowsCount int       `json:"rows_count"`
	CreatedAt time.Time `json:"created_at"`
	Columns   []string  `json:"columns"`
}

// ExcelError представляет ошибку при обработке Excel файла.
type ExcelError struct {
	SheetName string `json:"sheet_name"`
	Row       int    `json:"row,omitempty"`
	Column    string `json:"column,omitempty"`
	Message   string `json:"message"`
	Error     string `json:"error"`
}

// ExcelProcessingResult содержит результат обработки Excel файла.
type ExcelProcessingResult struct {
	Success        bool                 `json:"success"`
	TablesCreated  []ExcelTableMetadata `json:"tables_created"`
	Errors         []ExcelError         `json:"errors,omitempty"`
	TotalRows      int                  `json:"total_rows"`
	ProcessingTime time.Duration        `json:"processing_time"`
}

// DynamicTableColumn представляет колонку динамически созданной таблицы.
type DynamicTableColumn struct {
	Name     string `json:"name"`
	DataType string `json:"data_type"`
	Nullable bool   `json:"nullable"`
}

// DynamicTableSchema представляет схему динамически созданной таблицы.
type DynamicTableSchema struct {
	TableName string               `json:"table_name"`
	Columns   []DynamicTableColumn `json:"columns"`
}
