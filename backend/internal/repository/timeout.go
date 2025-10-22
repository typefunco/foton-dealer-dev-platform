package repository

import (
	"context"
	"time"

	"github.com/typefunco/dealer_dev_platform/internal/database"
)

// WithDBTimeout создает контекст с таймаутом для операций с БД
func WithDBTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	// Используем глобальную конфигурацию БД
	if database.DBConfig != nil {
		return context.WithTimeout(ctx, database.DBConfig.QueryTimeout)
	}
	// Fallback на 30 секунд
	return context.WithTimeout(ctx, 30*time.Second)
}

// Example использования в репозитории:
/*
func (r *userRepository) GetUserByID(ctx context.Context, id int64) (*model.User, error) {
	// Создаем контекст с таймаутом
	queryCtx, cancel := WithDBTimeout(ctx)
	defer cancel()

	// Используем queryCtx для запроса
	query := r.sq.Select("id", "login", "password", "is_admin", "role", "region").
		From("users").
		Where(squirrel.Eq{"id": id})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("UserRepository.GetUserByID: error building query: %w", err)
	}

	var user model.User
	err = r.pool.QueryRow(queryCtx, sql, args...).Scan(
		&user.ID, &user.Login, &user.Password, &user.IsAdmin, &user.Role, &user.Region,
	)
	if err != nil {
		return nil, fmt.Errorf("UserRepository.GetUserByID: error scanning: %w", err)
	}

	return &user, nil
}
*/
