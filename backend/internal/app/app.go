package app

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/typefunco/dealer_dev_platform/internal/config"
	"github.com/typefunco/dealer_dev_platform/internal/database"
	"github.com/typefunco/dealer_dev_platform/internal/delivery"
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

// RunApp запускает приложение.
func RunApp() {
	// Инициализация логгера
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	logger.Info("Starting Dealer Development Platform")

	// Загрузка конфигурации
	cfg, err := config.Load()
	if err != nil {
		logger.Error("Failed to load configuration", slog.String("error", err.Error()))
		os.Exit(1)
	}

	logger.Info("Configuration loaded",
		slog.String("server_port", cfg.ServerPort),
		slog.String("database_url", config.MaskDSN(cfg.DatabaseURL)),
	)

	// Подключение к базе данных
	ctx := context.Background()
	dbConfig := database.DefaultPostgresConfig()
	dbConfig.MaxConns = cfg.DBMaxConns
	pool, err := database.NewPostgresPool(ctx, cfg.DatabaseURL, dbConfig)
	if err != nil {
		logger.Error("Failed to connect to database", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer pool.Close()

	logger.Info("Database connection established")

	// Инициализация компонентов
	if err := run(ctx, pool, cfg, logger); err != nil {
		logger.Error("Application error", slog.String("error", err.Error()))
		os.Exit(1)
	}

	logger.Info("Application stopped gracefully")
}

// run содержит основную логику приложения.
func run(ctx context.Context, pool *pgxpool.Pool, cfg *config.Config, logger *slog.Logger) error {
	// Инициализация репозиториев
	dealerRepo := repository.NewDealerRepository(pool)
	dealerDevRepo := repository.NewDealerDevRepository(pool)
	salesRepo := repository.NewSalesRepository(pool)
	performanceRepo := repository.NewPerformanceRepository(pool)
	performanceSalesRepo := repository.NewPerformanceSalesRepository(pool, logger)
	performanceASRepo := repository.NewPerformanceAfterSalesRepository(pool, logger)
	afterSalesRepo := repository.NewAfterSalesRepository(pool)
	authRepo := repository.NewAuthRepository(pool, logger)
	userRepo := repository.NewUserRepository(pool, logger)
	dynamicRepo := repository.NewDynamicTableRepository(pool, logger)

	logger.Info("Repositories initialized")

	// Инициализация сервисов
	jwtService := jwt.NewService()
	authService := auth.NewService(authRepo, jwtService, logger)
	perfService := performance.NewService(performanceRepo, logger)
	perfSalesService := performance_sales.NewService(performanceSalesRepo, logger)
	perfASService := performance_aftersales.NewService(performanceASRepo, logger)
	userService := user.NewService(userRepo, logger)
	afterSalesService := aftersales.NewService(afterSalesRepo, logger)
	dealerService := dealer.NewService(dealerRepo, logger)
	salesService := sales.NewService(salesRepo, logger)
	dealerDevService := dealerdev.NewService(dealerDevRepo, logger)
	excelService := excel.NewService(dynamicRepo, logger)

	logger.Info("Services initialized")

	// Инициализация HTTP сервера
	server := delivery.NewServer(authService, perfService, perfSalesService, perfASService, userService, afterSalesService, dealerService, salesService, dealerDevService, excelService, dynamicRepo, pool, cfg.MaxFileSize, logger)
	logger.Info("HTTP server initialized", slog.String("port", cfg.ServerPort))

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Запуск сервера в отдельной горутине
	go func() {
		server.RunServer()
	}()

	logger.Info("Server started successfully")

	// Ожидание сигнала завершения
	sig := <-quit
	logger.Info("Shutdown signal received", slog.String("signal", sig.String()))

	// Graceful shutdown с таймаутом
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return shutdown(shutdownCtx, pool, logger)
}

// shutdown выполняет graceful shutdown приложения.
func shutdown(ctx context.Context, pool *pgxpool.Pool, logger *slog.Logger) error {
	logger.Info("Starting graceful shutdown...")

	// TODO: Добавить shutdown для HTTP сервера когда будет реализован

	// Закрываем пул соединений БД
	pool.Close()
	logger.Info("Database connections closed")

	return nil
}
