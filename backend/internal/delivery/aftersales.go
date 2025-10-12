package delivery

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// AfterSalesDealerResponse представляет дилера с данными After Sales для API.
type AfterSalesDealerResponse struct {
	ID                    string   `json:"id"`
	Name                  string   `json:"name"`
	City                  string   `json:"city"`
	RStockPercent         *float64 `json:"rStockPercent"`         // Recommended Stock %
	WStockPercent         *float64 `json:"wStockPercent"`         // Warranty Stock %
	FlhPercent            *float64 `json:"flhPercent"`            // Foton Labor Hours %
	WarrantyHours         *float64 `json:"warrantyHours"`         // Warranty Hours
	ServiceContractsHours *float64 `json:"serviceContractsHours"` // Service Contracts Hours
	AsTrainings           *string  `json:"asTrainings"`           // AS Trainings status
	SparePartsSalesQ      *float64 `json:"sparePartsSalesQ"`      // Spare Parts Sales Q
	SparePartsSalesYtdPct *float64 `json:"sparePartsSalesYtdPct"` // Spare Parts Sales YTD %
	AsDecision            *string  `json:"asDecision"`            // AS Decision
}

// GetAfterSalesData возвращает данные послепродажного обслуживания по региону и периоду.
// @Summary Get After Sales data
// @Description Получение данных послепродажного обслуживания с фильтрацией по региону
// @Tags aftersales
// @Accept json
// @Produce json
// @Param region query string false "Region filter" default("all-russia")
// @Param quarter query string false "Quarter" default("q1")
// @Param year query int false "Year" default(2024)
// @Success 200 {array} AfterSalesDealerResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/aftersales [get]
func (s *Server) GetAfterSalesData(c echo.Context) error {
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
	afterSalesList, err := s.afterSalesService.GetAfterSalesByPeriod(c.Request().Context(), quarter, year, region)
	if err != nil {
		s.logger.Error("GetAfterSalesData: failed to get after sales data",
			"region", region,
			"quarter", quarter,
			"year", year,
			"error", err,
		)
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to get after sales data",
		})
	}

	// Преобразование в API response
	response := make([]AfterSalesDealerResponse, 0, len(afterSalesList))
	for _, as := range afterSalesList {
		response = append(response, AfterSalesDealerResponse{
			ID:                    strconv.FormatInt(int64(as.DealerID), 10),
			Name:                  as.DealerNameRu,
			City:                  as.City,
			RStockPercent:         as.RecommendedStockPct,
			WStockPercent:         as.WarrantyStockPct,
			FlhPercent:            as.FotonLaborHoursPct,
			WarrantyHours:         as.WarrantyHours,
			ServiceContractsHours: as.ServiceContractsHours,
			AsTrainings:           (*string)(as.ASTrainings),
			SparePartsSalesQ:      as.SparePartsSalesQ,
			SparePartsSalesYtdPct: as.SparePartsSalesYtdPct,
			AsDecision:            as.ASRecommendation,
		})
	}

	return c.JSON(http.StatusOK, response)
}
