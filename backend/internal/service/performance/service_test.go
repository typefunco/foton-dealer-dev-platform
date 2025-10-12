package performance_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/typefunco/dealer_dev_platform/internal/repository"
	"github.com/typefunco/dealer_dev_platform/internal/service/performance"
	"github.com/typefunco/dealer_dev_platform/internal/testutil"
)

func TestPerformanceService_CreatePerformance(t *testing.T) {
	// Настройка тестовой базы данных
	testDB := testutil.SetupTestDB(t)
	defer testDB.Cleanup(t)
	testDB.RunMigrations(t)

	logger := testutil.GetTestLogger()
	repo := repository.NewPerformanceRepository(testDB.Pool)
	service := performance.NewService(repo, logger)

	ctx := context.Background()

	t.Run("successful creation", func(t *testing.T) {
		defer testDB.CleanupTable(t, "performance_sales")

		// Создаем тестового дилера
		dealerRepo := repository.NewDealerRepository(testDB.Pool)
		dealer := testutil.CreateTestDealer()
		dealerID, err := dealerRepo.Create(ctx, dealer)
		require.NoError(t, err)

		// Создаем тестовые данные производительности
		perfData := testutil.CreateTestPerformanceSales(dealerID)

		// Создаем запись
		id, err := service.CreatePerformance(ctx, perfData)
		require.NoError(t, err)
		assert.Greater(t, id, int64(0))

		// Проверяем, что запись создана
		created, err := service.GetPerformanceByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, dealerID, created.DealerID)
		assert.Equal(t, perfData.Period, created.Period)
		assert.Equal(t, *perfData.QuantitySold, *created.QuantitySold)
		assert.Equal(t, *perfData.SalesRevenue, *created.SalesRevenue)
	})

	t.Run("validation error - invalid dealer ID", func(t *testing.T) {
		perfData := testutil.CreateTestPerformanceSales(0) // Неверный ID дилера

		_, err := service.CreatePerformance(ctx, perfData)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "dealer_id is required")
	})

	t.Run("validation error - invalid period", func(t *testing.T) {
		perfData := testutil.CreateTestPerformanceSales(1)
		perfData.Period = time.Time{} // Пустой период

		_, err := service.CreatePerformance(ctx, perfData)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "period is required")
	})

	t.Run("validation error - negative revenue", func(t *testing.T) {
		perfData := testutil.CreateTestPerformanceSales(1)
		negativeRevenue := -1000.0
		perfData.SalesRevenue = &negativeRevenue

		_, err := service.CreatePerformance(ctx, perfData)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "sales_revenue cannot be negative")
	})
}

func TestPerformanceService_GetPerformanceByPeriod(t *testing.T) {
	// Настройка тестовой базы данных
	testDB := testutil.SetupTestDB(t)
	defer testDB.Cleanup(t)
	testDB.RunMigrations(t)

	logger := testutil.GetTestLogger()
	repo := repository.NewPerformanceRepository(testDB.Pool)
	service := performance.NewService(repo, logger)

	ctx := context.Background()

	t.Run("successful retrieval", func(t *testing.T) {
		defer testDB.CleanupTable(t, "performance_sales")

		// Создаем тестового дилера
		dealerRepo := repository.NewDealerRepository(testDB.Pool)
		dealer := testutil.CreateTestDealer()
		dealerID, err := dealerRepo.Create(ctx, dealer)
		require.NoError(t, err)

		// Создаем тестовые данные
		perfData := testutil.CreateTestPerformanceSales(dealerID)
		_, err = service.CreatePerformance(ctx, perfData)
		require.NoError(t, err)

		// Получаем данные за период
		results, err := service.GetPerformanceByPeriod(ctx, "q1", 2024, "all-russia")
		require.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, dealerID, results[0].DealerID)
	})

	t.Run("invalid quarter", func(t *testing.T) {
		_, err := service.GetPerformanceByPeriod(ctx, "invalid", 2024, "all-russia")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid quarter")
	})

	t.Run("invalid year", func(t *testing.T) {
		_, err := service.GetPerformanceByPeriod(ctx, "q1", 1999, "all-russia")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid year")
	})
}

func TestPerformanceService_UpdatePerformance(t *testing.T) {
	// Настройка тестовой базы данных
	testDB := testutil.SetupTestDB(t)
	defer testDB.Cleanup(t)
	testDB.RunMigrations(t)

	logger := testutil.GetTestLogger()
	repo := repository.NewPerformanceRepository(testDB.Pool)
	service := performance.NewService(repo, logger)

	ctx := context.Background()

	t.Run("successful update", func(t *testing.T) {
		defer testDB.CleanupTable(t, "performance_sales")

		// Создаем тестового дилера
		dealerRepo := repository.NewDealerRepository(testDB.Pool)
		dealer := testutil.CreateTestDealer()
		dealerID, err := dealerRepo.Create(ctx, dealer)
		require.NoError(t, err)

		// Создаем тестовые данные
		perfData := testutil.CreateTestPerformanceSales(dealerID)
		id, err := service.CreatePerformance(ctx, perfData)
		require.NoError(t, err)

		// Обновляем данные
		newQuantitySold := 30
		updates := map[string]interface{}{
			"quantity_sold": newQuantitySold,
		}

		err = service.UpdatePerformance(ctx, id, updates)
		require.NoError(t, err)

		// Проверяем обновление
		updated, err := service.GetPerformanceByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, newQuantitySold, *updated.QuantitySold)
	})

	t.Run("update non-existent record", func(t *testing.T) {
		updates := map[string]interface{}{
			"quantity_sold": 30,
		}

		err := service.UpdatePerformance(ctx, 99999, updates)
		require.Error(t, err)
	})

	t.Run("empty updates", func(t *testing.T) {
		err := service.UpdatePerformance(ctx, 1, map[string]interface{}{})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "no fields to update")
	})
}

func TestPerformanceService_DeletePerformance(t *testing.T) {
	// Настройка тестовой базы данных
	testDB := testutil.SetupTestDB(t)
	defer testDB.Cleanup(t)
	testDB.RunMigrations(t)

	logger := testutil.GetTestLogger()
	repo := repository.NewPerformanceRepository(testDB.Pool)
	service := performance.NewService(repo, logger)

	ctx := context.Background()

	t.Run("successful deletion", func(t *testing.T) {
		defer testDB.CleanupTable(t, "performance_sales")

		// Создаем тестового дилера
		dealerRepo := repository.NewDealerRepository(testDB.Pool)
		dealer := testutil.CreateTestDealer()
		dealerID, err := dealerRepo.Create(ctx, dealer)
		require.NoError(t, err)

		// Создаем тестовые данные
		perfData := testutil.CreateTestPerformanceSales(dealerID)
		id, err := service.CreatePerformance(ctx, perfData)
		require.NoError(t, err)

		// Удаляем запись
		err = service.DeletePerformance(ctx, id)
		require.NoError(t, err)

		// Проверяем, что запись удалена
		_, err = service.GetPerformanceByID(ctx, id)
		require.Error(t, err)
	})

	t.Run("delete non-existent record", func(t *testing.T) {
		err := service.DeletePerformance(ctx, 99999)
		require.Error(t, err)
	})

	t.Run("invalid ID", func(t *testing.T) {
		err := service.DeletePerformance(ctx, 0)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid ID")
	})
}

func TestPerformanceService_GetPerformanceByID(t *testing.T) {
	// Настройка тестовой базы данных
	testDB := testutil.SetupTestDB(t)
	defer testDB.Cleanup(t)
	testDB.RunMigrations(t)

	logger := testutil.GetTestLogger()
	repo := repository.NewPerformanceRepository(testDB.Pool)
	service := performance.NewService(repo, logger)

	ctx := context.Background()

	t.Run("successful retrieval", func(t *testing.T) {
		defer testDB.CleanupTable(t, "performance_sales")

		// Создаем тестового дилера
		dealerRepo := repository.NewDealerRepository(testDB.Pool)
		dealer := testutil.CreateTestDealer()
		dealerID, err := dealerRepo.Create(ctx, dealer)
		require.NoError(t, err)

		// Создаем тестовые данные
		perfData := testutil.CreateTestPerformanceSales(dealerID)
		id, err := service.CreatePerformance(ctx, perfData)
		require.NoError(t, err)

		// Получаем данные по ID
		retrieved, err := service.GetPerformanceByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, id, retrieved.ID)
		assert.Equal(t, dealerID, retrieved.DealerID)
		assert.Equal(t, perfData.Period, retrieved.Period)
	})

	t.Run("non-existent ID", func(t *testing.T) {
		_, err := service.GetPerformanceByID(ctx, 99999)
		require.Error(t, err)
	})

	t.Run("invalid ID", func(t *testing.T) {
		_, err := service.GetPerformanceByID(ctx, 0)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid ID")
	})
}

func TestPerformanceService_FindPerformances(t *testing.T) {
	// Настройка тестовой базы данных
	testDB := testutil.SetupTestDB(t)
	defer testDB.Cleanup(t)
	testDB.RunMigrations(t)

	logger := testutil.GetTestLogger()
	repo := repository.NewPerformanceRepository(testDB.Pool)
	service := performance.NewService(repo, logger)

	ctx := context.Background()

	t.Run("successful retrieval by region", func(t *testing.T) {
		defer testDB.CleanupTable(t, "performance_sales")

		// Создаем тестового дилера
		dealerRepo := repository.NewDealerRepository(testDB.Pool)
		dealer := testutil.CreateTestDealer()
		dealerID, err := dealerRepo.Create(ctx, dealer)
		require.NoError(t, err)

		// Создаем тестовые данные
		perfData := testutil.CreateTestPerformanceSales(dealerID)
		_, err = service.CreatePerformance(ctx, perfData)
		require.NoError(t, err)

		// Получаем данные по региону (deprecated метод)
		results, err := service.FindPerformances(ctx, "all-russia")
		require.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, dealerID, results[0].DealerID)
	})
}
