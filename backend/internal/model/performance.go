package model

// Performance отвечает за блок Performance.
type Performance struct {
	DealerName     string `json:"dealer_name"`
	Region         string `json:"region"`
	City           string `json:"city"`
	Manager        string `json:"manager"`
	SalesRevenues  int32  `json:"sales_revenues"` // in rub absolute
	SalesProfits   int32  `json:"sales_profits"`  // in %
	SalesMargin    int32  `json:"sales_margin"`   // in %
	ASRevenues     int32  `json:"as_revenues"`    // in rub absolute
	ASProfits      int32  `json:"as_profits"`     // in rub absolute
	ASMargin       int32  `json:"as_margin"`      // in %
	FotonRank      int16  `json:"foton_rank"`     // from 1 to 10
	Recommendation string `json:"dealer_dev_recommendation"`
}
