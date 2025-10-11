package delivery

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/typefunco/dealer_dev_platform/internal/model"
)

// AllDealerData представляет комплексные данные дилера со всех таблиц.
type AllDealerData struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	City         string `json:"city"`
	SalesManager string `json:"salesManager"`

	// Dealer Development fields
	Class                   string   `json:"class"`
	Checklist               int      `json:"checklist"`
	BrandsInPortfolio       []string `json:"brandsInPortfolio"`
	DealerDevRecommendation string   `json:"dealerDevRecommendation"`

	// Sales Team fields
	SalesTarget     string `json:"salesTarget"`
	StockHdtMdtLdt  string `json:"stockHdtMdtLdt"`
	BuyoutHdtMdtLdt string `json:"buyoutHdtMdtLdt"`
	FotonSalesmen   int    `json:"fotonSalesmen"`
	SalesTrainings  bool   `json:"salesTrainings"`
	SalesDecision   string `json:"salesDecision"`

	// Performance fields
	SrRub               string  `json:"srRub"`               // Sales Revenue
	SalesProfit         string  `json:"salesProfit"`         // Sales Profit formatted
	SalesMargin         float64 `json:"salesMargin"`         // Sales Margin %
	AutoSalesRevenue    string  `json:"autoSalesRevenue"`    // After Sales Revenue
	AutoSalesProfitsRap string  `json:"autoSalesProfitsRap"` // After Sales Profit
	AutoSalesMargin     float64 `json:"autoSalesMargin"`     // After Sales Margin %
	Ranking             int     `json:"ranking"`             // Foton Ranking
	AutoSalesDecision   string  `json:"autoSalesDecision"`   // Performance Decision

	// After Sales fields
	RStockPercent   float64 `json:"rStockPercent"`   // Recommended Stock %
	WStockPercent   float64 `json:"wStockPercent"`   // Warranty Stock %
	FlhPercent      float64 `json:"flhPercent"`      // FLH %
	ServiceContract string  `json:"serviceContract"` // Service Contract level
	AsTrainings     bool    `json:"asTrainings"`     // After Sales Trainings
	Csi             string  `json:"csi"`             // CSI
	AsDecision      string  `json:"asDecision"`      // After Sales Decision
}

// GetAllData возвращает комплексные данные всех дилеров за период.
// Объединяет данные из таблиц: dealer_dev, sales, performance, after_sales.
// @Summary Get all dealer data
// @Description Получение всех данных дилеров за период (dealer dev + sales + performance + after sales)
// @Tags all-data
// @Accept json
// @Produce json
// @Param region query string false "Region filter" default("all-russia")
// @Param quarter query string false "Quarter" default("q1")
// @Param year query int false "Year" default(2024)
// @Success 200 {array} AllDealerData
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/all-data [get]
func (s *Server) GetAllData(c echo.Context) error {
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

	ctx := c.Request().Context()

	// Получаем данные из всех сервисов
	ddList, err := s.dealerDevService.GetDealerDevByPeriod(ctx, quarter, year, region)
	if err != nil {
		s.logger.Error("GetAllData: failed to get dealer dev data",
			"error", err,
		)
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to get dealer development data",
		})
	}

	salesList, err := s.salesService.GetSalesByPeriod(ctx, quarter, year, region)
	if err != nil {
		s.logger.Error("GetAllData: failed to get sales data",
			"error", err,
		)
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to get sales data",
		})
	}

	perfList, err := s.perfService.GetPerformanceByPeriod(ctx, quarter, year, region)
	if err != nil {
		s.logger.Error("GetAllData: failed to get performance data",
			"error", err,
		)
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to get performance data",
		})
	}

	asList, err := s.afterSalesService.GetAfterSalesByPeriod(ctx, quarter, year, region)
	if err != nil {
		s.logger.Error("GetAllData: failed to get after sales data",
			"error", err,
		)
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to get after sales data",
		})
	}

	// Создаем map для быстрого поиска по dealer_id
	ddMap := make(map[int64]*model.DealerDevWithDetails)
	salesMap := make(map[int64]*model.SalesWithDetails)
	perfMap := make(map[int64]*model.PerformanceWithDetails)
	asMap := make(map[int64]*model.AfterSalesWithDetails)

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

	// Получаем уникальный список dealer_id
	dealerIDs := make(map[int64]bool)
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

	// Собираем комплексные данные
	response := make([]AllDealerData, 0, len(dealerIDs))

	for dealerID := range dealerIDs {
		allData := AllDealerData{
			ID: strconv.FormatInt(dealerID, 10),
		}

		// Dealer Development data
		if dd, ok := ddMap[dealerID]; ok {
			allData.Name = dd.DealerName
			allData.City = dd.City
			allData.SalesManager = dd.Manager
			allData.Class = string(dd.DealerShipClass)
			allData.Checklist = int(dd.CheckListScore)
			allData.BrandsInPortfolio = dd.Brands
			allData.DealerDevRecommendation = string(dd.Recommendation)
		}

		// Sales data
		if sales, ok := salesMap[dealerID]; ok {
			if allData.Name == "" {
				allData.Name = sales.DealerName
				allData.City = sales.City
				allData.SalesManager = sales.Manager
			}
			allData.SalesTarget = sales.SalesTarget
			allData.StockHdtMdtLdt = sales.StockHdtMdtLdt
			allData.BuyoutHdtMdtLdt = sales.BuyoutHdtMdtLdt
			allData.FotonSalesmen = int(sales.FotonSalesmen)
			allData.SalesTrainings = sales.SalesTrainings
			allData.SalesDecision = string(sales.SalesDecision)
		}

		// Performance data
		if perf, ok := perfMap[dealerID]; ok {
			if allData.Name == "" {
				allData.Name = perf.DealerName
				allData.City = perf.City
				allData.SalesManager = perf.Manager
			}
			allData.SrRub = perf.SalesRevenueFormatted
			allData.SalesProfit = perf.SalesProfitFormatted
			allData.SalesMargin = perf.SalesMarginPercent
			allData.AutoSalesRevenue = perf.AfterSalesRevenueFormatted
			allData.AutoSalesProfitsRap = perf.AfterSalesProfitFormatted
			allData.AutoSalesMargin = perf.AfterSalesMarginPercent
			allData.Ranking = int(perf.FotonRank)
			allData.AutoSalesDecision = string(perf.PerformanceDecision)
		}

		// After Sales data
		if as, ok := asMap[dealerID]; ok {
			if allData.Name == "" {
				allData.Name = as.DealerName
				allData.City = as.City
				allData.SalesManager = as.Manager
			}
			allData.RStockPercent = float64(as.RecommendedStock)
			allData.WStockPercent = float64(as.WarrantyStock)
			allData.FlhPercent = float64(as.FotonLaborHours)
			allData.ServiceContract = convertServiceContractsToLevel(as.ServiceContracts)
			allData.AsTrainings = as.ASTrainings
			allData.Csi = as.CSI
			allData.AsDecision = string(as.AfterSalesDecision)
		}

		response = append(response, allData)
	}

	s.logger.Info("GetAllData: successfully retrieved all data",
		"region", region,
		"quarter", quarter,
		"year", year,
		"count", len(response),
	)

	return c.JSON(http.StatusOK, response)
}

// convertServiceContractsToLevel преобразует количество контрактов в уровень (Gold, Silver, Bronze).
func convertServiceContractsToLevel(contracts int16) string {
	switch {
	case contracts >= 50:
		return "Gold"
	case contracts >= 20:
		return "Silver"
	case contracts >= 5:
		return "Bronze"
	default:
		return "None"
	}
}
