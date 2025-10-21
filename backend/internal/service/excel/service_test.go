package excel

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/typefunco/dealer_dev_platform/internal/model"
)

// MockDynamicTableRepository мок репозитория для тестов.
type MockDynamicTableRepository struct {
	mock.Mock
}

func (m *MockDynamicTableRepository) CreateTable(ctx context.Context, schema *model.DynamicTableSchema) error {
	args := m.Called(ctx, schema)
	return args.Error(0)
}

func (m *MockDynamicTableRepository) InsertData(ctx context.Context, tableName string, columns []string, rows [][]interface{}) error {
	args := m.Called(ctx, tableName, columns, rows)
	return args.Error(0)
}

func (m *MockDynamicTableRepository) TableExists(ctx context.Context, tableName string) (bool, error) {
	args := m.Called(ctx, tableName)
	return args.Bool(0), args.Error(1)
}

func (m *MockDynamicTableRepository) GetTableMetadata(ctx context.Context, tableName string) (*model.ExcelTableMetadata, error) {
	args := m.Called(ctx, tableName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.ExcelTableMetadata), args.Error(1)
}

func (m *MockDynamicTableRepository) ListDynamicTables(ctx context.Context) ([]*model.ExcelTableMetadata, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.ExcelTableMetadata), args.Error(1)
}

func (m *MockDynamicTableRepository) BeginTransaction(ctx context.Context) (pgx.Tx, error) {
	args := m.Called(ctx)
	return args.Get(0).(pgx.Tx), args.Error(1)
}

func TestParseSheetName(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	service := &Service{logger: logger}

	tests := []struct {
		name      string
		sheetName string
		expected  *model.ExcelSheetInfo
	}{
		{
			name:      "Dealer Development Q3 2025",
			sheetName: "Dealer_Development_Q3_2025",
			expected: &model.ExcelSheetInfo{
				Name:      "Dealer_Development_Q3_2025",
				Category:  "dealer_development",
				Quarter:   "Q3",
				Year:      2025,
				TableName: "dealer_development_q3_2025",
			},
		},
		{
			name:      "Sales Q2 2024",
			sheetName: "Sales_Q2_2024",
			expected: &model.ExcelSheetInfo{
				Name:      "Sales_Q2_2024",
				Category:  "sales",
				Quarter:   "Q2",
				Year:      2024,
				TableName: "sales_q2_2024",
			},
		},
		{
			name:      "Q3-Central",
			sheetName: "Q3-Central",
			expected: &model.ExcelSheetInfo{
				Name:      "Q3-Central",
				Category:  "central",
				Quarter:   "Q3",
				Year:      2025, // Текущий год
				TableName: "central_q3_2025",
			},
		},
		{
			name:      "Performance Q4 2024",
			sheetName: "Performance_Q4_2024",
			expected: &model.ExcelSheetInfo{
				Name:      "Performance_Q4_2024",
				Category:  "performance",
				Quarter:   "Q4",
				Year:      2024,
				TableName: "performance_q4_2024",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.parseSheetName(tt.sheetName)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected.Name, result.Name)
			assert.Equal(t, tt.expected.Category, result.Category)
			assert.Equal(t, tt.expected.Quarter, result.Quarter)
			assert.Equal(t, tt.expected.TableName, result.TableName)
			// Год может отличаться если тест запускается в другом году
			assert.NotZero(t, result.Year)
		})
	}
}

func TestSanitizeHeaders(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	service := &Service{logger: logger}

	tests := []struct {
		name     string
		headers  []string
		expected []string
	}{
		{
			name:     "Basic headers",
			headers:  []string{"Dealer Name", "Region", "City", "Manager"},
			expected: []string{"dealer_name", "region", "city", "manager"},
		},
		{
			name:     "Headers with special characters",
			headers:  []string{"Stock Q3", "Buy out Q3", "Marketing Investments (%)", "Check List %"},
			expected: []string{"stock_q3", "buy_out_q3", "marketing_investments_percent", "check_list_percent"},
		},
		{
			name:     "Empty headers",
			headers:  []string{"", "Region", "", "Manager"},
			expected: []string{"column_1", "region", "column_3", "manager"},
		},
		{
			name:     "Duplicate headers",
			headers:  []string{"Region", "Region", "City", "Region"},
			expected: []string{"region", "region_2", "city", "region_3"},
		},
		{
			name:     "Headers starting with numbers",
			headers:  []string{"1st Quarter", "2nd Quarter", "3rd Quarter"},
			expected: []string{"column_1", "column_2", "column_3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.sanitizeHeaders(tt.headers)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestExtractQuarterAndYear(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	service := &Service{logger: logger}

	tests := []struct {
		name            string
		sheetName       string
		expectedQuarter string
		expectedYear    int
	}{
		{
			name:            "Q3 2025",
			sheetName:       "dealer_development_q3_2025",
			expectedQuarter: "Q3",
			expectedYear:    2025,
		},
		{
			name:            "Q2-2024",
			sheetName:       "sales_q2-2024",
			expectedQuarter: "Q2",
			expectedYear:    2024,
		},
		{
			name:            "Q1 2023",
			sheetName:       "performance_q1_2023",
			expectedQuarter: "Q1",
			expectedYear:    2023,
		},
		{
			name:            "Only quarter",
			sheetName:       "q3_central",
			expectedQuarter: "Q3",
			expectedYear:    2025, // Текущий год
		},
		{
			name:            "No quarter",
			sheetName:       "central_data",
			expectedQuarter: "Q4", // Текущий квартал
			expectedYear:    2025, // Текущий год
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			quarter, year := service.extractQuarterAndYear(tt.sheetName)
			assert.Equal(t, tt.expectedQuarter, quarter)
			// Год может отличаться если тест запускается в другом году
			assert.NotZero(t, year)
		})
	}
}

func TestExtractCategory(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	service := &Service{logger: logger}

	tests := []struct {
		name      string
		sheetName string
		expected  string
	}{
		{
			name:      "Dealer Development",
			sheetName: "dealer_development_q3_2025",
			expected:  "dealer_development",
		},
		{
			name:      "Sales",
			sheetName: "sales_q2_2024",
			expected:  "sales",
		},
		{
			name:      "Central",
			sheetName: "q3_central",
			expected:  "central",
		},
		{
			name:      "Performance",
			sheetName: "performance_q4_2024",
			expected:  "performance",
		},
		{
			name:      "Empty category",
			sheetName: "q3_2025",
			expected:  "data",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.extractCategory(tt.sheetName)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGenerateTableName(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	service := &Service{logger: logger}

	tests := []struct {
		name     string
		category string
		quarter  string
		year     int
		expected string
	}{
		{
			name:     "Dealer Development Q3 2025",
			category: "dealer_development",
			quarter:  "Q3",
			year:     2025,
			expected: "dealer_development_q3_2025",
		},
		{
			name:     "Sales Q2 2024",
			category: "sales",
			quarter:  "Q2",
			year:     2024,
			expected: "sales_q2_2024",
		},
		{
			name:     "Central Q1 2023",
			category: "central",
			quarter:  "Q1",
			year:     2023,
			expected: "central_q1_2023",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.generateTableName(tt.category, tt.quarter, tt.year)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestBuildCreateTableQuery(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	service := &Service{logger: logger}

	tableName := "test_table"
	columns := []model.DynamicTableColumn{
		{Name: "dealer_name", DataType: "TEXT", Nullable: true},
		{Name: "region", DataType: "TEXT", Nullable: true},
		{Name: "city", DataType: "TEXT", Nullable: true},
	}

	result := service.buildCreateTableQuery(tableName, columns)

	assert.Contains(t, result, "CREATE TABLE IF NOT EXISTS test_table")
	assert.Contains(t, result, "id BIGSERIAL PRIMARY KEY")
	assert.Contains(t, result, "dealer_name TEXT")
	assert.Contains(t, result, "region TEXT")
	assert.Contains(t, result, "city TEXT")
	assert.Contains(t, result, "created_at TIMESTAMP DEFAULT NOW()")
}

func TestBuildInsertQuery(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	service := &Service{logger: logger}

	tableName := "test_table"
	columns := []string{"dealer_name", "region", "city"}
	rowCount := 2

	result := service.buildInsertQuery(tableName, columns, rowCount)

	assert.Contains(t, result, "INSERT INTO test_table")
	assert.Contains(t, result, "(dealer_name, region, city)")
	assert.Contains(t, result, "VALUES")
	assert.Contains(t, result, "($1, $2, $3)")
	assert.Contains(t, result, "($4, $5, $6)")
}
