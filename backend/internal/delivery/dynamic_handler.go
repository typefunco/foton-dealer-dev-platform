package delivery

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/typefunco/dealer_dev_platform/internal/model"
	"github.com/typefunco/dealer_dev_platform/internal/utils"
)

// DynamicDataResponse универсальный ответ для динамических данных
type DynamicDataResponse struct {
	TableType string      `json:"tableType"`
	Year      int         `json:"year"`
	Quarter   string      `json:"quarter"`
	Regions   []string    `json:"regions"`
	DealerIDs []int       `json:"dealer_ids"`
	Data      interface{} `json:"data"`
	Count     int         `json:"count"`
}

// GetDynamicData универсальный хендлер для получения данных из динамических таблиц
// @Summary Get table data
// @Description Получение данных из таблиц с поддержкой фильтрации по году, кварталу и региону
// @Tags tables
// @Accept json
// @Produce json
// @Param year query int false "Year filter" default(2024)
// @Param quarter query string false "Quarter filter (Q1, Q2, Q3, Q4)" default(Q1)
// @Param regions query string false "Comma-separated regions filter" default(all-russia)
// @Param region query string false "Single region filter (for backward compatibility)" default(all-russia)
// @Param dealer_ids query string false "Comma-separated dealer IDs"
// @Param limit query int false "Limit for pagination"
// @Param offset query int false "Offset for pagination"
// @Param sort_by query string false "Sort field"
// @Param sort_order query string false "Sort order (asc, desc)"
// @Success 200 {object} DynamicDataResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/{table_type} [get]
func (s *Server) GetDynamicData(c echo.Context) error {
	// Извлекаем тип таблицы из URL пути
	path := c.Request().URL.Path
	tableType := s.getTableTypeFromPath(path)

	// Валидируем тип таблицы
	if !isValidTableType(tableType) {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid table type. Supported types: dealer_dev, sales, after_sales, performance, sales_team",
		})
	}

	// Парсим динамические параметры
	params, err := model.ParseFromContext(c)
	if err != nil {
		s.logger.Error("GetDynamicData: failed to parse parameters", "error", err)
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid parameters: " + err.Error(),
		})
	}

	// Парсим дополнительные фильтры
	filters := utils.ParseFilterParamsFromContext(c)

	// Добавляем динамические параметры в фильтры
	filters.Year = params.Year
	filters.Quarter = params.Quarter
	filters.DealerIDs = params.DealerIDs

	// Обрабатываем множественные регионы
	if len(params.Regions) > 0 {
		// Если регионов несколько, используем первый для обратной совместимости
		// В будущем можно расширить сервисы для поддержки множественных регионов
		filters.Region = params.Regions[0]
	}

	// Получаем данные в зависимости от типа таблицы
	var data interface{}
	var count int

	switch tableType {
	case model.TableTypeDealerDev:
		data, err = s.getDealerDevData(c, filters)
	case model.TableTypeSales:
		data, err = s.getSalesData(c, filters)
	case model.TableTypeAfterSales:
		data, err = s.getAfterSalesData(c, filters)
	case model.TableTypePerformance:
		data, err = s.getPerformanceData(c, filters)
	case model.TableTypeSalesTeam:
		data, err = s.getSalesTeamData(c, filters)
	default:
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Unsupported table type",
		})
	}

	if err != nil {
		s.logger.Error("GetDynamicData: failed to get data",
			"table_type", tableType,
			"year", params.Year,
			"quarter", params.Quarter,
			"regions", params.Regions,
			"dealer_ids", params.DealerIDs,
			"error", err,
		)
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to get data from " + string(tableType),
		})
	}

	// Подсчитываем количество элементов
	if slice, ok := data.([]interface{}); ok {
		count = len(slice)
	} else {
		count = 1
	}

	response := DynamicDataResponse{
		TableType: string(tableType),
		Year:      params.Year,
		Quarter:   params.Quarter,
		Regions:   params.Regions,
		DealerIDs: params.DealerIDs,
		Data:      data,
		Count:     count,
	}

	return c.JSON(http.StatusOK, response)
}

// getTableTypeFromPath определяет тип таблицы по URL пути
func (s *Server) getTableTypeFromPath(path string) model.TableType {
	switch path {
	case "/api/dealer_dev":
		return model.TableTypeDealerDev
	case "/api/sales":
		return model.TableTypeSales
	case "/api/after_sales":
		return model.TableTypeAfterSales
	case "/api/performance":
		return model.TableTypePerformance
	case "/api/sales_team":
		return model.TableTypeSalesTeam
	default:
		return ""
	}
}

// isValidTableType проверяет валидность типа таблицы
func isValidTableType(tableType model.TableType) bool {
	validTypes := map[model.TableType]bool{
		model.TableTypeDealerDev:   true,
		model.TableTypeSales:       true,
		model.TableTypeAfterSales:  true,
		model.TableTypePerformance: true,
		model.TableTypeSalesTeam:   true,
	}
	return validTypes[tableType]
}

// getDealerDevData получает данные Dealer Development
func (s *Server) getDealerDevData(c echo.Context, filters *model.FilterParams) (interface{}, error) {
	region := filters.Region
	quarter := filters.Quarter
	year := filters.Year

	ddList, err := s.dealerDevService.GetDealerDevByPeriod(c.Request().Context(), quarter, year, region)
	if err != nil {
		return nil, err
	}

	// Преобразуем в API response
	response := make([]DealerDevResponse, 0, len(ddList))
	for _, dd := range ddList {
		response = append(response, DealerDevResponse{
			ID:                      strconv.FormatInt(int64(dd.DealerID), 10),
			Name:                    dd.DealerNameRu,
			City:                    dd.City,
			Class:                   dd.DealershipClass,
			Checklist:               dd.CheckListScore,
			BrandsInPortfolio:       []string{},
			BrandsCount:             0,
			Branding:                dd.Branding,
			BuySideBusiness:         []string{},
			DealerDevRecommendation: dd.DDRecommendation,
		})
	}

	return response, nil
}

// getSalesData получает данные Sales
func (s *Server) getSalesData(c echo.Context, filters *model.FilterParams) (interface{}, error) {
	// Используем существующий сервис sales
	salesList, err := s.salesService.GetSalesByPeriod(c.Request().Context(), filters.Quarter, filters.Year, filters.Region)
	if err != nil {
		return nil, err
	}

	// Преобразуем в API response с полными данными для Sales таблицы
	response := make([]interface{}, 0, len(salesList))
	for _, sale := range salesList {
		// Формируем строки для stock и buyout в формате "hdt/mdt/ldt"
		stockHdtMdtLdt := fmt.Sprintf("%d/%d/%d", sale.StockHDT, sale.StockMDT, sale.StockLDT)
		buyoutHdtMdtLdt := fmt.Sprintf("%d/%d/%d", sale.BuyoutHDT, sale.BuyoutMDT, sale.BuyoutLDT)

		response = append(response, map[string]interface{}{
			"id":              strconv.FormatInt(int64(sale.DealerID), 10),
			"name":            sale.DealerNameRu,
			"city":            sale.City,
			"salesManager":    sale.Manager,
			"salesTarget":     sale.SalesTarget,
			"stockHdtMdtLdt":  stockHdtMdtLdt,
			"buyoutHdtMdtLdt": buyoutHdtMdtLdt,
			"fotonSalesmen":   sale.FotonSalesmen,
			"salesTrainings":  sale.SalesTrainings,
			"salesDecision":   sale.SalesDecision,
		})
	}

	return response, nil
}

// getAfterSalesData получает данные After Sales
func (s *Server) getAfterSalesData(c echo.Context, filters *model.FilterParams) (interface{}, error) {
	afterSalesList, err := s.afterSalesService.GetAfterSalesWithFilters(c.Request().Context(), filters)
	if err != nil {
		return nil, err
	}

	// Преобразуем в API response
	response := make([]AfterSalesDealerResponse, 0, len(afterSalesList))
	for _, as := range afterSalesList {
		response = append(response, AfterSalesDealerResponse{
			ID:                    strconv.FormatInt(int64(as.DealerID), 10),
			Name:                  as.DealerNameRu,
			City:                  as.City,
			RStockPercent:         &as.RecommendedStock,
			WStockPercent:         &as.WarrantyStock,
			FlhPercent:            &as.FotonLaborHours,
			WarrantyHours:         &as.FotonWarrantyHours,
			ServiceContractsHours: &as.ServiceContracts,
			AsTrainings:           &as.ASTrainings,
			CSI:                   as.CSI,
			AsDecision:            &as.ASDecision,
		})
	}

	return response, nil
}

// getPerformanceData получает данные Performance
func (s *Server) getPerformanceData(c echo.Context, filters *model.FilterParams) (interface{}, error) {
	perfList, err := s.perfService.GetPerformanceWithFilters(c.Request().Context(), filters)
	if err != nil {
		return nil, err
	}

	// Преобразуем в API response
	response := make([]PerformanceDealerResponse, 0, len(perfList))
	for _, perf := range perfList {
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

	return response, nil
}

// getSalesTeamData получает данные Sales Team
func (s *Server) getSalesTeamData(c echo.Context, filters *model.FilterParams) (interface{}, error) {
	// Используем существующий сервис sales для получения данных команды продаж
	salesList, err := s.salesService.GetSalesByPeriod(c.Request().Context(), filters.Quarter, filters.Year, filters.Region)
	if err != nil {
		return nil, err
	}

	// Преобразуем в API response для Sales Team
	response := make([]interface{}, 0, len(salesList))
	for _, sale := range salesList {
		response = append(response, map[string]interface{}{
			"id":              strconv.FormatInt(int64(sale.DealerID), 10),
			"name":            sale.DealerNameRu,
			"city":            sale.City,
			"salesManager":    sale.Manager,
			"salesTarget":     sale.SalesTarget,
			"stockHdtMdtLdt":  fmt.Sprintf("%d/%d/%d", sale.StockHDT, sale.StockMDT, sale.StockLDT),
			"buyoutHdtMdtLdt": fmt.Sprintf("%d/%d/%d", sale.BuyoutHDT, sale.BuyoutMDT, sale.BuyoutLDT),
			"fotonSalesmen":   sale.FotonSalesmen,
			"salesTrainings":  sale.SalesTrainings,
			"salesDecision":   sale.SalesDecision,
		})
	}

	return response, nil
}
