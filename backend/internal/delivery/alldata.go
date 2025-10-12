package delivery

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/typefunco/dealer_dev_platform/internal/model"
)

// Импортируем типы из модели для использования в delivery
type DealershipClass = model.DealershipClass
type BrandingStatus = model.BrandingStatus
type SalesTrainingsStatus = model.SalesTrainingsStatus
type ASTrainingsStatus = model.ASTrainingsStatus

// AllDealerData представляет комплексные данные дилера со всех таблиц.
// Использует новую структуру БД с нормализованными таблицами.
type AllDealerData struct {
	// Базовая информация
	DealerID     int       `json:"dealer_id"`
	DealerNameRu string    `json:"dealer_name_ru"`
	City         string    `json:"city"`
	Region       string    `json:"region"`
	Manager      string    `json:"manager"`
	Period       time.Time `json:"period"`

	// Dealer Development данные
	CheckListScore       *int    `json:"check_list_score"`
	DealershipClass      *string `json:"dealership_class"`
	Branding             *bool   `json:"branding"`
	MarketingInvestments *int64  `json:"marketing_investments"`
	DDRecommendation     *string `json:"dd_recommendation"`

	// Sales данные
	StockHDT              *int                  `json:"stock_hdt"`
	StockMDT              *int                  `json:"stock_mdt"`
	StockLDT              *int                  `json:"stock_ldt"`
	BuyoutHDT             *int                  `json:"buyout_hdt"`
	BuyoutMDT             *int                  `json:"buyout_mdt"`
	BuyoutLDT             *int                  `json:"buyout_ldt"`
	FotonSalesPersonnel   *int                  `json:"foton_sales_personnel"`
	SalesTargetPlan       *int                  `json:"sales_target_plan"`
	SalesTargetFact       *int                  `json:"sales_target_fact"`
	ServiceContractsSales *float64              `json:"service_contracts_sales"`
	SalesTrainings        *SalesTrainingsStatus `json:"sales_trainings"`
	SalesRecommendation   *string               `json:"sales_recommendation"`

	// AfterSales данные
	RecommendedStock   *int    `json:"recommended_stock"`
	WarrantyStock      *int    `json:"warranty_stock"`
	FotonLaborHours    *int    `json:"foton_labor_hours"`
	FotonWarrantyHours *int    `json:"foton_warranty_hours"`
	ServiceContracts   *int    `json:"service_contracts"`
	ASTrainings        *bool   `json:"as_trainings"`
	CSI                *string `json:"csi"`
	ASDecision         *string `json:"as_decision"`

	// Performance данные
	SalesRevenueRub      *float64 `json:"sales_revenue_rub"`
	SalesProfitRub       *float64 `json:"sales_profit_rub"`
	SalesMarginPercent   *float64 `json:"sales_margin_percent"`
	AfterSalesRevenueRub *float64 `json:"after_sales_revenue_rub"`
	AfterSalesProfitRub  *float64 `json:"after_sales_profit_rub"`
	AfterSalesMarginPct  *float64 `json:"after_sales_margin_pct"`
	MarketingInvestment  *float64 `json:"marketing_investment"`
	FotonRank            *int     `json:"foton_rank"`
	PerformanceDecision  *string  `json:"performance_decision"`
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

	// Получаем данные производительности
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

	// Получаем уникальный список dealer_id
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

	// Собираем комплексные данные
	response := make([]AllDealerData, 0, len(dealerIDs))

	for dealerID := range dealerIDs {
		allData := AllDealerData{
			DealerID: dealerID,
			Period:   time.Date(year, getMonthForQuarter(quarter), 1, 0, 0, 0, 0, time.UTC),
		}

		// Dealer Development data
		if dd, ok := ddMap[dealerID]; ok {
			allData.DealerNameRu = dd.DealerNameRu
			allData.City = dd.City
			allData.Region = dd.Region
			allData.Manager = dd.Manager
			allData.CheckListScore = &dd.CheckListScore
			allData.DealershipClass = &dd.DealershipClass
			allData.Branding = &dd.Branding
			allData.MarketingInvestments = &dd.MarketingInvestments
			allData.DDRecommendation = &dd.DDRecommendation
		}

		// Sales data
		if sales, ok := salesMap[dealerID]; ok {
			if allData.DealerNameRu == "" {
				allData.DealerNameRu = sales.DealerNameRu
				allData.City = sales.City
				allData.Region = sales.Region
				allData.Manager = sales.Manager
			}
			allData.StockHDT = &sales.StockHDT
			allData.StockMDT = &sales.StockMDT
			allData.StockLDT = &sales.StockLDT
			allData.BuyoutHDT = &sales.BuyoutHDT
			allData.BuyoutMDT = &sales.BuyoutMDT
			allData.BuyoutLDT = &sales.BuyoutLDT
			allData.FotonSalesPersonnel = &sales.FotonSalesmen
			// SalesTarget в модели Sales - это строка, а в AllData ожидается int
			// Пока оставляем nil, так как нужно преобразование
			allData.SalesTargetPlan = nil
			allData.SalesTargetFact = nil
			serviceContractsFloat := float64(sales.ServiceContractsSales)
			allData.ServiceContractsSales = &serviceContractsFloat
			var salesTrainingsStatus model.SalesTrainingsStatus
			if sales.SalesTrainings {
				salesTrainingsStatus = model.SalesTrainingsYes
			} else {
				salesTrainingsStatus = model.SalesTrainingsNo
			}
			allData.SalesTrainings = &salesTrainingsStatus
			allData.SalesRecommendation = &sales.SalesDecision
		}

		// Performance data
		if perf, ok := perfMap[dealerID]; ok {
			allData.SalesRevenueRub = &perf.SalesRevenueRub
			allData.SalesProfitRub = &perf.SalesProfitRub
			allData.SalesMarginPercent = &perf.SalesMarginPercent
			allData.AfterSalesRevenueRub = &perf.AfterSalesRevenueRub
			allData.AfterSalesProfitRub = &perf.AfterSalesProfitRub
			allData.AfterSalesMarginPct = &perf.AfterSalesMarginPercent
			allData.MarketingInvestment = &perf.MarketingInvestment
			allData.FotonRank = &perf.FotonRank
			allData.PerformanceDecision = &perf.PerformanceDecision
		}

		// After Sales data
		if as, ok := asMap[dealerID]; ok {
			if allData.DealerNameRu == "" {
				allData.DealerNameRu = as.DealerNameRu
				allData.City = as.City
				allData.Region = as.Region
				allData.Manager = as.Manager
			}
			allData.RecommendedStock = &as.RecommendedStock
			allData.WarrantyStock = &as.WarrantyStock
			allData.FotonLaborHours = &as.FotonLaborHours
			allData.FotonWarrantyHours = &as.FotonWarrantyHours
			allData.ServiceContracts = &as.ServiceContracts
			allData.ASTrainings = &as.ASTrainings
			allData.CSI = as.CSI
			allData.ASDecision = &as.ASDecision
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

func getMonthForQuarter(quarter string) time.Month {
	switch quarter {
	case "Q1":
		return time.January
	case "Q2":
		return time.April
	case "Q3":
		return time.July
	case "Q4":
		return time.October
	default:
		return time.January
	}
}
