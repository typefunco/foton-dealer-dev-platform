package performance_aftersales_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/typefunco/dealer_dev_platform/internal/repository"
	"github.com/typefunco/dealer_dev_platform/internal/service/performance_aftersales"
	"github.com/typefunco/dealer_dev_platform/internal/testutil"
)

func TestPerformanceAfterSalesService_Create(t *testing.T) {
	// Настройка тестовой базы данных
	testDB := testutil.SetupTestDB(t)
	defer testDB.Cleanup(t)
	testDB.RunMigrations(t)

	logger := testutil.GetTestLogger()
	repo := repository.NewPerformanceAfterSalesRepository(testDB.Pool, logger)
	service := performance_aftersales.NewService(repo, logger)

	ctx := context.Background()

	t.Run("successful creation", func(t *testing.T) {
		defer testDB.CleanupTable(t, "performance_aftersales")

		// Создаем тестового дилера
		dealerRepo := repository.NewDealerRepository(testDB.Pool)
		dealer := testutil.CreateTestDealer()
		dealerID, err := dealerRepo.Create(ctx, dealer)
		require.NoError(t, err)

		// Создаем тестовые данные производительности послепродажного обслуживания
		perfData := testutil.CreateTestPerformanceAfterSales(dealerID)

		// Создаем запись
		id, err := service.Create(ctx, perfData)
		require.NoError(t, err)
		assert.Greater(t, id, 0)

		// Проверяем, что запись создана
		created, err := service.GetByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, dealerID, created.DealerID)
		assert.Equal(t, perfData.Period, created.Period)
		assert.Equal(t, *perfData.ASRevenue, *created.ASRevenue)
		assert.Equal(t, *perfData.ASMargin, *created.ASMargin)
	})

	t.Run("validation error - invalid dealer ID", func(t *testing.T) {
		perfData := testutil.CreateTestPerformanceAfterSales(0) // Неверный ID дилера

		_, err := service.Create(ctx, perfData)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid dealer ID")
	})

	t.Run("validation error - invalid period", func(t *testing.T) {
		perfData := testutil.CreateTestPerformanceAfterSales(1)
		perfData.Period = time.Time{} // Пустой период

		_, err := service.Create(ctx, perfData)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid period")
	})

	t.Run("successful creation with negative values", func(t *testing.T) {
		defer testDB.CleanupTable(t, "performance_aftersales")
		defer testDB.CleanupTable(t, "dealers")

		dealer := testutil.CreateTestDealerWithRUFT("TEST002")
		dealerID, err := repository.NewDealerRepository(testDB.Pool).Create(ctx, dealer)
		require.NoError(t, err)

		perfData := testutil.CreateTestPerformanceAfterSales(int(dealerID))
		// Устанавливаем отрицательные значения
		negativeRevenue := -1000.0
		negativeMargin := -500.0
		perfData.ASRevenue = &negativeRevenue
		perfData.ASMargin = &negativeMargin

		id, err := service.Create(ctx, perfData)
		require.NoError(t, err)
		assert.Greater(t, id, 0)

		// Проверяем, что запись создана с отрицательными значениями
		created, err := service.GetByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, -1000.0, *created.ASRevenue)
		assert.Equal(t, -500.0, *created.ASMargin)
	})
}

func TestPerformanceAfterSalesService_GetByDealerIDAndPeriod(t *testing.T) {
	// Настройка тестовой базы данных
	testDB := testutil.SetupTestDB(t)
	defer testDB.Cleanup(t)
	testDB.RunMigrations(t)

	logger := testutil.GetTestLogger()
	repo := repository.NewPerformanceAfterSalesRepository(testDB.Pool, logger)
	service := performance_aftersales.NewService(repo, logger)

	ctx := context.Background()

	t.Run("successful retrieval", func(t *testing.T) {
		defer testDB.CleanupTable(t, "performance_aftersales")

		// Создаем тестового дилера
		dealerRepo := repository.NewDealerRepository(testDB.Pool)
		dealer := testutil.CreateTestDealer()
		dealerID, err := dealerRepo.Create(ctx, dealer)
		require.NoError(t, err)

		// Создаем тестовые данные
		perfData := testutil.CreateTestPerformanceAfterSales(dealerID)
		_, err = service.Create(ctx, perfData)
		require.NoError(t, err)

		// Получаем данные по дилеру и периоду
		retrieved, err := service.GetByDealerIDAndPeriod(ctx, dealerID, perfData.Period)
		require.NoError(t, err)
		assert.Equal(t, dealerID, retrieved.DealerID)
		assert.Equal(t, perfData.Period, retrieved.Period)
	})

	t.Run("invalid dealer ID", func(t *testing.T) {
		_, err := service.GetByDealerIDAndPeriod(ctx, 0, time.Now())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid dealer ID")
	})

	t.Run("invalid period", func(t *testing.T) {
		_, err := service.GetByDealerIDAndPeriod(ctx, 1, time.Time{})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid period")
	})

	t.Run("non-existent record", func(t *testing.T) {
		_, err := service.GetByDealerIDAndPeriod(ctx, 99999, time.Now())
		require.Error(t, err)
	})
}

func TestPerformanceAfterSalesService_GetAllByPeriod(t *testing.T) {
	// Настройка тестовой базы данных
	testDB := testutil.SetupTestDB(t)
	defer testDB.Cleanup(t)
	testDB.RunMigrations(t)

	logger := testutil.GetTestLogger()
	repo := repository.NewPerformanceAfterSalesRepository(testDB.Pool, logger)
	service := performance_aftersales.NewService(repo, logger)

	ctx := context.Background()

	t.Run("successful retrieval", func(t *testing.T) {
		defer testDB.CleanupTable(t, "performance_aftersales")

		// Создаем тестового дилера
		dealerRepo := repository.NewDealerRepository(testDB.Pool)
		dealer := testutil.CreateTestDealer()
		dealerID, err := dealerRepo.Create(ctx, dealer)
		require.NoError(t, err)

		// Создаем тестовые данные
		perfData := testutil.CreateTestPerformanceAfterSales(dealerID)
		_, err = service.Create(ctx, perfData)
		require.NoError(t, err)

		// Получаем все данные за период
		results, err := service.GetAllByPeriod(ctx, perfData.Period)
		require.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, dealerID, results[0].DealerID)
	})

	t.Run("invalid period", func(t *testing.T) {
		_, err := service.GetAllByPeriod(ctx, time.Time{})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid period")
	})

	t.Run("empty result", func(t *testing.T) {
		defer testDB.CleanupTable(t, "performance_aftersales")

		results, err := service.GetAllByPeriod(ctx, time.Now())
		require.NoError(t, err)
		assert.Len(t, results, 0)
	})
}

func TestPerformanceAfterSalesService_Update(t *testing.T) {
	// Настройка тестовой базы данных
	testDB := testutil.SetupTestDB(t)
	defer testDB.Cleanup(t)
	testDB.RunMigrations(t)

	logger := testutil.GetTestLogger()
	repo := repository.NewPerformanceAfterSalesRepository(testDB.Pool, logger)
	service := performance_aftersales.NewService(repo, logger)

	ctx := context.Background()

	t.Run("successful update", func(t *testing.T) {
		defer testDB.CleanupTable(t, "performance_aftersales")

		// Создаем тестового дилера
		dealerRepo := repository.NewDealerRepository(testDB.Pool)
		dealer := testutil.CreateTestDealer()
		dealerID, err := dealerRepo.Create(ctx, dealer)
		require.NoError(t, err)

		// Создаем тестовые данные
		perfData := testutil.CreateTestPerformanceAfterSales(dealerID)
		id, err := service.Create(ctx, perfData)
		require.NoError(t, err)

		// Обновляем данные
		newRevenue := 900000.0
		perfData.ID = id
		perfData.ASRevenue = &newRevenue

		err = service.Update(ctx, perfData)
		require.NoError(t, err)

		// Проверяем обновление
		updated, err := service.GetByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, newRevenue, *updated.ASRevenue)
	})

	t.Run("update non-existent record", func(t *testing.T) {
		perfData := testutil.CreateTestPerformanceAfterSales(1)
		perfData.ID = 99999

		err := service.Update(ctx, perfData)
		require.Error(t, err)
	})

	t.Run("invalid ID", func(t *testing.T) {
		perfData := testutil.CreateTestPerformanceAfterSales(1)
		perfData.ID = 0

		err := service.Update(ctx, perfData)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid ID")
	})

	t.Run("invalid dealer ID", func(t *testing.T) {
		perfData := testutil.CreateTestPerformanceAfterSales(0)
		perfData.ID = 1

		err := service.Update(ctx, perfData)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid dealer ID")
	})
}

func TestPerformanceAfterSalesService_Delete(t *testing.T) {
	// Настройка тестовой базы данных
	testDB := testutil.SetupTestDB(t)
	defer testDB.Cleanup(t)
	testDB.RunMigrations(t)

	logger := testutil.GetTestLogger()
	repo := repository.NewPerformanceAfterSalesRepository(testDB.Pool, logger)
	service := performance_aftersales.NewService(repo, logger)

	ctx := context.Background()

	t.Run("successful deletion", func(t *testing.T) {
		defer testDB.CleanupTable(t, "performance_aftersales")

		// Создаем тестового дилера
		dealerRepo := repository.NewDealerRepository(testDB.Pool)
		dealer := testutil.CreateTestDealer()
		dealerID, err := dealerRepo.Create(ctx, dealer)
		require.NoError(t, err)

		// Создаем тестовые данные
		perfData := testutil.CreateTestPerformanceAfterSales(dealerID)
		id, err := service.Create(ctx, perfData)
		require.NoError(t, err)

		// Удаляем запись
		err = service.Delete(ctx, id)
		require.NoError(t, err)

		// Проверяем, что запись удалена
		_, err = service.GetByID(ctx, id)
		require.Error(t, err)
	})

	t.Run("delete non-existent record", func(t *testing.T) {
		err := service.Delete(ctx, 99999)
		require.Error(t, err)
	})

	t.Run("invalid ID", func(t *testing.T) {
		err := service.Delete(ctx, 0)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid ID")
	})
}

func TestPerformanceAfterSalesService_GetByID(t *testing.T) {
	// Настройка тестовой базы данных
	testDB := testutil.SetupTestDB(t)
	defer testDB.Cleanup(t)
	testDB.RunMigrations(t)

	logger := testutil.GetTestLogger()
	repo := repository.NewPerformanceAfterSalesRepository(testDB.Pool, logger)
	service := performance_aftersales.NewService(repo, logger)

	ctx := context.Background()

	t.Run("successful retrieval", func(t *testing.T) {
		defer testDB.CleanupTable(t, "performance_aftersales")

		// Создаем тестового дилера
		dealerRepo := repository.NewDealerRepository(testDB.Pool)
		dealer := testutil.CreateTestDealer()
		dealerID, err := dealerRepo.Create(ctx, dealer)
		require.NoError(t, err)

		// Создаем тестовые данные
		perfData := testutil.CreateTestPerformanceAfterSales(dealerID)
		id, err := service.Create(ctx, perfData)
		require.NoError(t, err)

		// Получаем данные по ID
		retrieved, err := service.GetByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, id, retrieved.ID)
		assert.Equal(t, dealerID, retrieved.DealerID)
		assert.Equal(t, perfData.Period, retrieved.Period)
	})

	t.Run("non-existent ID", func(t *testing.T) {
		_, err := service.GetByID(ctx, 99999)
		require.Error(t, err)
	})

	t.Run("invalid ID", func(t *testing.T) {
		_, err := service.GetByID(ctx, 0)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid ID")
	})
}
