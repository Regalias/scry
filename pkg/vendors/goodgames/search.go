package goodgames

import (
	"encoding/json"
)

// Request model

func UnmarshalSearchRequest(data []byte) (SearchRequest, error) {
	var r SearchRequest
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *SearchRequest) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type SearchRequest struct {
	Query            string           `json:"query"`
	Fields           []string         `json:"fields"`
	SearchFields     []string         `json:"searchFields"`
	Filter           string           `json:"filter"`
	Sort             []string         `json:"sort"`
	Skip             int64            `json:"skip"`
	Count            int64            `json:"count"`
	Collection       string           `json:"collection"`
	FacetCount       int64            `json:"facetCount"`
	GroupCount       int64            `json:"groupCount"`
	TypoTolerance    int64            `json:"typoTolerance"`
	TextFacetFilters TextFacetFilters `json:"textFacetFilters"`
}

type TextFacetFilters struct {
	ProductType []string `json:"product_type"`
}

// Response model

func UnmarshalSearchResponse(data []byte) (SearchResponse, error) {
	var r SearchResponse
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *SearchResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type SearchResponse struct {
	Results   []Result `json:"results"`
	TotalHits int64    `json:"totalHits"`

	// Query        Query    `json:"query"`
	// ResponseTime int64    `json:"responseTime"`
	// UniqueID     string   `json:"uniqueId"`
}

// type Query struct {
// 	Query            string           `json:"query"`
// 	Fields           []string         `json:"fields"`
// 	SearchFields     []string         `json:"searchFields"`
// 	Filter           string           `json:"filter"`
// 	Sort             []string         `json:"sort"`
// 	Skip             int64            `json:"skip"`
// 	Collection       string           `json:"collection"`
// 	FacetCount       int64            `json:"facetCount"`
// 	Count            int64            `json:"count"`
// 	GroupCount       int64            `json:"groupCount"`
// 	TypoTolerance    int64            `json:"typoTolerance"`
// 	TextFacetFilters TextFacetFilters `json:"textFacetFilters"`
// }

type Result struct {
	ID              string           `json:"id"`
	Title           string           `json:"title"`  // product name
	Handle          string           `json:"handle"` // used to build product URI slug
	InStockVariants []InStockVariant `json:"inStockVariants"`
	Image           Image            `json:"image"`
	Set             []string         `json:"set"`

	// Rank            int64            `json:"_rank"`
	// IsActive        int64            `json:"isActive"`
}

type Image struct {
	Src string `json:"src"` // product img href

	// ID int64 `json:"id"`
	// Width int64 `json:"width"`
	// Height int64 `json:"height"`
	// UpdatedAt         time.Time `json:"updated_at"`
	// ProductID         int64     `json:"product_id"`
	// AdminGraphqlAPIID string    `json:"admin_graphql_api_id"`
	// CreatedAt         time.Time `json:"created_at"`
	// Position          int64     `json:"position"`
}

type InStockVariant struct {
	ID                int64  `json:"id"`                 // SKU
	Title             string `json:"title"`              // Condition + Variant
	InventoryQuantity int64  `json:"inventory_quantity"` // qty
	Price             string `json:"price"`              // e.g. 8.80

	// InventoryManagement  string      `json:"inventory_management"`
	// RequiresShipping     bool        `json:"requires_shipping"`
	// OldInventoryQuantity int64       `json:"old_inventory_quantity"`
	// CreatedAt            time.Time   `json:"created_at"`

	// UpdatedAt            time.Time   `json:"updated_at"`

	// InventoryItemID int64 `json:"inventory_item_id"`
	// ProductID int64 `json:"product_id"`
	// Option1              string      `json:"option1"` // Condition + Variant

	// Grams                int64       `json:"grams"`
	// Sku string `json:"sku"` // Not the SKU we care about
	// Barcode              *string     `json:"barcode"`

	// Image                interface{} `json:"image"`
	// CompareAtPrice       interface{} `json:"compare_at_price"`
	// Taxable              bool        `json:"taxable"`
	// FulfillmentService   string      `json:"fulfillment_service"`
	// Weight               int64       `json:"weight"`
	// InventoryPolicy      string      `json:"inventory_policy"`
	// WeightUnit           string      `json:"weight_unit"`
	// AdminGraphqlAPIID    string      `json:"admin_graphql_api_id"`
	// Position             int64       `json:"position"`
	// ImageID              interface{} `json:"image_id"`
}
