package testutil

import (
	"fmt"
	"time"

	"github.com/typefunco/dealer_dev_platform/internal/model"
)

// CreateTestDealer создает тестового дилера
func CreateTestDealer() *model.Dealer {
	return CreateTestDealerWithRUFT("TEST001")
}

// CreateTestDealerWithRUFT создает тестового дилера с указанным RUFT
func CreateTestDealerWithRUFT(ruft string) *model.Dealer {
	jointDecision := "Planned Result"
	return &model.Dealer{
		Ruft:          ruft,
		DealerNameRu:  fmt.Sprintf("Тестовый Дилер %s", ruft),
		DealerNameEn:  fmt.Sprintf("Test Dealer %s", ruft),
		Region:        "Москва",
		City:          "Москва",
		Manager:       "Иван Иванов",
		JointDecision: &jointDecision,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
}

// CreateTestUser создает тестового пользователя
func CreateTestUser() *model.User {
	return &model.User{
		Login:     "testuser",
		Email:     "test@example.com",
		FirstName: "Тест",
		LastName:  "Пользователь",
		Role:      model.UserRoleAdmin,
		IsAdmin:   true,
		Region:    "Москва",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// CreateTestSales создает тестовые данные продаж
func CreateTestSales(dealerID int) *model.Sales {
	stockHDT := 10
	stockMDT := 15
	stockLDT := 20
	buyoutHDT := 5
	buyoutMDT := 8
	buyoutLDT := 12
	fotonPersonnel := 3
	targetPlan := 100
	targetFact := 85
	serviceContracts := 5.0
	trainings := model.SalesTrainingsYes
	recommendation := "Planned Result"

	return &model.Sales{
		DealerID:              dealerID,
		Period:                time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		StockHDT:              &stockHDT,
		StockMDT:              &stockMDT,
		StockLDT:              &stockLDT,
		BuyoutHDT:             &buyoutHDT,
		BuyoutMDT:             &buyoutMDT,
		BuyoutLDT:             &buyoutLDT,
		FotonSalesPersonnel:   &fotonPersonnel,
		SalesTargetPlan:       &targetPlan,
		SalesTargetFact:       &targetFact,
		ServiceContractsSales: &serviceContracts,
		SalesTrainings:        &trainings,
		SalesRecommendation:   &recommendation,
		CreatedAt:             time.Now(),
		UpdatedAt:             time.Now(),
	}
}

// // CreateTestAfterSales создает тестовые данные послепродажного обслуживания
// func CreateTestAfterSales(dealerID int) *model.AfterSales {
// 	recommendedStockPct := 80.0
// 	warrantyStockPct := 70.0
// 	fotonLaborHoursPct := 90.0
// 	warrantyHours := 120.0
// 	serviceContractsHours := 80.0
// 	trainings := model.ASTrainingsYes
// 	sparePartsSalesQ := 50.0
// 	sparePartsSalesYtdPct := 75.0
// 	recommendation := "Planned Result"

// 	return &model.AfterSales{
// 		DealerID:              dealerID,
// 		Period:                time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
// 		RecommendedStockPct:   &recommendedStockPct,
// 		WarrantyStockPct:      &warrantyStockPct,
// 		FotonLaborHoursPct:    &fotonLaborHoursPct,
// 		WarrantyHours:         &warrantyHours,
// 		ServiceContractsHours: &serviceContractsHours,
// 		ASTrainings:           &trainings,
// 		SparePartsSalesQ:      &sparePartsSalesQ,
// 		SparePartsSalesYtdPct: &sparePartsSalesYtdPct,
// 		ASRecommendation:      &recommendation,
// 		CreatedAt:             time.Now(),
// 		UpdatedAt:             time.Now(),
// 	}
// }

// // CreateTestDealerDevelopment создает тестовые данные развития дилера
// func CreateTestDealerDevelopment(dealerID int) *model.DealerDevelopment {
// 	checkListScore := 85.0
// 	dealershipClass := model.DealershipClassA
// 	brands := []string{"Foton", "Dongfeng"}
// 	branding := model.BrandingYes
// 	marketingInvestments := 100000.0
// 	bySideBusinesses := "Service, Parts"
// 	recommendation := "Planned Result"

// 	return &model.DealerDevelopment{
// 		DealerID:             dealerID,
// 		Period:               time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
// 		CheckListScore:       &checkListScore,
// 		DealershipClass:      &dealershipClass,
// 		Brands:               brands,
// 		Branding:             &branding,
// 		MarketingInvestments: &marketingInvestments,
// 		BySideBusinesses:     &bySideBusinesses,
// 		DDRecommendation:     &recommendation,
// 		CreatedAt:            time.Now(),
// 		UpdatedAt:            time.Now(),
// 	}
// }

// CreateTestPerformanceSales создает тестовые данные производительности продаж
func CreateTestPerformanceSales(dealerID int) *model.PerformanceSales {
	quantitySold := 25
	salesRevenue := 5000000.0
	salesRevenueNoVat := 4200000.0
	salesCost := 3500000.0
	salesMargin := 1500000.0
	salesMarginPct := 30.0
	salesProfitPct := 25.0

	return &model.PerformanceSales{
		DealerID:          dealerID,
		Period:            time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		QuantitySold:      &quantitySold,
		SalesRevenue:      &salesRevenue,
		SalesRevenueNoVat: &salesRevenueNoVat,
		SalesCost:         &salesCost,
		SalesMargin:       &salesMargin,
		SalesMarginPct:    &salesMarginPct,
		SalesProfitPct:    &salesProfitPct,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}
}

// CreateTestPerformanceAfterSales создает тестовые данные производительности послепродажного обслуживания
func CreateTestPerformanceAfterSales(dealerID int) *model.PerformanceAfterSales {
	asRevenue := 800000.0
	asRevenueNoVat := 670000.0
	asCost := 600000.0
	asMargin := 200000.0
	asMarginPct := 25.0
	asProfitPct := 30.0

	return &model.PerformanceAfterSales{
		DealerID:       dealerID,
		Period:         time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		ASRevenue:      &asRevenue,
		ASRevenueNoVat: &asRevenueNoVat,
		ASCost:         &asCost,
		ASMargin:       &asMargin,
		ASMarginPct:    &asMarginPct,
		ASProfitPct:    &asProfitPct,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
}

// CreateTestRegions создает тестовые регионы
func CreateTestRegions() []*model.Region {
	return []*model.Region{
		{
			Code: "MSK",
			Name: "Москва",
		},
		{
			Code: "SPB",
			Name: "Санкт-Петербург",
		},
		{
			Code: "EKB",
			Name: "Екатеринбург",
		},
	}
}

// CreateTestBrands создает тестовые бренды
func CreateTestBrands() []*model.Brand {
	return []*model.Brand{
		{
			Name:     "Foton",
			LogoPath: "/brands/Foton.png",
		},
		{
			Name:     "Dongfeng",
			LogoPath: "/brands/Dongfeng.png",
		},
		{
			Name:     "Shacman",
			LogoPath: "/brands/Shacman.png",
		},
	}
}
