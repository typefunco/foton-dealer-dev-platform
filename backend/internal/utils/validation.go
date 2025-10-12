package utils

import "fmt"

// Validation utilities for the dealer development platform

// IsValidQuarter проверяет валидность квартала.
// Поддерживает как верхний, так и нижний регистр.
func IsValidQuarter(quarter string) bool {
	validQuarters := map[string]bool{
		"Q1": true,
		"Q2": true,
		"Q3": true,
		"Q4": true,
		"q1": true,
		"q2": true,
		"q3": true,
		"q4": true,
	}
	return validQuarters[quarter]
}

// IsValidYear проверяет валидность года.
func IsValidYear(year int) bool {
	return year >= 2020 && year <= 2030
}

// IsValidRegion проверяет валидность региона.
func IsValidRegion(region string) bool {
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
	return validRegions[region]
}

// NormalizeQuarter нормализует квартал к верхнему регистру.
func NormalizeQuarter(quarter string) string {
	switch quarter {
	case "q1":
		return "Q1"
	case "q2":
		return "Q2"
	case "q3":
		return "Q3"
	case "q4":
		return "Q4"
	default:
		return quarter
	}
}

// ValidateFilters валидирует параметры фильтрации.
func ValidateFilters(quarter string, year int, region string) error {
	if !IsValidQuarter(quarter) {
		return fmt.Errorf("invalid quarter: %s", quarter)
	}

	if !IsValidYear(year) {
		return fmt.Errorf("invalid year: %d", year)
	}

	if !IsValidRegion(region) {
		return fmt.Errorf("invalid region: %s", region)
	}

	return nil
}
