package model

// QuarterMetrics представляет агрегированные метрики за квартал.
// Используется для страницы сравнения кварталов.
type QuarterMetrics struct {
	Quarter string `json:"quarter"` // q1, q2, q3, q4
	Year    string `json:"year"`    // 2024, 2025 и т.д.

	// Dealer Development метрики
	AverageClass      string         `json:"average_class"`
	AverageChecklist  float64        `json:"average_checklist"`
	ClassDistribution map[string]int `json:"class_distribution"` // {"A": 35, "B": 40, "C": 20, "D": 5}

	// Sales Team метрики
	AverageSalesTarget         string         `json:"average_sales_target"` // "45/100"
	AverageFotonSalesmen       float64        `json:"average_foton_salesmen"`
	SalesTrainingsPercentage   float64        `json:"sales_trainings_percentage"`
	SalesTrainingsDistribution map[string]int `json:"sales_trainings_distribution"` // {"Yes": 65, "No": 35}

	// Stocks and Buyout
	StocksData StockData `json:"stocks_data"` // {HDT: 50, MDT: 30, LDT: 20}
	BuyoutData StockData `json:"buyout_data"` // {HDT: 45, MDT: 25, LDT: 15}

	// Performance метрики
	AverageSalesRevenue string  `json:"average_sales_revenue"` // "5,200,000"
	AverageSalesProfit  float64 `json:"average_sales_profit"`
	AverageSalesMargin  float64 `json:"average_sales_margin"`
	AverageRanking      float64 `json:"average_ranking"`
	MarketingInvestment float64 `json:"marketing_investment"` // В млн рублей
	AutoSalesRevenue    float64 `json:"auto_sales_revenue"`   // В млн рублей
	AutoSalesProfit     float64 `json:"auto_sales_profit"`    // В млн рублей
	AutoSalesMargin     float64 `json:"auto_sales_margin"`

	// After Sales метрики
	AverageRStockPercent    float64        `json:"average_r_stock_percent"`
	AverageWStockPercent    float64        `json:"average_w_stock_percent"`
	AverageFlhPercent       float64        `json:"average_flh_percent"`
	AsTrainingsPercentage   float64        `json:"as_trainings_percentage"`
	AsTrainingsDistribution map[string]int `json:"as_trainings_distribution"` // {"Yes": 70, "No": 30}
	CsiPercentage           float64        `json:"csi_percentage"`
	FotonWarrantyHours      int            `json:"foton_warranty_hours"`

	// Decisions
	DecisionDistribution map[string]int `json:"decision_distribution"` // {"Planned Result": 30, "Needs Development": 45, ...}
}

// StockData представляет данные о запасах по категориям грузовиков.
type StockData struct {
	HDT int `json:"HDT"` // Heavy Duty Truck
	MDT int `json:"MDT"` // Medium Duty Truck
	LDT int `json:"LDT"` // Light Duty Truck
}
