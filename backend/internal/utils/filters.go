package utils

import (
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/typefunco/dealer_dev_platform/internal/model"
)

// ParseFilterParamsFromContext парсит параметры фильтрации из Echo контекста.
// Универсальная функция для всех эндпоинтов.
func ParseFilterParamsFromContext(c echo.Context) *model.FilterParams {
	filters := &model.FilterParams{
		Region:    c.QueryParam("region"),
		Quarter:   c.QueryParam("quarter"),
		SortBy:    c.QueryParam("sort_by"),
		SortOrder: c.QueryParam("sort_order"),
	}

	// Парсим год
	if yearStr := c.QueryParam("year"); yearStr != "" {
		if year, err := strconv.Atoi(yearStr); err == nil {
			filters.Year = year
		}
	}

	// Парсим лимит и оффсет
	if limitStr := c.QueryParam("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			filters.Limit = limit
		}
	}
	if offsetStr := c.QueryParam("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil {
			filters.Offset = offset
		}
	}

	// Парсим ID дилеров
	if dealerIDsStr := c.QueryParam("dealer_ids"); dealerIDsStr != "" {
		dealerIDs := strings.Split(dealerIDsStr, ",")
		for _, idStr := range dealerIDs {
			if id, err := strconv.Atoi(strings.TrimSpace(idStr)); err == nil {
				filters.DealerIDs = append(filters.DealerIDs, id)
			}
		}
	}

	return filters
}

// SetDefaultFilters устанавливает значения по умолчанию для фильтров.
// Используется для обратной совместимости.
func SetDefaultFilters(filters *model.FilterParams, defaults map[string]interface{}) {
	if filters.Region == "" && defaults["region"] != nil {
		filters.Region = defaults["region"].(string)
	}
	if filters.Quarter == "" && defaults["quarter"] != nil {
		filters.Quarter = defaults["quarter"].(string)
	}
	if filters.Year == 0 && defaults["year"] != nil {
		filters.Year = defaults["year"].(int)
	}
}
