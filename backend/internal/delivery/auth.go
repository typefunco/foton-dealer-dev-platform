package delivery

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

const (
	ttl = time.Second * 5
)

// LoginRequest представляет запрос на логин
type LoginRequest struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse представляет ответ на логин
type LoginResponse struct {
	Token string `json:"token"`
	User  struct {
		Login   string `json:"login"`
		IsAdmin bool   `json:"is_admin"`
		Role    string `json:"role"`
	} `json:"user"`
}

// Login - ручка логина.
func (s *Server) Login(c echo.Context) error {
	var req LoginRequest
	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid request body",
		})
	}

	ctx, deadline := context.WithTimeout(context.Background(), ttl)
	defer deadline()

	claims, err := s.authService.Login(ctx, req.Login, req.Password)
	if err != nil {
		s.logger.Error("Login failed", "login", req.Login, "error", err)
		return c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "Invalid credentials",
		})
	}

	// Генерируем новый токен для ответа
	token, err := s.authService.GenerateToken(claims.Login, claims.IsAdmin, claims.Role)
	if err != nil {
		s.logger.Error("Failed to generate JWT", "error", err)
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to generate token",
		})
	}

	response := LoginResponse{
		Token: token,
		User: struct {
			Login   string `json:"login"`
			IsAdmin bool   `json:"is_admin"`
			Role    string `json:"role"`
		}{
			Login:   claims.Login,
			IsAdmin: claims.IsAdmin,
			Role:    claims.Role,
		},
	}

	return c.JSON(http.StatusOK, response)
}
