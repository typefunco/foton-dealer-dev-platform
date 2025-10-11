package delivery

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

// Health - health check.
func (s *Server) Health(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := s.authService.PingDatabase(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}
