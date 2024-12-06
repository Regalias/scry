package models

// Summary
type BuylistSummary struct {
	Vendors         []string                  `json:"vendors"`
	SummaryByVendor map[string]*VendorSummary `json:"summaryByVendor"`
}

type VendorSummary struct {
	CardList   []CardSummary     `json:"cardList"`
	TotalQty   int64             `json:"totalQty"`
	TotalCost  int64             `json:"totalCost"`
	Selections ProductSelections `json:"selections"`
}

type CardSummary struct {
	Name string `json:"name"`
	Qty  int64  `json:"qty"`
}
