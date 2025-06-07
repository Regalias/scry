package gamescube

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/regalias/scry/pkg/models"
	"github.com/regalias/scry/pkg/scrape"
)

func (v *Vendor) parsePageProducts(ctx context.Context, cardName string, document *goquery.Document) (offerings []models.Offering, err error) {

	offerings = make([]models.Offering, 0)
	for _, selection := range document.Find("section.main ul.products > li.product").EachIter() {

		if err := ctx.Err(); err != nil {
			return nil, err
		}

		parsed, err := v.parseProducts(selection, cardName)
		if err != nil {
			// don't fail the whole thing on a single parsing error
			v.logger.With("cardName", cardName).Error(err.Error())
		}

		if offerings != nil {
			offerings = append(offerings, parsed...)
		}
		// else out of stock or no name match, do nothing
	}
	return offerings, nil
}

func (v *Vendor) parseProducts(s *goquery.Selection, cardName string) ([]models.Offering, error) {
	s = s.Find("div.inner")

	// Parse common metadata
	metadataNode := s.Find("div.image-meta > div.image > a")
	productURI, found := scrape.FindAttr(metadataNode, "href")
	if !found {
		return nil, fmt.Errorf("failed to find product URI in: %s", metadataNode.Text())
	}

	productName, found := scrape.FindAttr(metadataNode, "title")
	if !found {
		return nil, fmt.Errorf("failed to find product name in: %s", metadataNode.Text())
	}

	imgURI, exists := scrape.FindAttr(metadataNode.Find("img"), "src")
	if !exists {
		return nil, fmt.Errorf("failed to find image URI in: %s", metadataNode.Text())
	}

	// Check for exact card name match
	match, err := checkCardNameMatch(productName, cardName)
	if err != nil {
		return nil, err
	}

	if !match {
		v.logger.Debug("not exact match for card", "cardName", cardName, "productName", productName)
		return nil, nil
	}

	// Parse finish/properties
	propertyParts := strings.Split(productName[len(cardName):], "-")
	properties := []string{}
	for _, part := range propertyParts {
		part = strings.Trim(part, " ")
		if part != "" {
			properties = append(properties, part)
		}
	}

	offerings := []models.Offering{}
	// Parse variations
	for _, selection := range s.Find("div.variants > div.variant-row form").EachIter() {

		// QTY
		qtyText, _ := scrape.FindAttr(selection.Find("input.qty"), "max")
		if qtyText == "" || qtyText == "0" {
			// Out of stock
			continue
		}
		quantity, err := strconv.ParseInt(qtyText, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to convert qty to int64 for '%s'", productName)
		}

		// Price
		price, exists := scrape.FindAttr(selection, "data-price")
		if !exists {
			return nil, fmt.Errorf("failed to find product price")
		}
		price = strings.Trim(price, "AUD$ ")
		price_int, err := scrape.ParsePrice(price)
		if err != nil {
			return nil, fmt.Errorf("failed parsing price: %w", err)
		}

		// Other metadata
		set, _ := scrape.FindAttr(selection, "data-category")
		set = strings.TrimSpace(set)

		sku, exists := scrape.FindAttr(selection, "data-vid")
		if !exists {
			return nil, fmt.Errorf("failed to find product SKU")
		}

		condition, _ := scrape.FindAttr(selection, "data-variant")

		variant_name, _ := scrape.FindAttr(selection, "data-name")
		if variant_name != productName {
			v.logger.Warn("variant name is different to product name!", "variant_name", variant_name, "productName", productName)
		}

		offerings = append(offerings, models.Offering{
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
		})
	}

	return offerings, nil
}

// Checks if a product name matches the card name exactly, minus printing variations
// This function is not perfect and in rare cases may return true if the actual product card name contains parenthesis or dashes in a way that is identical to printing variation formats
func checkCardNameMatch(productName string, cardName string) (bool, error) {

	// Convert representation to base ASCII cases
	cardName, err := canonicalizeName(cardName)
	if err != nil {
		return false, fmt.Errorf("failed to canonicalize cardName: %w", err)
	}
	productName, err = canonicalizeName(productName)
	if err != nil {
		return false, fmt.Errorf("failed to canonicalize productName: %w", err)
	}

	// Match character by character from the start, because there's no well defined set of rules for product names
	// Another way to do this is check strings.Index(productName, cardName) == 0,
	// check len(productName) > or == len(cardName), then index directly on productName[len(cardName)]
	if len(cardName) > len(productName) {
		return false, nil
	}
	curr := 0
	for curr < len(cardName) {
		if cardName[curr] != productName[curr] {
			return false, nil
		}
		curr += 1
	}

	if len(productName) == curr {
		// Exact length match, strings are identical
		return true, nil
	}

	if len(productName) > curr+1 {
		// at least 2 extra characters present, check if it's a space followed by one of the known delimiters for printing variations
		if productName[curr] == ' ' {
			switch productName[curr+1] {
			case '-':
				fallthrough
			case '(':
				return true, nil
			default:
				return false, nil
			}
		}
	}

	return false, nil
}

// Get canonicalized name for a card in lowercase and ASCII representation
func canonicalizeName(name string) (string, error) {
	name, err := scrape.ToAscii(name)
	if err != nil {
		return "", fmt.Errorf("failed to convert cardname to ascii: %w", err)
	}
	name = strings.ToLower(name)
	return name, nil
}
