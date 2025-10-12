package dealer_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/typefunco/dealer_dev_platform/internal/repository"
	"github.com/typefunco/dealer_dev_platform/internal/service/dealer"
	"github.com/typefunco/dealer_dev_platform/internal/testutil"
)

func TestDealerService_CreateDealer(t *testing.T) {
	// Настройка тестовой базы данных
	testDB := testutil.SetupTestDB(t)
	defer testDB.Cleanup(t)
	testDB.RunMigrations(t)

	logger := testutil.GetTestLogger()
	repo := repository.NewDealerRepository(testDB.Pool)
	service := dealer.NewService(repo, logger)

	ctx := context.Background()

	t.Run("successful creation", func(t *testing.T) {
		defer testDB.CleanupTable(t, "dealers")

		dealerData := testutil.CreateTestDealerWithRUFT("TEST007")

		id, err := service.CreateDealer(ctx, dealerData)
		require.NoError(t, err)
		assert.Greater(t, id, 0)

		// Проверяем, что дилер создан
		created, err := service.GetDealerByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, dealerData.Ruft, created.Ruft)
		assert.Equal(t, dealerData.DealerNameRu, created.DealerNameRu)
		assert.Equal(t, dealerData.DealerNameEn, created.DealerNameEn)
		assert.Equal(t, dealerData.Region, created.Region)
		assert.Equal(t, dealerData.City, created.City)
		assert.Equal(t, dealerData.Manager, created.Manager)
	})

	t.Run("validation error - missing dealer name", func(t *testing.T) {
		dealerData := testutil.CreateTestDealerWithRUFT("TEST007")
		dealerData.DealerNameRu = ""

		_, err := service.CreateDealer(ctx, dealerData)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "dealer_name_ru is required")
	})

	t.Run("validation error - missing city", func(t *testing.T) {
		dealerData := testutil.CreateTestDealerWithRUFT("TEST007")
		dealerData.City = ""

		_, err := service.CreateDealer(ctx, dealerData)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "city is required")
	})
}

func TestDealerService_GetAllDealers(t *testing.T) {
	// Настройка тестовой базы данных
	testDB := testutil.SetupTestDB(t)
	defer testDB.Cleanup(t)
	testDB.RunMigrations(t)

	logger := testutil.GetTestLogger()
	repo := repository.NewDealerRepository(testDB.Pool)
	service := dealer.NewService(repo, logger)

	ctx := context.Background()

	t.Run("successful retrieval", func(t *testing.T) {
		defer testDB.CleanupTable(t, "dealers")

		// Создаем несколько дилеров
		dealer1 := testutil.CreateTestDealerWithRUFT("TEST007")
		dealer1.Ruft = "TEST001"
		dealer1.DealerNameRu = "Дилер 1"

		dealer2 := testutil.CreateTestDealerWithRUFT("TEST007")
		dealer2.Ruft = "TEST002"
		dealer2.DealerNameRu = "Дилер 2"

		_, err := service.CreateDealer(ctx, dealer1)
		require.NoError(t, err)

		_, err = service.CreateDealer(ctx, dealer2)
		require.NoError(t, err)

		// Получаем всех дилеров
		dealers, err := service.GetAllDealers(ctx)
		require.NoError(t, err)
		assert.Len(t, dealers, 2)
	})

	t.Run("empty result", func(t *testing.T) {
		defer testDB.CleanupTable(t, "dealers")

		dealers, err := service.GetAllDealers(ctx)
		require.NoError(t, err)
		assert.Len(t, dealers, 0)
	})
}

func TestDealerService_GetDealersByRegion(t *testing.T) {
	// Настройка тестовой базы данных
	testDB := testutil.SetupTestDB(t)
	defer testDB.Cleanup(t)
	testDB.RunMigrations(t)

	logger := testutil.GetTestLogger()
	repo := repository.NewDealerRepository(testDB.Pool)
	service := dealer.NewService(repo, logger)

	ctx := context.Background()

	t.Run("successful retrieval by region", func(t *testing.T) {
		defer testDB.CleanupTable(t, "dealers")

		// Создаем дилеров в разных регионах
		dealer1 := testutil.CreateTestDealerWithRUFT("TEST007")
		dealer1.Ruft = "TEST001"
		dealer1.Region = "Москва"

		dealer2 := testutil.CreateTestDealerWithRUFT("TEST007")
		dealer2.Ruft = "TEST002"
		dealer2.Region = "Санкт-Петербург"

		dealer3 := testutil.CreateTestDealerWithRUFT("TEST007")
		dealer3.Ruft = "TEST003"
		dealer3.Region = "Москва"

		_, err := service.CreateDealer(ctx, dealer1)
		require.NoError(t, err)

		_, err = service.CreateDealer(ctx, dealer2)
		require.NoError(t, err)

		_, err = service.CreateDealer(ctx, dealer3)
		require.NoError(t, err)

		// Получаем дилеров по региону
		dealers, err := service.GetDealersByRegion(ctx, "Москва")
		require.NoError(t, err)
		assert.Len(t, dealers, 2)

		for _, d := range dealers {
			assert.Equal(t, "Москва", d.Region)
		}
	})

	t.Run("no dealers in region", func(t *testing.T) {
		defer testDB.CleanupTable(t, "dealers")

		dealers, err := service.GetDealersByRegion(ctx, "Несуществующий регион")
		require.NoError(t, err)
		assert.Len(t, dealers, 0)
	})
}

func TestDealerService_UpdateDealer(t *testing.T) {
	// Настройка тестовой базы данных
	testDB := testutil.SetupTestDB(t)
	defer testDB.Cleanup(t)
	testDB.RunMigrations(t)

	logger := testutil.GetTestLogger()
	repo := repository.NewDealerRepository(testDB.Pool)
	service := dealer.NewService(repo, logger)

	ctx := context.Background()

	t.Run("successful update", func(t *testing.T) {
		defer testDB.CleanupTable(t, "dealers")

		// Создаем дилера
		dealerData := testutil.CreateTestDealerWithRUFT("TEST007")
		id, err := service.CreateDealer(ctx, dealerData)
		require.NoError(t, err)

		// Обновляем данные
		updates := map[string]interface{}{
			"city":    "Новый Город",
			"manager": "Новый Менеджер",
		}

		err = service.UpdateDealer(ctx, id, updates)
		require.NoError(t, err)

		// Проверяем обновление
		updated, err := service.GetDealerByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, "Новый Город", updated.City)
		assert.Equal(t, "Новый Менеджер", updated.Manager)
	})

	t.Run("update non-existent dealer", func(t *testing.T) {
		updates := map[string]interface{}{
			"city": "Новый Город",
		}

		err := service.UpdateDealer(ctx, 99999, updates)
		require.Error(t, err)
	})

	t.Run("empty updates", func(t *testing.T) {
		err := service.UpdateDealer(ctx, 1, map[string]interface{}{})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "no fields to update")
	})
}

func TestDealerService_DeleteDealer(t *testing.T) {
	// Настройка тестовой базы данных
	testDB := testutil.SetupTestDB(t)
	defer testDB.Cleanup(t)
	testDB.RunMigrations(t)

	logger := testutil.GetTestLogger()
	repo := repository.NewDealerRepository(testDB.Pool)
	service := dealer.NewService(repo, logger)

	ctx := context.Background()

	t.Run("successful deletion", func(t *testing.T) {
		defer testDB.CleanupTable(t, "dealers")

		// Создаем дилера
		dealerData := testutil.CreateTestDealerWithRUFT("TEST007")
		id, err := service.CreateDealer(ctx, dealerData)
		require.NoError(t, err)

		// Удаляем дилера
		err = service.DeleteDealer(ctx, id)
		require.NoError(t, err)

		// Проверяем, что дилер удален
		_, err = service.GetDealerByID(ctx, id)
		require.Error(t, err)
	})

	t.Run("delete non-existent dealer", func(t *testing.T) {
		err := service.DeleteDealer(ctx, 99999)
		require.Error(t, err)
	})

	t.Run("invalid ID", func(t *testing.T) {
		err := service.DeleteDealer(ctx, 0)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid ID")
	})
}

func TestDealerService_GetDealerCard(t *testing.T) {
	// Настройка тестовой базы данных
	testDB := testutil.SetupTestDB(t)
	defer testDB.Cleanup(t)
	testDB.RunMigrations(t)

	logger := testutil.GetTestLogger()
	repo := repository.NewDealerRepository(testDB.Pool)
	service := dealer.NewService(repo, logger)

	ctx := context.Background()

	t.Run("successful retrieval", func(t *testing.T) {
		defer testDB.CleanupTable(t, "dealers")

		// Создаем дилера
		dealerData := testutil.CreateTestDealerWithRUFT("TEST007")
		id, err := service.CreateDealer(ctx, dealerData)
		require.NoError(t, err)

		// Получаем карточку дилера
		card, err := service.GetDealerCard(ctx, int64(id), "q1", 2024)
		require.NoError(t, err)
		assert.Equal(t, id, card.DealerID)
		assert.Equal(t, dealerData.DealerNameRu, card.DealerNameRu)
		assert.Equal(t, dealerData.Region, card.Region)
		assert.Equal(t, dealerData.City, card.City)
	})

	t.Run("invalid quarter", func(t *testing.T) {
		_, err := service.GetDealerCard(ctx, 1, "invalid", 2024)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid quarter")
	})

	t.Run("invalid year", func(t *testing.T) {
		_, err := service.GetDealerCard(ctx, 1, "q1", 1999)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid year")
	})

	t.Run("invalid dealer ID", func(t *testing.T) {
		_, err := service.GetDealerCard(ctx, 0, "q1", 2024)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid dealer ID")
	})
}
