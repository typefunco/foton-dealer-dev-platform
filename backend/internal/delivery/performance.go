package delivery

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/typefunco/dealer_dev_platform/internal/utils"
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

// GetPerformanceData возвращает данные производительности с поддержкой фильтров.
// @Summary Get Performance data
// @Description Получение данных производительности с фильтрацией по региону, году, кварталу и дилерам
// @Tags performance
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
// @Success 200 {array} PerformanceDealerResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/performance [get]
func (s *Server) GetPerformanceData(c echo.Context) error {
	// Парсим параметры фильтрации с помощью универсальной функции
	filters := utils.ParseFilterParamsFromContext(c)

	// Устанавливаем значения по умолчанию для обратной совместимости
	defaults := map[string]interface{}{
		"region":  "all-russia",
		"quarter": "Q1",
		"year":    2024,
	}
	utils.SetDefaultFilters(filters, defaults)

	// Получение данных из сервиса с фильтрами
	perfList, err := s.perfService.GetPerformanceWithFilters(c.Request().Context(), filters)
	if err != nil {
		s.logger.Error("GetPerformanceData: failed to get performance data",
			"filters", filters,
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
		rap := calculateRAP(int16(perf.FotonRank))

		response = append(response, PerformanceDealerResponse{
			ID:                  strconv.FormatInt(int64(perf.DealerID), 10),
			Name:                perf.DealerNameRu,
			City:                perf.City,
			SrRub:               formatMoney(int64(perf.SalesRevenueRub)),
			SalesProfit:         perf.SalesProfitRub,
			SalesMargin:         perf.SalesMarginPercent,
			AutoSalesRevenue:    formatMoney(int64(perf.AfterSalesRevenueRub)),
			Rap:                 rap,
			AutoSalesProfitsRap: formatMoney(int64(perf.AfterSalesProfitRub)),
			AutoSalesMargin:     perf.AfterSalesMarginPercent,
			MarketingInvestment: perf.MarketingInvestment,
			Ranking:             perf.FotonRank,
			AutoSalesDecision:   perf.PerformanceDecision,
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
