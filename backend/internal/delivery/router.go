package delivery

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/typefunco/dealer_dev_platform/internal/service/auth"
	"log/slog"
	"net/http"
)

// Server структура сервера.
type Server struct {
	authService *auth.Service
	srv         *echo.Echo
	logger      *slog.Logger
}

// NewServer - конструктор сервера.
func NewServer(authService *auth.Service, logger *slog.Logger) *Server {
	return &Server{
		authService: authService,
		srv:         echo.New(),
		logger:      logger,
	}
}

// RunServer - команда запуска сервера.
func (s *Server) RunServer() {
	s.srv.Use(middleware.Logger())
	s.srv.Use(middleware.Recover())

	// TODO: переделать на http
	s.srv.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	}))

	s.srv.POST("/auth/login", s.Login)

	if err := s.srv.Start(":8080"); err != nil && !errors.Is(err, http.ErrServerClosed) {
		s.logger.Error("failed to start server", "error", err)
	}
}
