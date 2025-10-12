package model

import (
	"fmt"
	"strings"
	"time"
)

// FilterParams представляет параметры фильтрации для API запросов.
// Используется для унификации фильтрации по всем эндпоинтам.
type FilterParams struct {
	// Временные фильтры
	Quarter string `json:"quarter" form:"quarter"` // Q1, Q2, Q3, Q4
	Year    int    `json:"year" form:"year"`       // 2024, 2025, etc.

	// Географические фильтры
	Region string `json:"region" form:"region"` // central, north-west, volga, etc.

	// Фильтры по дилерам
	DealerIDs []int `json:"dealer_ids" form:"dealer_ids"` // Список ID дилеров

	// Пагинация
	Limit  int `json:"limit" form:"limit"`   // Количество записей на странице
	Offset int `json:"offset" form:"offset"` // Смещение для пагинации

	// Сортировка
	SortBy    string `json:"sort_by" form:"sort_by"`       // Поле для сортировки
	SortOrder string `json:"sort_order" form:"sort_order"` // asc, desc
}

// Validate проверяет корректность параметров фильтрации.
func (f *FilterParams) Validate() error {
	// Валидация квартала
	if f.Quarter != "" {
		validQuarters := map[string]bool{
			"Q1": true, "Q2": true, "Q3": true, "Q4": true,
			"q1": true, "q2": true, "q3": true, "q4": true,
		}
		if !validQuarters[f.Quarter] {
			return fmt.Errorf("invalid quarter: %s", f.Quarter)
		}
		// Нормализуем квартал к верхнему регистру
		f.Quarter = strings.ToUpper(f.Quarter)
	}

	// Валидация года
	if f.Year != 0 && (f.Year < 2020 || f.Year > 2030) {
		return fmt.Errorf("invalid year: %d", f.Year)
	}

	// Валидация региона
	if f.Region != "" {
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
		if !validRegions[f.Region] {
			return fmt.Errorf("invalid region: %s", f.Region)
		}
	}

	// Валидация пагинации
	if f.Limit < 0 {
		return fmt.Errorf("limit cannot be negative")
	}
	if f.Offset < 0 {
		return fmt.Errorf("offset cannot be negative")
	}
	if f.Limit == 0 {
		f.Limit = 100 // Значение по умолчанию
	}

	// Валидация сортировки
	if f.SortOrder != "" && f.SortOrder != "asc" && f.SortOrder != "desc" {
		return fmt.Errorf("invalid sort_order: %s", f.SortOrder)
	}
	if f.SortOrder == "" {
		f.SortOrder = "asc" // Значение по умолчанию
	}

	return nil
}

// GetPeriodFromFilters возвращает time.Time на основе квартала и года.
func (f *FilterParams) GetPeriodFromFilters() (time.Time, error) {
	if f.Quarter == "" || f.Year == 0 {
		return time.Time{}, fmt.Errorf("quarter and year are required to generate period")
	}

	var month int
	switch f.Quarter {
	case "Q1":
		month = 1 // Январь
	case "Q2":
		month = 4 // Апрель
	case "Q3":
		month = 7 // Июль
	case "Q4":
		month = 10 // Октябрь
	default:
		return time.Time{}, fmt.Errorf("invalid quarter: %s", f.Quarter)
	}

	return time.Date(f.Year, time.Month(month), 1, 0, 0, 0, 0, time.UTC), nil
}

// RegionMapping маппинг регионов для API.
var RegionMapping = map[string]string{
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

// QuarterMapping маппинг кварталов для API.
var QuarterMapping = map[string]string{
	"q1": "Q1",
	"q2": "Q2",
	"q3": "Q3",
	"q4": "Q4",
}

// GetMappedRegion возвращает маппированное значение региона.
func (f *FilterParams) GetMappedRegion() string {
	if mapped, exists := RegionMapping[f.Region]; exists {
		return mapped
	}
	return f.Region
}

// GetMappedQuarter возвращает маппированное значение квартала.
func (f *FilterParams) GetMappedQuarter() string {
	if mapped, exists := QuarterMapping[f.Quarter]; exists {
		return mapped
	}
	return f.Quarter
}

// HasRegionFilter проверяет, есть ли фильтр по региону.
func (f *FilterParams) HasRegionFilter() bool {
	return f.Region != "" && f.Region != "allAll"
}

// HasPeriodFilter проверяет, есть ли фильтр по периоду.
func (f *FilterParams) HasPeriodFilter() bool {
	return f.Quarter != "" && f.Year != 0
}

// HasDealerFilter проверяет, есть ли фильтр по дилерам.
func (f *FilterParams) HasDealerFilter() bool {
	return len(f.DealerIDs) > 0
}
