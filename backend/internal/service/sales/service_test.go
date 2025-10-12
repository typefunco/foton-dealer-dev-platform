package sales_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/typefunco/dealer_dev_platform/internal/repository"
	"github.com/typefunco/dealer_dev_platform/internal/service/sales"
	"github.com/typefunco/dealer_dev_platform/internal/testutil"
)

func TestSalesService_CreateSales(t *testing.T) {
	// Настройка тестовой базы данных
	testDB := testutil.SetupTestDB(t)
	defer testDB.Cleanup(t)
	testDB.RunMigrations(t)

	logger := testutil.GetTestLogger()
	repo := repository.NewSalesRepository(testDB.Pool)
	service := sales.NewService(repo, logger)

	ctx := context.Background()

	t.Run("successful creation", func(t *testing.T) {
		defer testDB.CleanupTable(t, "sales")

		// Создаем тестового дилера
		dealerRepo := repository.NewDealerRepository(testDB.Pool)
		dealer := testutil.CreateTestDealerWithRUFT("TEST005")
		dealerID, err := dealerRepo.Create(ctx, dealer)
		require.NoError(t, err)

		// Создаем тестовые данные продаж
		salesData := testutil.CreateTestSales(dealerID)

		// Создаем запись
		id, err := service.CreateSales(ctx, salesData)
		require.NoError(t, err)
		assert.Greater(t, id, int64(0))

		// Проверяем, что запись создана
		created, err := service.GetSalesByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, dealerID, created.DealerID)
		assert.Equal(t, salesData.Period, created.Period)
		assert.Equal(t, *salesData.StockHDT, *created.StockHDT)
		assert.Equal(t, *salesData.SalesTargetPlan, *created.SalesTargetPlan)
	})

	t.Run("validation error - invalid dealer ID", func(t *testing.T) {
		salesData := testutil.CreateTestSales(0) // Неверный ID дилера

		_, err := service.CreateSales(ctx, salesData)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "dealer_id is required")
	})

	t.Run("validation error - invalid period", func(t *testing.T) {
		salesData := testutil.CreateTestSales(1)
		salesData.Period = time.Time{} // Пустой период

		_, err := service.CreateSales(ctx, salesData)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "period is required")
	})

	t.Run("successful creation with negative values", func(t *testing.T) {
		defer testDB.CleanupTable(t, "sales")
		defer testDB.CleanupTable(t, "dealers")

		// Создаем тестового дилера
		dealerRepo := repository.NewDealerRepository(testDB.Pool)
		dealer := testutil.CreateTestDealerWithRUFT("TEST009")
		dealerID, err := dealerRepo.Create(ctx, dealer)
		require.NoError(t, err)

		salesData := testutil.CreateTestSales(dealerID)
		// Устанавливаем отрицательные значения
		negativeStock := -10
		negativeBuyout := -5
		negativePersonnel := -2
		salesData.StockHDT = &negativeStock
		salesData.BuyoutHDT = &negativeBuyout
		salesData.FotonSalesPersonnel = &negativePersonnel

		id, err := service.CreateSales(ctx, salesData)
		require.NoError(t, err)
		assert.Greater(t, id, int64(0))

		// Проверяем, что запись создана с отрицательными значениями
		created, err := service.GetSalesByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, -10, *created.StockHDT)
		assert.Equal(t, -5, *created.BuyoutHDT)
		assert.Equal(t, -2, *created.FotonSalesPersonnel)
	})
}

func TestSalesService_GetSalesByPeriod(t *testing.T) {
	// Настройка тестовой базы данных
	testDB := testutil.SetupTestDB(t)
	defer testDB.Cleanup(t)
	testDB.RunMigrations(t)

	logger := testutil.GetTestLogger()
	repo := repository.NewSalesRepository(testDB.Pool)
	service := sales.NewService(repo, logger)

	ctx := context.Background()

	t.Run("successful retrieval", func(t *testing.T) {
		defer testDB.CleanupTable(t, "sales")

		// Создаем тестового дилера
		dealerRepo := repository.NewDealerRepository(testDB.Pool)
		dealer := testutil.CreateTestDealerWithRUFT("TEST005")
		dealerID, err := dealerRepo.Create(ctx, dealer)
		require.NoError(t, err)

		// Создаем тестовые данные
		salesData := testutil.CreateTestSales(dealerID)
		_, err = service.CreateSales(ctx, salesData)
		require.NoError(t, err)

		// Получаем данные за период
		results, err := service.GetSalesByPeriod(ctx, "q1", 2024, "all-russia")
		require.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, dealerID, results[0].DealerID)
	})

	t.Run("invalid quarter", func(t *testing.T) {
		_, err := service.GetSalesByPeriod(ctx, "invalid", 2024, "all-russia")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid quarter")
	})

	t.Run("invalid year", func(t *testing.T) {
		_, err := service.GetSalesByPeriod(ctx, "q1", 1999, "all-russia")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid year")
	})
}

func TestSalesService_UpdateSales(t *testing.T) {
	// Настройка тестовой базы данных
	testDB := testutil.SetupTestDB(t)
	defer testDB.Cleanup(t)
	testDB.RunMigrations(t)

	logger := testutil.GetTestLogger()
	repo := repository.NewSalesRepository(testDB.Pool)
	service := sales.NewService(repo, logger)

	ctx := context.Background()

	t.Run("successful update", func(t *testing.T) {
		defer testDB.CleanupTable(t, "sales")

		// Создаем тестового дилера
		dealerRepo := repository.NewDealerRepository(testDB.Pool)
		dealer := testutil.CreateTestDealerWithRUFT("TEST005")
		dealerID, err := dealerRepo.Create(ctx, dealer)
		require.NoError(t, err)

		// Создаем тестовые данные
		salesData := testutil.CreateTestSales(dealerID)
		id, err := service.CreateSales(ctx, salesData)
		require.NoError(t, err)

		// Обновляем данные
		newStockHDT := 20
		updates := map[string]interface{}{
			"stock_hdt": newStockHDT,
		}

		err = service.UpdateSales(ctx, id, updates)
		require.NoError(t, err)

		// Проверяем обновление
		updated, err := service.GetSalesByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, newStockHDT, *updated.StockHDT)
	})

	t.Run("update non-existent record", func(t *testing.T) {
		updates := map[string]interface{}{
			"stock_hdt": 20,
		}

		err := service.UpdateSales(ctx, 99999, updates)
		require.Error(t, err)
	})

	t.Run("empty updates", func(t *testing.T) {
		err := service.UpdateSales(ctx, 1, map[string]interface{}{})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "no fields to update")
	})
}

func TestSalesService_DeleteSales(t *testing.T) {
	// Настройка тестовой базы данных
	testDB := testutil.SetupTestDB(t)
	defer testDB.Cleanup(t)
	testDB.RunMigrations(t)

	logger := testutil.GetTestLogger()
	repo := repository.NewSalesRepository(testDB.Pool)
	service := sales.NewService(repo, logger)

	ctx := context.Background()

	t.Run("successful deletion", func(t *testing.T) {
		defer testDB.CleanupTable(t, "sales")

		// Создаем тестового дилера
		dealerRepo := repository.NewDealerRepository(testDB.Pool)
		dealer := testutil.CreateTestDealerWithRUFT("TEST005")
		dealerID, err := dealerRepo.Create(ctx, dealer)
		require.NoError(t, err)

		// Создаем тестовые данные
		salesData := testutil.CreateTestSales(dealerID)
		id, err := service.CreateSales(ctx, salesData)
		require.NoError(t, err)

		// Удаляем запись
		err = service.DeleteSales(ctx, id)
		require.NoError(t, err)

		// Проверяем, что запись удалена
		_, err = service.GetSalesByID(ctx, id)
		require.Error(t, err)
	})

	t.Run("delete non-existent record", func(t *testing.T) {
		err := service.DeleteSales(ctx, 99999)
		require.Error(t, err)
	})

	t.Run("invalid ID", func(t *testing.T) {
		err := service.DeleteSales(ctx, 0)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid ID")
	})
}

func TestSalesService_GetSalesByID(t *testing.T) {
	// Настройка тестовой базы данных
	testDB := testutil.SetupTestDB(t)
	defer testDB.Cleanup(t)
	testDB.RunMigrations(t)

	logger := testutil.GetTestLogger()
	repo := repository.NewSalesRepository(testDB.Pool)
	service := sales.NewService(repo, logger)

	ctx := context.Background()

	t.Run("successful retrieval", func(t *testing.T) {
		defer testDB.CleanupTable(t, "sales")

		// Создаем тестового дилера
		dealerRepo := repository.NewDealerRepository(testDB.Pool)
		dealer := testutil.CreateTestDealerWithRUFT("TEST005")
		dealerID, err := dealerRepo.Create(ctx, dealer)
		require.NoError(t, err)

		// Создаем тестовые данные
		salesData := testutil.CreateTestSales(dealerID)
		id, err := service.CreateSales(ctx, salesData)
		require.NoError(t, err)

		// Получаем данные по ID
		retrieved, err := service.GetSalesByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, id, retrieved.ID)
		assert.Equal(t, dealerID, retrieved.DealerID)
		assert.Equal(t, salesData.Period, retrieved.Period)
	})

	t.Run("non-existent ID", func(t *testing.T) {
		_, err := service.GetSalesByID(ctx, 99999)
		require.Error(t, err)
	})

	t.Run("invalid ID", func(t *testing.T) {
		_, err := service.GetSalesByID(ctx, 0)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid ID")
	})
}
