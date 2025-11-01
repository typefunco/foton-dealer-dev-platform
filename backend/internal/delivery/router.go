package delivery

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	authMiddleware "github.com/typefunco/dealer_dev_platform/internal/middleware"
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
	"github.com/typefunco/dealer_dev_platform/internal/utils/jwt"
)

// Server структура сервера.
type Server struct {
	authService       *auth.Service
	jwtService        *jwt.Service
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
	jwtService *jwt.Service,
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
		jwtService:        jwtService,
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
			"http://localhost:3000",      // Development
			"http://frontend:3000",       // Docker frontend service
			"http://127.0.0.1:3000",      // Alternative localhost
			"http://localhost",           // Generic localhost
			"http://frontend",            // Docker service name
			"http://localhost:80",        // Production HTTP
			"http://localhost:443",       // Production HTTPS
			"https://localhost",          // Production HTTPS
			"https://localhost:443",      // Production HTTPS
			"http://77.223.123.147:3000", // Production server
			"http://77.223.123.147",      // Production server without port
		},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization,
			echo.HeaderXRequestedWith,
			echo.HeaderXRealIP,
			echo.HeaderXForwardedFor,
			echo.HeaderXForwardedProto,
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
		MaxAge:           86400, // Cache preflight for 24 hours
	}))

	// Auth routes (без middleware)
	s.srv.POST("/auth/login", s.Login)

	// Health check (без middleware)
	s.srv.GET("/health", s.Health)

	// API group с обязательной аутентификацией
	api := s.srv.Group("/api")
	api.Use(authMiddleware.AuthMiddleware(s.jwtService))

	// User management routes (только чтение для всех пользователей)
	api.GET("/users", s.GetUsers)           // Получить список пользователей с фильтрами
	api.GET("/users/stats", s.GetUserStats) // Получить статистику по регионам
	api.GET("/users/:id", s.GetUserByID)    // Получить пользователя по ID

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

	// Admin routes (только для администраторов)
	admin := api.Group("/admin")
	admin.Use(authMiddleware.AdminMiddleware())

	// Excel operations routes (только для админов)
	admin.POST("/excel/upload", s.UploadExcelFile)                  // Загрузка Excel файла
	admin.POST("/excel/brands/upload", s.UploadBrandsFile)          // Загрузка файла с брендами и побочными бизнесами
	admin.GET("/excel/tables", s.GetExcelTables)                    // Список созданных таблиц
	admin.GET("/excel/tables/:tableName", s.GetExcelTableMetadata)  // Метаданные таблицы
	admin.GET("/excel/tables/:tableName/data", s.GetExcelTableData) // Данные таблицы
	admin.DELETE("/excel/tables/:tableName", s.DeleteExcelTable)    // Удаление таблицы

	// User management routes (только для админов)
	admin.POST("/users", s.CreateUser)       // Создать пользователя
	admin.PUT("/users/:id", s.UpdateUser)    // Обновить пользователя
	admin.DELETE("/users/:id", s.DeleteUser) // Удалить пользователя

	// Bulk operations routes (только для админов)
	admin.POST("/bulk", s.BulkOperations)    // Массовые операции
	admin.POST("/bulk/update", s.BulkUpdate) // Массовое обновление
	admin.POST("/bulk/export", s.BulkExport) // Массовый экспорт

	if err := s.srv.Start(":8080"); err != nil && !errors.Is(err, http.ErrServerClosed) {
		s.logger.Error("failed to start server", "error", err)
	}
}
