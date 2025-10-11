package delivery

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// SalesTeamDealerResponse представляет дилера с данными Sales Team для API.
type SalesTeamDealerResponse struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	City            string `json:"city"`
	SalesManager    string `json:"salesManager"`
	SalesTarget     string `json:"salesTarget"`     // "40/100"
	StockHdtMdtLdt  string `json:"stockHdtMdtLdt"`  // "5/2/3"
	BuyoutHdtMdtLdt string `json:"buyoutHdtMdtLdt"` // "5/2/3"
	FotonSalesmen   int    `json:"fotonSalesmen"`   // Количество продавцов
	SalesTrainings  bool   `json:"salesTrainings"`  // Пройдены ли тренинги
	SalesDecision   string `json:"salesDecision"`   // Решение по продажам
}

// GetSalesTeamData возвращает данные команды продаж по региону и периоду.
// @Summary Get Sales Team data
// @Description Получение данных команды продаж с фильтрацией по региону
// @Tags sales
// @Accept json
// @Produce json
// @Param region query string false "Region filter" default("all-russia")
// @Param quarter query string false "Quarter" default("q1")
// @Param year query int false "Year" default(2024)
// @Success 200 {array} SalesTeamDealerResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/sales [get]
func (s *Server) GetSalesTeamData(c echo.Context) error {
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
	salesList, err := s.salesService.GetSalesByPeriod(c.Request().Context(), quarter, year, region)
	if err != nil {
		s.logger.Error("GetSalesTeamData: failed to get sales data",
			"region", region,
			"quarter", quarter,
			"year", year,
			"error", err,
		)
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to get sales team data",
		})
	}

	// Преобразование в API response
	response := make([]SalesTeamDealerResponse, 0, len(salesList))
	for _, sales := range salesList {
		response = append(response, SalesTeamDealerResponse{
			ID:              strconv.FormatInt(sales.DealerID, 10),
			Name:            sales.DealerName,
			City:            sales.City,
			SalesManager:    sales.Manager,
			SalesTarget:     sales.SalesTarget,
			StockHdtMdtLdt:  sales.StockHdtMdtLdt,
			BuyoutHdtMdtLdt: sales.BuyoutHdtMdtLdt,
			FotonSalesmen:   int(sales.FotonSalesmen),
			SalesTrainings:  sales.SalesTrainings,
			SalesDecision:   string(sales.SalesDecision),
		})
	}

	return c.JSON(http.StatusOK, response)
}
