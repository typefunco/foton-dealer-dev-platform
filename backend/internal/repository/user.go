package repository

import (
	"context"
	"fmt"
	"time"

	"log/slog"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/typefunco/dealer_dev_platform/internal/model"
)

// UserRepository интерфейс для работы с пользователями в базе данных.
type UserRepository interface {
	// GetUserByID возвращает пользователя по ID.
	GetUserByID(ctx context.Context, id int64) (*model.User, error)

	// GetUserByLogin возвращает пользователя по логину.
	GetUserByLogin(ctx context.Context, login string) (*model.User, error)

	// GetUsers возвращает список пользователей согласно фильтру.
	GetUsers(ctx context.Context, filter model.UserFilter) ([]*model.User, error)

	// CreateUser создает нового пользователя.
	CreateUser(ctx context.Context, user *model.User) (*model.User, error)

	// UpdateUser обновляет пользователя по ID.
	UpdateUser(ctx context.Context, id int64, update model.UserUpdate) (*model.User, error)

	// DeleteUser удаляет пользователя по ID.
	DeleteUser(ctx context.Context, id int64) error
}

// userRepository реализация UserRepository для работы с PostgreSQL.
type userRepository struct {
	pool   *pgxpool.Pool
	logger *slog.Logger
	sq     squirrel.StatementBuilderType
}

// NewUserRepository создает новый экземпляр репозитория пользователей.
func NewUserRepository(pool *pgxpool.Pool, logger *slog.Logger) UserRepository {
	return &userRepository{
		pool:   pool,
		logger: logger,
		sq:     squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

// GetUserByID возвращает пользователя по ID.
func (r *userRepository) GetUserByID(ctx context.Context, id int64) (*model.User, error) {
	query := r.sq.Select(
		"id", "login", "password", "is_admin", "role", "region",
		"first_name", "last_name", "email", "created_at", "updated_at",
	).From(usersTableName).
		Where(squirrel.Eq{"id": id})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("UserRepository.GetUserByID: failed to build query: %w", err)
	}

	var user model.User
	err = r.pool.QueryRow(ctx, sql, args...).Scan(
		&user.ID, &user.Login, &user.Password, &user.IsAdmin, &user.Role, &user.Region,
		&user.FirstName, &user.LastName, &user.Email, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("UserRepository.GetUserByID: user not found")
		}
		return nil, fmt.Errorf("UserRepository.GetUserByID: failed to scan user: %w", err)
	}

	return &user, nil
}

// GetUserByLogin возвращает пользователя по логину.
func (r *userRepository) GetUserByLogin(ctx context.Context, login string) (*model.User, error) {
	query := r.sq.Select(
		"id", "login", "password", "is_admin", "role", "region",
		"first_name", "last_name", "email", "created_at", "updated_at",
	).From(usersTableName).
		Where(squirrel.Eq{"login": login})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("UserRepository.GetUserByLogin: failed to build query: %w", err)
	}

	var user model.User
	err = r.pool.QueryRow(ctx, sql, args...).Scan(
		&user.ID, &user.Login, &user.Password, &user.IsAdmin, &user.Role, &user.Region,
		&user.FirstName, &user.LastName, &user.Email, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("UserRepository.GetUserByLogin: user not found")
		}
		return nil, fmt.Errorf("UserRepository.GetUserByLogin: failed to scan user: %w", err)
	}

	return &user, nil
}

// GetUsers возвращает список пользователей согласно фильтру.
func (r *userRepository) GetUsers(ctx context.Context, filter model.UserFilter) ([]*model.User, error) {
	query := r.sq.Select(
		"id", "login", "password", "is_admin", "role", "region",
		"first_name", "last_name", "email", "created_at", "updated_at",
	).From(usersTableName)

	// Применяем фильтры
	query = r.applyFilters(query, filter)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("UserRepository.GetUsers: failed to build query: %w", err)
	}

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("UserRepository.GetUsers: failed to execute query: %w", err)
	}
	defer rows.Close()

	var users []*model.User
	for rows.Next() {
		var user model.User
		err := rows.Scan(
			&user.ID, &user.Login, &user.Password, &user.IsAdmin, &user.Role, &user.Region,
			&user.FirstName, &user.LastName, &user.Email, &user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			r.logger.Error("UserRepository.GetUsers: failed to scan user", "error", err)
			continue
		}
		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("UserRepository.GetUsers: rows iteration error: %w", err)
	}

	return users, nil
}

// CreateUser создает нового пользователя.
func (r *userRepository) CreateUser(ctx context.Context, user *model.User) (*model.User, error) {
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	query := r.sq.Insert(usersTableName).
		Columns(
			"login", "password", "is_admin", "role", "region",
			"first_name", "last_name", "email", "created_at", "updated_at",
		).
		Values(
			user.Login, user.Password, user.IsAdmin, user.Role, user.Region,
			user.FirstName, user.LastName, user.Email, user.CreatedAt, user.UpdatedAt,
		).
		Suffix("RETURNING id, created_at, updated_at")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("UserRepository.CreateUser: failed to build query: %w", err)
	}

	err = r.pool.QueryRow(ctx, sql, args...).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("UserRepository.CreateUser: failed to create user: %w", err)
	}

	return user, nil
}

// UpdateUser обновляет пользователя по ID.
func (r *userRepository) UpdateUser(ctx context.Context, id int64, update model.UserUpdate) (*model.User, error) {
	query := r.sq.Update(usersTableName).
		Set("updated_at", time.Now())

	// Применяем только не-nil поля для обновления
	if update.Login != nil {
		query = query.Set("login", *update.Login)
	}
	if update.Password != nil {
		query = query.Set("password", *update.Password)
	}
	if update.IsAdmin != nil {
		query = query.Set("is_admin", *update.IsAdmin)
	}
	if update.Role != nil {
		query = query.Set("role", *update.Role)
	}
	if update.Region != nil {
		query = query.Set("region", *update.Region)
	}
	if update.FirstName != nil {
		query = query.Set("first_name", *update.FirstName)
	}
	if update.LastName != nil {
		query = query.Set("last_name", *update.LastName)
	}
	if update.Email != nil {
		query = query.Set("email", *update.Email)
	}

	query = query.Where(squirrel.Eq{"id": id}).
		Suffix("RETURNING id, login, password, is_admin, role, region, first_name, last_name, email, created_at, updated_at")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("UserRepository.UpdateUser: failed to build query: %w", err)
	}

	var user model.User
	err = r.pool.QueryRow(ctx, sql, args...).Scan(
		&user.ID, &user.Login, &user.Password, &user.IsAdmin, &user.Role, &user.Region,
		&user.FirstName, &user.LastName, &user.Email, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("UserRepository.UpdateUser: user not found")
		}
		return nil, fmt.Errorf("UserRepository.UpdateUser: failed to update user: %w", err)
	}

	return &user, nil
}

// DeleteUser удаляет пользователя по ID.
func (r *userRepository) DeleteUser(ctx context.Context, id int64) error {
	query := r.sq.Delete(usersTableName).Where(squirrel.Eq{"id": id})

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("UserRepository.DeleteUser: failed to build query: %w", err)
	}

	result, err := r.pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("UserRepository.DeleteUser: failed to delete user: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("UserRepository.DeleteUser: user not found")
	}

	return nil
}

// applyFilters применяет фильтры к запросу.
func (r *userRepository) applyFilters(query squirrel.SelectBuilder, filter model.UserFilter) squirrel.SelectBuilder {
	if filter.ID != nil {
		query = query.Where(squirrel.Eq{"id": *filter.ID})
	}
	if filter.Login != nil {
		query = query.Where(squirrel.Eq{"login": *filter.Login})
	}
	if filter.Email != nil {
		query = query.Where(squirrel.Eq{"email": *filter.Email})
	}
	if filter.IsAdmin != nil {
		query = query.Where(squirrel.Eq{"is_admin": *filter.IsAdmin})
	}
	if filter.Role != nil {
		query = query.Where(squirrel.Eq{"role": *filter.Role})
	}
	if filter.Region != nil {
		query = query.Where(squirrel.Eq{"region": *filter.Region})
	}
	if filter.FirstName != nil {
		query = query.Where(squirrel.Eq{"first_name": *filter.FirstName})
	}
	if filter.LastName != nil {
		query = query.Where(squirrel.Eq{"last_name": *filter.LastName})
	}

	return query
}
