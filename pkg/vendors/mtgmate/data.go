package mtgmate

import "encoding/json"

func unmarshalTableData(data []byte) (tableData, error) {
	var r tableData
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *tableData) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type tableData struct {
	// Title string          `json:"title"`
	Cards []Card          `json:"cards"`
	UUID  map[string]UUID `json:"uuid"`
	// CartData         string          `json:"cart_data"`
	// ChannelID        string          `json:"channelId"`
	// CurrentUserID    interface{}     `json:"currentUserId"`
	// DisplaySetName   bool            `json:"displaySetName"`
	// FeatureFlag      bool            `json:"featureFlag"`
	// PreorderSetCodes []interface{}   `json:"preorderSetCodes"`
}

type Card struct {
	UUID     string   `json:"uuid"`
	Variant  string   `json:"variant"`
	Name     string   `json:"name"`
	Price    string   `json:"price"`
	Set      string   `json:"set"`
	Rarity   string   `json:"rarity"`
	Quantity int64    `json:"quantity"`
	Colour   []string `json:"colour"`
	Type     string   `json:"type"`
	Finish   string   `json:"finish"`
}

type UUID struct {
	UUID      string   `json:"uuid"`
	Name      string   `json:"name"`
	Price     int64    `json:"price"`
	SetName   string   `json:"set_name"`
	SetCode   string   `json:"set_code"`
	Rarity    string   `json:"rarity"`
	Quantity  int64    `json:"quantity"`
	Colour    []string `json:"colour"`
	Type      string   `json:"type"`
	Image     string   `json:"image"`
	LinkPath  string   `json:"link_path"`
	Finish    string   `json:"finish"`
	Condition string   `json:"condition"`
}
