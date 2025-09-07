package model

// DealerDev отвечает за блок Dealer Development.
type DealerDev struct {
	DealerName           string   `json:"dealer_name"`
	Region               string   `json:"region"`
	City                 string   `json:"city"`
	Manager              string   `json:"manager"`
	Brands               []string `json:"brands"`
	CheckListScore       int16    `json:"check_list_score"`
	DealerShipClass      string   `json:"dealer_ship_class"`
	Branding             string   `json:"branding"`
	MarketingInvestments int32    `json:"marketing_investments"`
	BySideBusinesses     []string `json:"by_side_businesses"`
	Recommendation       string   `json:"dealer_dev_recommendation"`
}
