package delivery

import (
	"crypto/rand"
	"math/big"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/typefunco/dealer_dev_platform/internal/model"
)

// CreateUserRequest представляет запрос на создание пользователя через API.
type CreateUserRequest struct {
	Email     string `json:"email" validate:"required,email"`
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
	Region    string `json:"region" validate:"required"`
	Position  string `json:"position" validate:"required"`
}

// UpdateUserRequest представляет запрос на обновление пользователя через API.
type UpdateUserRequest struct {
	Email     *string `json:"email,omitempty" validate:"omitempty,email"`
	FirstName *string `json:"firstName,omitempty"`
	LastName  *string `json:"lastName,omitempty"`
	Region    *string `json:"region,omitempty"`
	Position  *string `json:"position,omitempty"`
	Status    *string `json:"status,omitempty"`
}

// UserFilterRequest представляет параметры фильтрации из query string.
type UserFilterRequest struct {
	SearchTerm string `query:"search"`
	Region     string `query:"region"`
	Position   string `query:"position"`
	Page       int    `query:"page"`
	Limit      int    `query:"limit"`
}

// UserAPIResponse представляет пользователя для API (с дополнительными полями для фронтенда).
type UserAPIResponse struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Region    string `json:"region"`
	Position  string `json:"position"`
	CreatedAt string `json:"createdAt"`
	Status    string `json:"status"`
}

// RegionStatsResponse представляет статистику по региону.
type RegionStatsResponse struct {
	Region    string            `json:"region"`
	UserCount int               `json:"userCount"`
	Users     []UserAPIResponse `json:"users"`
}

// CreateUserResponse возвращается при успешном создании пользователя.
type CreateUserResponse struct {
	User        UserAPIResponse  `json:"user"`
	Credentials *UserCredentials `json:"credentials,omitempty"`
}

// UserCredentials содержит сгенерированные учетные данные.
type UserCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// GetUsers возвращает список пользователей с фильтрацией.
// @Summary Get users list
// @Description Получение списка пользователей с возможностью фильтрации
// @Tags users
// @Accept json
// @Produce json
// @Param search query string false "Поиск по имени"
// @Param region query string false "Фильтр по региону"
// @Param position query string false "Фильтр по позиции"
// @Param page query int false "Номер страницы" default(1)
// @Param limit query int false "Количество на странице" default(10)
// @Success 200 {array} UserAPIResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/users [get]
func (s *Server) GetUsers(c echo.Context) error {
	var filterReq UserFilterRequest
	if err := c.Bind(&filterReq); err != nil {
		s.logger.Error("GetUsers: failed to bind filter request", "error", err)
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid filter parameters",
		})
	}

	// Установка значений по умолчанию для пагинации
	if filterReq.Page < 1 {
		filterReq.Page = 1
	}
	if filterReq.Limit < 1 {
		filterReq.Limit = 10
	}

	// Построение фильтра для сервиса
	filter := model.UserFilter{}

	if filterReq.Region != "" {
		filter.Region = &filterReq.Region
	}

	// Position мапится на Role в нашей модели
	if filterReq.Position != "" {
		role := model.UserRole(filterReq.Position)
		filter.Role = &role
	}

	// Получение пользователей из сервиса
	users, err := s.userService.GetUsers(c.Request().Context(), filter)
	if err != nil {
		s.logger.Error("GetUsers: failed to get users", "error", err)
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to get users",
		})
	}

	// Фильтрация по searchTerm (по имени) на уровне приложения
	filteredUsers := users
	if filterReq.SearchTerm != "" {
		filteredUsers = filterUsersByName(users, filterReq.SearchTerm)
	}

	// Применение пагинации
	totalUsers := len(filteredUsers)
	startIdx := (filterReq.Page - 1) * filterReq.Limit
	endIdx := startIdx + filterReq.Limit

	if startIdx >= totalUsers {
		filteredUsers = []*model.UserResponse{}
	} else {
		if endIdx > totalUsers {
			endIdx = totalUsers
		}
		filteredUsers = filteredUsers[startIdx:endIdx]
	}

	// Преобразование в API response
	apiUsers := make([]UserAPIResponse, 0, len(filteredUsers))
	for _, user := range filteredUsers {
		apiUsers = append(apiUsers, toUserAPIResponse(user))
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"users": apiUsers,
		"pagination": map[string]interface{}{
			"page":       filterReq.Page,
			"limit":      filterReq.Limit,
			"total":      totalUsers,
			"totalPages": (totalUsers + filterReq.Limit - 1) / filterReq.Limit,
		},
	})
}

// GetUserByID возвращает пользователя по ID.
// @Summary Get user by ID
// @Description Получение пользователя по ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} UserAPIResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/users/{id} [get]
func (s *Server) GetUserByID(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid user ID",
		})
	}

	user, err := s.userService.GetUserByID(c.Request().Context(), id)
	if err != nil {
		s.logger.Error("GetUserByID: failed to get user", "id", id, "error", err)
		return c.JSON(http.StatusNotFound, ErrorResponse{
			Error: "User not found",
		})
	}

	return c.JSON(http.StatusOK, toUserAPIResponse(user))
}

// CreateUser создает нового пользователя.
// @Summary Create new user
// @Description Создание нового пользователя с автоматической генерацией пароля
// @Tags users
// @Accept json
// @Produce json
// @Param user body CreateUserRequest true "User data"
// @Success 201 {object} CreateUserResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/users [post]
func (s *Server) CreateUser(c echo.Context) error {
	var req CreateUserRequest
	if err := c.Bind(&req); err != nil {
		s.logger.Error("CreateUser: failed to bind request", "error", err)
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid request body",
		})
	}

	// Валидация
	if req.Email == "" || req.FirstName == "" || req.LastName == "" || req.Region == "" || req.Position == "" {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "All fields are required: email, firstName, lastName, region, position",
		})
	}

	// Генерация пароля
	password := generatePassword()

	// Создание login из email (часть до @)
	login := req.Email

	// Маппинг position на role (упрощенная логика)
	role := mapPositionToRole(req.Position)

	// Создание пользователя через сервис
	createReq := model.UserCreateRequest{
		Login:     login,
		Password:  password,
		IsAdmin:   role == model.UserRoleAdmin,
		Role:      role,
		Region:    req.Region,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
	}

	user, err := s.userService.CreateUser(c.Request().Context(), createReq)
	if err != nil {
		s.logger.Error("CreateUser: failed to create user", "email", req.Email, "error", err)
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to create user",
		})
	}

	// Возврат пользователя с учетными данными
	response := CreateUserResponse{
		User: toUserAPIResponse(user),
		Credentials: &UserCredentials{
			Email:    user.Email,
			Password: password,
		},
	}

	return c.JSON(http.StatusCreated, response)
}

// UpdateUser обновляет пользователя.
// @Summary Update user
// @Description Обновление данных пользователя
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body UpdateUserRequest true "User data"
// @Success 200 {object} UserAPIResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/users/{id} [put]
func (s *Server) UpdateUser(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid user ID",
		})
	}

	var req UpdateUserRequest
	if err := c.Bind(&req); err != nil {
		s.logger.Error("UpdateUser: failed to bind request", "error", err)
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid request body",
		})
	}

	// Построение update модели
	update := model.UserUpdate{
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Region:    req.Region,
	}

	// Маппинг position на role если передан
	if req.Position != nil {
		role := mapPositionToRole(*req.Position)
		update.Role = &role
	}

	// Обновление через сервис
	user, err := s.userService.UpdateUser(c.Request().Context(), id, update)
	if err != nil {
		s.logger.Error("UpdateUser: failed to update user", "id", id, "error", err)
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to update user",
		})
	}

	return c.JSON(http.StatusOK, toUserAPIResponse(user))
}

// DeleteUser удаляет пользователя.
// @Summary Delete user
// @Description Удаление пользователя
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 204
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/users/{id} [delete]
func (s *Server) DeleteUser(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid user ID",
		})
	}

	err = s.userService.DeleteUser(c.Request().Context(), id)
	if err != nil {
		s.logger.Error("DeleteUser: failed to delete user", "id", id, "error", err)
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to delete user",
		})
	}

	return c.NoContent(http.StatusNoContent)
}

// GetUserStats возвращает статистику пользователей по регионам.
// @Summary Get user statistics
// @Description Получение статистики пользователей по регионам
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} ErrorResponse
// @Router /api/users/stats [get]
func (s *Server) GetUserStats(c echo.Context) error {
	// Получение всех пользователей
	users, err := s.userService.GetUsers(c.Request().Context(), model.UserFilter{})
	if err != nil {
		s.logger.Error("GetUserStats: failed to get users", "error", err)
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to get user statistics",
		})
	}

	// Список всех регионов
	regions := []string{
		"Central", "Caucasus", "Volga", "Ural",
		"Siberia", "Far East", "North-West", "South",
	}

	// Группировка по регионам
	regionMap := make(map[string][]*model.UserResponse)
	for _, user := range users {
		regionMap[user.Region] = append(regionMap[user.Region], user)
	}

	// Формирование статистики
	regionStats := make([]RegionStatsResponse, 0, len(regions))
	for _, region := range regions {
		usersInRegion := regionMap[region]
		apiUsers := make([]UserAPIResponse, 0, len(usersInRegion))
		for _, user := range usersInRegion {
			apiUsers = append(apiUsers, toUserAPIResponse(user))
		}

		regionStats = append(regionStats, RegionStatsResponse{
			Region:    region,
			UserCount: len(usersInRegion),
			Users:     apiUsers,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"totalUsers":  len(users),
		"regionStats": regionStats,
	})
}

// ErrorResponse представляет ошибку API.
type ErrorResponse struct {
	Error string `json:"error"`
}

// Вспомогательные функции

// toUserAPIResponse преобразует UserResponse в UserAPIResponse.
func toUserAPIResponse(user *model.UserResponse) UserAPIResponse {
	return UserAPIResponse{
		ID:        strconv.FormatInt(user.ID, 10),
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Region:    user.Region,
		Position:  string(user.Role), // Role мапится на Position
		CreatedAt: user.CreatedAt.Format("2006-01-02"),
		Status:    "active", // По умолчанию все пользователи активны
	}
}

// filterUsersByName фильтрует пользователей по имени или фамилии.
func filterUsersByName(users []*model.UserResponse, searchTerm string) []*model.UserResponse {
	if searchTerm == "" {
		return users
	}

	filtered := make([]*model.UserResponse, 0)
	searchLower := strings.ToLower(searchTerm)

	for _, user := range users {
		firstNameLower := strings.ToLower(user.FirstName)
		lastNameLower := strings.ToLower(user.LastName)

		if strings.Contains(firstNameLower, searchLower) || strings.Contains(lastNameLower, searchLower) {
			filtered = append(filtered, user)
		}
	}

	return filtered
}

// generatePassword генерирует криптографически стойкий случайный пароль длиной 12 символов.
func generatePassword() string {
	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!@#$%^&*"
	const passwordLength = 12

	password := make([]byte, passwordLength)
	charsLength := big.NewInt(int64(len(chars)))

	for i := 0; i < passwordLength; i++ {
		num, err := rand.Int(rand.Reader, charsLength)
		if err != nil {
			// Fallback на простую генерацию в случае ошибки
			password[i] = chars[i%len(chars)]
			continue
		}
		password[i] = chars[num.Int64()]
	}

	return string(password)
}

// mapPositionToRole маппит должность на роль в системе.
func mapPositionToRole(position string) model.UserRole {
	switch position {
	case "Sales Director", "Regional Director", "Head of Sales":
		return model.UserRoleAdmin
	case "Regional Manager", "Sales Manager", "Senior Sales Manager":
		return model.UserRoleManager
	case "Account Manager", "Account Executive", "Business Development Manager":
		return model.UserRoleSales
	case "Sales Representative":
		return model.UserRoleViewer
	default:
		return model.UserRoleViewer
	}
}
