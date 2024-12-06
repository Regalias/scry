package gamescube

import (
	"fmt"
	"net/url"

	"github.com/PuerkitoBio/goquery"
	"github.com/regalias/scry/pkg/scrape"
)

func (v *Vendor) getPageURL(cardName string) (string, error) {
	cardName, err := scrape.ToAscii(cardName)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(BASE_URL+"/advanced_search?utf8=%%E2%%9C%%93&search%%5Bfuzzy_search%%5D=&search%%5Btags_name_eq%%5D=&search%%5Bsell_price_gte%%5D=&search%%5Bsell_price_lte%%5D=&search%%5Bbuy_price_gte%%5D=&search%%5Bbuy_price_lte%%5D=&search%%5Bin_stock%%5D=0&search%%5Bin_stock%%5D=1&buylist_mode=0&search%%5Bcategory_ids_with_descendants%%5D%%5B%%5D=&search%%5Bcategory_ids_with_descendants%%5D%%5B%%5D=8&search%%5Bwith_descriptor_values%%5D%%5B6%%5D=&search%%5Bwith_descriptor_values%%5D%%5B7%%5D=&search%%5Bwith_descriptor_values%%5D%%5B9%%5D=&search%%5Bwith_descriptor_values%%5D%%5B10%%5D=&search%%5Bwith_descriptor_values%%5D%%5B11%%5D=&search%%5Bwith_descriptor_values%%5D%%5B13%%5D=&search%%5Bwith_descriptor_values%%5D%%5B348%%5D=%s&search%%5Bwith_descriptor_values%%5D%%5B361%%5D=&search%%5Bwith_descriptor_values%%5D%%5B1259%%5D=&search%%5Bwith_descriptor_values%%5D%%5B9703%%5D=&search%%5Bwith_descriptor_values%%5D%%5B9713%%5D=&search%%5Bwith_descriptor_values%%5D%%5B9723%%5D=&search%%5Bwith_descriptor_values%%5D%%5B9733%%5D=&search%%5Bwith_descriptor_values%%5D%%5B9743%%5D=&search%%5Bwith_descriptor_values%%5D%%5B9753%%5D=&search%%5Bwith_descriptor_values%%5D%%5B9763%%5D=&search%%5Bwith_descriptor_values%%5D%%5B9773%%5D=&search%%5Bwith_descriptor_values%%5D%%5B9783%%5D=&search%%5Bwith_descriptor_values%%5D%%5B9793%%5D=&search%%5Bwith_descriptor_values%%5D%%5B9803%%5D=&search%%5Bwith_descriptor_values%%5D%%5B9813%%5D=&search%%5Bwith_descriptor_values%%5D%%5B9823%%5D=&search%%5Bwith_descriptor_values%%5D%%5B9833%%5D=&search%%5Bwith_descriptor_values%%5D%%5B9843%%5D=&search%%5Bwith_descriptor_values%%5D%%5B9853%%5D=&search%%5Bwith_descriptor_values%%5D%%5B9863%%5D=&search%%5Bwith_descriptor_values%%5D%%5B9873%%5D=&search%%5Bwith_descriptor_values%%5D%%5B9883%%5D=&search%%5Bwith_descriptor_values%%5D%%5B9893%%5D=&search%%5Bwith_descriptor_values%%5D%%5B10853%%5D=&search%%5Bwith_descriptor_values%%5D%%5B10873%%5D=&search%%5Bwith_descriptor_values%%5D%%5B10893%%5D=&search%%5Bwith_descriptor_values%%5D%%5B10903%%5D=&search%%5Bwith_descriptor_values%%5D%%5B10913%%5D=&search%%5Bwith_descriptor_values%%5D%%5B10923%%5D=&search%%5Bwith_descriptor_values%%5D%%5B10933%%5D=&search%%5Bwith_descriptor_values%%5D%%5B10963%%5D=&search%%5Bwith_descriptor_values%%5D%%5B10973%%5D=&search%%5Bwith_descriptor_values%%5D%%5B10983%%5D=&search%%5Bwith_descriptor_values%%5D%%5B11133%%5D=&search%%5Bwith_descriptor_values%%5D%%5B11153%%5D=&search%%5Bwith_descriptor_values%%5D%%5B11163%%5D=&search%%5Bwith_descriptor_values%%5D%%5B11173%%5D=&search%%5Bwith_descriptor_values%%5D%%5B11183%%5D=&search%%5Bwith_descriptor_values%%5D%%5B11193%%5D=&search%%5Bwith_descriptor_values%%5D%%5B11276%%5D=&search%%5Bwith_descriptor_values%%5D%%5B11282%%5D=&search%%5Bwith_descriptor_values%%5D%%5B11285%%5D=&search%%5Bwith_descriptor_values%%5D%%5B11286%%5D=&search%%5Bwith_descriptor_values%%5D%%5B11287%%5D=&search%%5Bvariants_with_identifier%%5D%%5B14%%5D%%5B%%5D=&search%%5Bvariants_with_identifier%%5D%%5B15%%5D%%5B%%5D=&search%%5Bsort%%5D=name&search%%5Bdirection%%5D=ascend&commit=Search&search%%5Bcatalog_group_id_eq%%5D=", url.QueryEscape(cardName)), nil
}

func (v *Vendor) parseNextPageURL(document *goquery.Document) (url string, hasNextPage bool) {
	for _, s := range document.Find("a.next_page").EachIter() {
		// Should only be one result
		url, _ = scrape.FindAttr(s, "href")
		if url != "" {
			url = BASE_URL + url
			hasNextPage = true
		}
		break
	}
	return url, hasNextPage
}
