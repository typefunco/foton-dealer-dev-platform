package delivery

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// DealerDevResponse представляет дилера с данными Dealer Development для API.
type DealerDevResponse struct {
	ID                      string   `json:"id"`
	Name                    string   `json:"name"`
	City                    string   `json:"city"`
	Class                   string   `json:"class"`                   // A, B, C, D
	Checklist               int      `json:"checklist"`               // 0-100
	BrandsInPortfolio       []string `json:"brandsInPortfolio"`       // Массив брендов
	BrandsCount             int      `json:"brandsCount"`             // Количество брендов
	Branding                bool     `json:"branding"`                // Наличие брендинга
	BuySideBusiness         []string `json:"buySideBusiness"`         // Типы побочного бизнеса
	DealerDevRecommendation string   `json:"dealerDevRecommendation"` // Рекомендация
}

// GetDealerDevData возвращает данные Dealer Development по региону и периоду.
// @Summary Get Dealer Development data
// @Description Получение данных развития дилеров с фильтрацией по региону
// @Tags dealerdev
// @Accept json
// @Produce json
// @Param region query string false "Region filter" default("all-russia")
// @Param quarter query string false "Quarter" default("q1")
// @Param year query int false "Year" default(2024)
// @Success 200 {array} DealerDevResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/dealerdev [get]
func (s *Server) GetDealerDevData(c echo.Context) error {
	// Получение параметров из query string
	region := c.QueryParam("region")
	if region == "" {
		region = "all-russia"
	}

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
	ddList, err := s.dealerDevService.GetDealerDevByPeriod(c.Request().Context(), quarter, year, region)
	if err != nil {
		s.logger.Error("GetDealerDevData: failed to get dealer dev data",
			"region", region,
			"quarter", quarter,
			"year", year,
			"error", err,
		)
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to get dealer development data",
		})
	}

	// Преобразование в API response
	response := make([]DealerDevResponse, 0, len(ddList))
	for _, dd := range ddList {
		response = append(response, DealerDevResponse{
			ID:                      strconv.FormatInt(dd.DealerID, 10),
			Name:                    dd.DealerName,
			City:                    dd.City,
			Class:                   string(dd.DealerShipClass),
			Checklist:               int(dd.CheckListScore),
			BrandsInPortfolio:       dd.Brands,
			BrandsCount:             dd.BrandsCount,
			Branding:                dd.Branding,
			BuySideBusiness:         dd.BySideBusinesses,
			DealerDevRecommendation: string(dd.Recommendation),
		})
	}

	return c.JSON(http.StatusOK, response)
}
