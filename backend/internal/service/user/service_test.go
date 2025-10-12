package user_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/typefunco/dealer_dev_platform/internal/model"
	"github.com/typefunco/dealer_dev_platform/internal/repository"
	"github.com/typefunco/dealer_dev_platform/internal/service/user"
	"github.com/typefunco/dealer_dev_platform/internal/testutil"
)

// stringPtr возвращает указатель на строку
func stringPtr(s string) *string {
	return &s
}

func TestUserService_CreateUser(t *testing.T) {
	// Настройка тестовой базы данных
	testDB := testutil.SetupTestDB(t)
	defer testDB.Cleanup(t)
	testDB.RunMigrations(t)

	logger := testutil.GetTestLogger()
	repo := repository.NewUserRepository(testDB.Pool, logger)
	service := user.NewService(repo, logger)

	ctx := context.Background()

	t.Run("successful creation", func(t *testing.T) {
		defer testDB.CleanupTable(t, "users")

		userData := testutil.CreateTestUser()
		createReq := model.UserCreateRequest{
			Login:     userData.Login,
			Password:  "password123",
			IsAdmin:   userData.IsAdmin,
			Role:      userData.Role,
			Region:    userData.Region,
			FirstName: userData.FirstName,
			LastName:  userData.LastName,
			Email:     userData.Email,
		}

		created, err := service.CreateUser(ctx, createReq)
		require.NoError(t, err)
		assert.Greater(t, created.ID, int64(0))

		// Проверяем, что пользователь создан
		retrieved, err := service.GetUserByID(ctx, created.ID)
		require.NoError(t, err)
		assert.Equal(t, userData.Login, retrieved.Login)
		assert.Equal(t, userData.Email, retrieved.Email)
		assert.Equal(t, userData.FirstName, retrieved.FirstName)
		assert.Equal(t, userData.LastName, retrieved.LastName)
		assert.Equal(t, userData.Role, retrieved.Role)
		assert.Equal(t, userData.IsAdmin, retrieved.IsAdmin)
	})

	t.Run("validation error - missing login", func(t *testing.T) {
		createReq := model.UserCreateRequest{
			Login:    "",
			Password: "password123",
			Role:     model.UserRoleAdmin,
		}

		_, err := service.CreateUser(ctx, createReq)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "login is required")
	})

	t.Run("validation error - missing password", func(t *testing.T) {
		createReq := model.UserCreateRequest{
			Login: "testuser",
			Role:  model.UserRoleAdmin,
		}

		_, err := service.CreateUser(ctx, createReq)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "password is required")
	})

	t.Run("validation error - invalid email format", func(t *testing.T) {
		createReq := model.UserCreateRequest{
			Login:    "testuser",
			Password: "password123",
			Email:    "invalid-email",
			Role:     model.UserRoleAdmin,
		}

		_, err := service.CreateUser(ctx, createReq)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid email format")
	})
}

func TestUserService_GetUsers(t *testing.T) {
	// Настройка тестовой базы данных
	testDB := testutil.SetupTestDB(t)
	defer testDB.Cleanup(t)
	testDB.RunMigrations(t)

	logger := testutil.GetTestLogger()
	repo := repository.NewUserRepository(testDB.Pool, logger)
	service := user.NewService(repo, logger)

	ctx := context.Background()

	t.Run("successful retrieval", func(t *testing.T) {
		defer testDB.CleanupTable(t, "users")

		// Создаем нескольких пользователей
		user1Req := model.UserCreateRequest{
			Login:     "user1",
			Password:  "password123",
			Email:     "user1@example.com",
			Role:      model.UserRoleAdmin,
			FirstName: "User",
			LastName:  "One",
		}

		user2Req := model.UserCreateRequest{
			Login:     "user2",
			Password:  "password123",
			Email:     "user2@example.com",
			Role:      model.UserRoleManager,
			FirstName: "User",
			LastName:  "Two",
		}

		_, err := service.CreateUser(ctx, user1Req)
		require.NoError(t, err)

		_, err = service.CreateUser(ctx, user2Req)
		require.NoError(t, err)

		// Получаем всех пользователей
		filter := model.UserFilter{}
		users, err := service.GetUsers(ctx, filter)
		require.NoError(t, err)
		assert.Len(t, users, 2)
	})

	t.Run("empty result", func(t *testing.T) {
		defer testDB.CleanupTable(t, "users")

		filter := model.UserFilter{}
		users, err := service.GetUsers(ctx, filter)
		require.NoError(t, err)
		assert.Len(t, users, 0)
	})
}

func TestUserService_GetUserByLogin(t *testing.T) {
	// Настройка тестовой базы данных
	testDB := testutil.SetupTestDB(t)
	defer testDB.Cleanup(t)
	testDB.RunMigrations(t)

	logger := testutil.GetTestLogger()
	repo := repository.NewUserRepository(testDB.Pool, logger)
	service := user.NewService(repo, logger)

	ctx := context.Background()

	t.Run("successful retrieval", func(t *testing.T) {
		defer testDB.CleanupTable(t, "users")

		// Создаем пользователя
		createReq := model.UserCreateRequest{
			Login:     "testuser",
			Password:  "password123",
			Email:     "test@example.com",
			Role:      model.UserRoleAdmin,
			FirstName: "Test",
			LastName:  "User",
		}
		_, err := service.CreateUser(ctx, createReq)
		require.NoError(t, err)

		// Получаем пользователя по login
		retrieved, err := service.GetUserByLogin(ctx, createReq.Login)
		require.NoError(t, err)
		assert.Equal(t, createReq.Login, retrieved.Login)
		assert.Equal(t, createReq.Email, retrieved.Email)
	})

	t.Run("non-existent login", func(t *testing.T) {
		_, err := service.GetUserByLogin(ctx, "nonexistent")
		require.Error(t, err)
	})

	t.Run("empty login", func(t *testing.T) {
		_, err := service.GetUserByLogin(ctx, "")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "login cannot be empty")
	})
}

func TestUserService_UpdateUser(t *testing.T) {
	// Настройка тестовой базы данных
	testDB := testutil.SetupTestDB(t)
	defer testDB.Cleanup(t)
	testDB.RunMigrations(t)

	logger := testutil.GetTestLogger()
	repo := repository.NewUserRepository(testDB.Pool, logger)
	service := user.NewService(repo, logger)

	ctx := context.Background()

	t.Run("successful update", func(t *testing.T) {
		defer testDB.CleanupTable(t, "users")

		// Создаем пользователя
		createReq := model.UserCreateRequest{
			Login:     "testuser",
			Password:  "password123",
			Email:     "test@example.com",
			Role:      model.UserRoleAdmin,
			FirstName: "Test",
			LastName:  "User",
		}
		created, err := service.CreateUser(ctx, createReq)
		require.NoError(t, err)

		// Обновляем данные
		update := model.UserUpdate{
			FirstName: stringPtr("Новое Имя"),
			LastName:  stringPtr("Новая Фамилия"),
		}

		updated, err := service.UpdateUser(ctx, created.ID, update)
		require.NoError(t, err)
		assert.Equal(t, "Новое Имя", updated.FirstName)
		assert.Equal(t, "Новая Фамилия", updated.LastName)
	})

	t.Run("update non-existent user", func(t *testing.T) {
		update := model.UserUpdate{
			FirstName: stringPtr("Новое Имя"),
		}

		_, err := service.UpdateUser(ctx, 99999, update)
		require.Error(t, err)
	})
}

func TestUserService_DeleteUser(t *testing.T) {
	// Настройка тестовой базы данных
	testDB := testutil.SetupTestDB(t)
	defer testDB.Cleanup(t)
	testDB.RunMigrations(t)

	logger := testutil.GetTestLogger()
	repo := repository.NewUserRepository(testDB.Pool, logger)
	service := user.NewService(repo, logger)

	ctx := context.Background()

	t.Run("successful deletion", func(t *testing.T) {
		defer testDB.CleanupTable(t, "users")

		// Создаем пользователя
		createReq := model.UserCreateRequest{
			Login:     "testuser",
			Password:  "password123",
			Email:     "test@example.com",
			Role:      model.UserRoleAdmin,
			FirstName: "Test",
			LastName:  "User",
		}
		created, err := service.CreateUser(ctx, createReq)
		require.NoError(t, err)

		// Удаляем пользователя
		err = service.DeleteUser(ctx, created.ID)
		require.NoError(t, err)

		// Проверяем, что пользователь удален
		_, err = service.GetUserByID(ctx, created.ID)
		require.Error(t, err)
	})

	t.Run("delete non-existent user", func(t *testing.T) {
		err := service.DeleteUser(ctx, 99999)
		require.Error(t, err)
	})

	t.Run("invalid ID", func(t *testing.T) {
		err := service.DeleteUser(ctx, 0)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid user ID")
	})
}

func TestUserService_GetUserByID(t *testing.T) {
	// Настройка тестовой базы данных
	testDB := testutil.SetupTestDB(t)
	defer testDB.Cleanup(t)
	testDB.RunMigrations(t)

	logger := testutil.GetTestLogger()
	repo := repository.NewUserRepository(testDB.Pool, logger)
	service := user.NewService(repo, logger)

	ctx := context.Background()

	t.Run("successful retrieval", func(t *testing.T) {
		defer testDB.CleanupTable(t, "users")

		// Создаем пользователя
		createReq := model.UserCreateRequest{
			Login:     "testuser",
			Password:  "password123",
			Email:     "test@example.com",
			Role:      model.UserRoleAdmin,
			FirstName: "Test",
			LastName:  "User",
		}
		created, err := service.CreateUser(ctx, createReq)
		require.NoError(t, err)

		// Получаем пользователя по ID
		retrieved, err := service.GetUserByID(ctx, created.ID)
		require.NoError(t, err)
		assert.Equal(t, created.ID, retrieved.ID)
		assert.Equal(t, createReq.Login, retrieved.Login)
		assert.Equal(t, createReq.Email, retrieved.Email)
	})

	t.Run("non-existent ID", func(t *testing.T) {
		_, err := service.GetUserByID(ctx, 99999)
		require.Error(t, err)
	})

	t.Run("invalid ID", func(t *testing.T) {
		_, err := service.GetUserByID(ctx, 0)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid user ID")
	})
}
