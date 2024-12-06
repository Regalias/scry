package mtgmate

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/regalias/scry/pkg/models"
	"github.com/regalias/scry/pkg/scrape"
)

type Vendor struct {
	logger *slog.Logger
	client *scrape.ScrapeClient
}

func New(logger *slog.Logger, client *scrape.ScrapeClient) *Vendor {
	return &Vendor{
		logger: logger.With("vendor", VENDOR_ID),
		client: client,
	}
}

func (v *Vendor) ScrapeCard(ctx context.Context, cardName string) ([]*models.Offering, error) {

	url, err := getPageURL(cardName)
	if err != nil {
		return nil, fmt.Errorf("failed getting page url: %w", err)
	}

	vendorProducts := []*models.Offering{}

	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("cancelling scrape: %w", err)
	}

	doc, err := v.client.Visit(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("failed visiting page: %w", err)
	}
	products, err := v.parsePageProducts(cardName, doc)
	if err != nil {
		return nil, fmt.Errorf("failed to parse products: %w", err)
	}
	vendorProducts = append(vendorProducts, products...)

	return vendorProducts, nil
}
