package delivery

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// PerformanceDealerResponse представляет дилера с данными Performance для API.
type PerformanceDealerResponse struct {
	ID                  string  `json:"id"`
	Name                string  `json:"name"`
	City                string  `json:"city"`
	SrRub               string  `json:"srRub"`               // Sales Revenue
	SalesProfit         float64 `json:"salesProfit"`         // Sales Profit %
	SalesMargin         float64 `json:"salesMargin"`         // Sales Margin %
	AutoSalesRevenue    string  `json:"autoSalesRevenue"`    // After Sales Revenue
	Rap                 string  `json:"rap"`                 // RAP (можно вычислить позже)
	AutoSalesProfitsRap string  `json:"autoSalesProfitsRap"` // After Sales Profit
	AutoSalesMargin     float64 `json:"autoSalesMargin"`     // After Sales Margin %
	MarketingInvestment float64 `json:"marketingInvestment"` // Marketing Investment (M Rub)
	Ranking             int     `json:"ranking"`             // Foton Ranking
	AutoSalesDecision   string  `json:"autoSalesDecision"`   // Performance Decision
}

// GetPerformanceData возвращает данные производительности по региону и периоду.
// @Summary Get Performance data
// @Description Получение данных производительности с фильтрацией по региону
// @Tags performance
// @Accept json
// @Produce json
// @Param region query string false "Region filter" default("all-russia")
// @Param quarter query string false "Quarter" default("q1")
// @Param year query int false "Year" default(2024)
// @Success 200 {array} PerformanceDealerResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/performance [get]
func (s *Server) GetPerformanceData(c echo.Context) error {
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
	perfList, err := s.perfService.GetPerformanceByPeriod(c.Request().Context(), quarter, year, region)
	if err != nil {
		s.logger.Error("GetPerformanceData: failed to get performance data",
			"region", region,
			"quarter", quarter,
			"year", year,
			"error", err,
		)
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to get performance data",
		})
	}

	// Преобразование в API response
	response := make([]PerformanceDealerResponse, 0, len(perfList))
	for _, perf := range perfList {
		// RAP можно вычислить на основе различных метрик
		rap := calculateRAP(perf.FotonRank)

		response = append(response, PerformanceDealerResponse{
			ID:                  strconv.FormatInt(perf.DealerID, 10),
			Name:                perf.DealerName,
			City:                perf.City,
			SrRub:               perf.SalesRevenueFormatted,
			SalesProfit:         perf.SalesProfitPercent,
			SalesMargin:         perf.SalesMarginPercent,
			AutoSalesRevenue:    perf.AfterSalesRevenueFormatted,
			Rap:                 rap,
			AutoSalesProfitsRap: perf.AfterSalesProfitFormatted,
			AutoSalesMargin:     perf.AfterSalesMarginPercent,
			MarketingInvestment: perf.MarketingInvestment,
			Ranking:             int(perf.FotonRank),
			AutoSalesDecision:   string(perf.PerformanceDecision),
		})
	}

	return c.JSON(http.StatusOK, response)
}

// calculateRAP вычисляет RAP на основе рейтинга.
// Это упрощенная логика - в реальности может быть более сложная.
func calculateRAP(rank int16) string {
	switch {
	case rank >= 1 && rank <= 3:
		return "Gold"
	case rank >= 4 && rank <= 6:
		return "Silver"
	case rank >= 7 && rank <= 10:
		return "Bronze"
	default:
		return "None"
	}
}
