package delivery

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/typefunco/dealer_dev_platform/internal/repository"
	"github.com/typefunco/dealer_dev_platform/internal/service/aftersales"
	"github.com/typefunco/dealer_dev_platform/internal/service/auth"
	"github.com/typefunco/dealer_dev_platform/internal/service/dealer"
	"github.com/typefunco/dealer_dev_platform/internal/service/dealerdev"
	"github.com/typefunco/dealer_dev_platform/internal/service/excel"
	"github.com/typefunco/dealer_dev_platform/internal/service/performance"
	"github.com/typefunco/dealer_dev_platform/internal/service/performance_aftersales"
	"github.com/typefunco/dealer_dev_platform/internal/service/performance_sales"
	"github.com/typefunco/dealer_dev_platform/internal/service/sales"
	"github.com/typefunco/dealer_dev_platform/internal/service/user"
)

// Server структура сервера.
type Server struct {
	authService       *auth.Service
	perfService       *performance.Service
	perfSalesService  *performance_sales.Service
	perfASService     *performance_aftersales.Service
	userService       *user.Service
	afterSalesService *aftersales.Service
	dealerService     *dealer.Service
	salesService      *sales.Service
	dealerDevService  *dealerdev.Service
	excelService      *excel.Service
	dynamicRepo       repository.DynamicTableRepository
	pool              *pgxpool.Pool
	maxFileSize       int64
	srv               *echo.Echo
	logger            *slog.Logger
}

// NewServer - конструктор сервера.
func NewServer(
	authService *auth.Service,
	perfService *performance.Service,
	perfSalesService *performance_sales.Service,
	perfASService *performance_aftersales.Service,
	userService *user.Service,
	afterSalesService *aftersales.Service,
	dealerService *dealer.Service,
	salesService *sales.Service,
	dealerDevService *dealerdev.Service,
	excelService *excel.Service,
	dynamicRepo repository.DynamicTableRepository,
	pool *pgxpool.Pool,
	maxFileSize int64,
	logger *slog.Logger,
) *Server {
	return &Server{
		authService:       authService,
		perfService:       perfService,
		perfSalesService:  perfSalesService,
		perfASService:     perfASService,
		userService:       userService,
		afterSalesService: afterSalesService,
		dealerService:     dealerService,
		salesService:      salesService,
		dealerDevService:  dealerDevService,
		excelService:      excelService,
		dynamicRepo:       dynamicRepo,
		pool:              pool,
		maxFileSize:       maxFileSize,
		srv:               echo.New(),
		logger:            logger,
	}
}

// RunServer - команда запуска сервера.
func (s *Server) RunServer() {
	s.srv.Use(middleware.Logger())
	s.srv.Use(middleware.Recover())

	// CORS configuration for Docker environment
	s.srv.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{
			"http://localhost:3000", // Development
			"http://frontend:3000",  // Docker frontend service
			"http://127.0.0.1:3000", // Alternative localhost
		},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization,
			echo.HeaderXRequestedWith,
		},
		AllowMethods: []string{
			http.MethodGet,
			http.MethodHead,
			http.MethodPut,
			http.MethodPatch,
			http.MethodPost,
			http.MethodDelete,
			http.MethodOptions,
		},
		AllowCredentials: true,
	}))

	// Auth routes
	s.srv.POST("/auth/login", s.Login)

	// Health check
	s.srv.GET("/health", s.Health)

	// API group
	api := s.srv.Group("/api")

	// User management routes
	api.GET("/users", s.GetUsers)           // Получить список пользователей с фильтрами
	api.GET("/users/stats", s.GetUserStats) // Получить статистику по регионам
	api.GET("/users/:id", s.GetUserByID)    // Получить пользователя по ID
	api.POST("/users", s.CreateUser)        // Создать пользователя
	api.PUT("/users/:id", s.UpdateUser)     // Обновить пользователя
	api.DELETE("/users/:id", s.DeleteUser)  // Удалить пользователя

	// After Sales routes
	api.GET("/aftersales", s.GetAfterSalesData) // Получить данные After Sales по региону (legacy)

	// Dealer routes
	api.GET("/dealers/list", s.GetDealersList)    // Получить упрощенный список дилеров для UI
	api.GET("/dealers", s.GetDealers)             // Получить список дилеров
	api.GET("/dealers/:id", s.GetDealerByID)      // Получить базовую информацию о дилере
	api.GET("/dealers/:id/card", s.GetDealerCard) // Получить полную карточку дилера

	// Унифицированные маршруты для всех типов таблиц (без префикса dynamic)
	api.GET("/dealer_dev", s.GetDynamicData)  // Dealer Development
	api.GET("/sales", s.GetDynamicData)       // Sales Team
	api.GET("/after_sales", s.GetDynamicData) // After Sales
	api.GET("/performance", s.GetDynamicData) // Performance

	// Legacy routes (сохраняем для обратной совместимости)
	api.GET("/dealerdev", s.GetDealerDevData) // Получить данные Dealer Development

	// Quarter Comparison routes
	api.GET("/quarter-comparison", s.GetQuarterComparison) // Сравнение кварталов

	// All Data routes (комплексные данные всех таблиц)
	api.GET("/all-data", s.GetAllData) // Получить все данные дилеров (DealerDev + Sales + Performance + AfterSales)

	// Filter routes
	api.GET("/filters", s.GetAvailableFilters) // Получить доступные фильтры

	// Analytics routes
	api.GET("/analytics", s.GetAnalytics) // Получить аналитические данные

	// Bulk operations routes
	api.POST("/bulk", s.BulkOperations)    // Массовые операции
	api.POST("/bulk/update", s.BulkUpdate) // Массовое обновление
	api.POST("/bulk/export", s.BulkExport) // Массовый экспорт

	// Excel operations routes
	api.POST("/excel/upload", s.UploadExcelFile)                  // Загрузка Excel файла
	api.GET("/excel/tables", s.GetExcelTables)                    // Список созданных таблиц
	api.GET("/excel/tables/:tableName", s.GetExcelTableMetadata)  // Метаданные таблицы
	api.GET("/excel/tables/:tableName/data", s.GetExcelTableData) // Данные таблицы
	api.DELETE("/excel/tables/:tableName", s.DeleteExcelTable)    // Удаление таблицы

	if err := s.srv.Start(":8080"); err != nil && !errors.Is(err, http.ErrServerClosed) {
		s.logger.Error("failed to start server", "error", err)
	}
}
