package model

// AfterSales отвечает за блок AfterSales.
type AfterSales struct {
	DealerName       string `json:"dealer_name"`
	Region           string `json:"region"`
	City             string `json:"city"`
	Manager          string `json:"manager"`
	RecommendedStock int16  `json:"recommended_stock"` // in %
	WarrantyStock    int16  `json:"warranty_stock"`    // in %
	FotonLaborHours  int16  `json:"foton_labor_hours"` // in %
	ServiceContracts int16  `json:"service_contracts"` // absolute count
	ASTrainings      string `json:"as_trainings"`
	CSI              int16  `json:"marketing_investments"` // in %
	Recommendation   string `json:"as_recommendation"`
}
