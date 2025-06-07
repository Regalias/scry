package gamescube

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

	url, err := v.getPageURL(cardName)
	if err != nil {
		return nil, fmt.Errorf("failed getting page url: %w", err)
	}

	nextPage := url
	vendorProducts := []*models.Offering{}
	addedProductURIs := map[string]struct{}{}

	for nextPage != "" {
		if err := ctx.Err(); err != nil {
			return nil, fmt.Errorf("cancelling scrape: %w", err)
		}

		doc, err := v.client.Visit(ctx, nextPage)
		if err != nil {
			return nil, fmt.Errorf("failed visiting page: %w", err)
		}
		products, err := v.parsePageProducts(ctx, cardName, doc)
		if err != nil {
			return nil, fmt.Errorf("failed to parse products: %w", err)
		}

		for i := range products {
			// Add product only if it doesn't already exist
			// Needed because gamescube can return duplicate products between pages
			// which will break frontend rendering keys
			if _, exists := addedProductURIs[products[i].ProductURI]; !exists {
				addedProductURIs[products[i].ProductURI] = struct{}{}
				vendorProducts = append(vendorProducts, &products[i])
			}
		}
		nextPage, _ = v.parseNextPageURL(doc)
	}

	return vendorProducts, nil
}
