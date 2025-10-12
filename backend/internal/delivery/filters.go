package delivery

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/typefunco/dealer_dev_platform/internal/utils"
)

// FilterRequest представляет запрос с фильтрами
type FilterRequest struct {
	Region   string `json:"region" query:"region"`
	Quarter  string `json:"quarter" query:"quarter"`
	Year     int    `json:"year" query:"year"`
	Search   string `json:"search" query:"search"`
	Page     int    `json:"page" query:"page"`
	Limit    int    `json:"limit" query:"limit"`
	DealerID int    `json:"dealer_id" query:"dealer_id"`
}

// FilterResponse представляет ответ с фильтрами
type FilterResponse struct {
	Data       interface{} `json:"data"`
	Pagination Pagination  `json:"pagination"`
	Filters    FilterInfo  `json:"filters"`
}

// Pagination представляет информацию о пагинации
type Pagination struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

// FilterInfo представляет информацию о примененных фильтрах
type FilterInfo struct {
	Region  string `json:"region"`
	Quarter string `json:"quarter"`
	Year    int    `json:"year"`
	Search  string `json:"search"`
}

// ParseFilters парсит и валидирует параметры фильтрации из запроса
func ParseFilters(c echo.Context) (*FilterRequest, error) {
	filters := &FilterRequest{
		Region:  c.QueryParam("region"),
		Quarter: c.QueryParam("quarter"),
		Search:  c.QueryParam("search"),
		Page:    1,
		Limit:   10,
	}

	// Парсим год
	if yearStr := c.QueryParam("year"); yearStr != "" {
		year, err := strconv.Atoi(yearStr)
		if err != nil {
			return nil, echo.NewHTTPError(http.StatusBadRequest, "Invalid year parameter")
		}
		filters.Year = year
	} else {
		filters.Year = 2024 // Значение по умолчанию
	}

	// Парсим страницу
	if pageStr := c.QueryParam("page"); pageStr != "" {
		page, err := strconv.Atoi(pageStr)
		if err != nil {
			return nil, echo.NewHTTPError(http.StatusBadRequest, "Invalid page parameter")
		}
		filters.Page = page
	}

	// Парсим лимит
	if limitStr := c.QueryParam("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			return nil, echo.NewHTTPError(http.StatusBadRequest, "Invalid limit parameter")
		}
		filters.Limit = limit
	}

	// Парсим dealer_id
	if dealerIDStr := c.QueryParam("dealer_id"); dealerIDStr != "" {
		dealerID, err := strconv.Atoi(dealerIDStr)
		if err != nil {
			return nil, echo.NewHTTPError(http.StatusBadRequest, "Invalid dealer_id parameter")
		}
		filters.DealerID = dealerID
	}

	// Устанавливаем значения по умолчанию
	if filters.Region == "" {
		filters.Region = "all-russia"
	}
	if filters.Quarter == "" {
		filters.Quarter = "Q1"
	}

	// Валидируем параметры
	if err := utils.ValidateFilters(filters.Quarter, filters.Year, filters.Region); err != nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return filters, nil
}

// CreateFilterResponse создает ответ с фильтрами
func CreateFilterResponse(data interface{}, filters *FilterRequest, total int) *FilterResponse {
	totalPages := (total + filters.Limit - 1) / filters.Limit
	if totalPages == 0 {
		totalPages = 1
	}

	return &FilterResponse{
		Data: data,
		Pagination: Pagination{
			Page:       filters.Page,
			Limit:      filters.Limit,
			Total:      total,
			TotalPages: totalPages,
		},
		Filters: FilterInfo{
			Region:  filters.Region,
			Quarter: filters.Quarter,
			Year:    filters.Year,
			Search:  filters.Search,
		},
	}
}

// GetAvailableFilters возвращает доступные фильтры
// @Summary Get available filters
// @Description Получение списка доступных фильтров (регионы, кварталы, годы)
// @Tags filters
// @Accept json
// @Produce json
// @Success 200 {object} AvailableFiltersResponse
// @Router /api/filters [get]
func (s *Server) GetAvailableFilters(c echo.Context) error {
	response := map[string]interface{}{
		"regions": []map[string]string{
			{"id": "all-russia", "name": "All Russia"},
			{"id": "Central", "name": "Central"},
			{"id": "North West", "name": "North West"},
			{"id": "Volga", "name": "Volga"},
			{"id": "South", "name": "South"},
			{"id": "Kavkaz", "name": "Kavkaz"},
			{"id": "Ural", "name": "Ural"},
			{"id": "Siberia", "name": "Siberia"},
			{"id": "Far East", "name": "Far East"},
		},
		"quarters": []map[string]string{
			{"id": "Q1", "name": "Q1"},
			{"id": "Q2", "name": "Q2"},
			{"id": "Q3", "name": "Q3"},
			{"id": "Q4", "name": "Q4"},
		},
		"years": []map[string]interface{}{
			{"id": 2024, "name": "2024"},
			{"id": 2025, "name": "2025"},
			{"id": 2026, "name": "2026"},
			{"id": 2027, "name": "2027"},
		},
	}

	return c.JSON(http.StatusOK, response)
}
