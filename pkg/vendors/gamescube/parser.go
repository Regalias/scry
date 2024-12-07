package gamescube

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/regalias/scry/pkg/models"
	"github.com/regalias/scry/pkg/scrape"
)

func (v *Vendor) parsePageProducts(ctx context.Context, cardName string, document *goquery.Document) (offerings []*models.Offering, err error) {

	offerings = make([]*models.Offering, 0)
	for _, selection := range document.Find("section.main ul.products > li.product").EachIter() {

		if err := ctx.Err(); err != nil {
			return nil, err
		}

		offering, err := v.parseProduct(selection, cardName)
		if err != nil {
			// don't fail the whole thing on a single parsing error
			v.logger.Error(err.Error())
		}

		if offering != nil {
			offerings = append(offerings, offering)
		}
		// else out of stock or no name match, do nothing
	}
	return offerings, nil
}

func (v *Vendor) parseProduct(s *goquery.Selection, cardName string) (*models.Offering, error) {
	// Check for exact card name match
	// Formats:
	// - Swords to Plowshares
	// - Damn - Foil - Borderless
	// - Island (0022) - Foil Double Sided
	// - Forest (Full art)
	// https://regexr.com/880as
	re, err := regexp.Compile(fmt.Sprintf("(?i)^%s(\\s\\([\\w\\d\\s\\/]+\\))?$", regexp.QuoteMeta(cardName)))
	if err != nil {
		return nil, fmt.Errorf("failed to build regex for name: %w", err)
	}

	productName, exists := scrape.FindAttr(s.Find("h4.name"), "title")
	if !exists {
		return nil, fmt.Errorf("failed to find product name: %s", s.Text())
	}

	parts := strings.Split(productName, "-")

	// Check first part matches name exactly
	namePart := strings.TrimSpace(parts[0])
	if !re.MatchString(namePart) {

		// Try again with ascii encoding
		asciiName, err := scrape.ToAscii(cardName)
		if err != nil {
			return nil, fmt.Errorf("failed to convert cardname to ascii: %w", err)
		}
		if asciiName != cardName {
			re, err := regexp.Compile(fmt.Sprintf("(?i)^%s(\\s\\([\\w\\d\\s\\/]+\\))?$", regexp.QuoteMeta(asciiName)))
			if err != nil {
				return nil, fmt.Errorf("failed to build regex for name: %w", err)
			}

			if !re.MatchString(namePart) {
				v.logger.Debug("not exact match for card", "cardName", cardName, "ascii", asciiName, "productName", productName)
				return nil, nil
			}
		}

		v.logger.Debug("not exact match for card", "cardName", cardName, "productName", productName)
		return nil, nil
	}

	// Parse finish/properties
	properties := []string{}
	if len(parts) > 1 {
		for i, part := range parts {
			if i == 0 {
				continue
			}
			properties = append(properties, strings.Trim(part, " "))
		}
	}

	// Check for out of stock
	qtyText, _ := scrape.FindAttr(s.Find("select.qty"), "max")
	if qtyText == "" {
		// Out of stock
		return nil, nil
	}

	quantity, err := strconv.ParseInt(qtyText, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to convert qty to int64 for '%s'", productName)
	}

	price, exists := scrape.FindAttr(s.Find("form"), "data-price")
	if !exists {
		return nil, fmt.Errorf("failed to find product price")
	}
	price = strings.Trim(price, "AUD$ ")
	price_parts := strings.Split(price, ".")
	if len(price_parts) != 2 {
		return nil, errors.New("failed parsing price: multiple '.' chars found")
	}
	// Truncate off the cents to 2 digits
	cents := price_parts[1]
	if len(cents) > 2 {
		cents = cents[:2]
	}
	price_int, err := strconv.ParseInt(price_parts[0]+cents, 10, 64)
	if err != nil {
		price_int = -1
		return nil, fmt.Errorf("failed to convert price to int64 for '%s'", productName)
	}

	// Other metadata
	set := strings.TrimSpace(s.Find("span.category").Text())
	sku, exists := scrape.FindChildAttr(s, "form", "data-vid")
	if !exists {
		return nil, fmt.Errorf("failed to find product SKU")
	}

	condition, _ := scrape.FindChildAttr(s, "form", "data-variant")

	imgURI, exists := scrape.FindChildAttr(s, "img", "src")
	if !exists {
		return nil, fmt.Errorf("failed to find image URI")
	}

	productURI, exists := scrape.FindChildAttr(s, "a[itemprop=url]", "href")
	if !exists {
		return nil, fmt.Errorf("failed to find product URI")
	}

	offering := &models.Offering{
		Name:       productName,
		Price:      price_int,
		Quantity:   quantity,
		VendorID:   VENDOR_ID,
		Set:        set,
		Condition:  condition,
		Properties: properties,
		ImgURI:     imgURI,
		ProductURI: BASE_URL + productURI,
		StoreSKU:   sku,
		CreatedAt:  time.Now().UnixMilli(),
	}
	return offering, nil
}
