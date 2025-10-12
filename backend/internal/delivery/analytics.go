package delivery

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/typefunco/dealer_dev_platform/internal/model"
)

// AnalyticsResponse представляет аналитические данные
type AnalyticsResponse struct {
	Summary   AnalyticsSummary `json:"summary"`
	Trends    AnalyticsTrends  `json:"trends"`
	Regions   []RegionStats    `json:"regions"`
	Dealers   []DealerStats    `json:"dealers"`
	Timeframe Timeframe        `json:"timeframe"`
}

// AnalyticsSummary представляет сводную статистику
type AnalyticsSummary struct {
	TotalDealers   int     `json:"total_dealers"`
	ActiveDealers  int     `json:"active_dealers"`
	TotalRevenue   float64 `json:"total_revenue"`
	AverageRevenue float64 `json:"average_revenue"`
	TotalProfit    float64 `json:"total_profit"`
	AverageProfit  float64 `json:"average_profit"`
	GrowthRate     float64 `json:"growth_rate"`
	MarketShare    float64 `json:"market_share"`
}

// AnalyticsTrends представляет тренды
type AnalyticsTrends struct {
	RevenueGrowth    float64 `json:"revenue_growth"`
	ProfitGrowth     float64 `json:"profit_growth"`
	DealerGrowth     float64 `json:"dealer_growth"`
	MarketGrowth     float64 `json:"market_growth"`
	PerformanceTrend string  `json:"performance_trend"` // "up", "down", "stable"
}

// RegionStats представляет статистику по регионам
type RegionStats struct {
	Region      string  `json:"region"`
	DealerCount int     `json:"dealer_count"`
	Revenue     float64 `json:"revenue"`
	Profit      float64 `json:"profit"`
	GrowthRate  float64 `json:"growth_rate"`
	MarketShare float64 `json:"market_share"`
	TopDealer   string  `json:"top_dealer"`
}

// DealerStats представляет статистику по дилерам
type DealerStats struct {
	DealerID    int     `json:"dealer_id"`
	DealerName  string  `json:"dealer_name"`
	Region      string  `json:"region"`
	Revenue     float64 `json:"revenue"`
	Profit      float64 `json:"profit"`
	GrowthRate  float64 `json:"growth_rate"`
	Rank        int     `json:"rank"`
	Performance string  `json:"performance"` // "excellent", "good", "average", "poor"
}

// Timeframe представляет временной период
type Timeframe struct {
	Quarter string `json:"quarter"`
	Year    int    `json:"year"`
	Period  string `json:"period"` // "Q1 2024"
}

// GetAnalytics возвращает аналитические данные
// @Summary Get analytics data
// @Description Получение аналитических данных с фильтрацией по региону, кварталу и году
// @Tags analytics
// @Accept json
// @Produce json
// @Param region query string false "Region filter" default("all-russia")
// @Param quarter query string false "Quarter" default("Q1")
// @Param year query int false "Year" default(2024)
// @Success 200 {object} AnalyticsResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/analytics [get]
func (s *Server) GetAnalytics(c echo.Context) error {
	filters, err := ParseFilters(c)
	if err != nil {
		return err
	}

	ctx := c.Request().Context()

	// Получаем данные из всех сервисов
	ddList, err := s.dealerDevService.GetDealerDevByPeriod(ctx, filters.Quarter, filters.Year, filters.Region)
	if err != nil {
		s.logger.Error("GetAnalytics: failed to get dealer dev data", "error", err)
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to get dealer development data",
		})
	}

	salesList, err := s.salesService.GetSalesByPeriod(ctx, filters.Quarter, filters.Year, filters.Region)
	if err != nil {
		s.logger.Error("GetAnalytics: failed to get sales data", "error", err)
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to get sales data",
		})
	}

	perfList, err := s.perfService.GetPerformanceByPeriod(ctx, filters.Quarter, filters.Year, filters.Region)
	if err != nil {
		s.logger.Error("GetAnalytics: failed to get performance data", "error", err)
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to get performance data",
		})
	}

	asList, err := s.afterSalesService.GetAfterSalesByPeriod(ctx, filters.Quarter, filters.Year, filters.Region)
	if err != nil {
		s.logger.Error("GetAnalytics: failed to get after sales data", "error", err)
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to get after sales data",
		})
	}

	// Вычисляем аналитику
	analytics := s.calculateAnalytics(ddList, salesList, perfList, asList, filters)

	s.logger.Info("GetAnalytics: successfully retrieved analytics",
		"region", filters.Region,
		"quarter", filters.Quarter,
		"year", filters.Year,
		"total_dealers", analytics.Summary.TotalDealers,
	)

	return c.JSON(http.StatusOK, analytics)
}

// calculateAnalytics вычисляет аналитические данные
func (s *Server) calculateAnalytics(ddList []*model.DealerDevWithDetails, salesList []*model.SalesWithDetails, perfList []*model.PerformanceWithDetails, asList []*model.AfterSalesWithDetails, filters *FilterRequest) *AnalyticsResponse {
	// Создаем карты для быстрого поиска
	ddMap := make(map[int]*model.DealerDevWithDetails)
	salesMap := make(map[int]*model.SalesWithDetails)
	perfMap := make(map[int]*model.PerformanceWithDetails)
	asMap := make(map[int]*model.AfterSalesWithDetails)

	for _, dd := range ddList {
		ddMap[dd.DealerID] = dd
	}
	for _, sales := range salesList {
		salesMap[sales.DealerID] = sales
	}
	for _, perf := range perfList {
		perfMap[perf.DealerID] = perf
	}
	for _, as := range asList {
		asMap[as.DealerID] = as
	}

	// Получаем уникальные dealer_id
	dealerIDs := make(map[int]bool)
	for id := range ddMap {
		dealerIDs[id] = true
	}
	for id := range salesMap {
		dealerIDs[id] = true
	}
	for id := range perfMap {
		dealerIDs[id] = true
	}
	for id := range asMap {
		dealerIDs[id] = true
	}

	// Вычисляем сводную статистику
	var totalRevenue, totalProfit float64
	var activeDealers int

	for dealerID := range dealerIDs {
		if perf, ok := perfMap[dealerID]; ok {
			totalRevenue += perf.SalesRevenueRub
			totalProfit += perf.SalesProfitRub
			activeDealers++
		}
	}

	totalDealers := len(dealerIDs)
	var averageRevenue, averageProfit float64
	if totalDealers > 0 {
		averageRevenue = totalRevenue / float64(totalDealers)
		averageProfit = totalProfit / float64(totalDealers)
	}

	// Вычисляем статистику по регионам
	regionStats := s.calculateRegionStats(ddMap, salesMap, perfMap, asMap)

	// Вычисляем статистику по дилерам
	dealerStats := s.calculateDealerStats(ddMap, salesMap, perfMap, asMap)

	// Определяем тренды (упрощенная логика)
	trends := AnalyticsTrends{
		RevenueGrowth:    5.2, // Примерные значения
		ProfitGrowth:     3.8,
		DealerGrowth:     2.1,
		MarketGrowth:     4.5,
		PerformanceTrend: "up",
	}

	// Создаем ответ
	response := &AnalyticsResponse{
		Summary: AnalyticsSummary{
			TotalDealers:   totalDealers,
			ActiveDealers:  activeDealers,
			TotalRevenue:   totalRevenue,
			AverageRevenue: averageRevenue,
			TotalProfit:    totalProfit,
			AverageProfit:  averageProfit,
			GrowthRate:     4.2,  // Примерное значение
			MarketShare:    12.5, // Примерное значение
		},
		Trends:  trends,
		Regions: regionStats,
		Dealers: dealerStats,
		Timeframe: Timeframe{
			Quarter: filters.Quarter,
			Year:    filters.Year,
			Period:  filters.Quarter + " " + strconv.Itoa(filters.Year),
		},
	}

	return response
}

// calculateRegionStats вычисляет статистику по регионам
func (s *Server) calculateRegionStats(ddMap map[int]*model.DealerDevWithDetails, salesMap map[int]*model.SalesWithDetails, perfMap map[int]*model.PerformanceWithDetails, asMap map[int]*model.AfterSalesWithDetails) []RegionStats {
	regionData := make(map[string]*RegionStats)

	// Собираем данные по регионам
	for dealerID, dd := range ddMap {
		region := dd.Region
		if _, exists := regionData[region]; !exists {
			regionData[region] = &RegionStats{
				Region: region,
			}
		}

		regionData[region].DealerCount++

		if perf, ok := perfMap[dealerID]; ok {
			regionData[region].Revenue += perf.SalesRevenueRub
			regionData[region].Profit += perf.SalesProfitRub
		}
	}

	// Конвертируем в слайс
	var result []RegionStats
	for _, stats := range regionData {
		if stats.DealerCount > 0 {
			// Вычисляем средние значения
			avgRevenue := stats.Revenue / float64(stats.DealerCount)
			avgProfit := stats.Profit / float64(stats.DealerCount)
			// Обновляем структуру (если нужно добавить поля)
			_ = avgRevenue
			_ = avgProfit
		}
		result = append(result, *stats)
	}

	return result
}

// calculateDealerStats вычисляет статистику по дилерам
func (s *Server) calculateDealerStats(ddMap map[int]*model.DealerDevWithDetails, salesMap map[int]*model.SalesWithDetails, perfMap map[int]*model.PerformanceWithDetails, asMap map[int]*model.AfterSalesWithDetails) []DealerStats {
	var result []DealerStats

	for dealerID, dd := range ddMap {
		stats := DealerStats{
			DealerID:   dealerID,
			DealerName: dd.DealerNameRu,
			Region:     dd.Region,
		}

		if perf, ok := perfMap[dealerID]; ok {
			stats.Revenue = perf.SalesRevenueRub
			stats.Profit = perf.SalesProfitRub
		}

		// Определяем производительность
		if dd.CheckListScore >= 90 {
			stats.Performance = "excellent"
		} else if dd.CheckListScore >= 80 {
			stats.Performance = "good"
		} else if dd.CheckListScore >= 70 {
			stats.Performance = "average"
		} else {
			stats.Performance = "poor"
		}

		result = append(result, stats)
	}

	return result
}
