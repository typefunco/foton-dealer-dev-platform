package user

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/typefunco/dealer_dev_platform/internal/model"
)

// Repository интерфейс репозитория пользователей.
type Repository interface {
	GetUserByID(ctx context.Context, id int64) (*model.User, error)
	GetUserByLogin(ctx context.Context, login string) (*model.User, error)
	GetUsers(ctx context.Context, filter model.UserFilter) ([]*model.User, error)
	CreateUser(ctx context.Context, user *model.User) (*model.User, error)
	UpdateUser(ctx context.Context, id int64, update model.UserUpdate) (*model.User, error)
	DeleteUser(ctx context.Context, id int64) error
}

// Service сервис для работы с пользователями.
type Service struct {
	repo   Repository
	logger *slog.Logger
}

// NewService создает новый экземпляр сервиса пользователей.
func NewService(repo Repository, logger *slog.Logger) *Service {
	return &Service{
		repo:   repo,
		logger: logger,
	}
}

// GetUserByID возвращает пользователя по ID.
func (s *Service) GetUserByID(ctx context.Context, id int64) (*model.UserResponse, error) {
	if id <= 0 {
		return nil, fmt.Errorf("UserService.GetUserByID: invalid user ID")
	}

	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		s.logger.Error("UserService.GetUserByID: failed to get user", "id", id, "error", err)
		return nil, fmt.Errorf("UserService.GetUserByID: %w", err)
	}

	return s.toUserResponse(user), nil
}

// GetUserByLogin возвращает пользователя по логину.
func (s *Service) GetUserByLogin(ctx context.Context, login string) (*model.UserResponse, error) {
	if login == "" {
		return nil, fmt.Errorf("UserService.GetUserByLogin: login cannot be empty")
	}

	user, err := s.repo.GetUserByLogin(ctx, login)
	if err != nil {
		s.logger.Error("UserService.GetUserByLogin: failed to get user", "login", login, "error", err)
		return nil, fmt.Errorf("UserService.GetUserByLogin: %w", err)
	}

	return s.toUserResponse(user), nil
}

// GetUsers возвращает список пользователей согласно фильтру.
func (s *Service) GetUsers(ctx context.Context, filter model.UserFilter) ([]*model.UserResponse, error) {
	users, err := s.repo.GetUsers(ctx, filter)
	if err != nil {
		s.logger.Error("UserService.GetUsers: failed to get users", "error", err)
		return nil, fmt.Errorf("UserService.GetUsers: %w", err)
	}

	responses := make([]*model.UserResponse, 0, len(users))
	for _, user := range users {
		responses = append(responses, s.toUserResponse(user))
	}

	return responses, nil
}

// CreateUser создает нового пользователя.
func (s *Service) CreateUser(ctx context.Context, req model.UserCreateRequest) (*model.UserResponse, error) {
	// Валидация
	if err := s.validateCreateRequest(req); err != nil {
		return nil, fmt.Errorf("UserService.CreateUser: validation failed: %w", err)
	}

	user := &model.User{
		Login:     req.Login,
		Password:  req.Password, // В реальном приложении здесь должно быть хеширование пароля
		IsAdmin:   req.IsAdmin,
		Role:      req.Role,
		Region:    req.Region,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
	}

	createdUser, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		s.logger.Error("UserService.CreateUser: failed to create user", "login", req.Login, "error", err)
		return nil, fmt.Errorf("UserService.CreateUser: %w", err)
	}

	s.logger.Info("UserService.CreateUser: user created successfully", "id", createdUser.ID, "login", createdUser.Login)
	return s.toUserResponse(createdUser), nil
}

// UpdateUser обновляет пользователя.
func (s *Service) UpdateUser(ctx context.Context, id int64, update model.UserUpdate) (*model.UserResponse, error) {
	if id <= 0 {
		return nil, fmt.Errorf("UserService.UpdateUser: invalid user ID")
	}

	// Валидация обновления
	if err := s.validateUpdateRequest(update); err != nil {
		return nil, fmt.Errorf("UserService.UpdateUser: validation failed: %w", err)
	}

	// Если обновляется пароль, здесь должно быть хеширование
	// if update.Password != nil {
	//     hashedPassword := hashPassword(*update.Password)
	//     update.Password = &hashedPassword
	// }

	updatedUser, err := s.repo.UpdateUser(ctx, id, update)
	if err != nil {
		s.logger.Error("UserService.UpdateUser: failed to update user", "id", id, "error", err)
		return nil, fmt.Errorf("UserService.UpdateUser: %w", err)
	}

	s.logger.Info("UserService.UpdateUser: user updated successfully", "id", updatedUser.ID)
	return s.toUserResponse(updatedUser), nil
}

// DeleteUser удаляет пользователя.
func (s *Service) DeleteUser(ctx context.Context, id int64) error {
	if id <= 0 {
		return fmt.Errorf("UserService.DeleteUser: invalid user ID")
	}

	err := s.repo.DeleteUser(ctx, id)
	if err != nil {
		s.logger.Error("UserService.DeleteUser: failed to delete user", "id", id, "error", err)
		return fmt.Errorf("UserService.DeleteUser: %w", err)
	}

	s.logger.Info("UserService.DeleteUser: user deleted successfully", "id", id)
	return nil
}

// validateCreateRequest валидирует запрос на создание пользователя.
func (s *Service) validateCreateRequest(req model.UserCreateRequest) error {
	if req.Login == "" {
		return fmt.Errorf("login is required")
	}
	if req.Password == "" {
		return fmt.Errorf("password is required")
	}
	if len(req.Password) < 6 {
		return fmt.Errorf("password must be at least 6 characters")
	}
	if req.Role == "" {
		return fmt.Errorf("role is required")
	}
	// Проверка валидности роли
	validRoles := map[model.UserRole]bool{
		model.UserRoleAdmin:   true,
		model.UserRoleManager: true,
		model.UserRoleSales:   true,
		model.UserRoleViewer:  true,
	}
	if !validRoles[req.Role] {
		return fmt.Errorf("invalid role: %s", req.Role)
	}
	return nil
}

// validateUpdateRequest валидирует запрос на обновление пользователя.
func (s *Service) validateUpdateRequest(update model.UserUpdate) error {
	if update.Login != nil && *update.Login == "" {
		return fmt.Errorf("login cannot be empty")
	}
	if update.Password != nil && len(*update.Password) < 6 {
		return fmt.Errorf("password must be at least 6 characters")
	}
	if update.Role != nil {
		validRoles := map[model.UserRole]bool{
			model.UserRoleAdmin:   true,
			model.UserRoleManager: true,
			model.UserRoleSales:   true,
			model.UserRoleViewer:  true,
		}
		if !validRoles[*update.Role] {
			return fmt.Errorf("invalid role: %s", *update.Role)
		}
	}
	return nil
}

// toUserResponse преобразует User в UserResponse (без пароля).
func (s *Service) toUserResponse(user *model.User) *model.UserResponse {
	return &model.UserResponse{
		ID:        user.ID,
		Login:     user.Login,
		IsAdmin:   user.IsAdmin,
		Role:      user.Role,
		Region:    user.Region,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}
}
