package delivery

import (
	"net/http"
	"strconv"

	"log/slog"

	"github.com/labstack/echo/v4"
	"github.com/typefunco/dealer_dev_platform/internal/model"
	"github.com/typefunco/dealer_dev_platform/internal/utils"
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
		quarter = "Q1"
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
	var cardData *model.DealerCardData

	// Если указаны год и квартал, используем Excel данные
	if year > 0 && quarter != "" {
		s.logger.Info("Getting dealer card from Excel table",
			slog.Int64("id", id),
			slog.Int("year", year),
			slog.String("quarter", quarter),
		)

		// Получаем карточку дилера по ID из Excel таблицы
		cardData, err = s.dealerService.GetDealerByIDFromExcel(c.Request().Context(), year, quarter, int(id))
		if err != nil {
			s.logger.Error("GetDealerCard: failed to get dealer card from Excel",
				slog.Int64("id", id),
				slog.Int("year", year),
				slog.String("quarter", quarter),
				slog.String("error", err.Error()),
			)
			return c.JSON(http.StatusNotFound, ErrorResponse{
				Error: "Dealer data not found in Excel table",
			})
		}
	} else {
		// Используем обычные данные
		cardData, err = s.dealerService.GetDealerCard(c.Request().Context(), id, quarter, year)
		if err != nil {
			s.logger.Error("GetDealerCard: failed to get dealer card",
				slog.Int64("id", id),
				slog.String("quarter", quarter),
				slog.Int("year", year),
				slog.String("error", err.Error()),
			)
			return c.JSON(http.StatusNotFound, ErrorResponse{
				Error: "Dealer not found or data unavailable",
			})
		}
	}

	return c.JSON(http.StatusOK, cardData)
}

// GetDealers возвращает список дилеров.
// @Summary Get dealers list
// @Description Получение списка дилеров с возможностью фильтрации по региону, году, кварталу и дилерам
// @Tags dealers
// @Accept json
// @Produce json
// @Param region query string false "Region filter"
// @Param quarter query string false "Quarter filter (Q1, Q2, Q3, Q4)"
// @Param year query int false "Year filter"
// @Param dealer_ids query string false "Comma-separated dealer IDs"
// @Param limit query int false "Limit for pagination"
// @Param offset query int false "Offset for pagination"
// @Param sort_by query string false "Sort field"
// @Param sort_order query string false "Sort order (asc, desc)"
// @Success 200 {array} model.Dealer
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/dealers [get]
func (s *Server) GetDealers(c echo.Context) error {
	// Парсим параметры фильтрации с помощью универсальной функции
	filters := utils.ParseFilterParamsFromContext(c)

	var dealers []*model.Dealer
	var err error

	// Если указаны год и квартал, используем Excel данные
	if filters.Year > 0 && filters.Quarter != "" {
		s.logger.Info("Getting dealers from Excel table",
			slog.Int("year", filters.Year),
			slog.String("quarter", filters.Quarter),
			slog.String("region", filters.Region),
		)

		dealers, err = s.dealerService.GetDealersFromExcel(c.Request().Context(), filters.Year, filters.Quarter, filters)
		if err != nil {
			s.logger.Error("GetDealers: failed to get dealers from Excel",
				slog.Int("year", filters.Year),
				slog.String("quarter", filters.Quarter),
				slog.String("error", err.Error()),
			)
			return c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error: "Failed to get dealers from Excel data",
			})
		}
	} else {
		// Используем обычные данные из основной таблицы
		if filters.Region != "" || len(filters.DealerIDs) > 0 || filters.Limit > 0 || filters.Offset > 0 {
			dealers, err = s.dealerService.GetDealersWithFilters(c.Request().Context(), filters)
		} else {
			dealers, err = s.dealerService.GetAllDealers(c.Request().Context())
		}

		if err != nil {
			s.logger.Error("GetDealers: failed to get dealers",
				slog.String("error", err.Error()),
			)
			return c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error: "Failed to get dealers",
			})
		}
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

	dealer, err := s.dealerService.GetDealerByID(c.Request().Context(), int(id))
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

// GetDealersList возвращает упрощенный список дилеров для выпадающих меню.
// @Summary Get dealers list for dropdown
// @Description Получение упрощенного списка дилеров для использования в UI (выпадающие меню, селекторы)
// @Tags dealers
// @Accept json
// @Produce json
// @Param region query string false "Region filter"
// @Param limit query int false "Limit for pagination" default(100)
// @Param offset query int false "Offset for pagination" default(0)
// @Success 200 {array} DealerListItem
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/dealers/list [get]
func (s *Server) GetDealersList(c echo.Context) error {
	// Парсим параметры фильтрации
	filters := utils.ParseFilterParamsFromContext(c)

	// Устанавливаем разумные значения по умолчанию для списка дилеров
	if filters.Limit == 0 {
		filters.Limit = 100 // Загружаем до 100 дилеров сразу
	}
	if filters.Offset < 0 {
		filters.Offset = 0
	}

	// Получаем дилеров с фильтрами
	dealers, err := s.dealerService.GetDealersWithFilters(c.Request().Context(), filters)
	if err != nil {
		s.logger.Error("GetDealersList: failed to get dealers",
			"filters", filters,
			"error", err,
		)
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to get dealers list",
		})
	}

	// Преобразуем в упрощенный формат для фронтенда
	dealerList := make([]DealerListItem, 0, len(dealers))
	for _, dealer := range dealers {
		dealerList = append(dealerList, DealerListItem{
			ID:      dealer.DealerID,
			Name:    dealer.DealerNameRu,
			Region:  dealer.Region,
			City:    dealer.City,
			Manager: dealer.Manager,
		})
	}

	return c.JSON(http.StatusOK, dealerList)
}

// DealerListItem представляет упрощенную информацию о дилере для UI.
type DealerListItem struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Region  string `json:"region"`
	City    string `json:"city"`
	Manager string `json:"manager"`
}
