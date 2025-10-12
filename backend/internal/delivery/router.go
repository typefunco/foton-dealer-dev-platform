package delivery

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/typefunco/dealer_dev_platform/internal/service/aftersales"
	"github.com/typefunco/dealer_dev_platform/internal/service/auth"
	"github.com/typefunco/dealer_dev_platform/internal/service/dealer"
	"github.com/typefunco/dealer_dev_platform/internal/service/dealerdev"
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
		srv:               echo.New(),
		logger:            logger,
	}
}

// RunServer - команда запуска сервера.
func (s *Server) RunServer() {
	s.srv.Use(middleware.Logger())
	s.srv.Use(middleware.Recover())

	// TODO: переделать на http
	s.srv.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
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
	api.GET("/aftersales", s.GetAfterSalesData) // Получить данные After Sales по региону

	// Dealer routes
	api.GET("/dealers", s.GetDealers)             // Получить список дилеров
	api.GET("/dealers/:id", s.GetDealerByID)      // Получить базовую информацию о дилере
	api.GET("/dealers/:id/card", s.GetDealerCard) // Получить полную карточку дилера

	// Dealer Development routes
	api.GET("/dealerdev", s.GetDealerDevData) // Получить данные Dealer Development

	// Performance routes
	api.GET("/performance", s.GetPerformanceData) // Получить данные производительности

	// Sales Team routes
	api.GET("/sales", s.GetSalesTeamData) // Получить данные команды продаж

	// Quarter Comparison routes
	api.GET("/quarter-comparison", s.GetQuarterComparison) // Сравнение кварталов

	// All Data routes (комплексные данные всех таблиц)
	api.GET("/all-data", s.GetAllData) // Получить все данные дилеров (DealerDev + Sales + Performance + AfterSales)

	if err := s.srv.Start(":8080"); err != nil && !errors.Is(err, http.ErrServerClosed) {
		s.logger.Error("failed to start server", "error", err)
	}
}
