package delivery

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/typefunco/dealer_dev_platform/internal/utils"
)

// AfterSalesDealerResponse представляет дилера с данными After Sales для API.
type AfterSalesDealerResponse struct {
	ID                    string  `json:"id"`
	Name                  string  `json:"name"`
	City                  string  `json:"city"`
	RStockPercent         *int    `json:"rStockPercent"`         // Recommended Stock
	WStockPercent         *int    `json:"wStockPercent"`         // Warranty Stock
	FlhPercent            *int    `json:"flhPercent"`            // Foton Labor Hours
	FlhSharePercent       *string `json:"flhSharePercent"`       // Foton Labour Hours Share
	WarrantyHours         *int    `json:"warrantyHours"`         // Foton Warranty Hours
	ServiceContractsHours *int    `json:"serviceContractsHours"` // Service Contracts
	AsTrainings           *bool   `json:"asTrainings"`           // AS Trainings status
	CSI                   *string `json:"csi"`                   // Customer Satisfaction Index
	AsDecision            *string `json:"asDecision"`            // AS Decision
	SparePartsSalesQ3     *string `json:"sparePartsSalesQ3"`     // Spare Parts Sales Q3
	SparePartsSalesYtd    *string `json:"sparePartsSalesYtd"`    // Spare Parts Sales YTD %
}

// GetAfterSalesData возвращает данные послепродажного обслуживания с поддержкой фильтров.
// @Summary Get After Sales data
// @Description Получение данных послепродажного обслуживания с фильтрацией по региону, году, кварталу и дилерам
// @Tags aftersales
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
// @Success 200 {array} AfterSalesDealerResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/aftersales [get]
func (s *Server) GetAfterSalesData(c echo.Context) error {
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
	afterSalesList, err := s.afterSalesService.GetAfterSalesWithFilters(c.Request().Context(), filters)
	if err != nil {
		s.logger.Error("GetAfterSalesData: failed to get after sales data",
			"filters", filters,
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

	return c.JSON(http.StatusOK, response)
}
