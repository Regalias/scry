package goodgames

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"net/http"

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

	var currResults int64 = 0
	var totalResults int64 = 1

	offerings := []*models.Offering{}

	for currResults < totalResults {

		reqBody, err := buildRequest(cardName, currResults).Marshal()
		if err != nil {
			return nil, fmt.Errorf("failed to build request body: %w", err)
		}

		resp, err := v.client.VisitJson(ctx, SEARCH_URL, http.MethodPost, bytes.NewReader(reqBody), SEARCH_CLIENT_TOKEN)
		if err != nil {
			return nil, fmt.Errorf("remote returned error response: %w", err)
		}

		searchResp, err := UnmarshalSearchResponse(resp)
		if err != nil {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}

		page, err := v.parseProducts(&searchResp, cardName)
		if err != nil {
			return nil, fmt.Errorf("failed to parse products: %w", err)
		}

		offerings = append(offerings, page...)

		totalResults = searchResp.TotalHits
		currResults += int64(len(searchResp.Results))
	}

	return offerings, nil
}

func buildRequest(cardName string, skip int64) *SearchRequest {
	return &SearchRequest{
		Query: cardName,
		Fields: []string{
			"title",
			"image",
			"handle",
			"id",
			"isActive",
			"set",
			"inStockVariants",
		},
		SearchFields: []string{"title"},
		Filter:       "isSearchable = 1 AND isActive = 1",
		Sort: []string{
			"-isActive",
			"st_bestseller_position",
			"-_rank",
		},
		Skip:          skip,
		Count:         500,
		Collection:    "UFAEBTQUSK66VEW5EKLJQYH7",
		FacetCount:    99,
		GroupCount:    -1,
		TypoTolerance: 1,
		TextFacetFilters: TextFacetFilters{
			[]string{"mtg single"},
		},
	}
}
