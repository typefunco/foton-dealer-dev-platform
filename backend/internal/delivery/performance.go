package delivery

import (
	"context"
	"github.com/labstack/echo/v4"
	"net/http"
)

// GetPerformance - ручка логина.
func (s *Server) GetPerformance(c echo.Context) error {
	region := c.Param("region")
	if region == "" {
		c.JSON(http.StatusBadRequest, echo.Map{"error": "region query parameter is missing"})
	}
	ctx, deadline := context.WithTimeout(context.Background(), ttl)
	defer deadline()
	data, err := s.perfService.FindPerformances(ctx, region)
	if err != nil || len(data) == 0 {
		c.JSON(http.StatusInternalServerError, err)
	}

	c.JSON(http.StatusOK, data)
	return nil
}
