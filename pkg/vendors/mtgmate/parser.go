package mtgmate

import (
	"fmt"
	"regexp"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/regalias/scry/pkg/models"
	"github.com/regalias/scry/pkg/scrape"
)

func (v *Vendor) parsePageProducts(cardName string, document *goquery.Document) (offerings []*models.Offering, err error) {

	elem := document.Find("div[data-react-class=FilterableTable]").First()
	if elem == nil {
		return nil, fmt.Errorf("failed to find table element")
	}
	value, _ := scrape.FindAttr(elem, "data-react-props")
	if value == "" {
		return nil, fmt.Errorf("failed to find data-react-props attribute")
	}

	data, err := unmarshalTableData([]byte(value))
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshall table data: %w", err)
	}

	return v.getOfferingsFromTableData(cardName, &data)
}

func (v *Vendor) getOfferingsFromTableData(cardName string, tableData *tableData) ([]*models.Offering, error) {

	offerings := make([]*models.Offering, 0)

	for i := range tableData.Cards {
		card := tableData.Cards[i]
		data := tableData.UUID[card.UUID]

		// Filter name
		match, err := v.doesNameMatch(cardName, data.Name)
		if err != nil {
			return nil, fmt.Errorf("doesNameMatch failed: %w", err)
		}
		if !match {
			continue
		}

		// Filter out of stock
		if data.Quantity < 1 {
			continue
		}

		properties := []string{}
		if data.Finish != "Nonfoil" {
			properties = append(properties, data.Finish)
		}

		offerings = append(offerings, &models.Offering{
			Name:       data.Name,
			Quantity:   data.Quantity,
			Set:        data.SetName,
			Condition:  data.Condition,
			Properties: properties,
			ImgURI:     data.Image,
			ProductURI: getProductURL(data.LinkPath),
			StoreSKU:   card.UUID,
			VendorID:   VENDOR_ID,
			CreatedAt:  time.Now().UnixMilli(),
			Price:      data.Price,
		})

	}

	return offerings, nil
}

const nameRegex = "(?i)^(?:(?:[\\w\\d\\s]+ \\/\\/ )?%s)(?:\\s\\([\\w\\d\\s\\/-]+\\))?$"

func (v *Vendor) doesNameMatch(cardName string, productName string) (bool, error) {

	// Valid formats
	// Forest // Forest
	// Forest (Full Art 280)
	// Forest
	// Forest (Borderless 231)

	reg := fmt.Sprintf(nameRegex, regexp.QuoteMeta(cardName))
	re, err := regexp.Compile(reg)
	if err != nil {
		return false, fmt.Errorf("failed to build regex for name: %w", err)
	}

	if !re.MatchString(productName) {

		// Try again with ascii encoding
		asciiName, err := scrape.ToAscii(cardName)
		if err != nil {
			return false, fmt.Errorf("failed converting to ascii repr: %w", err)
		}

		if asciiName != cardName {
			re, err := regexp.Compile(fmt.Sprintf(nameRegex, regexp.QuoteMeta(asciiName)))
			if err != nil {
				return false, fmt.Errorf("failed to build regex for name: %w", err)
			}

			if !re.MatchString(productName) {
				v.logger.Debug("not exact match for card", "cardName", cardName, "ascii", asciiName, "productName", productName)
				return false, nil
			}
		}

		v.logger.Debug("not exact match for card", "cardName", cardName, "productName", productName)
		return false, nil
	}

	return true, nil
}
