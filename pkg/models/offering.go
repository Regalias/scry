package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type Offering struct {
	Name       string   `json:"name"`
	Price      int64    `json:"price"`    // in cents
	Quantity   int64    `json:"quantity"` // Vendor quantity available
	Set        string   `json:"set"`
	Condition  string   `json:"condition"`
	Properties []string `json:"properties"` // Foil, borderless, showcase, etc
	ImgURI     string   `json:"imgUri"`     // Image URI to load inline
	ProductURI string   `json:"productUri"` // Direct link to item on store
	StoreSKU   string   `json:"storeSKU"`   // For cart automation purposes
	VendorID   string   `json:"vendorId"`   // Vendor name
	CreatedAt  int64    `json:"createdAt"`  // Unix ms epoch, ms because JS interop
}

func (o *Offering) Scan(src any) error {
	switch src := src.(type) {
	case []byte:
		return json.Unmarshal(src, o)
	default:
		return errors.New("unsupported type")
	}
}

func (o *Offering) Value() (driver.Value, error) {
	return json.Marshal(o)
}
