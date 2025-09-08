package repository

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/typefunco/dealer_dev_platform/internal/model"
	"log/slog"
)

type AuthRepository struct {
	pool   *pgxpool.Pool
	logger *slog.Logger
	sq     squirrel.StatementBuilderType
}

const (
	usersTableName = "users"
)

func NewAuthRepository(pool *pgxpool.Pool, logger *slog.Logger) *AuthRepository {
	return &AuthRepository{
		pool:   pool,
		logger: logger,
		sq:     squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (repo *AuthRepository) CreateUser(ctx context.Context, user model.User) error {
	query := repo.sq.Insert(usersTableName).
		Columns("login", "password", "is_admin", "role", "region", "created_at", "updated_at").
		Values(user.Login, user.Password, user.IsAdmin, user.Role, user.Region, user.CreatedAt, user.UpdatedAt)

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("AuthRepository.CreateUser error creating query: %w", err)
	}

	_, err = repo.pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("AuthRepository.CreateUser error exec tx: %w", err)
	}

	return nil
}

func (repo *AuthRepository) DeleteUser(ctx context.Context, login string) error {
	query := repo.sq.Delete(usersTableName).Where(squirrel.Eq{"login": login})

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("AuthRepository.DeleteUser error creating query: %w", err)
	}

	_, err = repo.pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("AuthRepository.DeleteUser error exec tx: %w", err)
	}

	return nil
}

func (repo *AuthRepository) GetUser(ctx context.Context, login string) (*model.User, error) {
	var user model.User
	query := repo.sq.Select(usersTableName).Where(squirrel.Eq{"login": login}).Columns("id", "login", "password")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("AuthRepository.GetUser error creating query: %w", err)
	}

	rows, err := repo.pool.Query(ctx, sql, args...)
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&user.ID, &user.Login, &user.Password)
		if err != nil {
			repo.logger.Error("AuthRepository.GetUser error parse sql")
		}
	}

	return &user, rows.Err()
}
