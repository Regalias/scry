package models

// type Model struct {
// 	ID        uint `gorm:"primarykey"`
// 	CreatedAt time.Time
// 	UpdatedAt time.Time
// }

type Buylist struct {
	ID        int64  `json:"id" db:"id"`
	Name      string `json:"name" db:"name"`
	CreatedAt int64  `json:"createdAt" db:"created_at"` // in ms

	Cards []Card `json:"cards" db:"-"`

	// calculated values
	TotalPrice      int64 `json:"totalPrice" db:"-"`
	TotalSelections int64 `json:"totalSelections" db:"-"`
	TotalCards      int64 `json:"totalCards" db:"-"`
}

type Card struct {
	ID       int64  `json:"id" db:"id"`
	Name     string `json:"name" db:"name"`
	Quantity int64  `json:"quantity" db:"quantity"`

	Selections ProductSelections `json:"selections" db:"-"`

	// calculated values
	TotalSelectionPrice int64 `json:"totalPrice" db:"-"`
	TotalSelections     int64 `json:"totalSelections" db:"-"`

	BuylistID int64 `json:"-" db:"buylist_id"` // FK
}

func (c *Card) GetPurchasedCount() (purchased int64) {
	for _, sku := range c.Selections {
		if sku.IsPurchased {
			purchased += sku.Quantity
		}
	}
	return purchased
}

type ProductSelections []ProductSelection

type ProductSelection struct {
	ID          int64    `json:"id" db:"id"`
	Quantity    int64    `json:"quantity" db:"quantity"` // User selected quantity
	Offering    Offering `json:"offering" db:"offering"`
	IsPurchased bool     `json:"isPurchased" db:"is_purchased"`
	IsFlagged   bool     `json:"isFlagged" db:"is_flagged"`

	CardID int64 `json:"-" db:"card_id"` // FK
}
