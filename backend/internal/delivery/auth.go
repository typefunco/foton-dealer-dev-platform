package delivery

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/typefunco/dealer_dev_platform/internal/model"
	"net/http"
	"time"
)

const (
	ttl = time.Second * 5
)

// Login - ручка логина.
func (s *Server) Login(c echo.Context) error {
	var user model.User
	err := c.Bind(&user)
	if err != nil {
		_ = c.JSON(http.StatusBadRequest, err.Error())
		return err
	}
	ctx, deadline := context.WithTimeout(context.Background(), ttl)
	defer deadline()
	err = s.authService.Login(ctx, user.Login, user.Password)
	if err != nil {
		_ = c.JSON(http.StatusBadRequest, err.Error())
		return err
	}

	c.JSON(http.StatusOK, user)
	return nil
}
