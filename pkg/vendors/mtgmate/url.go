package mtgmate

import (
	"fmt"
	"net/url"

	"github.com/regalias/scry/pkg/scrape"
)

func getPageURL(cardName string) (string, error) {
	cardName, err := scrape.ToAscii(cardName)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(BASE_URL+"/cards/search?q=%s", url.QueryEscape(cardName)), nil
}

func getProductURL(productPath string) string {
	return BASE_URL + productPath
}
