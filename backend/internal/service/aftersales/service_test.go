package aftersales_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/typefunco/dealer_dev_platform/internal/repository"
	"github.com/typefunco/dealer_dev_platform/internal/service/aftersales"
	"github.com/typefunco/dealer_dev_platform/internal/testutil"
)

func TestAfterSalesService_CreateAfterSales(t *testing.T) {
	// Настройка тестовой базы данных
	testDB := testutil.SetupTestDB(t)
	defer testDB.Cleanup(t)
	testDB.RunMigrations(t)

	logger := testutil.GetTestLogger()
	repo := repository.NewAfterSalesRepository(testDB.Pool)
	service := aftersales.NewService(repo, logger)

	ctx := context.Background()

	t.Run("successful creation", func(t *testing.T) {
		defer testDB.CleanupTable(t, "aftersales")

		// Создаем тестового дилера
		dealerRepo := repository.NewDealerRepository(testDB.Pool)
		dealer := testutil.CreateTestDealerWithRUFT("TEST004")
		dealerID, err := dealerRepo.Create(ctx, dealer)
		require.NoError(t, err)

		// Создаем тестовые данные послепродажного обслуживания
		afterSales := testutil.CreateTestAfterSales(dealerID)

		// Создаем запись
		id, err := service.CreateAfterSales(ctx, afterSales)
		require.NoError(t, err)
		assert.Greater(t, id, int64(0))

		// Проверяем, что запись создана
		created, err := service.GetAfterSalesByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, dealerID, created.DealerID)
		assert.Equal(t, afterSales.Period, created.Period)
		assert.Equal(t, *afterSales.RecommendedStockPct, *created.RecommendedStockPct)
	})

	t.Run("validation error - invalid dealer ID", func(t *testing.T) {
		afterSales := testutil.CreateTestAfterSales(0) // Неверный ID дилера

		_, err := service.CreateAfterSales(ctx, afterSales)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "dealer_id is required")
	})

	t.Run("validation error - invalid period", func(t *testing.T) {
		afterSales := testutil.CreateTestAfterSales(1)
		afterSales.Period = time.Time{} // Пустой период

		_, err := service.CreateAfterSales(ctx, afterSales)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "period is required")
	})

	t.Run("successful creation with negative values", func(t *testing.T) {
		defer testDB.CleanupTable(t, "aftersales")
		defer testDB.CleanupTable(t, "dealers")

		// Создаем тестового дилера
		dealerRepo := repository.NewDealerRepository(testDB.Pool)
		dealer := testutil.CreateTestDealerWithRUFT("TEST008")
		dealerID, err := dealerRepo.Create(ctx, dealer)
		require.NoError(t, err)

		afterSales := testutil.CreateTestAfterSales(dealerID)
		// Устанавливаем отрицательные значения
		negativeHours := -50.0
		negativeSpareParts := -100.0
		afterSales.ServiceContractsHours = &negativeHours
		afterSales.SparePartsSalesQ = &negativeSpareParts

		id, err := service.CreateAfterSales(ctx, afterSales)
		require.NoError(t, err)
		assert.Greater(t, id, int64(0))

		// Проверяем, что запись создана с отрицательными значениями
		created, err := service.GetAfterSalesByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, -50.0, *created.ServiceContractsHours)
		assert.Equal(t, -100.0, *created.SparePartsSalesQ)
	})
}

func TestAfterSalesService_GetAfterSalesByPeriod(t *testing.T) {
	// Настройка тестовой базы данных
	testDB := testutil.SetupTestDB(t)
	defer testDB.Cleanup(t)
	testDB.RunMigrations(t)

	logger := testutil.GetTestLogger()
	repo := repository.NewAfterSalesRepository(testDB.Pool)
	service := aftersales.NewService(repo, logger)

	ctx := context.Background()

	t.Run("successful retrieval", func(t *testing.T) {
		defer testDB.CleanupTable(t, "aftersales")

		// Создаем тестового дилера
		dealerRepo := repository.NewDealerRepository(testDB.Pool)
		dealer := testutil.CreateTestDealerWithRUFT("TEST004")
		dealerID, err := dealerRepo.Create(ctx, dealer)
		require.NoError(t, err)

		// Создаем тестовые данные
		afterSales := testutil.CreateTestAfterSales(dealerID)
		_, err = service.CreateAfterSales(ctx, afterSales)
		require.NoError(t, err)

		// Получаем данные за период
		results, err := service.GetAfterSalesByPeriod(ctx, "q1", 2024, "all-russia")
		require.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, dealerID, results[0].DealerID)
	})

	t.Run("invalid quarter", func(t *testing.T) {
		_, err := service.GetAfterSalesByPeriod(ctx, "invalid", 2024, "all-russia")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid quarter")
	})

	t.Run("invalid year", func(t *testing.T) {
		_, err := service.GetAfterSalesByPeriod(ctx, "q1", 1999, "all-russia")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid year")
	})
}

func TestAfterSalesService_UpdateAfterSales(t *testing.T) {
	// Настройка тестовой базы данных
	testDB := testutil.SetupTestDB(t)
	defer testDB.Cleanup(t)
	testDB.RunMigrations(t)

	logger := testutil.GetTestLogger()
	repo := repository.NewAfterSalesRepository(testDB.Pool)
	service := aftersales.NewService(repo, logger)

	ctx := context.Background()

	t.Run("successful update", func(t *testing.T) {
		defer testDB.CleanupTable(t, "aftersales")

		// Создаем тестового дилера
		dealerRepo := repository.NewDealerRepository(testDB.Pool)
		dealer := testutil.CreateTestDealerWithRUFT("TEST004")
		dealerID, err := dealerRepo.Create(ctx, dealer)
		require.NoError(t, err)

		// Создаем тестовые данные
		afterSales := testutil.CreateTestAfterSales(dealerID)
		id, err := service.CreateAfterSales(ctx, afterSales)
		require.NoError(t, err)

		// Обновляем данные
		newStockPct := 90.0
		updates := map[string]interface{}{
			"recommended_stock_pct": newStockPct,
		}

		err = service.UpdateAfterSales(ctx, id, updates)
		require.NoError(t, err)

		// Проверяем обновление
		updated, err := service.GetAfterSalesByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, newStockPct, *updated.RecommendedStockPct)
	})

	t.Run("update non-existent record", func(t *testing.T) {
		updates := map[string]interface{}{
			"recommended_stock_pct": 90.0,
		}

		err := service.UpdateAfterSales(ctx, 99999, updates)
		require.Error(t, err)
	})

	t.Run("empty updates", func(t *testing.T) {
		err := service.UpdateAfterSales(ctx, 1, map[string]interface{}{})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "no fields to update")
	})
}

func TestAfterSalesService_DeleteAfterSales(t *testing.T) {
	// Настройка тестовой базы данных
	testDB := testutil.SetupTestDB(t)
	defer testDB.Cleanup(t)
	testDB.RunMigrations(t)

	logger := testutil.GetTestLogger()
	repo := repository.NewAfterSalesRepository(testDB.Pool)
	service := aftersales.NewService(repo, logger)

	ctx := context.Background()

	t.Run("successful deletion", func(t *testing.T) {
		defer testDB.CleanupTable(t, "aftersales")

		// Создаем тестового дилера
		dealerRepo := repository.NewDealerRepository(testDB.Pool)
		dealer := testutil.CreateTestDealerWithRUFT("TEST004")
		dealerID, err := dealerRepo.Create(ctx, dealer)
		require.NoError(t, err)

		// Создаем тестовые данные
		afterSales := testutil.CreateTestAfterSales(dealerID)
		id, err := service.CreateAfterSales(ctx, afterSales)
		require.NoError(t, err)

		// Удаляем запись
		err = service.DeleteAfterSales(ctx, id)
		require.NoError(t, err)

		// Проверяем, что запись удалена
		_, err = service.GetAfterSalesByID(ctx, id)
		require.Error(t, err)
	})

	t.Run("delete non-existent record", func(t *testing.T) {
		err := service.DeleteAfterSales(ctx, 99999)
		require.Error(t, err)
	})

	t.Run("invalid ID", func(t *testing.T) {
		err := service.DeleteAfterSales(ctx, 0)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid ID")
	})
}
