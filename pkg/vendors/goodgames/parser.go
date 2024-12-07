package goodgames

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/regalias/scry/pkg/models"
	"github.com/regalias/scry/pkg/scrape"
)

const productNameRegex = `(?i)^%s(?: \(([\w\s]+?)\))?(?: \[.+\])?$`

var etchedRegex = regexp.MustCompile(`(?i)^(?:.+?\s)?Etched(?:\s.+?)?$`)

func (v *Vendor) parseProducts(data *SearchResponse, cardName string) ([]*models.Offering, error) {
	products := []*models.Offering{}

	for _, result := range data.Results {

		// Filter name
		match, isEtched, err := v.doesNameMatch(cardName, result.Title)
		if err != nil {
			v.logger.Error("name match failure", "err", err, "productName", result.Title)
			continue
		}
		if !match {
			continue
		}

		for _, variant := range result.InStockVariants {
			// Shouldn't happen with our filters, but double check product is in stock
			if variant.InventoryQuantity < 1 {
				continue
			}

			// Parse price to cents in int
			price, err := scrape.ParsePrice(variant.Price)
			if err != nil {
				v.logger.Error("failed parsing price", "err", err, "variant", variant)
				continue
			}

			set := ""
			if result.Set == nil || len(result.Set) < 1 {
				// Sometimes goodgames doesn't specify a set for the card
				v.logger.Debug("result is missing set information", "result", result)
			} else {
				set = result.Set[0]
			}

			// Determine foil
			// e.g. Near Mint Foil
			condition, isFoil := strings.CutSuffix(variant.Title, " Foil")
			properties := []string{}
			if isFoil {
				properties = append(properties, "Foil")
			}
			if isEtched {
				properties = append(properties, "Etched")
			}

			variantSku := strconv.FormatInt(variant.ID, 10)

			products = append(products, &models.Offering{
				Name:       result.Title,
				Price:      price,
				Quantity:   variant.InventoryQuantity,
				Set:        set,
				Condition:  condition,
				Properties: properties,
				ImgURI:     result.Image.Src,
				StoreSKU:   variantSku,
				VendorID:   VENDOR_ID,
				CreatedAt:  time.Now().UnixMilli(),
				// https://tcg.goodgames.com.au/products/damnation-borderless-alternate-art-double-masters-2022?variant=41720986599603
				ProductURI: BASE_URL + PRODUCT_PATH + result.Handle + "?variant=" + variantSku,
			})
		}
	}

	return products, nil
}

func (v *Vendor) doesNameMatch(cardName string, productName string) (matched bool, isEtched bool, err error) {

	matched = false
	isEtched = false

	toMatch := []string{cardName}

	// Check if utf chars in string
	asciiName, err := scrape.ToAscii(cardName)
	if err != nil {
		return false, false, fmt.Errorf("failed converting to ascii: %w", err)
	}

	if asciiName != cardName {
		toMatch = append(toMatch, asciiName)
	}

	var matches []string
	for _, name := range toMatch {
		re, err := regexp.Compile(fmt.Sprintf(productNameRegex, regexp.QuoteMeta(name)))
		if err != nil {
			return false, false, fmt.Errorf("failed to compile regex: %w", err)
		}
		matches = re.FindStringSubmatch(productName)
		if matches != nil {
			matched = true
			break
		}
	}

	if !matched {
		v.logger.Debug("name mismatch", "cardName", toMatch, "productName", productName)
		return false, false, nil
	}

	if len(matches) != 2 {
		v.logger.Error("matches did not return expected length", "matches", matches)
	}

	// Check for etched because goodgames doesn't have an explicit property for it
	// e.g. Arid Mesa (Retro Foil Etched) [Modern Horizons 2]
	props := matches[1]
	if etchedRegex.MatchString(props) {
		isEtched = true
	}

	return true, isEtched, nil
}
