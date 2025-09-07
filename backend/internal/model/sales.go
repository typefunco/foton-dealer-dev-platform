package model

// Sales отвечает за блок Dealer Sales.
type Sales struct {
	DealerName            string   `json:"dealer_name"`
	Region                string   `json:"region"`
	City                  string   `json:"city"`
	Manager               string   `json:"manager"`
	StockHdtMdtLdt        [3]int16 `json:"stock_hdt_mdt_ldt"`
	BuyoutHdtMdtLdt       [3]int16 `json:"buyout_hdt_mdt_ldt"`
	SalesTarget           int32    `json:"sales_target"`
	FotonSalesmen         int16    `json:"foton_salesmen"`
	SalesTrainings        string   `json:"sales_trainings"`
	ServiceContractsSales int16    `json:"service_contracts_sales"`
	Recommendation        string   `json:"sales_recommendation"`
}
