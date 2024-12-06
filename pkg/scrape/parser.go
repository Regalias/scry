package scrape

import (
	"context"

	"github.com/regalias/scry/pkg/models"
)

type VendorParser interface {
	ScrapeCard(ctx context.Context, cardName string) ([]*models.Offering, error)
}
