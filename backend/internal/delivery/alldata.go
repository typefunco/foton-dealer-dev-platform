package delivery

import (
	"fmt"
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
	DealerID      int       `json:"dealer_id"`
	Ruft          string    `json:"ruft"`
	DealerNameRu  string    `json:"dealer_name_ru"`
	DealerNameEn  string    `json:"dealer_name_en"`
	City          string    `json:"city"`
	Region        string    `json:"region"`
	Manager       string    `json:"manager"`
	JointDecision *string   `json:"joint_decision"`
	Period        time.Time `json:"period"`

	// Dealer Development данные
	CheckListScore       *float64         `json:"check_list_score"`
	DealershipClass      *DealershipClass `json:"dealership_class"`
	Brands               []string         `json:"brands"`
	Branding             *BrandingStatus  `json:"branding"`
	MarketingInvestments *float64         `json:"marketing_investments"`
	BySideBusinesses     *string          `json:"by_side_businesses"`
	DDRecommendation     *string          `json:"dd_recommendation"`

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
	RecommendedStockPct   *float64           `json:"recommended_stock_pct"`
	WarrantyStockPct      *float64           `json:"warranty_stock_pct"`
	FotonLaborHoursPct    *float64           `json:"foton_labor_hours_pct"`
	WarrantyHours         *float64           `json:"warranty_hours"`
	ServiceContractsHours *float64           `json:"service_contracts_hours"`
	ASTrainings           *ASTrainingsStatus `json:"as_trainings"`
	SparePartsSalesQ      *float64           `json:"spare_parts_sales_q"`
	SparePartsSalesYtdPct *float64           `json:"spare_parts_sales_ytd_pct"`
	ASRecommendation      *string            `json:"as_recommendation"`

	// Performance Sales данные
	QuantitySold   *int     `json:"quantity_sold"`
	SalesRevenue   *float64 `json:"sales_revenue"`
	SalesMargin    *float64 `json:"sales_margin"`
	SalesMarginPct *float64 `json:"sales_margin_pct"`
	SalesProfitPct *float64 `json:"sales_profit_pct"`

	// Performance AfterSales данные
	ASRevenue   *float64 `json:"as_revenue"`
	ASMargin    *float64 `json:"as_margin"`
	ASMarginPct *float64 `json:"as_margin_pct"`
	ASProfitPct *float64 `json:"as_profit_pct"`
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

	// Преобразуем quarter/year в period
	period, err := parseQuarterToPeriod(quarter, year)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid quarter/year format",
		})
	}

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

	// Получаем данные из новых сервисов производительности
	perfSalesList, err := s.perfSalesService.GetAllByPeriod(ctx, period)
	if err != nil {
		s.logger.Error("GetAllData: failed to get performance sales data",
			"error", err,
		)
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to get performance sales data",
		})
	}

	perfASList, err := s.perfASService.GetAllByPeriod(ctx, period)
	if err != nil {
		s.logger.Error("GetAllData: failed to get performance aftersales data",
			"error", err,
		)
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to get performance aftersales data",
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
	perfSalesMap := make(map[int]*model.PerformanceSales)
	perfASMap := make(map[int]*model.PerformanceAfterSales)
	asMap := make(map[int]*model.AfterSalesWithDetails)

	for _, dd := range ddList {
		ddMap[dd.DealerID] = dd
	}

	for _, sales := range salesList {
		salesMap[sales.DealerID] = sales
	}

	for _, perf := range perfSalesList {
		perfSalesMap[perf.DealerID] = perf
	}

	for _, perf := range perfASList {
		perfASMap[perf.DealerID] = perf
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
	for id := range perfSalesMap {
		dealerIDs[id] = true
	}
	for id := range perfASMap {
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
			Period:   period,
		}

		// Dealer Development data
		if dd, ok := ddMap[dealerID]; ok {
			allData.DealerNameRu = dd.DealerNameRu
			allData.DealerNameEn = dd.DealerNameEn
			allData.City = dd.City
			allData.Region = dd.Region
			allData.Manager = dd.Manager
			allData.CheckListScore = dd.CheckListScore
			allData.DealershipClass = dd.DealershipClass
			allData.Brands = dd.Brands
			allData.Branding = dd.Branding
			allData.MarketingInvestments = dd.MarketingInvestments
			allData.BySideBusinesses = dd.BySideBusinesses
			allData.DDRecommendation = dd.DDRecommendation
		}

		// Sales data
		if sales, ok := salesMap[dealerID]; ok {
			if allData.DealerNameRu == "" {
				allData.DealerNameRu = sales.DealerNameRu
				allData.DealerNameEn = sales.DealerNameEn
				allData.City = sales.City
				allData.Region = sales.Region
				allData.Manager = sales.Manager
			}
			allData.StockHDT = sales.StockHDT
			allData.StockMDT = sales.StockMDT
			allData.StockLDT = sales.StockLDT
			allData.BuyoutHDT = sales.BuyoutHDT
			allData.BuyoutMDT = sales.BuyoutMDT
			allData.BuyoutLDT = sales.BuyoutLDT
			allData.FotonSalesPersonnel = sales.FotonSalesPersonnel
			allData.SalesTargetPlan = sales.SalesTargetPlan
			allData.SalesTargetFact = sales.SalesTargetFact
			allData.ServiceContractsSales = sales.ServiceContractsSales
			allData.SalesTrainings = sales.SalesTrainings
			allData.SalesRecommendation = sales.SalesRecommendation
		}

		// Performance Sales data
		if perf, ok := perfSalesMap[dealerID]; ok {
			allData.QuantitySold = perf.QuantitySold
			allData.SalesRevenue = perf.SalesRevenue
			allData.SalesMargin = perf.SalesMargin
			allData.SalesMarginPct = perf.SalesMarginPct
			allData.SalesProfitPct = perf.SalesProfitPct
		}

		// Performance AfterSales data
		if perf, ok := perfASMap[dealerID]; ok {
			allData.ASRevenue = perf.ASRevenue
			allData.ASMargin = perf.ASMargin
			allData.ASMarginPct = perf.ASMarginPct
			allData.ASProfitPct = perf.ASProfitPct
		}

		// After Sales data
		if as, ok := asMap[dealerID]; ok {
			if allData.DealerNameRu == "" {
				allData.DealerNameRu = as.DealerNameRu
				allData.DealerNameEn = as.DealerNameEn
				allData.City = as.City
				allData.Region = as.Region
				allData.Manager = as.Manager
			}
			allData.RecommendedStockPct = as.RecommendedStockPct
			allData.WarrantyStockPct = as.WarrantyStockPct
			allData.FotonLaborHoursPct = as.FotonLaborHoursPct
			allData.WarrantyHours = as.WarrantyHours
			allData.ServiceContractsHours = as.ServiceContractsHours
			allData.ASTrainings = as.ASTrainings
			allData.SparePartsSalesQ = as.SparePartsSalesQ
			allData.SparePartsSalesYtdPct = as.SparePartsSalesYtdPct
			allData.ASRecommendation = as.ASRecommendation
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

// parseQuarterToPeriod преобразует quarter и year в time.Time.
func parseQuarterToPeriod(quarter string, year int) (time.Time, error) {
	var month int
	switch quarter {
	case "q1":
		month = 1 // Январь
	case "q2":
		month = 4 // Апрель
	case "q3":
		month = 7 // Июль
	case "q4":
		month = 10 // Октябрь
	default:
		return time.Time{}, fmt.Errorf("invalid quarter: %s", quarter)
	}

	return time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC), nil
}
