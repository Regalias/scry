package vendors

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/regalias/scry/pkg/models"
	"github.com/regalias/scry/pkg/scrape"
	"github.com/regalias/scry/pkg/vendors/gamescube"
	"github.com/regalias/scry/pkg/vendors/mtgmate"
)

type Manager struct {
	logger   *slog.Logger
	client   *scrape.ScrapeClient
	scrapers map[string]scrape.VendorParser
	vendors  []string
}

func NewManager(logger *slog.Logger) (*Manager, error) {

	client, err := scrape.NewScrapeClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create scrape client: %w", err)
	}

	scrapers := map[string]scrape.VendorParser{
		gamescube.VENDOR_ID: gamescube.New(logger, client),
		mtgmate.VENDOR_ID:   mtgmate.New(logger, client),
	}

	vendors := make([]string, 0, len(scrapers))
	for k := range scrapers {
		vendors = append(vendors, k)
	}

	return &Manager{
		vendors:  vendors,
		scrapers: scrapers,
		logger:   logger,
		client:   client,
	}, nil
}

func (m *Manager) ListVendors() []string {
	return m.vendors
}

func (m *Manager) ScrapeProductForVendor(ctx context.Context, cardName string, vendor string) ([]*models.Offering, error) {

	parser, exists := m.scrapers[vendor]
	if !exists {
		return nil, fmt.Errorf("no parser for vendor: %s", vendor)
	}
	return parser.ScrapeCard(ctx, cardName)
}
