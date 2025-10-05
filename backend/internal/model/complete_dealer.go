package model

// CompleteDealerData представляет полные данные о дилере для страницы AllTable.
// Это объединение всех данных о дилере из разных блоков.
type CompleteDealerData struct {
	// Базовая информация
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	City         string `json:"city"`
	Region       string `json:"region"`
	SalesManager string `json:"sales_manager"`

	// Dealer Development данные
	Class                   string   `json:"class"`                     // A, B, C, D
	Checklist               int      `json:"checklist"`                 // 0-100
	BrandsInPortfolio       []string `json:"brands_in_portfolio"`       // Массив брендов
	BrandsCount             int      `json:"brands_count"`              // Количество брендов
	Branding                bool     `json:"branding"`                  // Наличие брендинга
	BuySideBusiness         []string `json:"buy_side_business"`         // Побочный бизнес
	DealerDevRecommendation string   `json:"dealer_dev_recommendation"` // Рекомендация

	// Sales Team данные
	SalesTarget     string `json:"sales_target"`       // "40/100"
	StockHdtMdtLdt  string `json:"stock_hdt_mdt_ldt"`  // "5/2/3"
	BuyoutHdtMdtLdt string `json:"buyout_hdt_mdt_ldt"` // "5/2/3"
	FotonSalesmen   int    `json:"foton_salesmen"`     // Количество продавцов
	SalesTrainings  bool   `json:"sales_trainings"`    // Обучение
	SalesDecision   string `json:"sales_decision"`     // Решение

	// Performance данные
	SrRub               string  `json:"sr_rub"`                 // Выручка от продаж "5 555 555"
	SalesProfit         string  `json:"sales_profit"`           // Прибыль от продаж "5 000 000"
	SalesMargin         float64 `json:"sales_margin"`           // Маржа продаж в %
	AutoSalesRevenue    string  `json:"auto_sales_revenue"`     // Выручка от послепродажного обслуживания
	AutoSalesProfitsRap string  `json:"auto_sales_profits_rap"` // Прибыль от послепродажного обслуживания
	AutoSalesMargin     float64 `json:"auto_sales_margin"`      // Маржа послепродажного обслуживания в %
	MarketingInvestment float64 `json:"marketing_investment"`   // Инвестиции в маркетинг (M Rub)
	Ranking             int     `json:"ranking"`                // Рейтинг 1-10
	AutoSalesDecision   string  `json:"auto_sales_decision"`    // Решение по производительности

	// After Sales данные
	RStockPercent   int    `json:"r_stock_percent"`  // Рекомендованный запас в %
	WStockPercent   int    `json:"w_stock_percent"`  // Гарантийный запас в %
	FlhPercent      int    `json:"flh_percent"`      // Foton Labor Hours в %
	ServiceContract string `json:"service_contract"` // Сервисные контракты
	AsTrainings     bool   `json:"as_trainings"`     // Обучение
	Csi             string `json:"csi"`              // Customer Satisfaction Index
	AsDecision      string `json:"as_decision"`      // Решение по послепродажному обслуживанию
}

// DealerCardData представляет детальные данные о дилере для карточки дилера.
// Используется на странице /dealer/:id
type DealerCardData struct {
	// Базовая информация
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	City         string `json:"city"`
	Region       string `json:"region"`
	SalesManager string `json:"sales_manager"`

	// Dealer Development
	Class               string   `json:"class"`
	Checklist           int      `json:"checklist"`
	BrandsInPortfolio   []string `json:"brands_in_portfolio"`
	BrandsCount         int      `json:"brands_count"`
	Branding            bool     `json:"branding"`
	BuySideBusiness     []string `json:"buy_side_business"`
	MarketingInvestment int      `json:"marketing_investment"` // В рублях

	// Sales
	SalesTarget     string `json:"sales_target"`
	StockHdtMdtLdt  string `json:"stock_hdt_mdt_ldt"`
	BuyoutHdtMdtLdt string `json:"buyout_hdt_mdt_ldt"`
	FotonSalesmen   string `json:"foton_salesmen"`
	SalesTrainings  bool   `json:"sales_trainings"`

	// Performance
	SrRub                string `json:"sr_rub"`
	SalesProfit          string `json:"sales_profit"`
	SalesMargin          int    `json:"sales_margin"`
	AfterSalesRevenue    string `json:"after_sales_revenue"`
	AfterSalesProfitsRap string `json:"after_sales_profits_rap"`
	AfterSalesMargin     int    `json:"after_sales_margin"`
	Ranking              int    `json:"ranking"`

	// After Sales
	RecommendedStock int    `json:"recommended_stock"`
	WarrantyStock    int    `json:"warranty_stock"`
	FotonLaborHours  int    `json:"foton_labor_hours"`
	ServiceContract  bool   `json:"service_contract"`
	AsTrainings      bool   `json:"as_trainings"`
	Csi              string `json:"csi"`
}
