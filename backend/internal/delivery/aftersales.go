package delivery

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// AfterSalesDealerResponse представляет дилера с данными After Sales для API.
type AfterSalesDealerResponse struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	City            string `json:"city"`
	RStockPercent   int    `json:"rStockPercent"`   // Recommended Stock %
	WStockPercent   int    `json:"wStockPercent"`   // Warranty Stock %
	FlhPercent      int    `json:"flhPercent"`      // Foton Labor Hours %
	ServiceContract int    `json:"serviceContract"` // Service Contracts count
	AsTrainings     bool   `json:"asTrainings"`     // AS Trainings completed
	Csi             bool   `json:"csi"`             // CSI available
	AsDecision      string `json:"asDecision"`      // AS Decision
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
			ID:              strconv.FormatInt(as.DealerID, 10),
			Name:            as.DealerName,
			City:            as.City,
			RStockPercent:   int(as.RecommendedStock),
			WStockPercent:   int(as.WarrantyStock),
			FlhPercent:      int(as.FotonLaborHours),
			ServiceContract: int(as.ServiceContracts),
			AsTrainings:     as.ASTrainings,
			Csi:             as.CSI != "" && as.CSI != "0",
			AsDecision:      string(as.AfterSalesDecision),
		})
	}

	return c.JSON(http.StatusOK, response)
}
