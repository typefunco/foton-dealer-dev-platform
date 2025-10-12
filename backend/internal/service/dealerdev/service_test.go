package dealerdev_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/typefunco/dealer_dev_platform/internal/model"
	"github.com/typefunco/dealer_dev_platform/internal/repository"
	"github.com/typefunco/dealer_dev_platform/internal/service/dealerdev"
	"github.com/typefunco/dealer_dev_platform/internal/testutil"
)

func TestDealerDevService_CreateDealerDev(t *testing.T) {
	// Настройка тестовой базы данных
	testDB := testutil.SetupTestDB(t)
	defer testDB.Cleanup(t)
	testDB.RunMigrations(t)

	logger := testutil.GetTestLogger()
	repo := repository.NewDealerDevRepository(testDB.Pool)
	service := dealerdev.NewService(repo, logger)

	ctx := context.Background()

	t.Run("successful creation", func(t *testing.T) {
		defer testDB.CleanupTable(t, "dealer_development")

		// Создаем тестового дилера
		dealerRepo := repository.NewDealerRepository(testDB.Pool)
		dealer := testutil.CreateTestDealerWithRUFT("TEST006")
		dealerID, err := dealerRepo.Create(ctx, dealer)
		require.NoError(t, err)

		// Создаем тестовые данные развития дилера
		devData := testutil.CreateTestDealerDevelopment(dealerID)

		// Создаем запись
		id, err := service.CreateDealerDev(ctx, devData)
		require.NoError(t, err)
		assert.Greater(t, id, 0)

		// Проверяем, что запись создана
		created, err := service.GetDealerDevByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, dealerID, created.DealerID)
		assert.Equal(t, devData.Period, created.Period)
		assert.Equal(t, *devData.CheckListScore, *created.CheckListScore)
		assert.Equal(t, *devData.DealershipClass, *created.DealershipClass)
	})

	t.Run("validation error - invalid dealer ID", func(t *testing.T) {
		devData := testutil.CreateTestDealerDevelopment(0) // Неверный ID дилера

		_, err := service.CreateDealerDev(ctx, devData)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "dealer_id is required")
	})

	t.Run("validation error - invalid period", func(t *testing.T) {
		devData := testutil.CreateTestDealerDevelopment(1)
		devData.Period = time.Time{} // Пустой период

		_, err := service.CreateDealerDev(ctx, devData)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "period is required")
	})

	t.Run("validation error - invalid check list score", func(t *testing.T) {
		devData := testutil.CreateTestDealerDevelopment(1)
		invalidScore := 150.0 // Больше 100
		devData.CheckListScore = &invalidScore

		_, err := service.CreateDealerDev(ctx, devData)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "check_list_score must be between 0 and 100")
	})

	t.Run("validation error - invalid dealership class", func(t *testing.T) {
		devData := testutil.CreateTestDealerDevelopment(1)
		invalidClass := model.DealershipClass("Z") // Неверный класс
		devData.DealershipClass = &invalidClass

		_, err := service.CreateDealerDev(ctx, devData)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid dealership_class")
	})
}

func TestDealerDevService_GetDealerDevByPeriod(t *testing.T) {
	// Настройка тестовой базы данных
	testDB := testutil.SetupTestDB(t)
	defer testDB.Cleanup(t)
	testDB.RunMigrations(t)

	logger := testutil.GetTestLogger()
	repo := repository.NewDealerDevRepository(testDB.Pool)
	service := dealerdev.NewService(repo, logger)

	ctx := context.Background()

	t.Run("successful retrieval", func(t *testing.T) {
		defer testDB.CleanupTable(t, "dealer_development")

		// Создаем тестового дилера
		dealerRepo := repository.NewDealerRepository(testDB.Pool)
		dealer := testutil.CreateTestDealerWithRUFT("TEST006")
		dealerID, err := dealerRepo.Create(ctx, dealer)
		require.NoError(t, err)

		// Создаем тестовые данные
		devData := testutil.CreateTestDealerDevelopment(dealerID)
		_, err = service.CreateDealerDev(ctx, devData)
		require.NoError(t, err)

		// Получаем данные за период
		results, err := service.GetDealerDevByPeriod(ctx, "q1", 2024, "all-russia")
		require.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, dealerID, results[0].DealerID)
	})

	t.Run("invalid quarter", func(t *testing.T) {
		_, err := service.GetDealerDevByPeriod(ctx, "invalid", 2024, "all-russia")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid quarter")
	})

	t.Run("invalid year", func(t *testing.T) {
		_, err := service.GetDealerDevByPeriod(ctx, "q1", 1999, "all-russia")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid year")
	})
}

func TestDealerDevService_UpdateDealerDev(t *testing.T) {
	// Настройка тестовой базы данных
	testDB := testutil.SetupTestDB(t)
	defer testDB.Cleanup(t)
	testDB.RunMigrations(t)

	logger := testutil.GetTestLogger()
	repo := repository.NewDealerDevRepository(testDB.Pool)
	service := dealerdev.NewService(repo, logger)

	ctx := context.Background()

	t.Run("successful update", func(t *testing.T) {
		defer testDB.CleanupTable(t, "dealer_development")

		// Создаем тестового дилера
		dealerRepo := repository.NewDealerRepository(testDB.Pool)
		dealer := testutil.CreateTestDealerWithRUFT("TEST006")
		dealerID, err := dealerRepo.Create(ctx, dealer)
		require.NoError(t, err)

		// Создаем тестовые данные
		devData := testutil.CreateTestDealerDevelopment(dealerID)
		id, err := service.CreateDealerDev(ctx, devData)
		require.NoError(t, err)

		// Обновляем данные
		newScore := 90.0
		updates := map[string]interface{}{
			"check_list_score": newScore,
		}

		err = service.UpdateDealerDev(ctx, id, updates)
		require.NoError(t, err)

		// Проверяем обновление
		updated, err := service.GetDealerDevByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, newScore, *updated.CheckListScore)
	})

	t.Run("update non-existent record", func(t *testing.T) {
		updates := map[string]interface{}{
			"check_list_score": 90.0,
		}

		err := service.UpdateDealerDev(ctx, 99999, updates)
		require.Error(t, err)
	})

	t.Run("empty updates", func(t *testing.T) {
		err := service.UpdateDealerDev(ctx, 1, map[string]interface{}{})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "no fields to update")
	})
}

func TestDealerDevService_DeleteDealerDev(t *testing.T) {
	// Настройка тестовой базы данных
	testDB := testutil.SetupTestDB(t)
	defer testDB.Cleanup(t)
	testDB.RunMigrations(t)

	logger := testutil.GetTestLogger()
	repo := repository.NewDealerDevRepository(testDB.Pool)
	service := dealerdev.NewService(repo, logger)

	ctx := context.Background()

	t.Run("successful deletion", func(t *testing.T) {
		defer testDB.CleanupTable(t, "dealer_development")

		// Создаем тестового дилера
		dealerRepo := repository.NewDealerRepository(testDB.Pool)
		dealer := testutil.CreateTestDealerWithRUFT("TEST006")
		dealerID, err := dealerRepo.Create(ctx, dealer)
		require.NoError(t, err)

		// Создаем тестовые данные
		devData := testutil.CreateTestDealerDevelopment(dealerID)
		id, err := service.CreateDealerDev(ctx, devData)
		require.NoError(t, err)

		// Удаляем запись
		err = service.DeleteDealerDev(ctx, id)
		require.NoError(t, err)

		// Проверяем, что запись удалена
		_, err = service.GetDealerDevByID(ctx, id)
		require.Error(t, err)
	})

	t.Run("delete non-existent record", func(t *testing.T) {
		err := service.DeleteDealerDev(ctx, 99999)
		require.Error(t, err)
	})

	t.Run("invalid ID", func(t *testing.T) {
		err := service.DeleteDealerDev(ctx, 0)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid ID")
	})
}

func TestDealerDevService_GetDealerDevByID(t *testing.T) {
	// Настройка тестовой базы данных
	testDB := testutil.SetupTestDB(t)
	defer testDB.Cleanup(t)
	testDB.RunMigrations(t)

	logger := testutil.GetTestLogger()
	repo := repository.NewDealerDevRepository(testDB.Pool)
	service := dealerdev.NewService(repo, logger)

	ctx := context.Background()

	t.Run("successful retrieval", func(t *testing.T) {
		defer testDB.CleanupTable(t, "dealer_development")

		// Создаем тестового дилера
		dealerRepo := repository.NewDealerRepository(testDB.Pool)
		dealer := testutil.CreateTestDealerWithRUFT("TEST006")
		dealerID, err := dealerRepo.Create(ctx, dealer)
		require.NoError(t, err)

		// Создаем тестовые данные
		devData := testutil.CreateTestDealerDevelopment(dealerID)
		id, err := service.CreateDealerDev(ctx, devData)
		require.NoError(t, err)

		// Получаем данные по ID
		retrieved, err := service.GetDealerDevByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, id, retrieved.ID)
		assert.Equal(t, dealerID, retrieved.DealerID)
		assert.Equal(t, devData.Period, retrieved.Period)
	})

	t.Run("non-existent ID", func(t *testing.T) {
		_, err := service.GetDealerDevByID(ctx, 99999)
		require.Error(t, err)
	})

	t.Run("invalid ID", func(t *testing.T) {
		_, err := service.GetDealerDevByID(ctx, 0)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid ID")
	})
}
