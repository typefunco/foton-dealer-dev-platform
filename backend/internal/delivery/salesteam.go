package delivery

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// SalesTeamDealerResponse представляет дилера с данными Sales Team для API.
type SalesTeamDealerResponse struct {
	ID                    string   `json:"id"`
	Name                  string   `json:"name"`
	City                  string   `json:"city"`
	SalesManager          string   `json:"salesManager"`
	SalesTargetPlan       *int     `json:"salesTargetPlan"`       // План продаж
	SalesTargetFact       *int     `json:"salesTargetFact"`       // Факт продаж
	StockHDT              *int     `json:"stockHDT"`              // Stock HDT
	StockMDT              *int     `json:"stockMDT"`              // Stock MDT
	StockLDT              *int     `json:"stockLDT"`              // Stock LDT
	BuyoutHDT             *int     `json:"buyoutHDT"`             // Buyout HDT
	BuyoutMDT             *int     `json:"buyoutMDT"`             // Buyout MDT
	BuyoutLDT             *int     `json:"buyoutLDT"`             // Buyout LDT
	FotonSalesPersonnel   *int     `json:"fotonSalesPersonnel"`   // Количество продавцов
	ServiceContractsSales *float64 `json:"serviceContractsSales"` // Service Contracts Sales
	SalesTrainings        *string  `json:"salesTrainings"`        // Статус тренингов
	SalesDecision         *string  `json:"salesDecision"`         // Решение по продажам
}

// GetSalesTeamData возвращает данные команды продаж по региону и периоду.
// @Summary Get Sales Team data
// @Description Получение данных команды продаж с фильтрацией по региону
// @Tags sales
// @Accept json
// @Produce json
// @Param region query string false "Region filter" default("all-russia")
// @Param quarter query string false "Quarter" default("q1")
// @Param year query int false "Year" default(2024)
// @Success 200 {array} SalesTeamDealerResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/sales [get]
func (s *Server) GetSalesTeamData(c echo.Context) error {
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

	// Получение данных из сервиса
	salesList, err := s.salesService.GetSalesByPeriod(c.Request().Context(), quarter, year, region)
	if err != nil {
		s.logger.Error("GetSalesTeamData: failed to get sales data",
			"region", region,
			"quarter", quarter,
			"year", year,
			"error", err,
		)
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to get sales team data",
		})
	}

	// Преобразование в API response
	response := make([]SalesTeamDealerResponse, 0, len(salesList))
	for _, sales := range salesList {
		response = append(response, SalesTeamDealerResponse{
			ID:                  strconv.FormatInt(int64(sales.DealerID), 10),
			Name:                sales.DealerNameRu,
			City:                sales.City,
			SalesManager:        sales.Manager,
			SalesTargetPlan:     nil, // SalesTarget теперь строка, а не int
			SalesTargetFact:     nil, // SalesTarget теперь строка, а не int
			StockHDT:            &sales.StockHDT,
			StockMDT:            &sales.StockMDT,
			StockLDT:            &sales.StockLDT,
			BuyoutHDT:           &sales.BuyoutHDT,
			BuyoutMDT:           &sales.BuyoutMDT,
			BuyoutLDT:           &sales.BuyoutLDT,
			FotonSalesPersonnel: &sales.FotonSalesmen,
			ServiceContractsSales: func() *float64 {
				f := float64(sales.ServiceContractsSales)
				return &f
			}(),
			SalesTrainings: func() *string {
				if sales.SalesTrainings {
					return &[]string{"Yes"}[0]
				}
				return &[]string{"No"}[0]
			}(),
			SalesDecision: &sales.SalesDecision,
		})
	}

	return c.JSON(http.StatusOK, response)
}
