package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/typefunco/dealer_dev_platform/internal/model"
	"github.com/typefunco/dealer_dev_platform/internal/repository"
	"github.com/typefunco/dealer_dev_platform/internal/service/aftersales"
	"github.com/typefunco/dealer_dev_platform/internal/service/dealer"
	"github.com/typefunco/dealer_dev_platform/internal/service/performance"
)

func main() {
	fmt.Println("🚀 Quick Filter Test")
	fmt.Println("===================")

	// Подключение к базе данных
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL environment variable is required for testing")
	}

	pool, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	fmt.Println("✅ Database connected successfully")

	// Создаем логгер
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	// Инициализируем репозитории
	dealerRepo := repository.NewDealerRepository(pool)
	excelDealerRepo := repository.NewExcelDealerRepository(pool, logger)
	afterSalesRepo := repository.NewAfterSalesRepository(pool)
	performanceRepo := repository.NewPerformanceRepository(pool)

	// Инициализируем сервисы
	dealerService := dealer.NewService(dealerRepo, excelDealerRepo, logger)
	afterSalesService := aftersales.NewService(afterSalesRepo, excelDealerRepo, logger)
	performanceService := performance.NewService(performanceRepo, logger)

	ctx := context.Background()

	// Тест 1: Фильтр по региону
	fmt.Println("\n=== Test 1: Region Filter ===")
	filters := &model.FilterParams{
		Region: "central",
		Limit:  5,
	}

	dealers, err := dealerService.GetDealersWithFilters(ctx, filters)
	if err != nil {
		fmt.Printf("❌ Dealers error: %v\n", err)
		return
	}
	fmt.Printf("✅ Found %d dealers in central region\n", len(dealers))

	// Тест 2: Фильтр по году и кварталу
	fmt.Println("\n=== Test 2: Year and Quarter Filter ===")
	filters2 := &model.FilterParams{
		Quarter: "Q1",
		Year:    2025,
		Limit:   3,
	}

	afterSales, err := afterSalesService.GetAfterSalesWithFilters(ctx, filters2)
	if err != nil {
		fmt.Printf("❌ AfterSales error: %v\n", err)
		return
	}
	fmt.Printf("✅ Found %d after sales records for Q1 2025\n", len(afterSales))

	// Тест 3: Комбинированные фильтры
	fmt.Println("\n=== Test 3: Combined Filters ===")
	filters3 := &model.FilterParams{
		Region:  "central",
		Quarter: "Q2",
		Year:    2025,
		Limit:   2,
	}

	performance, err := performanceService.GetPerformanceWithFilters(ctx, filters3)
	if err != nil {
		fmt.Printf("❌ Performance error: %v\n", err)
		return
	}
	fmt.Printf("✅ Found %d performance records in central region for Q2 2025\n", len(performance))

	// Тест 4: Пагинация
	fmt.Println("\n=== Test 4: Pagination ===")
	filters4 := &model.FilterParams{
		Limit:  3,
		Offset: 6,
	}

	dealers2, err := dealerService.GetDealersWithFilters(ctx, filters4)
	if err != nil {
		fmt.Printf("❌ Dealers error: %v\n", err)
		return
	}
	fmt.Printf("✅ Found %d dealers with pagination (limit=3, offset=6)\n", len(dealers2))

	// Тест 5: Фильтр по ID дилеров
	fmt.Println("\n=== Test 5: Dealer IDs Filter ===")
	filters5 := &model.FilterParams{
		DealerIDs: []int{1, 2, 3, 4},
		Limit:     10,
	}

	dealers3, err := dealerService.GetDealersWithFilters(ctx, filters5)
	if err != nil {
		fmt.Printf("❌ Dealers error: %v\n", err)
		return
	}
	fmt.Printf("✅ Found %d dealers with specific IDs [1,2,3,4]\n", len(dealers3))

	fmt.Println("\n🎉 All filter tests completed successfully!")
}
