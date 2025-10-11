package model

// QuarterMetrics представляет агрегированные метрики за квартал.
// Используется для сравнения кварталов.
type QuarterMetrics struct {
	Quarter string `json:"quarter"` // q1, q2, q3, q4
	Year    string `json:"year"`    // 2024, 2025

	// Dealer Development
	AverageClass      string             `json:"averageClass"`
	AverageChecklist  float64            `json:"averageChecklist"`
	ClassDistribution map[string]float64 `json:"classDistribution"` // {"A": 35, "B": 40, ...}

	// Sales Team
	AverageSalesTarget         string             `json:"averageSalesTarget"`
	AverageFotonSalesmen       float64            `json:"averageFotonSalesmen"`
	SalesTrainingsPercentage   float64            `json:"salesTrainingsPercentage"`
	SalesTrainingsDistribution map[string]float64 `json:"salesTrainingsDistribution"` // {"Yes": 65, "No": 35}

	// Stocks and Buyout
	StocksData StockDistribution `json:"stocksData"`
	BuyoutData StockDistribution `json:"buyoutData"`

	// Performance
	AverageSalesRevenue string  `json:"averageSalesRevenue"` // Форматированная строка
	AverageSalesProfit  float64 `json:"averageSalesProfit"`  // Процент
	AverageSalesMargin  float64 `json:"averageSalesMargin"`  // Процент
	AverageRanking      float64 `json:"averageRanking"`
	MarketingInvestment float64 `json:"marketingInvestment"` // M Rub
	AutoSalesRevenue    float64 `json:"autoSalesRevenue"`    // M Rub
	AutoSalesProfit     float64 `json:"autoSalesProfit"`     // M Rub
	AutoSalesMargin     float64 `json:"autoSalesMargin"`     // Процент

	// After Sales
	AverageRStockPercent    float64            `json:"averageRStockPercent"`
	AverageWStockPercent    float64            `json:"averageWStockPercent"`
	AverageFlhPercent       float64            `json:"averageFlhPercent"`
	AsTrainingsPercentage   float64            `json:"asTrainingsPercentage"`
	AsTrainingsDistribution map[string]float64 `json:"asTrainingsDistribution"` // {"Yes": 70, "No": 30}
	CsiPercentage           float64            `json:"csiPercentage"`
	FotonWarrantyHours      float64            `json:"fotonWarrantyHours"`

	// Decisions
	DecisionDistribution map[string]float64 `json:"decisionDistribution"` // {"Planned Result": 30, ...}
}

// StockDistribution представляет распределение по типам грузовиков.
type StockDistribution struct {
	HDT int `json:"HDT"` // Heavy Duty Truck
	MDT int `json:"MDT"` // Medium Duty Truck
	LDT int `json:"LDT"` // Light Duty Truck
}
