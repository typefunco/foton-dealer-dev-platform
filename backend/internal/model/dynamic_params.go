package model

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

// DynamicTableParams содержит параметры для динамического определения таблицы
type DynamicTableParams struct {
	Year      int      `json:"year"`       // Год (2024, 2025, etc.)
	Quarter   string   `json:"quarter"`    // Квартал (Q1, Q2, Q3, Q4)
	Regions   []string `json:"regions"`    // Список регионов (central, north-west, etc.)
	DealerIDs []int    `json:"dealer_ids"` // Список ID дилеров
}

// TableType определяет тип таблицы
type TableType string

const (
	TableTypeDealerDev   TableType = "dealer_dev"
	TableTypeSales       TableType = "sales"
	TableTypeAfterSales  TableType = "after_sales"
	TableTypePerformance TableType = "performance"
	TableTypeSalesTeam   TableType = "sales_team"
)

// GetTableName возвращает полное имя таблицы с учетом года и квартала
func (p *DynamicTableParams) GetTableName(tableType TableType) string {
	return string(tableType) + "_" + p.Quarter + "_" + string(rune(p.Year))
}

// Validate проверяет корректность параметров
func (p *DynamicTableParams) Validate() error {
	if p.Year < 2020 || p.Year > 2030 {
		return fmt.Errorf("invalid year: %d", p.Year)
	}

	validQuarters := map[string]bool{
		"Q1": true, "Q2": true, "Q3": true, "Q4": true,
	}
	if !validQuarters[p.Quarter] {
		return fmt.Errorf("invalid quarter: %s", p.Quarter)
	}

	// Валидируем регионы
	validRegions := map[string]bool{
		"all-russia": true,
		"Central":    true,
		"North West": true,
		"Volga":      true,
		"South":      true,
		"Kavkaz":     true,
		"Ural":       true,
		"Siberia":    true,
		"Far East":   true,
	}

	for _, region := range p.Regions {
		if !validRegions[region] {
			return fmt.Errorf("invalid region: %s", region)
		}
	}

	return nil
}

// ParseFromContext извлекает параметры из контекста Echo
func ParseFromContext(c echo.Context) (*DynamicTableParams, error) {
	yearStr := c.QueryParam("year")
	year := 2024 // значение по умолчанию
	if yearStr != "" {
		var err error
		year, err = strconv.Atoi(yearStr)
		if err != nil {
			return nil, fmt.Errorf("invalid year parameter: %s", yearStr)
		}
	}

	quarter := c.QueryParam("quarter")
	if quarter == "" {
		quarter = "Q1" // значение по умолчанию
	}

	// Парсим регионы (может быть один или несколько через запятую)
	regionsStr := c.QueryParam("regions")
	if regionsStr == "" {
		regionsStr = c.QueryParam("region") // обратная совместимость
	}
	if regionsStr == "" {
		regionsStr = "all-russia" // значение по умолчанию
	}

	regions := []string{}
	if regionsStr != "" {
		// Разделяем регионы по запятой и очищаем от пробелов
		regionParts := strings.Split(regionsStr, ",")
		for _, part := range regionParts {
			region := strings.TrimSpace(part)
			if region != "" {
				// Применяем маппинг регионов
				if mappedRegion := mapRegion(region); mappedRegion != "" {
					regions = append(regions, mappedRegion)
				} else {
					regions = append(regions, region)
				}
			}
		}
	}

	// Парсим ID дилеров (опционально)
	dealerIDsStr := c.QueryParam("dealer_ids")
	dealerIDs := []int{}
	if dealerIDsStr != "" {
		dealerParts := strings.Split(dealerIDsStr, ",")
		for _, part := range dealerParts {
			idStr := strings.TrimSpace(part)
			if idStr != "" {
				if id, err := strconv.Atoi(idStr); err == nil {
					dealerIDs = append(dealerIDs, id)
				}
			}
		}
	}

	params := &DynamicTableParams{
		Year:      year,
		Quarter:   quarter,
		Regions:   regions,
		DealerIDs: dealerIDs,
	}

	if err := params.Validate(); err != nil {
		return nil, err
	}

	return params, nil
}

// mapRegion применяет маппинг регионов для совместимости с фронтендом
func mapRegion(region string) string {
	regionMapping := map[string]string{
		"all":        "all-russia",
		"central":    "Central",
		"north-west": "North West",
		"volga":      "Volga",
		"south":      "South",
		"n-caucasus": "Kavkaz",
		"ural":       "Ural",
		"siberia":    "Siberia",
		"far-east":   "Far East",
	}

	if mapped, exists := regionMapping[region]; exists {
		return mapped
	}
	return ""
}
