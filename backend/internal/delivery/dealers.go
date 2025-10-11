package delivery

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/typefunco/dealer_dev_platform/internal/model"
)

// GetDealerCard возвращает полную карточку дилера.
// @Summary Get dealer card
// @Description Получение полной информации о дилере (карточка дилера)
// @Tags dealers
// @Accept json
// @Produce json
// @Param id path int true "Dealer ID"
// @Param quarter query string false "Quarter" default("q1")
// @Param year query int false "Year" default(2024)
// @Success 200 {object} model.DealerCardData
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/dealers/{id}/card [get]
func (s *Server) GetDealerCard(c echo.Context) error {
	// Получение ID дилера из параметров пути
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid dealer ID",
		})
	}

	// Получение параметров квартала и года
	quarter := c.QueryParam("quarter")
	if quarter == "" {
		quarter = "q1"
	}

	yearStr := c.QueryParam("year")
	year := 2024
	if yearStr != "" {
		parsedYear, err := strconv.Atoi(yearStr)
		if err == nil {
			year = parsedYear
		}
	}

	// Получение данных из сервиса
	cardData, err := s.dealerService.GetDealerCard(c.Request().Context(), id, quarter, year)
	if err != nil {
		s.logger.Error("GetDealerCard: failed to get dealer card",
			"id", id,
			"quarter", quarter,
			"year", year,
			"error", err,
		)
		return c.JSON(http.StatusNotFound, ErrorResponse{
			Error: "Dealer not found or data unavailable",
		})
	}

	return c.JSON(http.StatusOK, cardData)
}

// GetDealers возвращает список дилеров.
// @Summary Get dealers list
// @Description Получение списка дилеров с возможностью фильтрации по региону
// @Tags dealers
// @Accept json
// @Produce json
// @Param region query string false "Region filter"
// @Success 200 {array} model.Dealer
// @Failure 500 {object} ErrorResponse
// @Router /api/dealers [get]
func (s *Server) GetDealers(c echo.Context) error {
	region := c.QueryParam("region")

	var dealers []*model.Dealer
	var err error

	if region != "" {
		dealers, err = s.dealerService.GetDealersByRegion(c.Request().Context(), region)
	} else {
		dealers, err = s.dealerService.GetAllDealers(c.Request().Context())
	}

	if err != nil {
		s.logger.Error("GetDealers: failed to get dealers",
			"region", region,
			"error", err,
		)
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to get dealers",
		})
	}

	return c.JSON(http.StatusOK, dealers)
}

// GetDealerByID возвращает дилера по ID.
// @Summary Get dealer by ID
// @Description Получение базовой информации о дилере
// @Tags dealers
// @Accept json
// @Produce json
// @Param id path int true "Dealer ID"
// @Success 200 {object} model.Dealer
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/dealers/{id} [get]
func (s *Server) GetDealerByID(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid dealer ID",
		})
	}

	dealer, err := s.dealerService.GetDealerByID(c.Request().Context(), id)
	if err != nil {
		s.logger.Error("GetDealerByID: failed to get dealer",
			"id", id,
			"error", err,
		)
		return c.JSON(http.StatusNotFound, ErrorResponse{
			Error: "Dealer not found",
		})
	}

	return c.JSON(http.StatusOK, dealer)
}
