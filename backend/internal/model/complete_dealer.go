package model

import "time"

// CompleteDealerData представляет полные данные о дилере для страницы AllTable.
// Это объединение всех данных о дилере из разных блоков.
type CompleteDealerData struct {
	// Базовая информация
	DealerID      int       `json:"dealer_id"`
	Ruft          string    `json:"ruft"`
	DealerNameRu  string    `json:"dealer_name_ru"`
	DealerNameEn  string    `json:"dealer_name_en"`
	City          string    `json:"city"`
	Region        string    `json:"region"`
	Manager       string    `json:"manager"`
	JointDecision *string   `json:"joint_decision"`
	Period        time.Time `json:"period"`

	// Dealer Development данные
	CheckListScore       *float64         `json:"check_list_score"`      // Check List Score % (0-100)
	DealershipClass      *DealershipClass `json:"dealership_class"`      // A, B, C, D
	Brands               []string         `json:"brands"`                // Массив брендов ["Foton", "Kamaz", "Sitrak"]
	Branding             *BrandingStatus  `json:"branding"`              // Y, N, Yes, No
	MarketingInvestments *float64         `json:"marketing_investments"` // Marketing Investments Rub
	BySideBusinesses     *string          `json:"by_side_businesses"`    // By-side businesses описание
	DDRecommendation     *string          `json:"dd_recommendation"`     // Recommendation (из Excel)

	// Sales данные
	StockHDT              *int                  `json:"stock_hdt"`               // Heavy Duty Truck
	StockMDT              *int                  `json:"stock_mdt"`               // Medium Duty Truck
	StockLDT              *int                  `json:"stock_ldt"`               // Light Duty Truck
	BuyoutHDT             *int                  `json:"buyout_hdt"`              // Heavy Duty Truck
	BuyoutMDT             *int                  `json:"buyout_mdt"`              // Medium Duty Truck
	BuyoutLDT             *int                  `json:"buyout_ldt"`              // Light Duty Truck
	FotonSalesPersonnel   *int                  `json:"foton_sales_personnel"`   // Количество продавцов Foton
	SalesTargetPlan       *int                  `json:"sales_target_plan"`       // План продаж
	SalesTargetFact       *int                  `json:"sales_target_fact"`       // Факт продаж
	ServiceContractsSales *float64              `json:"service_contracts_sales"` // Service contracts sales
	SalesTrainings        *SalesTrainingsStatus `json:"sales_trainings"`         // Y, N, Yes, No
	SalesRecommendation   *string               `json:"sales_recommendation"`    // Recommendation (из Excel)

	// AfterSales данные
	RecommendedStockPct   *float64           `json:"recommended_stock_pct"`     // В процентах
	WarrantyStockPct      *float64           `json:"warranty_stock_pct"`        // В процентах
	FotonLaborHoursPct    *float64           `json:"foton_labor_hours_pct"`     // В процентах
	WarrantyHours         *float64           `json:"warranty_hours"`            // Гарантийные часы
	ServiceContractsHours *float64           `json:"service_contracts_hours"`   // Часы сервисных контрактов
	ASTrainings           *ASTrainingsStatus `json:"as_trainings"`              // Y, N, Yes, No
	SparePartsSalesQ      *float64           `json:"spare_parts_sales_q"`       // За квартал
	SparePartsSalesYtdPct *float64           `json:"spare_parts_sales_ytd_pct"` // YTD динамика %
	ASRecommendation      *string            `json:"as_recommendation"`         // Recommendation (из Excel)

	// Performance Sales данные
	QuantitySold   *int     `json:"quantity_sold"`    // Количество проданных автомобилей
	SalesRevenue   *float64 `json:"sales_revenue"`    // Выручка (с НДС)
	SalesMargin    *float64 `json:"sales_margin"`     // Валовая прибыль (в рублях)
	SalesMarginPct *float64 `json:"sales_margin_pct"` // Маржа % = (margin / revenue) * 100
	SalesProfitPct *float64 `json:"sales_profit_pct"` // Рентабельность %

	// Performance AfterSales данные
	ASRevenue   *float64 `json:"as_revenue"`    // Выручка (с НДС)
	ASMargin    *float64 `json:"as_margin"`     // Валовая прибыль (в рублях)
	ASMarginPct *float64 `json:"as_margin_pct"` // Маржа % = (margin / revenue) * 100
	ASProfitPct *float64 `json:"as_profit_pct"` // Рентабельность %
}

// DealerCardData представляет детальные данные о дилере для карточки дилера.
// Используется на странице /dealer/:id
type DealerCardData struct {
	// Базовая информация
	DealerID      int       `json:"dealer_id"`
	Ruft          string    `json:"ruft"`
	DealerNameRu  string    `json:"dealer_name_ru"`
	DealerNameEn  string    `json:"dealer_name_en"`
	City          string    `json:"city"`
	Region        string    `json:"region"`
	Manager       string    `json:"manager"`
	JointDecision *string   `json:"joint_decision"`
	Period        time.Time `json:"period"`

	// Dealer Development
	CheckListScore       *float64         `json:"check_list_score"`
	DealershipClass      *DealershipClass `json:"dealership_class"`
	Brands               []string         `json:"brands"`
	Branding             *BrandingStatus  `json:"branding"`
	MarketingInvestments *float64         `json:"marketing_investments"`
	BySideBusinesses     *string          `json:"by_side_businesses"`
	DDRecommendation     *string          `json:"dd_recommendation"`

	// Sales
	StockHDT              *int                  `json:"stock_hdt"`
	StockMDT              *int                  `json:"stock_mdt"`
	StockLDT              *int                  `json:"stock_ldt"`
	BuyoutHDT             *int                  `json:"buyout_hdt"`
	BuyoutMDT             *int                  `json:"buyout_mdt"`
	BuyoutLDT             *int                  `json:"buyout_ldt"`
	FotonSalesPersonnel   *int                  `json:"foton_sales_personnel"`
	SalesTargetPlan       *int                  `json:"sales_target_plan"`
	SalesTargetFact       *int                  `json:"sales_target_fact"`
	ServiceContractsSales *float64              `json:"service_contracts_sales"`
	SalesTrainings        *SalesTrainingsStatus `json:"sales_trainings"`
	SalesRecommendation   *string               `json:"sales_recommendation"`

	// Performance Sales
	QuantitySold   *int     `json:"quantity_sold"`
	SalesRevenue   *float64 `json:"sales_revenue"`
	SalesMargin    *float64 `json:"sales_margin"`
	SalesMarginPct *float64 `json:"sales_margin_pct"`
	SalesProfitPct *float64 `json:"sales_profit_pct"`

	// Performance AfterSales
	ASRevenue   *float64 `json:"as_revenue"`
	ASMargin    *float64 `json:"as_margin"`
	ASMarginPct *float64 `json:"as_margin_pct"`
	ASProfitPct *float64 `json:"as_profit_pct"`

	// AfterSales
	RecommendedStockPct   *float64           `json:"recommended_stock_pct"`
	WarrantyStockPct      *float64           `json:"warranty_stock_pct"`
	FotonLaborHoursPct    *float64           `json:"foton_labor_hours_pct"`
	WarrantyHours         *float64           `json:"warranty_hours"`
	ServiceContractsHours *float64           `json:"service_contracts_hours"`
	ASTrainings           *ASTrainingsStatus `json:"as_trainings"`
	SparePartsSalesQ      *float64           `json:"spare_parts_sales_q"`
	SparePartsSalesYtdPct *float64           `json:"spare_parts_sales_ytd_pct"`
	ASRecommendation      *string            `json:"as_recommendation"`
}
