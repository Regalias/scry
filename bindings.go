package main

import (
	"context"
	"time"

	"github.com/regalias/scry/pkg/buylist"
	"github.com/regalias/scry/pkg/models"
)

const DEFAULT_TIMEOUT = time.Second * 15

// =============================
// Buylist operations
// =============================

func (a *App) GetBuylistSummary(buylistId int64) (*models.BuylistSummary, error) {
	ctx, cancel := context.WithTimeout(a.ctx, DEFAULT_TIMEOUT)
	defer cancel()

	// Fetch buylist
	buylist, err := a.buylists.GetBuylist(ctx, buylistId)
	if err != nil {
		return nil, err
	}

	summary := &models.BuylistSummary{
		Vendors:         []string{},
		SummaryByVendor: map[string]*models.VendorSummary{},
	}

	// Break down selections by vendor
	for _, card := range buylist.Cards {
		for _, sel := range card.Selections {

			vendorId := sel.Offering.VendorID

			// Init stuff
			if summary.SummaryByVendor[vendorId] == nil {
				summary.SummaryByVendor[vendorId] = &models.VendorSummary{
					CardList:   []models.CardSummary{},
					Selections: models.ProductSelections{},
				}
			}
			vendorSummary := summary.SummaryByVendor[vendorId]

			vendorSummary.TotalCost += sel.Offering.Price
			vendorSummary.TotalQty += sel.Quantity

			vendorSummary.Selections = append(vendorSummary.Selections, sel)
			vendorSummary.CardList = append(vendorSummary.CardList, models.CardSummary{
				Name: card.Name,
				Qty:  sel.Quantity,
			})
		}
	}

	// Get existing vendor keys
	keys := make([]string, 0, len(summary.SummaryByVendor))
	for k := range summary.SummaryByVendor {
		keys = append(keys, k)
	}
	summary.Vendors = keys

	return summary, nil
}

func (a *App) NewBuylist(name string) (*models.Buylist, error) {
	ctx, cancel := context.WithTimeout(a.ctx, DEFAULT_TIMEOUT)
	defer cancel()
	return a.buylists.NewBuylist(ctx, name)
}

func (a *App) GetBuylist(id int64) (*models.Buylist, error) {
	ctx, cancel := context.WithTimeout(a.ctx, DEFAULT_TIMEOUT)
	defer cancel()
	return a.buylists.GetBuylist(ctx, id)
}

func (a *App) ListBuylists() ([]models.Buylist, error) {
	ctx, cancel := context.WithTimeout(a.ctx, DEFAULT_TIMEOUT)
	defer cancel()
	return a.buylists.ListBuylists(ctx)
}

func (a *App) DeleteBuylist(id int64) error {
	ctx, cancel := context.WithTimeout(a.ctx, DEFAULT_TIMEOUT)
	defer cancel()
	return a.buylists.DeleteBuylist(ctx, id)
}

func (a *App) UpdateBuylistName(id int64, name string) error {
	ctx, cancel := context.WithTimeout(a.ctx, DEFAULT_TIMEOUT)
	defer cancel()
	return a.buylists.UpdateBuylistName(ctx, id, name)
}

// =============================
// Card operations
// =============================

func (a *App) GetCardsForBuylist(buylistId int64) ([]models.Card, error) {
	ctx, cancel := context.WithTimeout(a.ctx, DEFAULT_TIMEOUT)
	defer cancel()
	return a.buylists.GetCardsForBuylist(ctx, buylistId)
}

func (a *App) AddCardsToBuylist(buylistId int64, cards []*buylist.AddCardsRequest) error {
	ctx, cancel := context.WithTimeout(a.ctx, DEFAULT_TIMEOUT)
	defer cancel()
	return a.buylists.AddCardsToBuylist(ctx, buylistId, cards)
}

func (a *App) DeleteCardsForBuylist(buylistId int64) error {
	ctx, cancel := context.WithTimeout(a.ctx, DEFAULT_TIMEOUT)
	defer cancel()
	return a.buylists.DeleteCardsForBuylist(ctx, nil, buylistId)
}

func (a *App) DeleteCards(cardIds []int64) error {
	ctx, cancel := context.WithTimeout(a.ctx, DEFAULT_TIMEOUT)
	defer cancel()
	return a.buylists.DeleteCards(ctx, cardIds)
}

func (a *App) UpdateCardQty(cardId int64, quantity int64) error {
	ctx, cancel := context.WithTimeout(a.ctx, DEFAULT_TIMEOUT)
	defer cancel()
	return a.buylists.UpdateCardQty(ctx, cardId, quantity)
}

// =============================
// Selection Opertions
// =============================

func (a *App) GetSelections(cardId int64) ([]models.ProductSelection, error) {
	ctx, cancel := context.WithTimeout(a.ctx, DEFAULT_TIMEOUT)
	defer cancel()
	return a.buylists.GetSelections(ctx, cardId)
}

func (a *App) AddSelection(cardId int64, offering *models.Offering, quantity int64) (*models.ProductSelection, error) {
	ctx, cancel := context.WithTimeout(a.ctx, DEFAULT_TIMEOUT)
	defer cancel()
	return a.buylists.AddSelection(ctx, cardId, offering, quantity)
}

func (a *App) DeleteSelection(selectionId int64) error {
	ctx, cancel := context.WithTimeout(a.ctx, DEFAULT_TIMEOUT)
	defer cancel()
	return a.buylists.DeleteSelection(ctx, selectionId)
}

func (a *App) DeleteSelectionsForCardId(cardId int64) error {
	ctx, cancel := context.WithTimeout(a.ctx, DEFAULT_TIMEOUT)
	defer cancel()
	return a.buylists.DeleteSelectionsForCardId(ctx, cardId)
}

func (a *App) UpdateSelection(req *buylist.UpdateSelectionRequest) error {
	ctx, cancel := context.WithTimeout(a.ctx, DEFAULT_TIMEOUT)
	defer cancel()
	return a.buylists.UpdateSelection(ctx, req)
}

// =============================
// Scraper operations
// =============================

func (a *App) GetCardOfferingForVendor(cardName string, vendor string) ([]*models.Offering, error) {
	ctx, cancel := context.WithCancel(a.ctx)
	defer cancel()
	return a.vendors.ScrapeProductForVendor(ctx, cardName, vendor)
}

func (a *App) ListVendors() []string {
	return a.vendors.ListVendors()
}
