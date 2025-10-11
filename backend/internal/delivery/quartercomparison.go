package delivery

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/typefunco/dealer_dev_platform/internal/model"
)

// GetQuarterComparison возвращает сравнение двух кварталов.
// @Summary Compare two quarters
// @Description Сравнение метрик двух кварталов
// @Tags quarter-comparison
// @Accept json
// @Produce json
// @Param quarter1 query string false "First quarter" default("q1")
// @Param year1 query int false "First year" default(2024)
// @Param quarter2 query string false "Second quarter" default("q2")
// @Param year2 query int false "Second year" default(2024)
// @Success 200 {array} model.QuarterMetrics
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/quarter-comparison [get]
func (s *Server) GetQuarterComparison(c echo.Context) error {
	// Получение параметров квартала 1
	quarter1 := c.QueryParam("quarter1")
	if quarter1 == "" {
		quarter1 = "q1"
	}

	year1Str := c.QueryParam("year1")
	year1 := 2024
	if year1Str != "" {
		parsed, err := strconv.Atoi(year1Str)
		if err == nil {
			year1 = parsed
		}
	}

	// Получение параметров квартала 2
	quarter2 := c.QueryParam("quarter2")
	if quarter2 == "" {
		quarter2 = "q2"
	}

	year2Str := c.QueryParam("year2")
	year2 := 2024
	if year2Str != "" {
		parsed, err := strconv.Atoi(year2Str)
		if err == nil {
			year2 = parsed
		}
	}

	// Вычисляем метрики для обоих кварталов
	metrics1, err := s.calculateQuarterMetrics(c.Request().Context(), quarter1, year1)
	if err != nil {
		s.logger.Error("GetQuarterComparison: failed to calculate metrics for quarter 1",
			"quarter", quarter1,
			"year", year1,
			"error", err,
		)
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to calculate metrics for first quarter",
		})
	}

	metrics2, err := s.calculateQuarterMetrics(c.Request().Context(), quarter2, year2)
	if err != nil {
		s.logger.Error("GetQuarterComparison: failed to calculate metrics for quarter 2",
			"quarter", quarter2,
			"year", year2,
			"error", err,
		)
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to calculate metrics for second quarter",
		})
	}

	// Возвращаем массив из двух метрик
	response := []model.QuarterMetrics{*metrics1, *metrics2}

	return c.JSON(http.StatusOK, response)
}

// calculateQuarterMetrics вычисляет агрегированные метрики за квартал.
func (s *Server) calculateQuarterMetrics(ctx context.Context, quarter string, year int) (*model.QuarterMetrics, error) {
	metrics := &model.QuarterMetrics{
		Quarter: quarter,
		Year:    fmt.Sprintf("%d", year),
	}

	// Получаем данные Dealer Dev
	dealerDevList, err := s.dealerDevService.GetDealerDevByPeriod(ctx, quarter, year, "all-russia")
	if err == nil && len(dealerDevList) > 0 {
		// Вычисляем средний checklist и распределение по классам
		var totalChecklist float64
		classCount := make(map[string]int)

		for _, dd := range dealerDevList {
			totalChecklist += float64(dd.CheckListScore)
			classCount[string(dd.DealerShipClass)]++
		}

		metrics.AverageChecklist = totalChecklist / float64(len(dealerDevList))

		// Распределение по классам в процентах
		metrics.ClassDistribution = make(map[string]float64)
		for class, count := range classCount {
			metrics.ClassDistribution[class] = float64(count) * 100 / float64(len(dealerDevList))
		}

		// Определяем средний класс
		metrics.AverageClass = getMostCommonClass(classCount)
	}

	// Получаем данные Sales
	salesList, err := s.salesService.GetSalesByPeriod(ctx, quarter, year, "all-russia")
	if err == nil && len(salesList) > 0 {
		var totalSalesmen float64
		var trainingCount int
		var totalStockHDT, totalStockMDT, totalStockLDT int
		var totalBuyoutHDT, totalBuyoutMDT, totalBuyoutLDT int

		for _, sales := range salesList {
			totalSalesmen += float64(sales.FotonSalesmen)
			if sales.SalesTrainings {
				trainingCount++
			}
			totalStockHDT += int(sales.StockHDT)
			totalStockMDT += int(sales.StockMDT)
			totalStockLDT += int(sales.StockLDT)
			totalBuyoutHDT += int(sales.BuyoutHDT)
			totalBuyoutMDT += int(sales.BuyoutMDT)
			totalBuyoutLDT += int(sales.BuyoutLDT)
		}

		metrics.AverageFotonSalesmen = totalSalesmen / float64(len(salesList))
		metrics.SalesTrainingsPercentage = float64(trainingCount) * 100 / float64(len(salesList))

		metrics.SalesTrainingsDistribution = map[string]float64{
			"Yes": metrics.SalesTrainingsPercentage,
			"No":  100 - metrics.SalesTrainingsPercentage,
		}

		// Средние значения stock и buyout
		count := len(salesList)
		metrics.StocksData = model.StockDistribution{
			HDT: totalStockHDT / count,
			MDT: totalStockMDT / count,
			LDT: totalStockLDT / count,
		}
		metrics.BuyoutData = model.StockDistribution{
			HDT: totalBuyoutHDT / count,
			MDT: totalBuyoutMDT / count,
			LDT: totalBuyoutLDT / count,
		}

		// Среднее sales target (упрощенно - берем первое значение)
		if len(salesList) > 0 {
			metrics.AverageSalesTarget = salesList[0].SalesTarget
		}
	}

	// Получаем данные Performance
	perfList, err := s.perfService.GetPerformanceByPeriod(ctx, quarter, year, "all-russia")
	if err == nil && len(perfList) > 0 {
		var totalRevenue int64
		var totalProfit, totalMargin, totalRanking, totalMarketing float64
		var totalAsRevenue int64
		var totalAsProfit int64
		var totalAsMargin float64

		for _, perf := range perfList {
			totalRevenue += perf.SalesRevenueRub
			totalProfit += perf.SalesProfitPercent
			totalMargin += perf.SalesMarginPercent
			totalRanking += float64(perf.FotonRank)
			totalMarketing += perf.MarketingInvestment
			totalAsRevenue += perf.AfterSalesRevenueRub
			totalAsProfit += perf.AfterSalesProfitRub
			totalAsMargin += perf.AfterSalesMarginPercent
		}

		count := float64(len(perfList))
		avgRevenue := totalRevenue / int64(len(perfList))
		metrics.AverageSalesRevenue = formatMoney(avgRevenue)
		metrics.AverageSalesProfit = totalProfit / count
		metrics.AverageSalesMargin = totalMargin / count
		metrics.AverageRanking = totalRanking / count
		metrics.MarketingInvestment = totalMarketing
		metrics.AutoSalesRevenue = float64(totalAsRevenue) / 1000000 // В миллионах
		metrics.AutoSalesProfit = float64(totalAsProfit) / 1000000   // В миллионах
		metrics.AutoSalesMargin = totalAsMargin / count
	}

	// Получаем данные After Sales
	asList, err := s.afterSalesService.GetAfterSalesByPeriod(ctx, quarter, year, "all-russia")
	if err == nil && len(asList) > 0 {
		var totalRStock, totalWStock, totalFlh float64
		var asTrainingCount int
		var csiCount int
		var totalWarrantyHours float64

		for _, as := range asList {
			totalRStock += float64(as.RecommendedStock)
			totalWStock += float64(as.WarrantyStock)
			totalFlh += float64(as.FotonLaborHours)
			if as.ASTrainings {
				asTrainingCount++
			}
			if as.CSI != "" && as.CSI != "0" {
				csiCount++
			}
			totalWarrantyHours += float64(as.FotonWarrantyHours)
		}

		count := float64(len(asList))
		metrics.AverageRStockPercent = totalRStock / count
		metrics.AverageWStockPercent = totalWStock / count
		metrics.AverageFlhPercent = totalFlh / count
		metrics.AsTrainingsPercentage = float64(asTrainingCount) * 100 / count
		metrics.CsiPercentage = float64(csiCount) * 100 / count
		metrics.FotonWarrantyHours = totalWarrantyHours / count

		metrics.AsTrainingsDistribution = map[string]float64{
			"Yes": metrics.AsTrainingsPercentage,
			"No":  100 - metrics.AsTrainingsPercentage,
		}
	}

	// Агрегируем решения из всех источников
	metrics.DecisionDistribution = make(map[string]float64)
	// Можно собрать решения из dealer_dev, sales, performance, after_sales
	// Упрощенно - используем данные из performance
	if perfList != nil && len(perfList) > 0 {
		decisionCount := make(map[string]int)
		for _, perf := range perfList {
			decisionCount[string(perf.PerformanceDecision)]++
		}

		for decision, count := range decisionCount {
			metrics.DecisionDistribution[decision] = float64(count) * 100 / float64(len(perfList))
		}
	}

	return metrics, nil
}

// getMostCommonClass возвращает наиболее распространенный класс.
func getMostCommonClass(classCount map[string]int) string {
	maxCount := 0
	mostCommon := "B" // По умолчанию

	for class, count := range classCount {
		if count > maxCount {
			maxCount = count
			mostCommon = class
		}
	}

	return mostCommon
}

// formatMoney форматирует денежное значение в строку вида "5,200,000".
func formatMoney(amount int64) string {
	if amount == 0 {
		return "0"
	}

	str := fmt.Sprintf("%d", amount)
	n := len(str)
	if n <= 3 {
		return str
	}

	var result []byte
	for i, digit := range str {
		if i > 0 && (n-i)%3 == 0 {
			result = append(result, ',')
		}
		result = append(result, byte(digit))
	}

	return string(result)
}
